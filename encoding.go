package memcache

import (
	"encoding/json"
	"reflect"
)

type Encoder func(in any) (out string, err error)

type Decoder func(in string, out any) (err error)

func JSONEncoder(in any) (string, error) {
	b, err := json.Marshal(in)
	return string(b), err
}
func JSONDecoder(in string, out any) (err error) {
	return json.Unmarshal([]byte(in), out)
}

func StringEncoder(in any) (out string, err error) {
	return in.(string), nil
}
func StringDecoder(in string, out any) (err error) {
	reflect.ValueOf(out).Elem().Set(reflect.ValueOf(in))
	return nil
}
