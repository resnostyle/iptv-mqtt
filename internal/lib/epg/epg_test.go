package epg

import (
	"testing"
	"time"
)

const sampleXML = `<?xml version="1.0" encoding="UTF-8"?>
<tv>
  <programme start="20260101000000 +0000" stop="20260101120000 +0000" channel="channel_a.us">
    <title>Morning Show</title>
  </programme>
  <programme start="20260101120000 +0000" stop="20260101150000 +0000" channel="channel_a.us">
    <title>Afternoon Show</title>
  </programme>
  <programme start="20260101120000 +0000" stop="20260101130000 +0000" channel="channel_b.us">
    <title lang="en">News Hour</title>
  </programme>
  <programme start="20260101120000 +0000" stop="20260101130000 +0000" channel="channel_c.de">
    <title lang="de">Abendnachrichten</title>
    <desc lang="de">Nachrichtenbeschreibung</desc>
  </programme>
</tv>
`

func TestParseXMLTVTime(t *testing.T) {
	got, err := ParseXMLTVTime("20260101123000 +0000")
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 1, 1, 12, 30, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestParseNowPlaying(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 45, 0, 0, time.UTC)
	result, err := ParseNowPlaying(sampleXML, now)
	if err != nil {
		t.Fatal(err)
	}
	if result["channel_a.us"].Title != "Afternoon Show" {
		t.Fatalf("channel_a title: got %q", result["channel_a.us"].Title)
	}
	if result["channel_b.us"].Title != "News Hour" {
		t.Fatalf("channel_b title: got %q", result["channel_b.us"].Title)
	}
	if result["channel_b.us"].Lang != "en" {
		t.Fatalf("channel_b lang: got %q", result["channel_b.us"].Lang)
	}
	if result["channel_c.de"].Lang != "de" {
		t.Fatalf("channel_c lang: got %q", result["channel_c.de"].Lang)
	}
	if result["channel_c.de"].Desc != "Nachrichtenbeschreibung" {
		t.Fatalf("channel_c desc: got %q", result["channel_c.de"].Desc)
	}
}
