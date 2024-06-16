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

var AppReaper *reaper.Reaper

type ModuleData struct {
	Name   string
	Mod    nerv.Module
	Topics []*nerv.TopicCfg
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
		buildModReaper(config.Reaper), // Must come first (sets global)
		buildModWebUi(config.WebUi),
	} {
		engine.UseModule(
			mod.Mod,
			mod.Topics,
		)
		slog.Debug("loaded module", "name", mod.Name, "idx", idx)
	}
	return nil
}

func buildModWebUi(config webui.Config) ModuleData {

	moduleName := "webui"

	channel := strings.Join([]string{moduleName, "command"}, ".")

	publishingTopic := formatGroup(channel)

	config.DesignatedTopic = publishingTopic

	mod := webui.New(config)

	topic := nerv.NewTopic(publishingTopic).
		UsingBroadcast()

	AppReaper.AddListener(mod.ShutdownWarning)

	return ModuleData{
		Name:   moduleName,
		Mod:    mod,
		Topics: []*nerv.TopicCfg{topic},
	}
}

func buildModReaper(config reaper.Config) ModuleData {

	moduleName := "reaper"

	channel := strings.Join([]string{moduleName, "kill"}, ".")

	publishingTopic := formatGroup(channel)

	config.DesignatedTopic = publishingTopic

	AppReaper = reaper.New(config)

	mod := AppReaper

	topic := nerv.NewTopic(publishingTopic).
		UsingBroadcast()

	return ModuleData{
		Name:   moduleName,
		Mod:    mod,
		Topics: []*nerv.TopicCfg{topic},
	}
}
