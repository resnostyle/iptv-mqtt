// Package m3u parses IPTV M3U playlists.
package m3u

import (
	"regexp"
	"strconv"
	"strings"
)

// Channel is a parsed playlist entry.
type Channel struct {
	ID    string
	Name  string
	TVGID string
	Logo  string
	Group string
	URL   string
}

var attrRE = regexp.MustCompile(`([a-z-]+)="([^"]*)"`)

func getAttr(line, name string) string {
	for _, match := range attrRE.FindAllStringSubmatch(line, -1) {
		if match[1] == name {
			return match[2]
		}
	}
	return ""
}

// Parse reads an M3U playlist and returns channel entries.
func Parse(text string) []Channel {
	lines := strings.Split(text, "\n")
	var channels []Channel
	channelID := 1

	for i, raw := range lines {
		line := strings.TrimSpace(raw)
		if !strings.HasPrefix(line, "#EXTINF:") {
			continue
		}

		url := ""
		for _, candidate := range lines[i+1:] {
			candidate = strings.TrimSpace(candidate)
			if candidate == "" || strings.HasPrefix(candidate, "#") {
				continue
			}
			url = candidate
			break
		}
		if url == "" {
			continue
		}

		name := "Unknown"
		if idx := strings.LastIndex(line, ","); idx >= 0 {
			name = strings.TrimSpace(line[idx+1:])
		}

		group := getAttr(line, "group-title")
		if group == "" {
			group = "Uncategorized"
		}

		channels = append(channels, Channel{
			ID:    strconv.Itoa(channelID),
			Name:  name,
			TVGID: getAttr(line, "tvg-id"),
			Logo:  getAttr(line, "tvg-logo"),
			Group: group,
			URL:   url,
		})
		channelID++
	}

	return channels
}
