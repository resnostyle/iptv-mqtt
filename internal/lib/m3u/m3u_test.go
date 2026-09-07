package m3u

import "testing"

const sampleM3U = `#EXTM3U
#EXTINF:-1 tvg-id="channel_a.us" tvg-name="Channel A" tvg-logo="https://example.com/channel_a.png" group-title="Example Group",Channel A
http://example.com/stream/1
#EXTINF:-1 tvg-id="channel_b.us" tvg-name="Channel B" group-title="Example Group",Channel B
http://example.com/stream/2
`

func TestParseExtractsChannels(t *testing.T) {
	channels := Parse(sampleM3U)
	if len(channels) != 2 {
		t.Fatalf("got %d channels want 2", len(channels))
	}
	if channels[0].TVGID != "channel_a.us" {
		t.Fatalf("tvgId: got %q", channels[0].TVGID)
	}
	if channels[0].Name != "Channel A" {
		t.Fatalf("name: got %q", channels[0].Name)
	}
	if channels[0].Group != "Example Group" {
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
