package unixdomainsockets

import (
	"context"
	"net"
	"os"
)

func streamingEchoServer(ctx context.Context, network, addr string) (net.Addr, error) {
	s, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}

	go func() {

		// Once the caller cancels the context, the server closes.
		go func() {
			<-ctx.Done()
			_ = s.Close()
		}()

		for {

			// Once the server accepts a connection, we start a goroutine to
			// echo incoming messages.
			conn, err := s.Accept()
			if err != nil {
				return
			}

			go func() {
				defer func() { _ = conn.Close() }()

				// Since we are using the net.Conn interface, we can use the
				// Read and Write methods to communicate with the client no
				// matter whether the server is communicating over a network
				// socket or a Unix domain socket.
				for {
					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if err != nil {
						return
					}

					_, err = conn.Write(buf[:n])
					if err != nil {
						return
					}
				}
			}()
		}
	}()

	return s.Addr(), nil
}

func datagramEchoServer(ctx context.Context, network, addr string) (net.Addr, error) {
	s, err := net.ListenPacket(network, addr)
	if err != nil {
		return nil, err
	}

	go func() {
		go func() {
			<-ctx.Done()
			_ = s.Close()
			if network == "unixgram" {
				_ = os.Remove(addr)
			}
		}()

		buf := make([]byte, 1024)
		for {
			n, clientAddr, err := s.ReadFrom(buf)
			if err != nil {
				return
			}

			_, err = s.WriteTo(buf[:n], clientAddr)
			if err != nil {
				return
			}
		}
	}()

	return s.LocalAddr(), nil
}
