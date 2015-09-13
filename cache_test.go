package cache

import (
	"log"
	"testing"
	"time"
)

func TestNormal(t *testing.T) {
	c := NewCache()
	k := c.Put(1234, 1*time.Second)
	v, ok := c.Pull(k)
	if !ok {
		t.Fatal("Must exists")
	}
	if v.(int) != 1234 {
		t.Fatal("Value broken")
	}
}

func TestExpired(t *testing.T) {
	c := NewCache()
	k := c.Put(1234, 1*time.Second)
	time.Sleep(2 * time.Second)
	_, ok := c.Pull(k)
	if ok {
		t.Fatal("Must not exists")
	}
}

func ExampleCache_Put() {
	cache := NewCache()
	key := cache.Put(1234, 1*time.Second)
	log.Println("Generated key:", key)
}

func ExampleCache_Pull() {
	cache := NewCache()
	key := "key1"
	if value, ok := cache.Pull(key); ok {
		log.Println(key, "=", value)
	} else {
		log.Println("Key", key, "is not exists")
	}
}
