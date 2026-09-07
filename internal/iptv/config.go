package iptv

import (
	"fmt"
	"time"

	"github.com/resnostyle/mqttkit/env"
)

type Settings struct {
	env.MQTT
	M3UURL                string
	EPGXMLURL             string
	PollIntervalSeconds   int
	ChannelRefreshSeconds int
	EPGCacheTTLSeconds    int
}

func FromEnv() (Settings, error) {
	mqtt, err := env.LoadMQTT("home/iptv", "iptv-mqtt")
	if err != nil {
		return Settings{}, err
	}
	m3uURL, err := env.Require("M3U_URL")
	if err != nil {
		return Settings{}, err
	}
	epgXMLURL, err := env.Require("EPG_XML_URL")
	if err != nil {
		return Settings{}, err
	}
	poll, err := env.Int("IPTV_POLL_INTERVAL_SECONDS", 60)
	if err != nil {
		return Settings{}, err
	}
	if poll < 15 {
		return Settings{}, fmt.Errorf("IPTV_POLL_INTERVAL_SECONDS must be >= 15")
	}
	refresh, err := env.Int("IPTV_CHANNEL_REFRESH_SECONDS", 3600)
	if err != nil {
		return Settings{}, err
	}
	if refresh < 60 {
		return Settings{}, fmt.Errorf("IPTV_CHANNEL_REFRESH_SECONDS must be >= 60")
	}
	cacheTTL, err := env.Int("EPG_CACHE_TTL_SECONDS", 900)
	if err != nil {
		return Settings{}, err
	}
	if cacheTTL < 60 {
		return Settings{}, fmt.Errorf("EPG_CACHE_TTL_SECONDS must be >= 60")
	}
	return Settings{
		MQTT:                  mqtt,
		M3UURL:                m3uURL,
		EPGXMLURL:             epgXMLURL,
		PollIntervalSeconds:   poll,
		ChannelRefreshSeconds: refresh,
		EPGCacheTTLSeconds:    cacheTTL,
	}, nil
}

func (s Settings) EPGCacheTTL() time.Duration {
	return time.Duration(s.EPGCacheTTLSeconds) * time.Second
}
