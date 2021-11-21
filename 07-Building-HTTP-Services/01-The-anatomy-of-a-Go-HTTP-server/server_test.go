package theanatomyofagohttpserver

import (
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/petrostrak/network-programming-with-go/07-Building-HTTP-Services/handlers"
)

func TestSimpleHTTPServer(t *testing.T) {
	srv := &http.Server{
		Addr: "127.0.0.1:8081",
		Handler: http.TimeoutHandler(
			handlers.DefaultHandler(), 2*time.Minute, ""),
		IdleTimeout:       5 * time.Minute,
		ReadHeaderTimeout: time.Minute,
	}

	l, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		err := srv.Serve(l)
		if err != http.ErrServerClosed {
			t.Error(err)
		}
	}()
}
