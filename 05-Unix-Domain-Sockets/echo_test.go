package unixdomainsockets

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestEchoServerUnix(t *testing.T) {

	// Create a subdirectory in your OS temp folder named echo_unix
	// that will contain the echo serve's socket file.
	dir, err := ioutil.TempDir("", "echo_unix")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {

		// The deferred call to os.RemoveAll cleans up after the server
		// test completes.
		if rErr := os.RemoveAll(dir); rErr != nil {
			t.Error(rErr)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Pass a socket file named #.sock where # is the server's process ID.
	socket := filepath.Join(dir, fmt.Sprintf("%d.sock", os.Getpid()))
	rAddr, err := streamingEchoServer(ctx, "unix", socket)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure everyone has read and write access to the socket.
	err = os.Chmod(socket, os.ModeSocket|0666)
	if err != nil {
		t.Fatal(err)
	}

	conn, err := net.Dial("unix", rAddr.String())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()

	msg := []byte("ping")
	for i := 0; i < 3; i++ { // write 3 "ping" messages
		_, err = conn.Write(msg)
		if err != nil {
			t.Fatal(err)
		}
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf) // read once from the server
	if err != nil {
		t.Fatal(err)
	}

	expected := bytes.Repeat(msg, 3)
	if !bytes.Equal(expected, buf[:n]) {
		t.Fatalf("expected reply %q; actual reply %q", expected, buf[:n])
	}
}
