# iptv-mqtt

Publishes IPTV now-playing data to MQTT for Home Assistant. Mirrors the [racing-mqtt](https://github.com/resnostyle/racing-mqtt) architecture.

## What it does

1. Loads all channels from your M3U playlist
2. Fetches bulk XMLTV EPG from the URL you configure
3. Publishes retained JSON per channel to MQTT
4. Registers Home Assistant MQTT discovery sensors (grouped by M3U `group-title`)

## Quick start

```bash
cp .env.example .env
# Edit M3U_URL and EPG_XML_URL for your provider

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
| `EPG_XML_URL` | (required) | XMLTV guide URL |
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

Slug is the channel `tvg-id` with `.` replaced by `_` (e.g. `channel.us` → `channel_us`).

### Per-channel payload

```json
{
  "tvg_id": "channel.us",
  "name": "Example Channel",
  "group": "Example Group",
  "logo": "https://example.com/logos/channel.png",
  "title": "Example Show",
  "desc": "Show description.",
  "lang": "en",
  "start": "2026-09-01T20:30:00Z",
  "stop": "2026-09-02T01:00:00Z",
  "on_air": true,
  "published_at": "2026-09-01T21:30:00Z"
}
```

## Home Assistant

Discovery creates one `sensor` per channel, grouped under devices named by M3U group. The sensor state is the show title; full JSON is available as attributes.

Example automation:

```yaml
trigger:
  - platform: mqtt
    topic: home/iptv/channels/channel_us/current
    value_template: "{{ value_json.title }}"
    payload: "Example Show"
```

## Development

```bash
mise run test
mise run build
```
