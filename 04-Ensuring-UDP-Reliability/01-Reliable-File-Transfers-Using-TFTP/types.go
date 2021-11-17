package reliablefiletransfersusingtftp

import (
	"bytes"
	"encoding/binary"
)

// Your read-only TFTP server will accept requests from the client, send data
// packets, transmit error packets, and accept acknowledgments from the client.

const (
	DatagramSize = 516              // The maximum supported datagram size (to avoid fragmentation)
	BlockSize    = DatagramSize - 4 // The DatagramSize minus a 4-byte header
)

// The first 2 bytes of a TFTP packet's header is an operation code.
// Each OpCode is an 2-byte unsigned integer. This server supports
// four operations: a read request (RPQ), a data operation, an
// acknowledgment, and an error.Since this server is a read-only,
// we skip the write request (WRQ) definition.
type OpCode uint16

const (
	OpRRQ OpCode = iota + 1
	_            // no WRQ support
	OpData
	OpAck
	OpErr
)

// Define a series of unsigned 16-bit integer error codes per the RFC.
type ErrCode uint16

const (
	ErrUnknown ErrCode = iota
	ErrNotFound
	ErrAccessViolation
	ErrDiskFull
	ErrIllegalOp
	ErrUnknownID
	ErrFileExists
	ErrNoUser
)

// The read request that needs to keep track of the filename and the mode.
type ReadReq struct {
	Filename string
	Mode     string
}

// Although not used by our server, a client would make use of this method.
func (q ReadReq) MarshalBinary() ([]byte, error) {
	mode := "octet"
	if q.Mode != "" {
		mode = q.Mode
	}

	// operation code + filename + 0 byte + mode + 0 byte
	cap := 2 + 2 + len(q.Filename) + 1 + len(q.Mode) + 1

	b := new(bytes.Buffer)
	b.Grow(cap)

	// Insert the operation code...
	err := binary.Write(b, binary.BigEndian, OpRRQ) // write operation code
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(q.Filename) // write Filename
	if err != nil {
		return nil, err
	}

	// ...and null bytes into the buffer while marshaling the packet to a
	// byte slice.
	err = b.WriteByte(0) // write 0 byte
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(mode) // write Mode
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0) // write 0 byte
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
