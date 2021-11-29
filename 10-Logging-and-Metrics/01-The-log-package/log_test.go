package thelogpackage

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func Example_log() {
	l := log.New(os.Stdout, "example:", log.Lshortfile)
	l.Print("logging to standard output")
}

func Example_logMultiWriter() {
	logFile := new(bytes.Buffer)

	w := SustainedMultiWriter(os.Stdout, logFile)
	l := log.New(w, "example:", log.Lshortfile|log.Lmsgprefix)

	fmt.Println("standard output:")
	l.Print("This is Peter")

	fmt.Print("\nlog file contents:\n", logFile.String())
}
