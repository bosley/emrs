package api

import (
	"time"
)

const (
	HttpV1SubmitEvent = "/submit/event"
	HttpV1Stat        = "/stat"

	HttpV1CNCShutdown = "/cnc/shutdown"
)

type Options struct {
	Binding     string
	AssetId     string
	AccessToken string
}

type CNCApi interface {
	Shutdown() error
}

type SubmissionApi interface {
	Submit(route string, data []byte) error
}

type StatsApi interface {
	GetUptime() (time.Duration, error)
}
