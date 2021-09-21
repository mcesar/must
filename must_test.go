package must

import (
	"errors"
	"testing"
)

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
		t.Run(test.name, func(t *testing.T) {
			x, err := func() (x int, err error) {
				defer Handle(&err)
				x = Do(nonNegativeOnly(test.x))
				x = Do(nonPositiveOnly(test.x))
				return x, nil
			}()
			if !test.err && x != test.x {
				t.Errorf("expected %v, got %v", test.x, x)
			}
			if (err != nil) != test.err {
				t.Errorf("expected %v, got %v", test.err, err)
			}
		})
	}
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
