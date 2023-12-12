package memcache

import (
	"errors"
	"reflect"

	"github.com/memcachier/mc/v3"
)

// ErrNoSet is to tell GetSet() not to Set()
var ErrNoSet = errors.New("")

type Config = mc.Config

type Client struct {
	client    *mc.Client
	namespace string
	encoder   Encoder
	decoder   Decoder
	servers   string // comma, semicolon or space seperated
	username  string
	password  string
	config    *Config
}

func NewClient(servers string, options ...Option) *Client {

	client := &Client{
		servers: servers,
		config:  mc.DefaultConfig(),
		encoder: JSONEncoder,
		decoder: JSONDecoder,
	}

	for _, option := range options {
		option(client)
	}

	client.client = mc.NewMCwithConfig(client.servers, client.username, client.password, client.config)

	return client
}

// GetClient gives you access to many other Memcache calls, inc, dec etc
func GetClient(c *Client) *mc.Client {
	return c.client
}

// Exists does not return an error when nothing found
func Exists(c *Client, key string) (exists bool, err error) {

	_, _, _, err = c.client.Get(c.namespace + key)
	if err != nil && !errors.Is(err, mc.ErrNotFound) {
		return false, err
	}
	if errors.Is(err, mc.ErrNotFound) {
		return false, nil
	}
	return true, nil
}

func Get[T any](c *Client, key string, out *T) (err error) {

	val, _, _, err := c.client.Get(c.namespace + key)
	if err != nil {
		return err
	}

	return c.decoder(val, out)
}

func Set(c *Client, key string, value any, seconds uint32) (err error) {

	encoded, err := c.encoder(value)
	if err != nil {
		return err
	}

	_, err = c.client.Set(c.namespace+key, encoded, 0, seconds, 0)
	return err
}

func GetSet[T any](c *Client, key string, seconds uint32, out *T, callback func() (*T, error)) (err error) {

	err = Get(c, key, out)
	if errors.Is(err, mc.ErrNotFound) {

		var s any
		var set = true

		s, err = callback()
		if errors.Is(err, ErrNoSet) {
			set = false
			err = nil
		}
		if err != nil {
			return err
		}

		// If s is nil it panics
		// todo, set out to empty value if s = nil
		if s != nil {
			reflect.ValueOf(out).Elem().Set(reflect.ValueOf(s).Elem())
		}

		if !set {
			return nil
		}

		return Set(c, key, s, seconds)
	}

	return err
}

// Delete does not error on missing keys
func Delete(c *Client, keys ...string) (err error) {

	for _, key := range keys {
		err = c.client.Del(c.namespace + key)
		if err != nil && !errors.Is(err, mc.ErrNotFound) {
			return err
		}
	}

	return nil
}

func Inc(c *Client, key string, delta uint64, seconds uint32) (new uint64, err error) {

	new, _, err = c.client.Incr(c.namespace+key, delta, 1, seconds, 0)
	return new, err
}

func Dec(c *Client, key string, delta uint64, seconds uint32) (new uint64, err error) {

	new, _, err = c.client.Decr(c.namespace+key, delta, 1, seconds, 0)
	return new, err
}

// DeleteAll does not delete keys, but expires them
func DeleteAll(c *Client) error {
	return c.client.Flush(0)
}

func Ping(c *Client) error {
	return c.client.NoOp()
}

func Close(c *Client) {
	c.client.Quit()
}
