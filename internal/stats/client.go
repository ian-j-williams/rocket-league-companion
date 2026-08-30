package stats

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

// Snapshot is the subset of Rocket League Stats API data the companion UI renders.
type Snapshot struct {
	Arena        string
	TimeSeconds  int
	PlaylistName string
	ScoreBlue    int
	ScoreOrange  int
}

// MessageEnvelope is the top-level message shape the game broadcasts.
type MessageEnvelope struct {
	Event string          `json:"Event"`
	Data  json.RawMessage `json:"Data"`
}

type updateStateData struct {
	Game struct {
		Arena       string `json:"Arena"`
		PlaylistId  int    `json:"PlaylistId"`
		TimeSeconds int    `json:"TimeSeconds"`
		Teams       []struct {
			TeamNum int    `json:"TeamNum"`
			Name    string `json:"Name"`
			Score   int    `json:"Score"`
		} `json:"Teams"`
	} `json:"Game"`
}

// Client connects to the local Rocket League stats socket and streams parsed snapshots.
type Client struct {
	port     int
	listener net.Listener
	updates  chan Snapshot
}

func NewClient(port int) *Client {
	return &Client{
		port:    port,
		updates: make(chan Snapshot, 32),
	}
}

func (c *Client) Start() error {
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(c.port))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}
	c.listener = listener
	go c.acceptLoop()
	return nil
}

func (c *Client) acceptLoop() {
	defer close(c.updates)

	for {
		conn, err := c.listener.Accept()
		if err != nil {
			return
		}
		go c.readLoop(conn)
	}
}

func (c *Client) Close() error {
	if c.listener == nil {
		return nil
	}
	return c.listener.Close()
}

func (c *Client) Updates() <-chan Snapshot {
	return c.updates
}

func (c *Client) Addr() net.Addr {
	if c.listener == nil {
		return nil
	}
	return c.listener.Addr()
}

func (c *Client) readLoop(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	for {
		var envelope MessageEnvelope
		if err := decoder.Decode(&envelope); err != nil {
			return
		}
		if envelope.Event != "UpdateState" {
			continue
		}

		var state updateStateData
		if err := json.Unmarshal(envelope.Data, &state); err != nil {
			continue
		}

		payload := Snapshot{
			Arena:        state.Game.Arena,
			TimeSeconds:  state.Game.TimeSeconds,
			PlaylistName: playlistName(state.Game.PlaylistId),
		}
		for _, team := range state.Game.Teams {
			switch team.TeamNum {
			case 0:
				payload.ScoreBlue = team.Score
			case 1:
				payload.ScoreOrange = team.Score
			}
		}
		c.updates <- payload
	}
}

func playlistName(id int) string {
	switch id {
	case 11:
		return "Duel"
	case 12:
		return "Doubles"
	case 13:
		return "Solo Standard"
	case 28:
		return "Ranked Duel"
	case 29:
		return "Ranked Doubles"
	case 30:
		return "Ranked Solo Standard"
	case 34:
		return "Chaos"
	case 35:
		return "Private Match"
	default:
		return fmt.Sprintf("Playlist %d", id)
	}
}
