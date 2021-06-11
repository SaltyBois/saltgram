package data

import (
	"bytes"
	"encoding/json"
)

type CertType int

const (
	Root         CertType = 10
	Intermediary CertType = 5
	EndEntity    CertType = 1
)

func (ct CertType) String() string {
	return toString[ct]
}

var toString = map[CertType]string{
	Root:         "Root",
	Intermediary: "Intermediary",
	EndEntity:    "EndEntity",
}

var toID = map[string]CertType{
	"Root":         Root,
	"Intermediary": Intermediary,
	"EndEntity":    EndEntity,
}

func (ct CertType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toString[ct])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (ct *CertType) UnmarshalJSON(b []byte) error {
	var jsonString string
	err := json.Unmarshal(b, &jsonString)
	if err != nil {
		return err
	}
	*ct = toID[jsonString]
	return nil
}