package thelogpackage

import (
	"log"
	"os"
)

func Example_log() {
	l := log.New(os.Stdout, "example:", log.Lshortfile)
	l.Print("logging to standard output")
}
