package protobuf

import (
	"io"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"github.com/petrostrak/network-programming-with-go/09-Data-Serialization/housework/v1"
)

func Load(r io.Reader) ([]*housework.Chore, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var chores housework.Chores

	return chores.Chores, proto.Unmarshal(b, &chores)
}

func Flush(w io.Writer, chores []*housework.Chore) error {
	b, err := proto.Marshal(&housework.Chores{Chores: chores})
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	if err != nil {
		return nil
	}

	return nil
}
