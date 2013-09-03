package util

import (
	"sync"
)

type CircularCounter struct {
	mu  sync.Mutex
	x   int
	Max int
}

func NewCircularCounter(max int) *CircularCounter {
	c := new(CircularCounter)
	c.Max = max
	c.x = 0
	return c
}

func (this *CircularCounter) Add(x int) {
	this.mu.Lock()
	this.x += x
	if this.x > (this.Max - 1) {
		this.x = 0
	}
	this.mu.Unlock()
}

func (this *CircularCounter) Value() (x int) {
	this.mu.Lock()
	x = this.x
	this.mu.Unlock()
	return
}

func (this *CircularCounter) Next() int {
	this.mu.Lock()
	defer this.mu.Unlock()
	current := this.x
	this.x += 1
	if this.x > (this.Max - 1) {
		this.x = 0
	}
	return current
}
