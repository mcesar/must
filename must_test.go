package must

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"
)

type testErrorEmitFunc func(a, delay int) (x int, err error)

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func TestDo(t *testing.T) {
	for _, test := range []struct {
		name string
		x    int
		err  bool
	}{
		{
			"no error",
			0,
			false,
		},
		{
			"error",
			-1,
			true,
		},
		{
			"error",
			1,
			true,
		},
	} {
		for _, f := range []testErrorEmitFunc{withMustErrorHandling, withMustErrorHandlingV2, withMustErrorHandlingV3} {
			t.Run(test.name+fmt.Sprintf("(%s)", getFunctionName(f)), func(t *testing.T) {
				x, err := f(test.x, 0)
				if !test.err && x != test.x {
					t.Errorf("expected %v, got %v", test.x, x)
				}
				if (err != nil) != test.err {
					t.Errorf("expected %v, got %v", test.err, err)
				}
				t.Logf("%+v", err)
			})
		}
	}
}

func BenchmarkMustErrorHandlingWithoutDelay(b *testing.B) {
	for n := 0; n < b.N; n++ {
		withMustErrorHandling(n, 0)
	}
}

func BenchmarkRegularErrorHandlingWithoutDelay(b *testing.B) {
	for n := 0; n < b.N; n++ {
		withRegularErrorHandling(n, 0)
	}
}

func BenchmarkMustErrorHandlingWith10msDelay(b *testing.B) {
	for n := 0; n < b.N; n++ {
		withMustErrorHandling(n, 10)
	}
}

func BenchmarkRegularErrorHandlingWith10msDelay(b *testing.B) {
	for n := 0; n < b.N; n++ {
		withRegularErrorHandling(n, 10)
	}
}

func withMustErrorHandling(a, delay int) (x int, err error) {
	defer Handle(&err)
	time.Sleep(time.Duration(delay) * time.Millisecond)
	x = Do(nonNegativeOnly(a))
	x = Do(nonPositiveOnly(a))
	return x, nil
}

func withMustErrorHandlingV2(a, delay int) (x int, rErr error) {
	defer Handlef(&rErr, "error: %w")
	time.Sleep(time.Duration(delay) * time.Millisecond)
	x = Do(nonNegativeOnly(a))
	x = Do(nonPositiveOnly(a))
	return x, nil
}

func withMustErrorHandlingV3(a, delay int) (x int, rErr error) {
	defer HandleFunc(func(err error) {
		if err != nil {
			rErr = fmt.Errorf("formatted: %+v", err)
		}
	})
	time.Sleep(time.Duration(delay) * time.Millisecond)
	x = Do(nonNegativeOnly(a))
	x = Do(nonPositiveOnly(a))
	return x, nil
}

func withRegularErrorHandling(a, delay int) (x int, err error) {
	time.Sleep(time.Duration(delay) * time.Millisecond)
	x, err = nonNegativeOnly(a)
	if err != nil {
		return 0, err
	}
	x, err = nonPositiveOnly(a)
	if err != nil {
		return 0, err
	}
	return x, nil
}

func nonNegativeOnly(x int) (int, error) {
	if x < 0 {
		return 0, errors.New("must be positive")
	}
	return x, nil
}

func nonPositiveOnly(x int) (int, error) {
	if x > 0 {
		return 0, errors.New("must be negative")
	}
	return x, nil
}

func TestDo_RethrowsPanic(t *testing.T) {
	err := func() (rErr error) {
		defer func() {
			e := recover()
			if e == nil {
				t.Errorf("panic must be rethrown")
			}
		}()
		defer Handle(&rErr)
		Do0(func() error {
			panic(fmt.Errorf("some non-wrapped error"))
		}())
		return nil
	}()
	if err != nil {
		t.Errorf("panic was interpreted as wrapped error")
	}
}
