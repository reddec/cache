package cache

import (
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
