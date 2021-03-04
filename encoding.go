package memcache

import (
	"encoding/json"
	"reflect"
)

type Encoder func(in interface{}) (out string, err error)

type Decoder func(in string, out interface{}) (err error)

func JSONEncoder(in interface{}) (string, error) {
	b, err := json.Marshal(in)
	return string(b), err
}
func JSONDecoder(in string, out interface{}) (err error) {
	return json.Unmarshal([]byte(in), out)
}

func StringEncoder(in interface{}) (out string, err error) {
	return in.(string), nil
}
func StringDecoder(in string, out interface{}) (err error) {
	reflect.ValueOf(out).Elem().Set(reflect.ValueOf(in))
	return nil
}
