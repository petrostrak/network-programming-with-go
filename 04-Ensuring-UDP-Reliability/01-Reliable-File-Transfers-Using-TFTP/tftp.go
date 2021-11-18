package reliablefiletransfersusingtftp

import (
	"flag"
	"io/ioutil"
	"log"
)

var (
	address = flag.String("a", "127.0.0.1:69", "listen address")
	payload = flag.String("p", "payload.svg", "file to serve to clients")
)

func main() {
	flag.Parse()

	p, err := ioutil.ReadFile(*payload)
	if err != nil {
		log.Fatal(err)
	}

	s := Server{Payload: p}
	log.Fatal(s.ListenAndServe(*address))
}
