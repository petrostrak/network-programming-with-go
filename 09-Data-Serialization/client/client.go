package client

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	addr, caCertFn string
)

func init() {
	flag.StringVar(&addr, "address", "localhost:34443", "server address")
	flag.StringVar(&caCertFn, "ca-cert", "cert.pem", "CA certificate")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			`Usage: %s [flags] [add chore, ...|complete #]
	add			add coma-separated chores
	complete	complete designated chore
	
Flags:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}
