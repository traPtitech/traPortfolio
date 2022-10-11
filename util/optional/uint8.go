package optional

import (
	"bytes"
	"database/sql"

	jsoniter "github.com/json-iterator/go"
)

type Uint8 struct {
	sql.NullByte
}

func Uint8From(i *uint8) Uint8 {
	if i == nil {
		return Uint8{}
	}

	return NewUint8(*i, true)
}

func NewUint8(i uint8, valid bool) Uint8 {
	return Uint8{
		NullByte: sql.NullByte{
			Byte:  byte(i),
			Valid: valid,
		},
	}
}

func (n *Uint8) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		n.Byte, n.Valid = 0, false
		return nil
	}

	if err := jsoniter.ConfigFastest.Unmarshal(data, &n.Byte); err != nil {
		return err
	}

	n.Valid = true
	return nil
}

func (n Uint8) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return jsoniter.ConfigFastest.Marshal(n.Byte)
	}
	return jsoniter.ConfigFastest.Marshal(nil)
}
