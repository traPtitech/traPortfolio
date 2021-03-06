package optional

import (
	"bytes"
	"database/sql"

	jsoniter "github.com/json-iterator/go"
)

type String struct {
	sql.NullString
}

func (s *String) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		s.String, s.Valid = "", false
		return nil
	}

	if err := jsoniter.ConfigFastest.Unmarshal(data, &s.String); err != nil {
		return err
	}

	s.Valid = true
	return nil
}

func (s String) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return jsoniter.ConfigFastest.Marshal(s.String)
	}
	return jsoniter.ConfigFastest.Marshal(nil)
}
