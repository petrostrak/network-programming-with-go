package proxyingdatabetweenconnections

import (
	"io"
	"net"
)

// One of the most useful functions from the io package, the io.Copy
// function reads data from an io.Reader and writes it to an io.Writer.
// This is useful for creating a proxy, which, in this context, is an
// intermediary that transfers data between two nodes.

func proxyConn(source, destination string) error {

	// Create a connection to the source.
	connSource, err := net.Dial("tcp", source)
	if err != nil {
		return err
	}
	defer connSource.Close()

	// Create a connection to the destination.
	connDestination, err := net.Dial("tcp", destination)
	if err != nil {
		return err
	}
	defer connDestination.Close()

	// connDestination replies to connSource.
	go func() { _, _ = io.Copy(connSource, connDestination) }()

	// connSource messages to connDestination.
	_, err = io.Copy(connDestination, connSource)

	return err
}
