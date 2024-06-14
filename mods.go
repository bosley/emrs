/*
  nerv-go modules supporting EMRS
*/

package main

import (
	"github.com/bosley/nerv-go"
	"internal/reaper"
	"internal/webui"
	"log/slog"
	"strings"
)

const (
	RootAppName         = "emrs"
	InternalTopicPrefix = "internal"
	TopicGroupWatchdog  = "watchdog"
	TopicGroupProducer  = "producer"
)

type ModuleData struct {
	Name      string
	Mod       nerv.Module
	Consumers []nerv.Consumer
	Topic     *nerv.TopicCfg
}

func formatGroup(name string) string {
	return strings.Join([]string{
		RootAppName,
		InternalTopicPrefix,
		name,
	}, ".")
}

func formatWatchdog(name string) string {
	return strings.Join([]string{
		formatGroup(TopicGroupWatchdog),
		name,
	}, ".")
}

func formatProducer(name string) string {
	return strings.Join([]string{
		formatGroup(TopicGroupProducer),
		name,
	}, ".")
}

func PopulateModules(engine *nerv.Engine, config *AppConfig) error {
	slog.Debug("populating nerv modules")
	for idx, mod := range []ModuleData{
		buildModReaper(config.Reaper),
		buildModWebUi(config.WebUi),
	} {
		if err := engine.UseModule(
			mod.Mod,
			mod.Topic,
			mod.Consumers); err != nil {
			return err
		}
		slog.Debug("loaded module", "name", mod.Name, "idx", idx)
	}
	return nil
}

func buildLoggingConsumer(channelName string, action nerv.EventRecvr) nerv.Consumer {

	consumerId := formatWatchdog(
		strings.Join([]string{
			channelName,
			"logger",
		}, "."))

	return nerv.Consumer{
		Id: consumerId,
		Fn: func(event *nerv.Event) {
			slog.Debug("logging-consumer", "id", consumerId)
			action(event)
		},
	}
}

func buildModWebUi(config webui.Config) ModuleData {

	moduleName := "webui"

	channel := strings.Join([]string{moduleName, "command"}, ".")

	publishingTopic := formatGroup(channel)

	mod := webui.New(config)

	topic := nerv.NewTopic(publishingTopic).
		UsingBroadcast()

	consumers := []nerv.Consumer{
		nerv.Consumer{
			Id: formatWatchdog(moduleName),
			Fn: mod.ReceiveEvent,
		},
	}

	return ModuleData{
		Name:      moduleName,
		Mod:       mod,
		Consumers: consumers,
		Topic:     topic,
	}
}

func buildModReaper(config reaper.Config) ModuleData {

	moduleName := "reaper"

	channel := strings.Join([]string{moduleName, "kill"}, ".")

	publishingTopic := formatGroup(channel)

	mod := reaper.New(config)

	topic := nerv.NewTopic(publishingTopic).
		UsingBroadcast()

	consumers := []nerv.Consumer{
		buildLoggingConsumer(channel, func(event *nerv.Event) {
			slog.Debug("shutdown imminent",
				"sec-remaining",
				event.Data.(*reaper.ReaperMsg).SecondsRemaining)
		}),
	}

	return ModuleData{
		Name:      moduleName,
		Mod:       mod,
		Consumers: consumers,
		Topic:     topic,
	}
}
