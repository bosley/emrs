/*
  nerv-go modules supporting EMRS
*/

package main

import (
	"github.com/bosley/nerv-go"
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

	return ModuleData{
		Name:   moduleName,
		Mod:    mod,
		Topics: []*nerv.TopicCfg{topic},
	}
}
