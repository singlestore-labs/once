package once

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	cases := []struct {
		f func() error
		e string
	}{
		{
			f: func() error {
				panic("hi, I'm a panic")
			},
			e: "I'm a panic",
		},
		{
			f: func() error {
				panic(fmt.Errorf("I'm an error and a panic"))
			},
			e: "I'm an error and a panic",
		},
		{
			f: func() error {
				return nil
			},
			e: "",
		},
		{
			f: func() error {
				return fmt.Errorf("a regular error")
			},
			e: "a regular error",
		},
	}

	for _, tc := range cases {
		c := make(chan error)
		go func() {
			var err error
			e := ReliableError(c, &err)
			defer e.Catch()
			err = tc.f()
			e.Do()
		}()
		timer := time.NewTimer(time.Second)
		select {
		case err := <-c:
			if tc.e == "" {
				assert.NoError(t, err, "no error expected")
			} else if assert.Error(t, err, tc.e) {
				assert.Contains(t, err.Error(), tc.e)
			}
			timer.Stop()
		case <-timer.C:
			t.Errorf("timeout for %s", tc.e)
		}
	}
}
