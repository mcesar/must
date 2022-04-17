// Package must provides generic error handling with panic, recover, and defer.
// Usage:
//  import "github.com/mcesar/must"
//  func f() (err error) {
//      must.Handle(&err)
//      f := must.Do(os.Open("file"))
//      defer f.close()
//      // ...
//  }
package must

import (
	"errors"
	"fmt"
)

type wrappedError struct{ error }

func (w wrappedError) Unwrap() error { return w.error }

// Do returns a or panics if err != nil
func Do[A any](a A, err error) A {
	if err != nil {
		panic(wrappedError{err})
	}
	return a
}

// Do0 panics if err != nil
func Do0(err error) {
	if err != nil {
		panic(wrappedError{err})
	}
}

// Do2 returns a and b or panics if err != nil
func Do2[A, B any](a A, b B, err error) (A, B) {
	if err != nil {
		panic(wrappedError{err})
	}
	return a, b
}

// Handle sets pErr to recovered value if it is an error emitted inside of a Do call
func Handle(pErr *error) {
	e := recover()
	handle(pErr, e)
}

func handle(pErr *error, e interface{}) {
	switch {
	case e == nil:
		return
	default:
		var errTyped wrappedError
		if eError, ok := e.(error); ok && errors.As(eError, &errTyped) {
			if pErr != nil {
				*pErr = errTyped.error
			}
			return
		}
	}
	panic(e)
}

// Handlef sets err to recovered value if it is an error, wrapped according to
// the formatting string specified
func Handlef(pErr *error, str string) {
	e := recover()
	handle(pErr, e)
	if pErr != nil && *pErr != nil {
		*pErr = fmt.Errorf(str, *pErr)
	}
}

// HandleFunc recovers error and passes it to the handler function
func HandleFunc(f func(err error)) {
	var err error
	e := recover()
	handle(&err, e)
	f(err)
}
