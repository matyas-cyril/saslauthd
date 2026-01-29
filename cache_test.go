package saslauthd_test

import (
	"fmt"
	"testing"

	cache "github.com/matyas-cyril/saslauthd/cache"
)

var PATH string = "/tmp"

var DATA = map[string][]byte{
	"hello": []byte("world"),
}

var HASH []byte = []byte("68b329da9893e34099c7d8ad5cb9c940")

// go test -timeout 5s -run ^TestLocal$
func TestLocal(t *testing.T) {

	opt := map[string]any{
		"path": PATH,
	}

	c, err := cache.New("LOCAL", nil, 3, 3, opt)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	if err := c.SetSucces(DATA, []byte(HASH)); err != nil {
		t.Fatal(err)
	}

	rst, err := c.GetCache(HASH)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Cache -> ", rst)
}

// go test -timeout 5s -run ^TestMemcache$
func TestMemcache(t *testing.T) {

	opt := map[string]any{
		"host":    "127.0.0.1",
		"port":    11211,
		"timeout": 10,
	}

	c, err := cache.New("MEMCACHE", nil, 3, 3, opt)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.SetSucces(DATA, HASH); err != nil {
		t.Fatal(err)
	}

	rst, err := c.GetCache(HASH)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Cache -> ", rst)

}

// go test -timeout 5s -run ^TestRedis$
func TestRedis(t *testing.T) {

	opt := map[string]any{
		"host":    "127.0.0.1",
		"port":    6379,
		"timeout": 10,
		"db":      0,
	}

	c, err := cache.New("REDIS", nil, 3, 3, opt)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.SetSucces(DATA, HASH); err != nil {
		t.Fatal(err)
	}

	rst, err := c.GetCache(HASH)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Cache -> ", rst)
}
