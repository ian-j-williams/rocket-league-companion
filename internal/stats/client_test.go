package stats

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestClientAcceptsStatsStream(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen test server: %v", err)
	}
	defer listener.Close()

	client := NewClient(listener.Addr().(*net.TCPAddr).Port)
	if err := client.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	defer client.Close()

	conn, err := listener.Accept()
	if err != nil {
		t.Fatalf("accept client connection: %v", err)
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
