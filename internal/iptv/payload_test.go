package iptv

import (
	"testing"

	"github.com/resnostyle/iptv-mqtt/internal/lib/m3u"
)

func TestSlug(t *testing.T) {
	ch := m3u.Channel{ID: "1", TVGID: "espn.us"}
	if got := Slug(ch); got != "espn_us" {
		t.Fatalf("got %q want espn_us", got)
	}
	ch = m3u.Channel{ID: "42"}
	if got := Slug(ch); got != "ch_42" {
		t.Fatalf("got %q want ch_42", got)
	}
}

func TestBuildDiscoveryConfigs(t *testing.T) {
	channels := []m3u.Channel{
		{ID: "1", Name: "US: ESPN", TVGID: "espn.us", Group: "US Sports"},
	}
	configs := BuildDiscoveryConfigs("home/iptv", channels)
	if len(configs) != 1 {
		t.Fatalf("got %d configs want 1", len(configs))
	}
	if configs[0].ObjectID != "iptv_mqtt_espn_us" {
		t.Fatalf("object id: %q", configs[0].ObjectID)
	}
	device := configs[0].Payload["device"].(map[string]any)
	if device["name"] != "US Sports" {
		t.Fatalf("device name: %v", device["name"])
	}
}
