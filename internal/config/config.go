package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Settings stores the app's user preferences for the Rocket League Stats API.
type Settings struct {
	PacketSendRate float64 `json:"packet_send_rate"`
	Port           int     `json:"port"`
	WebPort        int     `json:"web_port"`
}

func DefaultSettings() Settings {
	return Settings{
		PacketSendRate: 30,
		Port:           49123,
		WebPort:        49124,
	}
}

func Load(path string) (Settings, error) {
	settings := DefaultSettings()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return settings, nil
		}
		return Settings{}, err
	}

	if err := json.Unmarshal(data, &settings); err != nil {
		return Settings{}, err
	}

	return settings, nil
}

func Save(path string, settings Settings) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
