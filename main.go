package main

import (
	"flag"
	"fmt"
		"os"
)

const usageMessage = "" +
	`Usage of 'go tool cover':
Given a coverage profile produced by 'go test':
	   go test -coverprofile=c.out

Open a terminal ui  displaying annotated source code:
       goconver-cui -cui=c.out
`

func usage() {
	fmt.Fprintln(os.Stderr, usageMessage)
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\n  Only one -cui setting.")
	os.Exit(2)
}

var filePath = flag.String("f", "", "generate terminal ui representation of coverage profile")

var profile string // The profile to read; the value of -f

func main() {
	flag.Usage = usage
	flag.Parse()

	// Usage information when no arguments.
	if flag.NFlag() == 0 && flag.NArg() == 0 {
		flag.Usage()
	}

	err := parseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, `For usage information, run "goconver-cui -help"`)
		os.Exit(2)
	}

	if err := cuiOutput(profile); err != nil {
		panic(err)
	}
}

func parseFlags() error {
	profile = *filePath

	fmt.Println(flag.NArg(), profile)

	if profile == "" {
		return fmt.Errorf("missing source file")
	}

	return nil
}
