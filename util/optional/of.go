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

func (o Of[T]) Valid() bool {
	return o.valid
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
	if src == nil {
		var t T
		o.v, o.valid = t, false
		return nil
	}

	var scanner sql.Scanner
	switch v := any(o.v).(type) {
	case sql.Scanner:
		scanner = v
	case string:
		scanner = &sql.NullString{}
	case int64:
		scanner = &sql.NullInt64{}
	case int32:
		scanner = &sql.NullInt32{}
	case int16:
		scanner = &sql.NullInt16{}
	case byte:
		scanner = &sql.NullByte{}
	case float64:
		scanner = &sql.NullFloat64{}
	case bool:
		scanner = &sql.NullBool{}
	case time.Time:
		scanner = &sql.NullTime{}
	default:
		return fmt.Errorf("unsupported type for Scan: %T", v)
	}

	if err := scanner.Scan(src); err != nil {
		return err
	}

	switch v := scanner.(type) {
	case *sql.NullString:
		o.v, o.valid = any(v.String).(T), v.Valid
	case *sql.NullInt64:
		o.v, o.valid = any(v.Int64).(T), v.Valid
	case *sql.NullInt32:
		o.v, o.valid = any(v.Int32).(T), v.Valid
	case *sql.NullInt16:
		o.v, o.valid = any(v.Int16).(T), v.Valid
	case *sql.NullByte:
		o.v, o.valid = any(v.Byte).(T), v.Valid
	case *sql.NullFloat64:
		o.v, o.valid = any(v.Float64).(T), v.Valid
	case *sql.NullBool:
		o.v, o.valid = any(v.Bool).(T), v.Valid
	case *sql.NullTime:
		o.v, o.valid = any(v.Time).(T), v.Valid
	default:
		o.valid = true
	}

	return nil
}

func (o Of[T]) Value() (driver.Value, error) {
	if !o.valid {
		return nil, nil
	}

	var valuer driver.Valuer
	switch v := any(o.v).(type) {
	case driver.Valuer:
		valuer = v
	case string:
		valuer = sql.NullString{String: v, Valid: true}
	case int64:
		valuer = sql.NullInt64{Int64: v, Valid: true}
	case int32:
		valuer = sql.NullInt32{Int32: v, Valid: true}
	case int16:
		valuer = sql.NullInt16{Int16: v, Valid: true}
	case byte:
		valuer = sql.NullByte{Byte: v, Valid: true}
	case float64:
		valuer = sql.NullFloat64{Float64: v, Valid: true}
	case bool:
		valuer = sql.NullBool{Bool: v, Valid: true}
	case time.Time:
		valuer = sql.NullTime{Time: v, Valid: true}
	default:
		return nil, fmt.Errorf("unsupported type for Value: %T", v)
	}

	return valuer.Value()
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
