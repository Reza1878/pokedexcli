package internal

import (
	"testing"
	"time"
)

func TestAddCache(t *testing.T) {
	c := NewCache(time.Second * 5)
	c.Add("test", []byte{})

	_, ok := c.Get("test")
	if !ok {
		t.Error("Failed to add value to cache")
	}
}

func TestGetCache(t *testing.T) {
	c := NewCache(time.Second * 5)
	c.Add("test", []byte{})

	v, _ := c.Get("test")

	if v == nil {
		t.Error("Failed to get value from cache")
	}
}

func TestReapLoopCache(t *testing.T) {
	c := NewCache(time.Second * 1)
	c.Add("test", []byte{})

	time.Sleep(1 * time.Second)

	_, ok := c.Get("test")

	if ok {
		t.Error("Failed to clear data from cache")
	}
}
