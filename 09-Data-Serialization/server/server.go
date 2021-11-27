package server

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/petrostrak/network-programming-with-go/09-Data-Serialization/housework/v1"
	"google.golang.org/grpc"
)

var (
	addr, certF, keyFn string
)

func init() {
	flag.StringVar(&addr, "address", "localhost:3443", "listen address")
	flag.StringVar(&certF, "cert", "cert.pem", "certificate file")
	flag.StringVar(&keyFn, "key", "key.pem", "private key file")
}

func main() {
	flag.Parse()

	server := grpc.NewServer()
	rosie := new(Rosie)
	housework.RegisterRobotMaidServer(server, *rosie.Service())

	cert, err := tls.LoadX509KeyPair(certF, keyFn)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening for TLS connections on %s ...", addr)
	log.Fatal(err)
	server.Serve(
		tls.NewListener(
			listener,
			&tls.Config{
				Certificates:             []tls.Certificate{cert},
				CurvePreferences:         []tls.CurveID{tls.CurveP256},
				MinVersion:               tls.VersionTLS12,
				PreferServerCipherSuites: true,
			},
		),
	)
}
