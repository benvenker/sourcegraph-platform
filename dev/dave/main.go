package main

import (
	"fmt"
	"os"

	"github.com/moby/term"

	"github.com/sourcegraph/sourcegraph/lib/errors"
)

func main() {

	fd := os.Stdin.Fd()
	ws, err := term.GetWinsize(fd)
	if err != nil {
		err = errors.Wrap(err, "term.GetWinsize:")
	} else {
		fmt.Println(ws)
	}
}
