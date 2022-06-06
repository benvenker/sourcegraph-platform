package main

import (
	"fmt"
	"os"

	"github.com/moby/term"

	"github.com/sourcegraph/sourcegraph/lib/errors"
)

func main() {

	size, err := term.GetWinsize(os.Stdout.Fd())
	if err != nil {
		err = errors.Wrap(err, "GetWinsize")
	} else {
		fmt.Printf("%d %d\n", size.Height, size.Width)
	}
}
