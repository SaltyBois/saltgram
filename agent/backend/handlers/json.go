package handlers

import (
	"encoding/json"
	"io"
)

// NOTE(Jovan): Serializing to JSON
// NewEncoder provides better perf than json.Unmarshal
// https://golang.org/pkg/encoding/json/#NewEncoder
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

// NOTE(Jovan): Deserialize from JSON
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}
