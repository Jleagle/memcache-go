package memcache

import (
	"encoding/json"
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
	"io"
	"reflect"
)

//noinspection GoUnusedGlobalVariable
var (
	ErrCacheMiss    = memcache.ErrCacheMiss
	ErrNotPointer   = errors.New("value ust be a pointer")
	ErrInvalidTypes = errors.New("types must match")
)

type Item = memcache.Item

func New(namespace string, servers ...string) *Memcache {

	// Fallback DSN
	if len(servers) == 0 {
		servers = []string{"localhost:11211"}
	}

	mc := new(Memcache)
	mc.client = memcache.New(servers...)
	mc.namespace = namespace

	return mc
}

type Memcache struct {
	namespace string
	client    *memcache.Client
}

// Returns []byte
func (mc Memcache) Get(key string, i interface{}) error {

	item, err := mc.client.Get(mc.namespace + key)
	if err != nil {
		return err
	}

	return json.Unmarshal(item.Value, i)
}

func (mc Memcache) Set(key string, value interface{}, expiration int32) error {

	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	item := new(memcache.Item)
	item.Key = mc.namespace + key
	item.Value = bytes
	item.Expiration = expiration

	return mc.client.Set(item)
}

func (mc Memcache) SetItem(item memcache.Item) error {
	return mc.Set(item.Key, item.Value, item.Expiration)
}

func (mc Memcache) GetSet(key string, expiration int32, value interface{}, f func() (interface{}, error)) error {

	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		return ErrNotPointer
	}

	err := mc.Get(key, value)

	if err == memcache.ErrCacheMiss || err == io.EOF {

		s, err := f()
		if err != nil {
			return err
		}

		if reflect.TypeOf(s) != reflect.TypeOf(value).Elem() {
			return ErrInvalidTypes
		}

		err = setToPointer(s, value)

		return mc.Set(key, s, expiration)
	}

	return err
}

func (mc Memcache) Delete(item memcache.Item) (err error) {

	return mc.client.Delete(mc.namespace + item.Key)
}

func (mc Memcache) DeleteAll() (err error) {

	return mc.client.DeleteAll()
}

func (mc Memcache) Increment(key string, delta ...uint64) (err error) {

	if len(delta) == 0 {
		delta = []uint64{1}
	}

	_, err = mc.client.Increment(mc.namespace+key, delta[0])

	return err
}

func (mc Memcache) Decrement(key string, delta ...uint64) (err error) {

	if len(delta) == 0 {
		delta = []uint64{1}
	}

	_, err = mc.client.Decrement(mc.namespace+key, delta[0])

	return err
}

func setToPointer(in interface{}, out interface{}) error {

	b, err := json.Marshal(in)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, out)
}
