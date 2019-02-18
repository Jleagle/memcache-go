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
	ErrCacheMiss = memcache.ErrCacheMiss
	ErrNotStored = memcache.ErrNotStored

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

// Get gets the item for the given key. ErrCacheMiss is returned for a
// memcache cache miss. The key must be at most 250 bytes in length.
func (mc Memcache) Get(key string) (*Item, error) {
	return mc.client.Get(key)
}

// Set writes the given item, unconditionally.
func (mc Memcache) Set(item *Item) error {
	return mc.client.Set(item)
}

// Touch updates the expiry for the given key. The seconds parameter is either
// a Unix timestamp or, if seconds is less than 1 month, the number of seconds
// into the future at which time the item will expire. Zero means the item has
// no expiration time. ErrCacheMiss is returned if the key is not in the cache.
// The key must be at most 250 bytes in length.
func (mc Memcache) Touch(key string, seconds int32) error {

	operation := func() (err error) {
		err = mc.client.Touch(mc.namespace+key, seconds)
		if err == ErrCacheMiss {
			return backoff.Permanent(err)
		}
		return err
	}

	return backoff.Retry(operation, mc.backoff)
}

// Add writes the given item, if no value already exists for its
// key. ErrNotStored is returned if that condition is not met.
func (mc Memcache) Add(item *Item) error {

	operation := func() (err error) {
		err = mc.client.Add(item)
		if err == ErrNotStored {
			return backoff.Permanent(err)
		}
		return err
	}

	return backoff.Retry(operation, mc.backoff)
}

// Replace writes the given item, but only if the server *does*
// already hold data for this key
func (mc Memcache) Replace(item *Item) error {

	operation := func() (err error) {
		return mc.client.Replace(item)
	}

	return backoff.Retry(operation, mc.backoff)
}

// Delete deletes the item with the provided key. The error ErrCacheMiss is
// returned if the item didn't already exist in the cache.
func (mc Memcache) Delete(key string) (err error) {

	operation := func() (err error) {
		err = mc.client.Delete(mc.namespace + key)
		if err == ErrCacheMiss {
			return backoff.Permanent(err)
		}
		return err
	}

	return backoff.Retry(operation, mc.backoff)
}

// DeleteAll deletes all items in the cache.
func (mc Memcache) DeleteAll() (err error) {

	operation := func() (err error) {
		return mc.client.DeleteAll()
	}

	return backoff.Retry(operation, mc.backoff)
}

// Increment atomically increments key by delta. The return value is
// the new value after being incremented or an error. If the value
// didn't exist in memcached the error is ErrCacheMiss. The value in
// memcached must be an decimal number, or an error will be returned.
// On 64-bit overflow, the new value wraps around.
func (mc Memcache) Increment(key string, delta uint64) (newValue uint64, err error) {

	operation := func() (err error) {
		newValue, err = mc.client.Increment(mc.namespace+key, delta)
		if err == ErrCacheMiss {
			return backoff.Permanent(err)
		}
		return err
	}

	return newValue, backoff.Retry(operation, mc.backoff)
}

// Decrement atomically decrements key by delta. The return value is
// the new value after being decremented or an error. If the value
// didn't exist in memcached the error is ErrCacheMiss. The value in
// memcached must be an decimal number, or an error will be returned.
// On underflow, the new value is capped at zero and does not wrap
// around.
func (mc Memcache) Decrement(key string, delta uint64) (newValue uint64, err error) {

	operation := func() (err error) {
		newValue, err = mc.client.Decrement(mc.namespace+key, delta)
		if err == ErrCacheMiss {
			return backoff.Permanent(err)
		}
		return err
	}

	return newValue, backoff.Retry(operation, mc.backoff)
}

// Get gets the item for the given key. ErrCacheMiss is returned for a
// memcache cache miss. The key must be at most 250 bytes in length.
func (mc Memcache) GetInterface(key string, i interface{}) (err error) {

	var item *Item

	operation := func() (err error) {

		item, err = mc.client.Get(mc.namespace + key)
		if err == ErrCacheMiss {
			return backoff.Permanent(err)
		}
		return err
	}

	err = backoff.Retry(operation, mc.backoff)
	if err != nil {
		return err
	}

	return json.Unmarshal(item.Value, i)
}

// Set writes the given item, unconditionally.
func (mc Memcache) SetInterface(key string, value interface{}, expiration int32) error {

	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	item := new(Item)
	item.Key = mc.namespace + key
	item.Value = bytes
	item.Expiration = expiration

	operation := func() (err error) {
		return mc.client.Set(item)
	}

	return backoff.Retry(operation, mc.backoff)
}

func (mc Memcache) GetSetInterface(key string, expiration int32, value interface{}, f func() (interface{}, error)) error {

	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		return ErrNotPointer
	}

	err := mc.GetInterface(key, value)

	if err == memcache.ErrCacheMiss || err == io.EOF {

		s, err := f()
		if err != nil {
			return err
		}

		if reflect.TypeOf(s) != reflect.TypeOf(value).Elem() {
			return ErrInvalidTypes
		}

		err = setToPointer(s, value)

		return mc.SetInterface(key, s, expiration)
	}

	return err
}

func setToPointer(in interface{}, out interface{}) error {

	b, err := json.Marshal(in)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, out)
}
