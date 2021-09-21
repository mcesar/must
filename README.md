# must
Generic error handling with panic, recover, and defer.

Usage:
```go
import "github.com/mcesar/must"
func f() (err error) {
  must.Handle(&err)
  f := must.Do(os.Open("file"))
  defer f.close()
  // ...
}
```
