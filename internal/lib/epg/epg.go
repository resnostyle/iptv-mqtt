// Package epg parses XMLTV programme guides for now-playing lookups.
package epg

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const defaultCacheTTL = 60 * time.Second

// Programme is the currently airing show for a channel.
type Programme struct {
	Title string
	Desc  string
	Lang  string
	Start time.Time
	Stop  time.Time
}

type cacheEntry struct {
	fetchedAt  time.Time
	programmes map[string]Programme
}

var (
	cacheMu sync.Mutex
	cache   = make(map[string]cacheEntry)
)

// ParseXMLTVTime parses an XMLTV datetime string.
func ParseXMLTVTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if len(value) < 14 {
		return time.Time{}, fmt.Errorf("invalid xmltv time: %q", value)
	}
	stamp := value[:14]
	tzPart := strings.TrimSpace(value[15:])
	if tzPart == "" {
		tzPart = "+0000"
	}
	loc, err := parseXMLTVOffset(tzPart)
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.ParseInLocation("20060102150405", stamp, loc)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

func parseXMLTVOffset(tzPart string) (*time.Location, error) {
	if len(tzPart) < 5 {
		return time.UTC, nil
	}
	sign := 1
	if tzPart[0] == '-' {
		sign = -1
	}
	var hours, minutes int
	if _, err := fmt.Sscanf(tzPart[1:], "%2d%2d", &hours, &minutes); err != nil {
		return nil, err
	}
	offset := sign * ((hours * 60) + minutes)
	return time.FixedZone("xmltv", offset*60), nil
}

type textXML struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

type programmeXML struct {
	Start   string  `xml:"start,attr"`
	Stop    string  `xml:"stop,attr"`
	Channel string  `xml:"channel,attr"`
	Title   textXML `xml:"title"`
	Desc    textXML `xml:"desc"`
}

// ParseNowPlaying returns currently airing programmes keyed by channel ID.
func ParseNowPlaying(xmlText string, now time.Time) (map[string]Programme, error) {
	decoder := xml.NewDecoder(strings.NewReader(xmlText))
	programmes := make(map[string]Programme)

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		start, ok := tok.(xml.StartElement)
		if !ok || start.Name.Local != "programme" {
			continue
		}

		var prog programmeXML
		if err := decoder.DecodeElement(&prog, &start); err != nil {
			continue
		}
		channel := strings.TrimSpace(prog.Channel)
		if channel == "" {
			continue
		}
		if _, exists := programmes[channel]; exists {
			continue
		}

		startTime, err := ParseXMLTVTime(prog.Start)
		if err != nil {
			continue
		}
		stopTime, err := ParseXMLTVTime(prog.Stop)
		if err != nil {
			continue
		}
		if !startTime.After(now) && now.Before(stopTime) {
			programmes[channel] = Programme{
				Title: strings.TrimSpace(prog.Title.Value),
				Desc:  strings.TrimSpace(prog.Desc.Value),
				Lang:  firstNonEmpty(prog.Title.Lang, prog.Desc.Lang),
				Start: startTime,
				Stop:  stopTime,
			}
		}
	}

	return programmes, nil
}

// FetchNowPlaying downloads and parses the XMLTV guide with a short-lived cache.
func FetchNowPlaying(ctx context.Context, client *http.Client, xmlURL string, cacheTTL time.Duration) (map[string]Programme, error) {
	if cacheTTL <= 0 {
		cacheTTL = defaultCacheTTL
	}
	now := time.Now().UTC()
	cacheMu.Lock()
	if entry, ok := cache[xmlURL]; ok && now.Sub(entry.fetchedAt) < cacheTTL {
		programmes := copyProgrammes(entry.programmes)
		cacheMu.Unlock()
		return programmes, nil
	}
	cacheMu.Unlock()

	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, xmlURL, nil)
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
		return nil, fmt.Errorf("epg fetch %s: %s", xmlURL, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	programmes, err := ParseNowPlaying(string(body), now)
	if err != nil {
		return nil, err
	}

	cacheMu.Lock()
	cache[xmlURL] = cacheEntry{fetchedAt: now, programmes: programmes}
	cacheMu.Unlock()
	return copyProgrammes(programmes), nil
}

func copyProgrammes(src map[string]Programme) map[string]Programme {
	out := make(map[string]Programme, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
