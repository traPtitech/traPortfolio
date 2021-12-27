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

func TimeFrom(t *time.Time) Time {
	if t == nil {
		return Time{}
	}

	return NewTime(*t, true)
}

func NewTime(t time.Time, valid bool) Time {
	return Time{NullTime: sql.NullTime{Time: t, Valid: valid}}
}

func (t *Time) ValueOrZero() (zero time.Time) {
	if t.Valid {
		return t.Time
	}

	return
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
