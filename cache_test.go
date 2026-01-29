package saslauthd_test

import (
	"fmt"
	"testing"

	"github.com/matyas-cyril/saslauthd/cache_generic"
)

var PATH string = "/tmp"

// go test -timeout 5s -run ^TestLocal$
func TestLocal(t *testing.T) {
	opt := []any{PATH}
	c, err := cache_generic.New("LOCAL", nil, 3, 3, opt)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string][]byte{
		"hello": []byte("world"),
	}

	hash := "68b329da9893e34099c7d8ad5cb9c940"

	if err := c.SetSucces(data, []byte(hash)); err != nil {
		t.Fatal(err)
	}

}

// go test -timeout 5s -run ^TestMemcache$
func TestMemcache(t *testing.T) {
	opt := []any{"127.0.0.1", 11211, 10}

	c, err := cache_generic.New("MEMCACHE", nil, 3, 3, opt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c)
}

// go test -timeout 5s -run ^TestRedis$
func TestRedis(t *testing.T) {
	opt := []any{"127.0.0.1", 6379, 10}

	c, err := cache_generic.New("REDIS", nil, 3, 3, opt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c)
}
