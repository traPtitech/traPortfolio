package optional

import (
	"bytes"
	"database/sql"

	jsoniter "github.com/json-iterator/go"
)

type Int64 struct {
	sql.NullInt64
}

func Int64From(i int64) Int64 {
	return NewInt64(i, true)
}

func NewInt64(i int64, valid bool) Int64 {
	return Int64{
		NullInt64: sql.NullInt64{
			Int64: i,
			Valid: valid,
		},
	}
}

func (n *Int64) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		n.Int64, n.Valid = 0, false
		return nil
	}

	if err := jsoniter.ConfigFastest.Unmarshal(data, &n.Int64); err != nil {
		return err
	}

	n.Valid = true
	return nil
}

func (n Int64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return jsoniter.ConfigFastest.Marshal(n.Int64)
	}
	return jsoniter.ConfigFastest.Marshal(nil)
}
