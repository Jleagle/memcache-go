package memcache

import (
	"encoding/json"
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/cenkalti/backoff"
	"io"
	"reflect"
	"time"
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

	policy := backoff.NewExponentialBackOff()
	policy.InitialInterval = 100 * time.Millisecond

	mc := new(Memcache)
	mc.client = memcache.New(servers...)
	mc.namespace = namespace
	mc.backoff = backoff.WithMaxRetries(policy, 5)

	return mc
}

type Memcache struct {
	namespace string
	client    *memcache.Client
	backoff   backoff.BackOff
}

func (mc Memcache) SetBackoff(backoff backoff.BackOff) {
	mc.backoff = backoff
}

// Returns []byte
func (mc Memcache) Get(key string, i interface{}) (err error) {

	var item *memcache.Item

	operation := func() (err error) {

		item, err = mc.client.Get(mc.namespace + key)
		return err
	}

	err = backoff.Retry(operation, mc.backoff)
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

	operation := func() (err error) {
		return mc.client.Set(item)
	}

	return backoff.Retry(operation, mc.backoff)
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

	operation := func() (err error) {
		return mc.client.Delete(mc.namespace + item.Key)
	}

	return backoff.Retry(operation, mc.backoff)
}

func (mc Memcache) DeleteAll() (err error) {

	operation := func() (err error) {
		return mc.client.DeleteAll()
	}

	return backoff.Retry(operation, mc.backoff)
}

func (mc Memcache) Increment(key string, delta uint64) (newValue uint64, err error) {

	operation := func() (err error) {
		newValue, err = mc.client.Increment(mc.namespace+key, delta)
		return err
	}

	return newValue, backoff.Retry(operation, mc.backoff)
}

func (mc Memcache) Decrement(key string, delta uint64) (newValue uint64, err error) {

	operation := func() (err error) {
		newValue, err = mc.client.Decrement(mc.namespace+key, delta)
		return err
	}

	return newValue, backoff.Retry(operation, mc.backoff)
}

func setToPointer(in interface{}, out interface{}) error {

	b, err := json.Marshal(in)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, out)
}
