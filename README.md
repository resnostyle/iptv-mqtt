# iptv-mqtt

Publishes IPTV now-playing data to MQTT for Home Assistant. Mirrors the [racing-mqtt](https://github.com/resnostyle/racing-mqtt) architecture.

## What it does

1. Loads all channels from your M3U playlist (mybunny.tv)
2. Fetches bulk XMLTV EPG from the public mybunny.tv guide (`https://mybunny.tv/epg.xml` by default)
3. Publishes retained JSON per channel to MQTT
4. Registers Home Assistant MQTT discovery sensors (grouped by M3U `group-title`)

## Quick start

```bash
cp .env.example .env
# Edit M3U_URL with your playlist credentials

mise run iptv
```

Or with Docker:

```bash
docker compose up -d
```

## Configuration

| Variable | Default | Description |
|---|---|---|
| `M3U_URL` | (required) | IPTV playlist URL (channel list) |
| `EPG_XML_URL` | `https://mybunny.tv/epg.xml` | Public master XMLTV guide |
| `EPG_CACHE_TTL_SECONDS` | `900` | How long to cache the EPG download |
| `MQTT_HOST` | `127.0.0.1` | MQTT broker |
| `MQTT_PORT` | `1883` | MQTT port |
| `MQTT_TOPIC_PREFIX` | `home/iptv` | Topic prefix |
| `MQTT_DISCOVERY_ENABLED` | `true` | HA auto-discovery |
| `IPTV_POLL_INTERVAL_SECONDS` | `60` | EPG poll interval |
| `IPTV_CHANNEL_REFRESH_SECONDS` | `3600` | M3U refresh interval |

## MQTT topics

| Topic | Description |
|---|---|
| `home/iptv/channels/{slug}/current` | Per-channel now playing |
| `home/iptv/summary` | Channel and on-air counts |

Slug is the channel `tvg-id` with `.` replaced by `_` (e.g. `espn.us` → `espn_us`).

### Per-channel payload

```json
{
  "tvg_id": "espn.us",
  "name": "US: ESPN",
  "group": "US Sports",
  "logo": "https://logo.m3uassets.com/espn.png",
  "title": "SportsCenter",
  "desc": "Live sports news and highlights.",
  "lang": "en",
  "start": "2026-09-01T20:30:00Z",
  "stop": "2026-09-02T01:00:00Z",
  "on_air": true,
  "published_at": "2026-09-01T21:30:00Z"
}
```

## Home Assistant

Discovery creates one `sensor` per channel, grouped under devices named by M3U group (e.g. "US Sports"). The sensor state is the show title; full JSON is available as attributes.

Example automation:

```yaml
trigger:
  - platform: mqtt
    topic: home/iptv/channels/espn_us/current
    value_template: "{{ value_json.title }}"
    payload: "Monday Night Football"
```

## Development

```bash
mise run test
mise run build
```
