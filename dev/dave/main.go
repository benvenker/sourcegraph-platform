package main

import (
	"fmt"
	"os"

	"github.com/sourcegraph/sourcegraph/lib/output"
)

func main() {
	out := output.NewOutput(os.Stdout, output.OutputOpts{})

	fmt.Println(out)

}
