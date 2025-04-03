// Package once exists because sync.Once is annoying to use because each place you
// use it must provide the function to execute.  That's great if you want to do
// different things, but if you always want to do the same thing, it's repetative.
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

func New(f func()) *Once {
	return &Once{f: f}
}

func (o *Once) Do() {
	o.o.Do(o.f)
}
