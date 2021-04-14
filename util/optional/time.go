package optional

import (
	"bytes"
	"database/sql"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Time struct {
	sql.NullTime
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}

	if err := jsoniter.ConfigFastest.Unmarshal(data, &t.Time); err != nil {
		return err
	}

	t.Valid = true
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return t.Time.MarshalJSON()
	}
	return jsoniter.ConfigFastest.Marshal(nil)
}
