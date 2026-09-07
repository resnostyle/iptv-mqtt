package m3u

import "testing"

const sampleM3U = `#EXTM3U
#EXTINF:-1 tvg-id="espn.us" tvg-name="US: ESPN" tvg-logo="https://example.com/espn.png" group-title="US Sports",US: ESPN
http://example.com/stream/1
#EXTINF:-1 tvg-id="foxsports1.us" tvg-name="US: Fox Sports 1" group-title="US Sports",US: Fox Sports 1
http://example.com/stream/2
`

func TestParseExtractsChannels(t *testing.T) {
	channels := Parse(sampleM3U)
	if len(channels) != 2 {
		t.Fatalf("got %d channels want 2", len(channels))
	}
	if channels[0].TVGID != "espn.us" {
		t.Fatalf("tvgId: got %q", channels[0].TVGID)
	}
	if channels[0].Name != "US: ESPN" {
		t.Fatalf("name: got %q", channels[0].Name)
	}
	if channels[0].Group != "US Sports" {
		t.Fatalf("group: got %q", channels[0].Group)
	}
	if channels[0].URL != "http://example.com/stream/1" {
		t.Fatalf("url: got %q", channels[0].URL)
	}
}

func TestParseSkipsEntriesWithoutURL(t *testing.T) {
	text := `#EXTINF:-1 tvg-name="Broken",Broken
#EXTM3U
`
	if got := Parse(text); len(got) != 0 {
		t.Fatalf("got %d channels want 0", len(got))
	}
}
