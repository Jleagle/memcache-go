package memcache

import (
	"testing"

	"github.com/memcachier/mc/v3"
)

type test struct {
	Val1 int
	Val2 string
}

func TestSetGet(t *testing.T) {

	client := NewClient("localhost:11211")

	test1 := test{
		Val1: 1,
		Val2: "1",
	}

	// Set
	err := client.Set("key", test1, 10)
	if err != nil {
		t.Error(err)
	}

	// Get
	test2 := test{}
	err = client.Get("key", &test2)
	if err != nil {
		t.Error(err)
	}

	if test1.Val1 != test2.Val1 {
		t.Error("val1")
	}
	if test1.Val2 != test2.Val2 {
		t.Error("val2")
	}

}

func TestGetSet(t *testing.T) {

	client := NewClient("localhost:11211")

	test3 := test{}
	callback := func() (interface{}, error) {
		return test{Val1: 2, Val2: "2"}, nil
	}

	err := client.GetSet("key2", 10, &test3, callback)
	if err != nil {
		t.Error(err)
	}

	if test3.Val1 != 1 {
		t.Error("val1")
	}

	// Get
	test4 := test{}
	err = client.Get("key2", &test4)
	if err != nil {
		t.Error(err)
	}

	if test4.Val1 != 2 {
		t.Error("val1")
	}
}

func TestGetSetNoSet(t *testing.T) {

	client := NewClient("localhost:11211")

	test3 := test{}
	callback := func() (interface{}, error) {
		return test{Val1: 3, Val2: "3"}, ErrNoSet
	}

	err := client.GetSet("key3", 10, &test3, callback)
	if err != nil {
		t.Error(err)
	}

	if test3.Val1 != 1 {
		t.Error("val1")
	}

	// Get
	test4 := test{}
	err = client.Get("key3", &test4)
	if err != mc.ErrNotFound {
		t.Error(err)
	}
}

func TestGetDeleteGet(t *testing.T) {

	client := NewClient("localhost:11211")

	test := "test"

	// Set
	err := client.Set("key4", test, 10)
	if err != nil {
		t.Error(err)
	}

	// Get
	test2 := ""
	err = client.Get("key4", &test2)
	if err != nil {
		t.Error(err)
	}

	if test2 != "test" {
		t.Error("val1")
	}

	// Exists
	exists, err := client.Exists("key4")
	if err != nil {
		t.Error(err)
	}
	if !exists {
		t.Error("exists")
	}

	// Delete
	err = client.Delete("key4")
	if err != nil {
		t.Error(err)
	}

	// Get
	err = client.Get("key4", &test2)
	if err != mc.ErrNotFound {
		t.Error(err)
	}

	// Exists
	exists, err = client.Exists("key4")
	if err != nil {
		t.Error(err)
	}
	if exists {
		t.Error("exists")
	}
}
