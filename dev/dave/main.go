package main

import (
	"fmt"
	"os"

	"github.com/sourcegraph/sourcegraph/lib/output"
)

func main() {
	fmt.Print(output.NewOutput(os.Stdout, output.OutputOpts{}))
}
