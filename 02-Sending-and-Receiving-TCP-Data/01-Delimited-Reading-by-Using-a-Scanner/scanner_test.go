package delimitedreadingbyusingascanner

import (
	"bufio"
	"net"
	"reflect"
	"testing"
)

const (
	payload = "The bigger the interface, the weaker the abstraction."
)

func TestScanner(t *testing.T) {

	// All it's meant to do is serve up the payload.
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		_, err = conn.Read([]byte(payload))
		if err != nil {
			t.Error(err)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Create a bufio.Scanner that reads from the network connection.
	scanner := bufio.NewScanner(conn)

	// By default the scanner withh split data read from the network
	// connection when it encounters a newline character '\n' in the
	// stream of data. Instead we use the ScanWords.
	scanner.Split(bufio.ScanWords)

	var words []string

	// Every call to Scan() can result in multiple calls to the network
	// connection's Read() until the scanner finds its delimiter or reads
	// an error.
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		t.Error(err)
	}

	expected := []string{"The", "bigger", "the", "interface", "the", "weaker", "the", "abstraction."}

	if !reflect.DeepEqual(words, expected) {
		t.Fatal("inaccurate scanned word list")
	}

	t.Logf("Scanned words: %#v", words)

}
