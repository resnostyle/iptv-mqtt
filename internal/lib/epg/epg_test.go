package epg

import (
	"testing"
	"time"
)

const sampleXML = `<?xml version="1.0" encoding="UTF-8"?>
<tv>
  <programme start="20260101000000 +0000" stop="20260101120000 +0000" channel="espn.us">
    <title>Morning Show</title>
  </programme>
  <programme start="20260101120000 +0000" stop="20260101150000 +0000" channel="espn.us">
    <title>Afternoon Game</title>
  </programme>
  <programme start="20260101120000 +0000" stop="20260101130000 +0000" channel="fox.us">
    <title lang="en">Fox News Hour</title>
  </programme>
  <programme start="20260101120000 +0000" stop="20260101130000 +0000" channel="zdf.de">
    <title lang="de">heute journal</title>
    <desc lang="de">Nachrichten aus Deutschland</desc>
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
	if result["espn.us"].Title != "Afternoon Game" {
		t.Fatalf("espn title: got %q", result["espn.us"].Title)
	}
	if result["fox.us"].Title != "Fox News Hour" {
		t.Fatalf("fox title: got %q", result["fox.us"].Title)
	}
	if result["fox.us"].Lang != "en" {
		t.Fatalf("fox lang: got %q", result["fox.us"].Lang)
	}
	if result["zdf.de"].Lang != "de" {
		t.Fatalf("zdf lang: got %q", result["zdf.de"].Lang)
	}
	if result["zdf.de"].Desc != "Nachrichten aus Deutschland" {
		t.Fatalf("zdf desc: got %q", result["zdf.de"].Desc)
	}
}
