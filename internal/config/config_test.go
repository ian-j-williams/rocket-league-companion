package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveRocketLeagueWritesStatsAPIIni(t *testing.T) {
	path := filepath.Join(t.TempDir(), "TAGame", "Config", "TAStatsAPI.ini")
	settings := Settings{PacketSendRate: 45, Port: 49123, WebPort: 49124}

	if err := SaveRocketLeague(path, settings); err != nil {
		t.Fatalf("SaveRocketLeague() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated config: %v", err)
	}
	want := "[StatsAPI]\nbEnabled=true\nPacketSendRate=45\nPort=49123\nWebPort=49124\n"
	if string(data) != want {
		t.Fatalf("generated config = %q, want %q", data, want)
	}
	if !strings.Contains(string(data), "PacketSendRate=45") {
		t.Fatal("generated config does not contain the selected packet rate")
	}
}
