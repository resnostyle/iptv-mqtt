package iptv

import (
	"fmt"
	"strings"

	"github.com/resnostyle/iptv-mqtt/internal/lib/m3u"
	"github.com/resnostyle/mqttkit/hadisc"
	"github.com/resnostyle/mqttkit/mqttpub"
)

const (
	deviceManufacturer = "iptv-mqtt"
	deviceModel        = "IPTV"
)

func deviceUID(group string) string {
	slug := strings.ToLower(group)
	slug = strings.ReplaceAll(slug, " ", "_")
	slug = strings.ReplaceAll(slug, "/", "_")
	slug = strings.ReplaceAll(slug, ":", "_")
	if slug == "" {
		slug = "uncategorized"
	}
	return "iptv_mqtt_" + slug
}

func deviceBlock(group string) map[string]any {
	return hadisc.Device([]string{deviceUID(group)}, group, deviceManufacturer, deviceModel)
}

// BuildDiscoveryConfigs returns HA MQTT discovery configs for all channels.
func BuildDiscoveryConfigs(topicPrefix string, channels []m3u.Channel) []mqttpub.Config {
	configs := make([]mqttpub.Config, 0, len(channels))
	for _, ch := range channels {
		slug := Slug(ch)
		stateTopic := topicPrefix + "/channels/" + slug + "/current"
		objectID := "iptv_mqtt_" + slug
		name := ch.Name
		if name == "" {
			name = slug
		}
		payload := map[string]any{
			"name":                  name,
			"unique_id":             objectID,
			"state_topic":           stateTopic,
			"value_template":        "{{ value_json.title }}",
			"device":                deviceBlock(ch.Group),
			"object_id":             objectID,
			"icon":                  "mdi:television",
			"json_attributes_topic": stateTopic,
			"availability": []map[string]any{
				{
					"topic":                 stateTopic,
					"value_template":        "{{ value_json.on_air }}",
					"payload_available":     "true",
					"payload_not_available": "false",
				},
			},
		}
		configs = append(configs, mqttpub.Config{
			ObjectID:  objectID,
			Component: "sensor",
			Payload:   payload,
		})
	}
	return configs
}

func PublishDiscovery(settings Settings, mqtt mqttpub.Sink, channels []m3u.Channel) error {
	if !settings.MQTTDiscoveryEnabled {
		return nil
	}
	configs := BuildDiscoveryConfigs(settings.MQTTTopicPrefix, channels)
	if err := mqtt.PublishDiscovery(configs, settings.MQTTDiscoveryPrefix); err != nil {
		return err
	}
	return nil
}

func channelFingerprint(channels []m3u.Channel) string {
	parts := make([]string, 0, len(channels))
	for _, ch := range channels {
		parts = append(parts, fmt.Sprintf("%s:%s:%s", ch.ID, ch.TVGID, ch.Name))
	}
	return strings.Join(parts, "|")
}
