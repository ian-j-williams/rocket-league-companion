package state

// SessionStats tracks wins and losses for the current session.
type SessionStats struct {
	Wins   int
	Losses int
}

func NewSessionStats() *SessionStats {
	return &SessionStats{}
}

func (s *SessionStats) RecordWin() {
	s.Wins++
}

func (s *SessionStats) RecordLoss() {
	s.Losses++
}

func (s *SessionStats) Reset() {
	s.Wins = 0
	s.Losses = 0
}
