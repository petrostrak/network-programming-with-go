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

	// Pass a socket file named #.sock where # is the server's process ID.
	socket := filepath.Join(dir, fmt.Sprintf("%d.sock", os.Getpid()))
	rAddr, err := streamingEchoServer(ctx, "unix", socket)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

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

func TestEchoServerUnixDatagram(t *testing.T) {
	dir, err := ioutil.TempDir("", "echo_unixgram")
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

	// Pass a socket file named #.sock where # is the server's process ID.
	sSocket := filepath.Join(dir, fmt.Sprintf("%d.sock", os.Getpid()))
	serverAddr, err := datagramEchoServer(ctx, "unixgram", sSocket)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	// Make sure everyone has read and write access to the socket.
	err = os.Chmod(sSocket, os.ModeSocket|0666)
	if err != nil {
		t.Fatal(err)
	}

	cSocket := filepath.Join(dir, fmt.Sprintf("c%d.sock", os.Getpid()))
	client, err := net.ListenPacket("unixgram", cSocket)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = client.Close() }()

	err = os.Chmod(cSocket, os.ModeSocket|0622)
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("ping")
	for i := 0; i < 3; i++ { // write 3 "ping" messages
		_, err = client.WriteTo(msg, serverAddr)
		if err != nil {
			t.Fatal(err)
		}
	}

	buf := make([]byte, 1024)
	for i := 0; i < 3; i++ { // read 3 "ping" messages
		n, addr, err := client.ReadFrom(buf)
		if err != nil {
			t.Fatal(err)
		}

		if addr.String() != serverAddr.String() {
			t.Fatalf("received reply from %q instead of %q", addr, serverAddr)
		}

		if !bytes.Equal(msg, buf[:n]) {
			t.Fatalf("expected reply %q; actual reply %q", msg, buf[:n])
		}
	}
}
