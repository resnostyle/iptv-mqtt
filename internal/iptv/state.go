package iptv

import (
	"time"

	"github.com/resnostyle/iptv-mqtt/internal/lib/epg"
)

const (
	nearBoundaryWindow = 5 * time.Minute
	nearBoundaryPoll   = 15 * time.Second
)

// PollInterval returns how long to wait before the next EPG fetch.
func PollInterval(settings Settings, programmes map[string]epg.Programme, now time.Time) time.Duration {
	defaultInterval := time.Duration(settings.PollIntervalSeconds) * time.Second
	nearest := nextProgrammeStop(programmes, now)
	if nearest != nil && nearest.Sub(now) <= nearBoundaryWindow {
		return nearBoundaryPoll
	}
	return defaultInterval
}

func nextProgrammeStop(programmes map[string]epg.Programme, now time.Time) *time.Time {
	var nearest *time.Time
	for _, prog := range programmes {
		if !prog.Stop.After(now) {
			continue
		}
		if nearest == nil || prog.Stop.Before(*nearest) {
			copy := prog.Stop
			nearest = &copy
		}
	}
	return nearest
}
