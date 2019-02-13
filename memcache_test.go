package main

import (
	"github.com/Jleagle/memcache-go/memcache"
	"testing"
)

func Test(t *testing.T) {

	// Set
	item := memcache.Item{
		Key:        "test",
		Value:      []byte("value"),
		Expiration: 10,
	}

	mc := memcache.New("")

	err := mc.SetItem(&item)
	if err != nil {
		t.Error(err)
	}

	// Get
	var b []byte

	err = mc.Get(item.Key, &b)
	if err != nil {
		t.Error(err)
	}

	if string(b) != string(item.Value) {
		t.Error("wrong value")
	}

	// Get Set
	item2 := memcache.Item{
		Key:        "test2",
		Value:      []byte("value2"),
		Expiration: 10,
	}

	var b2 []byte

	err = mc.GetSet(item2.Key, item2.Expiration, &b2, func() (k interface{}, err error) {

		return []byte("value from database"), nil

	})

	if err != nil {
		t.Error(err)
	}

	if string(b2) != "value from database" {
		t.Error("wrong value")
	}
}
