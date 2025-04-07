package once

import (
	"sync"

	"github.com/memsql/errors"
)

type Error struct {
	c chan error
	e *error
	o sync.Once
}

// ReliableError provides a way to guarantee that a go routine returns an error
// over a channel, even if there is a panic.  This is important because a common
// idiom is to wait for a go-routine to finish by reading a error from a channel.
//
// Without something like ReliableError, that is unsafe.
//
// The other idiom, to use a sync.WaitGroup and then defer wg.Done() doesn't
// provide a reliable way to return an error.  If you're setting an error value,
// that error may not be set if the go-routine panics.
//
// Example:
//
//	  echan := make(chan error)
//
//	  go func() {
//		var err error
//		re := ReliableError(echan, &err)
//		defer re.Catch()
//
//		... do stuff
//		err = something()
//		if err != nil {
//		   re.Do()
//		   return
//		}
//		... repeat as needed
//	  }()
//
//	  err :- <-echan
func ReliableError(c chan error, e *error) *Error {
	return &Error{
		c: c,
		e: e,
	}
}

// Catch should be deferred right after calling ReliableError
func (e *Error) Catch() {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			e.o.Do(func() { e.c <- err })
		} else {
			e.o.Do(func() { e.c <- errors.Errorf("panic! %s", r) })
		}
	}
	e.o.Do(func() { e.c <- *e.e })
}

// Do is optional but is useful to when returning from the function
// that is generating an error.
func (e *Error) Do() {
	e.o.Do(func() { e.c <- *e.e })
}
