package main

import (
	"github.com/bosley/nerv-go"
	"internal/reaper"
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
		}, ","))

	return nerv.Consumer{
		Id: consumerId,
		Fn: func(event *nerv.Event) {
			slog.Debug("logging-consumer", "id", consumerId)
			action(event)
		},
	}
}

func buildModReaper(config reaper.Config) ModuleData {

	channel := "shutdown"

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
		Name:      "reaper",
		Mod:       mod,
		Consumers: consumers,
		Topic:     topic,
	}
}
