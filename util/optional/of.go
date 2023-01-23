package optional

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// Of nullableなjsonフィールドとして使用できます。
// json.Unmarshaler, json.Marshaler, encoding.TextUnmarshaler, encoding.TextMarshaler, sql.Scanner, driver.Valuer
// を実装する型と、一部の型についてこれらのメソッドを使用できます。
type Of[T any] struct {
	V     T
	Valid bool
}

func New[T any](v T, valid bool) Of[T] {
	return Of[T]{
		V:     v,
		Valid: valid,
	}
}

func From[T any](v T) Of[T] {
	return Of[T]{
		V:     v,
		Valid: true,
	}
}

func FromPtr[T any](v *T) Of[T] {
	if v == nil {
		return Of[T]{}
	}

	return Of[T]{
		V:     *v,
		Valid: true,
	}
}

func (o Of[T]) ValueOrZero() T {
	if o.Valid {
		return o.V
	}
	var t T
	return t
}

func (o *Of[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		var t T
		o.V, o.Valid = t, false
		return nil
	}

	if u, ok := any(&o.V).(json.Unmarshaler); ok {
		if err := u.UnmarshalJSON(data); err != nil {
			return err
		}
		o.Valid = true
		return nil
	}
	if err := jsoniter.ConfigFastest.Unmarshal(data, &o.V); err != nil {
		return err
	}
	o.Valid = true
	return nil
}

func (o Of[T]) MarshalJSON() ([]byte, error) {
	if !o.Valid {
		return jsoniter.ConfigFastest.Marshal(nil)
	}
	if m, ok := any(o.V).(json.Marshaler); ok {
		return m.MarshalJSON()
	}
	return jsoniter.ConfigFastest.Marshal(o.V)
}

func (o *Of[T]) UnmarshalText(data []byte) error {
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		var t T
		o.V = t
		o.Valid = false
		return nil
	}

	var valid bool
	switch any(o.V).(type) {
	case bool:
		var v bool
		s := string(data)
		switch s {
		case "", "null":
			v, valid = false, false
		case "true":
			v, valid = true, true
		case "false":
			v, valid = false, true
		default:
			return fmt.Errorf("invalid bool value: %s", s)
		}
		o.V, o.Valid = any(v).(T), valid
		return nil
	case int:
		var v int
		s := string(data)
		switch s {
		case "", "null":
			v, valid = 0, false
		default:
			n, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			v, valid = n, true
		}
		o.V, o.Valid = any(v).(T), valid
		return nil
	case string:
		var v string
		s := string(data)
		switch s {
		case "", "null":
			v, valid = "", false
		default:
			v, valid = s, true
		}
		o.V, o.Valid = any(v).(T), valid
		return nil
	default:
		t, ok := any(&o.V).(encoding.TextUnmarshaler)
		if !ok {
			return fmt.Errorf("unsupported type for UnmarshalText: %T", t)
		}
		if err := t.UnmarshalText(data); err != nil {
			return err
		}
		o.Valid = true
		return nil
	}
}

func (o Of[T]) MarshalText() ([]byte, error) {
	if !o.Valid {
		return []byte{}, nil
	}

	switch v := any(o.V).(type) {
	case bool:
		if v {
			return []byte("true"), nil
		}
		return []byte("false"), nil
	case int:
		return []byte(strconv.FormatInt(int64(v), 10)), nil
	case string:
		return []byte(v), nil
	default:
		t, ok := v.(encoding.TextMarshaler)
		if !ok {
			return nil, fmt.Errorf("unsupported type for MarshalText: %T", t)
		}
		return t.MarshalText()
	}
}

func (o *Of[T]) Scan(src any) error {
	switch any(o.V).(type) {
	case bool:
		var b sql.NullBool
		if err := b.Scan(src); err != nil {
			return err
		}
		o.V, o.Valid = any(b.Bool).(T), b.Valid
		return nil
	case int:
		var i sql.NullInt64
		if err := i.Scan(src); err != nil {
			return err
		}
		o.V, o.Valid = any(int(i.Int64)).(T), i.Valid
		return nil
	case string:
		var s sql.NullString
		if err := s.Scan(src); err != nil {
			return err
		}
		o.V, o.Valid = any(s.String).(T), s.Valid
		return nil
	case time.Time:
		var t sql.NullTime
		if err := t.Scan(src); err != nil {
			return err
		}
		o.V, o.Valid = any(t.Time).(T), t.Valid
		return nil
	default:
		s, ok := any(&o.V).(sql.Scanner)
		if !ok {
			return fmt.Errorf("unsupported type for Scan: %T", o.V)
		}
		if err := s.Scan(src); err != nil {
			return err
		}
		o.Valid = true
		return nil
	}
}

func (o Of[T]) Value() (driver.Value, error) {
	if !o.Valid {
		return nil, nil
	}
	switch v := any(o.V).(type) {
	case int:
		return int64(v), nil
	case driver.Valuer:
		return v.Value()
	}
	return o.V, nil
}

// interface guards
var (
	_ json.Marshaler           = (*Of[any])(nil)
	_ json.Unmarshaler         = (*Of[any])(nil)
	_ encoding.TextMarshaler   = (*Of[any])(nil)
	_ encoding.TextUnmarshaler = (*Of[any])(nil)
	_ sql.Scanner              = (*Of[any])(nil)
	_ driver.Valuer            = (*Of[any])(nil)
)
