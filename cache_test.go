package cache

import (
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	c := New(5, time.Second * 4)

	c.Add("Monday", 1)

	val, ok := c.Get("Monday")
	if !ok {
		t.Fatalf("c.Add() failed: can't take value from cache")
	}

	if val.(int) != 1 {
		t.Fatalf("c.Add() unexpected behaviour: want 1, have %d", val.(int))
	}

	time.Sleep(8 * time.Second)

	_, ok = c.Get("Monday")
	if ok {
		t.Fatalf("c.Add() unexpected behaviour: expected no value, but it is")
	}
}

func TestAddWithTTL(t *testing.T) {
	c := New(5, time.Second * 30)

	c.AddWithTTL("Monday", 1, time.Second * 4)

	val, ok := c.Get("Monday")
	if !ok {
		t.Fatalf("c.AddWithTTL() failed: can't take value from cache")
	}

	if val.(int) != 1 {
		t.Fatalf("c.AddWithTTL() unexpected behaviour: want 1, have %d", val.(int))
	}

	time.Sleep(8 * time.Second)

	_, ok = c.Get("Monday")
	if ok != false {
		t.Fatalf("c.AddWithTTL() unexpected behaviour: expected no value, but it is")
	}
}

func TestGet(t *testing.T) {
	c := New(5, time.Second * 30)

	c.Add("Monday", 1)

	val, ok := c.Get("Monday")
	if !ok {
		t.Fatalf("c.Get() failed: can't take value from cache")
	} else if val != 1 {
		t.Fatalf("c.Get() unexpected result: want 1, have %d", val)
	}
}

func TestRemove(t *testing.T) {
	c := New(5, time.Second * 30)
	
	c.Add("Monday", 1)

	ok := c.Remove("Monday")
	if !ok {
		t.Fatalf("c.Remove() unexpected behaviour: want \"ok\" to be true, but it's %v", ok)
	}

	_, ok = c.Get("Monday")
	if ok {
		t.Fatalf("c.Remove() unexpected behaviour: haven't deleted the \"Monday\" value")
	}
}

func TestClear(t *testing.T) {
	c := New(5, time.Second * 30)

	c.Add("Monday", 1)
	c.Add("Tuesday", 2)
	c.Add("Wednesday", 3)

	c.Clear()

	_, ok := c.Get("Monday")
	if ok {
		t.Fatalf("c.Clear() failed")
	}

	_, ok = c.Get("Tuesday")
	if ok {
		t.Fatalf("c.Clear() failed")
	}

	_, ok = c.Get("Wednesday")
	if ok {
		t.Fatalf("c.Clear() failed")
	}
}

func TestCap(t *testing.T) {
	c := New(5, time.Second * 10)

	cap := c.Cap()
	if cap != 5 {
		t.Fatalf("Cap() failed: expected 5, have %d", cap)
	}
}

func TestLRU(t *testing.T) {
	c := New(4, 15 * time.Second)

	c.Add("1", 1)
	c.Add("2", 2)
	c.Add("3", 3)
	c.Add("4", 4)

	c.Get("1")
	c.Get("2")
	c.Get("3")

	c.Add("5", 5)

	_, ok := c.Get("4")
	if ok {
		t.Fatalf("LRU failed: the least used value is not removed when a new one is added to a busy queue")
	}
}