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
