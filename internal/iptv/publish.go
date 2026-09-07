package iptv

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/resnostyle/iptv-mqtt/internal/lib/epg"
	"github.com/resnostyle/iptv-mqtt/internal/lib/m3u"
	"github.com/resnostyle/mqttkit/mqttpub"
)

// FetchChannels downloads and parses the M3U playlist.
func FetchChannels(ctx context.Context, client *http.Client, m3uURL string) ([]m3u.Channel, error) {
	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m3uURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "iptv-mqtt/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("m3u fetch %s: %s", m3uURL, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return m3u.Parse(string(body)), nil
}

type Publisher struct {
	settings Settings
	client   *http.Client
	mqtt     mqttpub.Sink
	channels []m3u.Channel
	lastFP   string
}

func NewPublisher(settings Settings, mqtt mqttpub.Sink) *Publisher {
	return &Publisher{
		settings: settings,
		client:   &http.Client{Timeout: 90 * time.Second},
		mqtt:     mqtt,
	}
}

func (p *Publisher) RefreshChannelsIfNeeded(ctx context.Context, force bool) error {
	if !force && len(p.channels) > 0 {
		return nil
	}
	channels, err := FetchChannels(ctx, p.client, p.settings.M3UURL)
	if err != nil {
		return err
	}
	fp := channelFingerprint(channels)
	if fp != p.lastFP {
		p.channels = channels
		p.lastFP = fp
		if err := PublishDiscovery(p.settings, p.mqtt, channels); err != nil {
			return fmt.Errorf("publish discovery: %w", err)
		}
		slog.Info("published mqtt discovery configs", "count", len(channels))
	} else {
		p.channels = channels
	}
	slog.Info("loaded channels", "count", len(channels))
	return nil
}

func (p *Publisher) FetchAndPublish(ctx context.Context) (map[string]epg.Programme, error) {
	if len(p.channels) == 0 {
		return nil, fmt.Errorf("no channels loaded")
	}
	programmes, err := epg.FetchNowPlaying(ctx, p.client, p.settings.EPGXMLURL, p.settings.EPGCacheTTL())
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	onAirCount := 0
	for _, ch := range p.channels {
		var prog *epg.Programme
		if ch.TVGID != "" {
			if found, ok := programmes[ch.TVGID]; ok {
				copy := found
				prog = &copy
				onAirCount++
			}
		}
		suffix := "channels/" + Slug(ch) + "/current"
		payload := BuildChannelPayload(ch, prog, now)
		if err := p.mqtt.PublishQuiet(suffix, payload, true); err != nil {
			return programmes, fmt.Errorf("publish %s: %w", suffix, err)
		}
	}

	summary := BuildSummaryPayload(len(p.channels), onAirCount, now)
	if err := p.mqtt.Publish("summary", summary, true); err != nil {
		return programmes, fmt.Errorf("publish summary: %w", err)
	}

	slog.Info("published iptv topics",
		"channels", len(p.channels),
		"on_air", onAirCount,
	)
	return programmes, nil
}

func (p *Publisher) ShouldRefreshChannels(lastRefresh time.Time, now time.Time) bool {
	if lastRefresh.IsZero() {
		return true
	}
	return now.Sub(lastRefresh) >= time.Duration(p.settings.ChannelRefreshSeconds)*time.Second
}
