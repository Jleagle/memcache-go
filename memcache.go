package memcache

import (
	"errors"
	"reflect"

	"github.com/memcachier/mc/v3"
)

var (
	ErrInvalidType = errors.New("value must be a pointer")
	ErrNoSet       = errors.New("") // Use this to tell GetSet() not to Set()
)

type Config = mc.Config

type Client struct {
	client     *mc.Client
	namespace  string
	encoder    Encoder
	decoder    Decoder
	servers    string // comma, semicolon or space seperated
	username   string
	password   string
	typeChecks bool
	config     *Config
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

// Client gives you access to many other memcache calls, inc, dec etc
func (c Client) Client() *mc.Client {
	return c.client
}

// Exists does not return an error when nothing found
func (c Client) Exists(key string) (exists bool, err error) {

	_, _, _, err = c.client.Get(c.namespace + key)
	if err != nil && err != mc.ErrNotFound {
		return false, err
	}

	return err != mc.ErrNotFound, nil
}

func (c Client) Get(key string, out any) (err error) {

	if c.typeChecks && reflect.TypeOf(out).Kind() != reflect.Ptr {
		return ErrInvalidType
	}

	val, _, _, err := c.client.Get(c.namespace + key)
	if err != nil {
		return err
	}

	return c.decoder(val, out)
}

func (c Client) Set(key string, value any, seconds uint32) (err error) {

	encoded, err := c.encoder(value)
	if err != nil {
		return err
	}

	_, err = c.client.Set(c.namespace+key, encoded, 0, seconds, 0)
	return err
}

func (c Client) GetSet(key string, seconds uint32, out any, callback func() (any, error)) (err error) {

	if c.typeChecks && reflect.TypeOf(out).Kind() != reflect.Ptr {
		return ErrInvalidType
	}

	err = c.Get(key, out)
	if err == mc.ErrNotFound {

		var s any
		var set = true

		s, err = callback()
		if err == ErrNoSet {
			set = false
			err = nil
		}
		if err != nil {
			return err
		}

		if c.typeChecks && reflect.TypeOf(s) != reflect.TypeOf(out).Elem() {
			return errors.New(reflect.TypeOf(s).String() + " does not match " + reflect.TypeOf(out).Elem().String())
		}

		// If s is nil it panics
		// todo, set out to empty value if s = nil
		if s != nil {
			reflect.ValueOf(out).Elem().Set(reflect.ValueOf(s))
		}

		if !set {
			return nil
		}

		return c.Set(key, s, seconds)
	}

	return err
}

// Delete does not error on missing keys
func (c Client) Delete(keys ...string) (err error) {

	for _, key := range keys {
		err = c.client.Del(c.namespace + key)
		if err != nil && err != mc.ErrNotFound {
			return err
		}
	}

	return nil
}

// DeleteAll does not delete keys, but expires them
func (c Client) DeleteAll() error {

	return c.client.Flush(0)
}

func (c Client) Ping() error {
	return c.client.NoOp()
}

func (c Client) Close() {
	c.client.Quit()
}
