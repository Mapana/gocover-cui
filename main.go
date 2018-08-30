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

Open a terminal ui  displaying annotated test log:
       goconver-cui -log=test.log

Open the terminal ui that displays the annotated source code and test log:
       goconver-cui -cui=c.out -log=test.log

`

func usage() {
	fmt.Fprintln(os.Stderr, usageMessage)
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\n  -cui and -log at least one.")
	os.Exit(2)
}

var (
	cui = flag.String("cui", "", "generate terminal ui representation of coverage profile")
	log = flag.String("log", "", "generate terminal ui representation of testting logfile")
)
var (
	profile string // The profile to read; the value of -cui
	logfile string // The logfile to read; the value of -log
)

func main() {
	flag.Usage = usage
	flag.Parse()

	// Usage information when no arguments.
	if flag.NFlag() == 0 && flag.NArg() == 0 {
		flag.Usage()
	}

	err := parseFlags()
	if err != nil {
		errExit(err)
	}

	if err := cuiOutput(profile, logfile); err != nil {
		errExit(err)
	}
}

func parseFlags() error {
	profile, logfile = *cui, *log

	if len(profile) <= 0 && len(logfile) <= 0 {
		return fmt.Errorf("missing source file")
	}

	return nil
}

func errExit(err error) {
	fmt.Fprintln(os.Stderr, err)
	fmt.Fprintln(os.Stderr, `For usage information, run "goconver-cui -help"`)
	os.Exit(2)
}
