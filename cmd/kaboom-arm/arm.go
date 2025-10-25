package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"github.com/naus3a/kaboom/payload"
	"github.com/naus3a/kaboom/fs"
)

const version = "0.0.1"

const usage = `Usage:
	kaboom-arm -p plaintext.file [-l localEncrypted.file] [-s 3] [-t 2]

Options:
	-h, --help				this help screen
	-v, --version			prints version
	-p, --payload		the file you want to encrypt
	-l, --local		the local version of the encrypted payload output
	-s, --shares		number of shares (default: 3)
	-t, --threshold		share threshold (default: 2)
	-n, --notes		extra notes for your payload
`

func main(){
	log.SetFlags(0)
	flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }

	var pFlag string
	var lFlag string
	var nFlag string
	var sFlag uint
	var tFlag uint
	var vFlag bool
	var hFlag bool

	flag.BoolVar(&hFlag, "h", false, "this help screen")
	flag.BoolVar(&hFlag, "help", false, "this help screen")
	flag.BoolVar(&vFlag, "v", false, "prints version")
	flag.BoolVar(&vFlag, "version", false, "prints version")
	flag.StringVar(&pFlag, "p", "", "the payload file to encrypt")
	flag.StringVar(&lFlag, "local", "", "the local encryoted output")
	flag.StringVar(&lFlag, "l", "", "the local encryoted output")
	flag.StringVar(&pFlag, "payload", "", "the payload file to encrypt")
	flag.StringVar(&nFlag, "notes", "", "estra notes")
	flag.StringVar(&nFlag, "n", "", "estra notes")
	flag.UintVar(&sFlag, "s", 3, "shamir shared secrer shares")
	flag.UintVar(&sFlag, "shares", 3, "shamir shared secrer shares")
	flag.UintVar(&tFlag, "threshold", 2, "shamir shared secret threshold")
	flag.UintVar(&tFlag, "t", 2, "shamir shared secret threshold")

	flag.Parse()

	if hFlag {
		flag.Usage()
		os.Exit(0)
	}

	if vFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	var hasPayloadOutput = false

	if pFlag==""{
		fmt.Println("You need to specify a payload file")
		flag.Usage()
		os.Exit(1)
	}

	if lFlag!=""{
		hasPayloadOutput = true
	}

	if !hasPayloadOutput {
		fmt.Println("You need at least 1 payload output")
		flag.Usage()
		os.Exit(1)
	}

	plaPayload, err := fs.LoadFile(pFlag)
	reportErrorAndExit(err)

	key, err := payload.NewArmoredPayloadKey("TODO", nFlag)
	reportErrorAndExit(err)

	encPayload, err := key.Encrypt(plaPayload)
	reportErrorAndExit(err)
	
	if lFlag!="" {
		err = fs.SaveFile(encPayload, lFlag)
		reportErrorAndExit(err)
	}

	shares, err := key.Split(int(tFlag), int(sFlag))
	reportErrorAndExit(err)
	for i:=0; i<len(shares); i++{
		fmt.Printf("%s\n\n", string(shares[i]))
	}
}

func reportErrorAndExit(err error){
	if err != nil {
                fmt.Printf("%w", err)
                os.Exit(1)
        }
}
