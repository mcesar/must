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

import "fmt"

// Do returns a or panics if err != nil
func Do[A any](a A, err error) A {
    if err != nil {
        panic(err)
    }
    return a
}

// Do0 panics if err != nil
func Do0(err error) {
	if err != nil {
		panic(err)
	}
}

// Do2 returns a and b or panics if err != nil
func Do2[A, B any](a A, b B, err error) (A, B) {
	if err != nil {
		panic(err)
	}
	return a, b
}


// Handle sets err to recovered value if it is an error
func Handle(err *error) {
	v := recover()
	if v == nil {
		return
	}
	if err == nil {
		panic(v)
	}
	if e, ok := v.(error); ok && e != nil {
		*err = e
		return
	}
	panic(v)
}

// Handlef sets err to recovered value if it is an error,
// wrapped according to the formatting string specified
func Handlef(err *error, str string) {
    Handle(err)
    if err != nil && *err != nil {
        *err = fmt.Errorf(str, err)
    }
}
