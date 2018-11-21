package memcache

import (
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
)

// Error alias'
var ErrCacheMiss = memcache.ErrCacheMiss

// Type alias'
type Item = memcache.Item

func New(namespace string, dsns ...string) *Memcache {

	mc := new(Memcache)
	mc.client = memcache.New(dsns...)
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

func (mc Memcache) GetSetInt(item memcache.Item, f func() (j int, err error)) (count int, err error) {

	err = mc.Get(item.Key, &count)

	if err != nil && (err == memcache.ErrCacheMiss || err.Error() == "EOF") {

		count, err := f()
		if err != nil {
			return count, err
		}

		err = mc.Set(item.Key, count, item.Expiration)
		return count, err
	}

	return count, err
}

func (mc Memcache) GetSetString(item memcache.Item, f func() (j string, err error)) (s string, err error) {

	err = mc.Get(item.Key, &s)

	if err != nil && (err == memcache.ErrCacheMiss || err.Error() == "EOF") {

		s, err := f()
		if err != nil {
			return s, err
		}

		err = mc.Set(item.Key, s, item.Expiration)
		return s, err
	}

	return s, err
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
