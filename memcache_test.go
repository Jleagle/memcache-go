package memcache

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/memcachier/mc/v3"
)

type test struct {
	Val1 int
	Val2 string
}

func TestSetGet(t *testing.T) {

	client := NewClient("localhost:11211")
	key := "TestSetGet-" + fmt.Sprintf("%d", time.Now().UnixNano())

	test1 := test{
		Val1: 1,
		Val2: "1",
	}

	// Set
	err := client.Set(key, test1, 10)
	if err != nil {
		t.Error(err)
	}

	// Get
	test2 := test{}
	err = client.Get(key, &test2)
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
	key := "TestGetSet-" + fmt.Sprintf("%d", time.Now().UnixNano())

	test3 := test{}
	callback := func() (test, error) {
		return test{Val1: 2, Val2: "2"}, nil
	}

	err := GetSet(client, key, 10, &test3, callback)
	if err != nil {
		t.Error(err)
	}

	if test3.Val1 != 2 {
		t.Error("val1")
	}

	// Get
	test4 := test{}
	err = client.Get(key, &test4)
	if err != nil {
		t.Error(err)
	}

	if test4.Val1 != 2 {
		t.Error("val1")
	}
}

func TestGetSetNoSet(t *testing.T) {

	client := NewClient("localhost:11211")
	key := "TestGetSetNoSet-" + fmt.Sprintf("%d", time.Now().UnixNano())

	test3 := test{}
	callback := func() (test, error) {
		return test{Val1: 3, Val2: "3"}, ErrNoSet
	}

	err := GetSet(client, key, 10, &test3, callback)
	if err != nil {
		t.Error(err)
	}

	if test3.Val1 != 3 {
		t.Error("val1")
	}

	// Get
	test4 := test{}
	err = client.Get(key, &test4)
	if !errors.Is(err, mc.ErrNotFound) {
		t.Error(err)
	}
}

func TestGetDeleteGet(t *testing.T) {

	client := NewClient("localhost:11211")
	key := "TestGetDeleteGet-" + fmt.Sprintf("%d", time.Now().UnixNano())

	test := "test"

	// Set
	err := client.Set(key, test, 10)
	if err != nil {
		t.Error(err)
	}

	// Get
	test2 := ""
	err = client.Get(key, &test2)
	if err != nil {
		t.Error(err)
	}

	if test2 != "test" {
		t.Error("val1")
	}

	// Exists
	exists, err := client.Exists(key)
	if err != nil {
		t.Error(err)
	}
	if !exists {
		t.Error("exists")
	}

	// Delete
	err = client.Delete(key)
	if err != nil {
		t.Error(err)
	}

	// Get
	err = client.Get(key, &test2)
	if !errors.Is(err, mc.ErrNotFound) {
		t.Error(err)
	}

	// Exists
	exists, err = client.Exists(key)
	if err != nil {
		t.Error(err)
	}
	if exists {
		t.Error("exists")
	}
}

func TestNils(t *testing.T) {

	client := NewClient("localhost:9002")
	key := "TestTestNils-" + fmt.Sprintf("%d", time.Now().UnixNano())

	var test3 []byte
	callback := func() ([]byte, error) {
		return nil, nil
	}

	err := GetSet(client, key, 10, &test3, callback)
	if err != nil {
		t.Error(err)
	}

	if test3 != nil {
		t.Error("val1")
	}
}
