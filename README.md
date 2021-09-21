# must
Generic error handling with panic, recover, and defer.

For now, it must be used with gotip:
```sh
$ go install golang.org/dl/gotip@latest
$ gotip download
```

Usage:
```go
// main.go
package main

import (
	"fmt"
	"os"

	"github.com/mcesar/must"
)

func main() {
	fmt.Println(f())
}

func f() (err error) {
	defer must.Handle(&err)
	f := must.Do(os.Open("file"))
	defer f.Close()
	// ...
	return nil
}

```
To run:
```sh
$ gotip run main.go
```

## Documentation:

**func Do**
```go
func Do[A any](a A, err error) A
```
Do returns a or panics if err != nil

**func Do0**
```go
func Do0(err error)
```
Do panics if err != nil

**func Do2**
```go
func Do2[A, B any](a A, b B, err error)
```
Do2 returns a and b or panics if err != nil

**func Handle**
```go
func Handle(err *error)
```
Handle sets err to recovered value if it is an error

**func Handlef**
```go
func Handlef(err *error, str string)
```
Handlef sets err to recovered value if it is an error,
wrapped according to the formatting string specified
