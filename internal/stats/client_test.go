package stats

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestClientAcceptsStatsStream(t *testing.T) {
	client := NewClient(0)
	if err := client.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	defer client.Close()

	conn, err := net.Dial("tcp", client.Addr().String())
	if err != nil {
		t.Fatalf("dial listener: %v", err)
	}
	defer conn.Close()

	message := `{"Event":"UpdateState","Data":{"Game":{"Arena":"DFH Stadium","PlaylistId":12,"TimeSeconds":87,"Teams":[{"TeamNum":0,"Score":2},{"TeamNum":1,"Score":1}]}}}`
	if _, err := fmt.Fprintln(conn, message); err != nil {
		t.Fatalf("write stats message: %v", err)
	}

	select {
	case snapshot := <-client.Updates():
		if snapshot.Arena != "DFH Stadium" || snapshot.PlaylistName != "Doubles" || snapshot.ScoreBlue != 2 || snapshot.ScoreOrange != 1 {
			t.Fatalf("unexpected snapshot: %+v", snapshot)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for stats snapshot")
	}
}
