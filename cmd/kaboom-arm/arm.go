package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const version = "0.0.1"

const usage = `Usage:
	kaboom-arm -p PAYLOAD

Options:
	-h, --help				this help screen
	-v, --version			prints version
	-p, --payload PAYLOAD	the file you want to encrypt
`

func main() {
	log.SetFlags(0)
	flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }

	var pFlag string
	var vFlag bool
	var hFlag bool

	flag.BoolVar(&hFlag, "h", false, "this help screen")
	flag.BoolVar(&hFlag, "help", false, "this help screen")
	flag.BoolVar(&vFlag, "v", false, "prints version")
	flag.BoolVar(&vFlag, "version", false, "prints version")
	flag.StringVar(&pFlag, "p", "", "the payload file to encrypt")
	flag.StringVar(&pFlag, "payload", "", "the payload file to encrypt")

	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(0)
	}

	if hFlag {
		flag.Usage()
		os.Exit(0)
	}

	if vFlag {
		fmt.Println(version)
		os.Exit(0)
	}
}
