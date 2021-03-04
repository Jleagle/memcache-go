package memcache

import (
	"testing"

	"github.com/memcachier/mc/v3"
)

func TestOptions(t *testing.T) {

	config := mc.DefaultConfig()
	config.Retries = 100

	client := NewClient(
		"test",
		WithAuth("user", "pass"),
		WithConfig(config),
		WithTypeChecks(true),
		WithNamespace("test_"),
		WithEncoding(StringEncoder, StringDecoder),
	)

	if client.username != "user" {
		t.Errorf("username = %s; want 'user'", client.username)
	}
	if client.password != "pass" {
		t.Errorf("password = %s; want 'pass'", client.password)
	}
	if client.config.Retries != 100 {
		t.Errorf("retries = %d; want 100", client.config.Retries)
	}
	if client.typeChecks != true {
		t.Errorf("typeChecks = %t; want true", client.typeChecks)
	}
	if client.namespace != "test_" {
		t.Errorf("namespace = %s; want 'test_'", client.namespace)
	}
	out, _ := client.encoder("abc")
	if out != "abc" {
		t.Errorf("encoder = %s; want 'abc'", out)
	}
	out2 := ""
	_ = client.decoder("def", &out2)
	if out2 != "def" {
		t.Errorf("decoder = %s; want 'def'", "def")
	}
}
