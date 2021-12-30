package main

import (
	"sync/atomic"
)

type Counter struct {
	create int64
	lookup int64
	submit int64
	notify int64
}

func NewCounter() *Counter {
	c := new(Counter)
	c.create = 0
	c.lookup = 0
	c.submit = 0
	c.notify = 0
	return c
}

func (c *Counter) OnCreate() {
	atomic.AddInt64(&c.create, 1)
}

func (c *Counter) OnLookup() {
	atomic.AddInt64(&c.lookup, 1)
}

func (c *Counter) OnSubmit() {
	atomic.AddInt64(&c.submit, 1)
}

func (c *Counter) OnNotify() {
	atomic.AddInt64(&c.notify, 1)
}

func (c *Counter) Reset() {
	atomic.StoreInt64(&c.create, 0)
	atomic.StoreInt64(&c.lookup, 0)
	atomic.StoreInt64(&c.submit, 0)
	atomic.StoreInt64(&c.notify, 0)
}

func (c *Counter) GetCreate() int64 {
	return atomic.LoadInt64(&c.create)
}
func (c *Counter) GetLookup() int64 {
	return atomic.LoadInt64(&c.lookup)
}
func (c *Counter) GetSubmit() int64 {
	return atomic.LoadInt64(&c.submit)
}
func (c *Counter) GetNotify() int64 {
	return atomic.LoadInt64(&c.notify)
}
