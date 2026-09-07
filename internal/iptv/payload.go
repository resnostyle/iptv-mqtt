package iptv

import (
	"strings"
	"time"

	"github.com/resnostyle/iptv-mqtt/internal/lib/epg"
	"github.com/resnostyle/iptv-mqtt/internal/lib/m3u"
	"github.com/resnostyle/mqttkit/payload"
)

// Slug returns a stable MQTT/HA identifier for a channel.
func Slug(ch m3u.Channel) string {
	if ch.TVGID != "" {
		return strings.ReplaceAll(ch.TVGID, ".", "_")
	}
	return "ch_" + ch.ID
}

func BuildChannelPayload(ch m3u.Channel, prog *epg.Programme, publishedAt time.Time) map[string]any {
	out := map[string]any{
		"tvg_id":       payload.NilIfEmpty(ch.TVGID),
		"name":         ch.Name,
		"group":        ch.Group,
		"logo":         payload.NilIfEmpty(ch.Logo),
		"title":        "",
		"stop":         nil,
		"on_air":       false,
		"published_at": publishedAt.UTC().Format(time.RFC3339),
	}
	if prog != nil && prog.Title != "" {
		out["title"] = prog.Title
		out["desc"] = payload.NilIfEmpty(prog.Desc)
		out["lang"] = payload.NilIfEmpty(prog.Lang)
		out["start"] = prog.Start.UTC().Format(time.RFC3339)
		out["stop"] = prog.Stop.UTC().Format(time.RFC3339)
		out["on_air"] = true
	}
	return out
}

func BuildSummaryPayload(channelCount, onAirCount int, publishedAt time.Time) map[string]any {
	return map[string]any{
		"channel_count": channelCount,
		"on_air_count":  onAirCount,
		"published_at":  publishedAt.UTC().Format(time.RFC3339),
	}
}
