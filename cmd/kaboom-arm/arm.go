package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/naus3a/kaboom/payload"
	"github.com/naus3a/kaboom/fs"
	"github.com/naus3a/kaboom/cmd"
)

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
	var pFlag string
	var lFlag string
	var nFlag string
	var sFlag uint
	var tFlag uint
	var vFlag bool
	var hFlag bool

	cmd.InitCli(usage)
	cmd.AddArg(&hFlag, false, "h", "help")
	cmd.AddArg(&vFlag, false, "v", "version")
	cmd.AddArg(&pFlag, "", "p", "payload")
	cmd.AddArg(&lFlag, "", "l", "local")
	cmd.AddArg(&nFlag, "", "n", "notes")
	cmd.AddArg(&sFlag, 3, "s", "shares")
	cmd.AddArg(&tFlag, 2, "t", "threshold")

	flag.Parse()

	if hFlag {
		flag.Usage()
		os.Exit(0)
	}

	if vFlag {
		fmt.Println(cmd.Version)
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
                fmt.Printf("%v", err)
                os.Exit(1)
        }
}
