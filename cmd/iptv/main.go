package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/resnostyle/iptv-mqtt/internal/iptv"
	"github.com/resnostyle/mqttkit/logx"
	"github.com/resnostyle/mqttkit/mqttpub"
	"github.com/resnostyle/mqttkit/poll"
)

func main() {
	settings, err := iptv.FromEnv()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	logx.Configure(settings.LogLevel, false)

	slog.Info("starting iptv-mqtt",
		"poll_interval", settings.PollIntervalSeconds,
		"channel_refresh", settings.ChannelRefreshSeconds,
		"mqtt", settings.MQTTHost,
		"port", settings.MQTTPort,
		"discovery", settings.MQTTDiscoveryEnabled,
	)

	ctx, cancel := poll.NotifyContext()
	defer cancel()

	mqtt, err := mqttpub.New(
		settings.MQTTHost,
		settings.MQTTPort,
		settings.MQTTClientID,
		settings.MQTTUsername,
		settings.MQTTPassword,
		settings.MQTTTopicPrefix,
	)
	if err != nil {
		slog.Error("mqtt connect failed", "err", err)
		os.Exit(1)
	}
	defer mqtt.Close()

	publisher := iptv.NewPublisher(settings, mqtt)
	if err := publisher.RefreshChannelsIfNeeded(ctx, true); err != nil {
		slog.Error("initial channel load failed", "err", err)
		os.Exit(1)
	}
	lastChannelRefresh := time.Now().UTC()

	for ctx.Err() == nil {
		now := time.Now().UTC()
		if publisher.ShouldRefreshChannels(lastChannelRefresh, now) {
			if err := publisher.RefreshChannelsIfNeeded(ctx, true); err != nil {
				slog.Error("channel refresh failed", "err", err)
			} else {
				lastChannelRefresh = now
			}
		}

		programmes, err := publisher.FetchAndPublish(ctx)
		if err != nil {
			slog.Error("fetch/publish failed", "err", err)
		}

		wait := iptv.PollInterval(settings, programmes, time.Now().UTC())
		slog.Debug("waiting for next poll", "seconds", wait.Seconds())
		poll.Wait(ctx, wait)
	}
	slog.Info("exited")
}
