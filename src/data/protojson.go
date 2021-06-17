package data

import (
	"io"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func ToPtorJSON(m proto.Message, w io.Writer) (int, error) {
	msopt := protojson.MarshalOptions{
		Indent:          " ",
		EmitUnpopulated: true,
	}
	json, err := msopt.Marshal(m)
	if err != nil {
		return 0, err
	}
	return w.Write(json)
}
