// Package once exists because sync.Once is annoying to use because each place you
// use it must provide the function to execute.  That's great if you want to do
// different things, but if you always want to do the same thing, it's repetitive.
//
// Further, ReliableError provides a way to use Once to make go-routines that return
// error easy to write safely.
package once

import (
	"sync"
)

type Once struct {
	o sync.Once
	f func()
}

// New provides the function that Do will will call. The function
// will only be called once no matter how many times Do is invoked.
func New(f func()) *Once {
	return &Once{f: f}
}

// Call the function provided in New
func (o *Once) Do() {
	o.o.Do(o.f)
}
