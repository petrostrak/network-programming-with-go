package usingaudpechoserver

import (
	"context"
	"fmt"
	"net"
)

func echoServerUDP(ctx context.Context, addr string) (net.Addr, error) {
	s, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("binding to udp %s: %w", addr, err)
	}

	// Once the caller cancels the context, receiving on the Done channel
	// unblocks and the server is closed, tearing down both this goroutine
	// and the parent goroutine.
	go func() {

		// The second goroutine blocks on the context's Done channel.
		go func() {
			<-ctx.Done()
			_ = s.Close()
		}()

		buf := make([]byte, 1024)

		for {

			// To read from the UDP connection, pass a byte of slice to the ReadFrom().
			// This returns the number of bytes read, the sender's address and an error
			// interface.
			n, clientAddr, err := s.ReadFrom(buf) // client to server
			if err != nil {
				return
			}

			// To write a UDP packet, pass a byte slice and a destination address to the
			// connection's WriteTo(). It returns the number of bytes written and an error
			// interface.
			_, err = s.WriteTo(buf[:n], clientAddr) // server to client
			if err != nil {
				return
			}
		}
	}()

	return s.LocalAddr(), nil
}
