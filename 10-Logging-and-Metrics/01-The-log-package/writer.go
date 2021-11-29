package thelogpackage

import (
	"io"

	"go.uber.org/multierr"
)

type sustainedMultiWriter struct {
	writers []io.Writer
}

func (s *sustainedMultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range s.writers {
		i, wErr := w.Write(p)
		n += i
		err = multierr.Append(err, wErr)
	}

	return n, err
}
