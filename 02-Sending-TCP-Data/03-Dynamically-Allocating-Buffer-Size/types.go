package dynamicallyallocatingbuffersize

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	BinaryType uint8 = iota + 1
	StringType
	MaxPayloadSize uint32 = 10 << 20 // 10 MB
)

var ErrMaxPayloadSize = errors.New("maximum payload size exceeded")

type Payload interface {
	fmt.Stringer
	io.ReaderFrom
	io.WriterTo
	Bytes() []byte
}

type Binary []byte

// Bytes returns itself.
func (m Binary) Bytes() []byte { return m }

// String cast itself as a string before returning.
func (m Binary) String() string { return string(m) }

// WriteTo returns the number of bytes written to the writer.
func (m Binary) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, BinaryType)
	if err != nil {
		return 0, err
	}

	var n int64 = 1

	err = binary.Write(w, binary.BigEndian, uint32(len(m))) // 4-byte size.
	if err != nil {
		return n, err
	}

	n += 4

	o, err := w.Write(m) // payload

	return n + int64(o), err
}
