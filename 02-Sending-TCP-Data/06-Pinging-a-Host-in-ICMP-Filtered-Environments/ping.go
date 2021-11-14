package pingingahostinicmpfilteredenvironments

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var (
	count    = flag.Int("c", 3, "number of pings: <= 0 means forever")
	interval = flag.Duration("i", time.Second, "interval between pings")
	timeout  = flag.Duration("W", 5*time.Second, "time to wait for a reply")
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] host:port\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Print("host:port is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	target := flag.Arg(0)
	fmt.Println("PING", target)

	if *count <= 0 {
		fmt.Println("CTRL+C to stop.")
	}

	msg := 0

	for (*count <= 0) || (msg < *count) {
		msg++
		fmt.Print(msg, " ")

		start := time.Now()

		// Attempt to establish a connection to a remote host's TCP port.
		c, err := net.DialTimeout("tcp", target, *timeout)

		// Keep track on the time it takes to complete the TCP handshare and
		// consider this duration the ping interval between the host and the
		// remote host.
		dur := time.Since(start)
		if err != nil {
			fmt.Printf("fail in %s: %v\n", dur, err)
			if nErr, ok := err.(net.Error); !ok || !nErr.Temporary() {
				os.Exit(1)
			} else {
				_ = c.Close()
				fmt.Println(dur)
			}

			time.Sleep(*interval)
		}
	}
}
