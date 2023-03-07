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
	v     T
	valid bool
}

func New[T any](v T, valid bool) Of[T] {
	return Of[T]{
		v:     v,
		valid: valid,
	}
}

func From[T any](v T) Of[T] {
	if fmt.Sprintf("%T", v)[0] == '*' {
		panic("optional: From[T](v T): T must not be a pointer type, use FromPtr[T](v *T) instead")
	}

	return Of[T]{
		v:     v,
		valid: true,
	}
}

func FromPtr[T any](v *T) Of[T] {
	if v == nil {
		return Of[T]{}
	}

	return Of[T]{
		v:     *v,
		valid: true,
	}
}

func (o Of[T]) ValueOrZero() T {
	if o.valid {
		return o.v
	}
	var t T
	return t
}

func (o Of[T]) V() (T, bool) {
	return o.v, o.valid
}

func (o *Of[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		var t T
		o.v, o.valid = t, false
		return nil
	}

	if err := jsoniter.ConfigFastest.Unmarshal(data, &o.v); err != nil {
		return err
	}
	o.valid = true
	return nil
}

func (o Of[T]) MarshalJSON() ([]byte, error) {
	if !o.valid {
		return jsoniter.ConfigFastest.Marshal(nil)
	}

	return jsoniter.ConfigFastest.Marshal(o.v)
}

func (o *Of[T]) UnmarshalText(data []byte) error {
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		var t T
		o.v = t
		o.valid = false
		return nil
	}

	var valid bool
	switch any(o.v).(type) {
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
		o.v, o.valid = any(v).(T), valid
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
		o.v, o.valid = any(v).(T), valid
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
		o.v, o.valid = any(v).(T), valid
		return nil
	default:
		t, ok := any(&o.v).(encoding.TextUnmarshaler)
		if !ok {
			return fmt.Errorf("unsupported type for UnmarshalText: %T", t)
		}
		if err := t.UnmarshalText(data); err != nil {
			return err
		}
		o.valid = true
		return nil
	}
}

func (o Of[T]) MarshalText() ([]byte, error) {
	if !o.valid {
		return []byte{}, nil
	}

	switch v := any(o.v).(type) {
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
	switch any(o.v).(type) {
	case bool:
		var b sql.NullBool
		if err := b.Scan(src); err != nil {
			return err
		}
		o.v, o.valid = any(b.Bool).(T), b.Valid
		return nil
	case int:
		var i sql.NullInt64
		if err := i.Scan(src); err != nil {
			return err
		}
		o.v, o.valid = any(int(i.Int64)).(T), i.Valid
		return nil
	case string:
		var s sql.NullString
		if err := s.Scan(src); err != nil {
			return err
		}
		o.v, o.valid = any(s.String).(T), s.Valid
		return nil
	case time.Time:
		var t sql.NullTime
		if err := t.Scan(src); err != nil {
			return err
		}
		o.v, o.valid = any(t.Time).(T), t.Valid
		return nil
	default:
		s, ok := any(&o.v).(sql.Scanner)
		if !ok {
			return fmt.Errorf("unsupported type for Scan: %T", o.v)
		}
		if err := s.Scan(src); err != nil {
			return err
		}
		o.valid = true
		return nil
	}
}

func (o Of[T]) Value() (driver.Value, error) {
	if !o.valid {
		return nil, nil
	}
	switch v := any(o.v).(type) {
	case int:
		return int64(v), nil
	case driver.Valuer:
		return v.Value()
	}
	return o.v, nil
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
