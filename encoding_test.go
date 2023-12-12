package memcache

import (
	"encoding/json"
	"errors"
	"testing"
)

type testStruct struct {
	A string
	B int
	C *func()
}

func TestEncoders(t *testing.T) {

	s := testStruct{A: "a", B: 2}
	r, err := JSONEncoder(s)
	if err != nil {
		t.Error(err)
	}
	if r != `{"A":"a","B":2,"C":null}` {
		t.Error("JSONEncoder")
	}
	x := func() {}
	s.C = &x
	r, err = JSONEncoder(s)
	var jsonErr *json.UnsupportedTypeError
	if !errors.As(err, &jsonErr) {
		t.Error("JSONEncoder")
	}
}

func TestDecoders(t *testing.T) {

	in := `{"A":"a","B":2}`
	out := testStruct{}
	err := JSONDecoder(in, &out)
	if err != nil {
		t.Error(err)
	}
	if out.A != "a" {
		t.Error("JSONDecoder")
	}
	if out.B != 2 {
		t.Error("JSONDecoder")
	}
}
