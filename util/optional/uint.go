package optional

import (
	"bytes"
	"database/sql"

	jsoniter "github.com/json-iterator/go"
)

type Uint struct {
	sql.NullByte
}

func UintFrom(i *uint) Uint {
	if i == nil {
		return Uint{}
	}

	return NewUint(*i, true)
}

func NewUint(i uint, valid bool) Uint {
	return Uint{
		NullByte: sql.NullByte{
			Byte:  byte(i),
			Valid: valid,
		},
	}
}

func (n *Uint) UnmarshalJSON(data []byte) error {
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

func (n Uint) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return jsoniter.ConfigFastest.Marshal(n.Byte)
	}
	return jsoniter.ConfigFastest.Marshal(nil)
}
