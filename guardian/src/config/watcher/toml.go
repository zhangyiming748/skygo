package watcher

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"reflect"
	"strconv"
)

// TOML is toml map.
type TOML = Map

// Set set the map by value.
func (m *TOML) Set(text string) error {
	if err := m.UnmarshalText([]byte(text)); err != nil {
		return err
	}
	return nil
}

// UnmarshalText implemented toml.
func (m *TOML) UnmarshalText(text []byte) error {
	raws := map[string]interface{}{}
	if err := toml.Unmarshal(text, &raws); err != nil {
		return err
	}
	values := map[string]*Value{}
	for k, v := range raws {
		k = keyNamed(k)
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Map:
			buf := bytes.NewBuffer(nil)
			err := toml.NewEncoder(buf).Encode(v)
			// b, err := toml.Marshal(v)
			if err != nil {
				return err
			}
			// NOTE: value is map[string]interface{}
			values[k] = &Value{val: v, raw: buf.String()}
		case reflect.Slice:
			raw := map[string]interface{}{
				k: v,
			}
			buf := bytes.NewBuffer(nil)
			err := toml.NewEncoder(buf).Encode(raw)
			// b, err := toml.Marshal(raw)
			if err != nil {
				return err
			}
			// NOTE: value is []interface{}
			values[k] = &Value{val: v, raw: buf.String()}
		case reflect.Bool:
			b := v.(bool)
			values[k] = &Value{val: b, raw: strconv.FormatBool(b)}
		case reflect.Int64:
			i := v.(int64)
			values[k] = &Value{val: i, raw: strconv.FormatInt(i, 10)}
		case reflect.Float64:
			f := v.(float64)
			values[k] = &Value{val: f, raw: strconv.FormatFloat(f, 'f', -1, 64)}
		case reflect.String:
			s := v.(string)
			values[k] = &Value{val: s, raw: s}
		default:
			return errors.New(fmt.Sprintf("UnmarshalTOML: unknown kind(%v)", rv.Kind()))
		}
	}
	m.Store(values)
	return nil
}
