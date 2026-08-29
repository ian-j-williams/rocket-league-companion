package state

import "testing"

func TestSessionStatsTracksWinsAndLosses(t *testing.T) {
	stats := NewSessionStats()

	stats.RecordWin()
	stats.RecordWin()
	stats.RecordLoss()

	if stats.Wins != 2 {
		t.Fatalf("expected 2 wins, got %d", stats.Wins)
	}
	if stats.Losses != 1 {
		t.Fatalf("expected 1 loss, got %d", stats.Losses)
	}
}

func TestSessionStatsResetClearsState(t *testing.T) {
	stats := NewSessionStats()
	stats.RecordWin()
	stats.RecordLoss()

	stats.Reset()

	if stats.Wins != 0 || stats.Losses != 0 {
		t.Fatalf("expected stats to be reset, got wins=%d losses=%d", stats.Wins, stats.Losses)
	}
}
