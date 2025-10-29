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
	-d, --delete		secure-delete plaintext
`

func main(){
	var pFlag string
	var lFlag string
	var nFlag string
	var sFlag uint
	var tFlag uint
	var vFlag bool
	var hFlag bool
	var dFlag bool

	cmd.InitCli(usage)
	cmd.AddArg(&hFlag, false, "h", "help")
	cmd.AddArg(&vFlag, false, "v", "version")
	cmd.AddArg(&dFlag, false, "d","delete")
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
	cmd.ReportErrorAndExit(err)

	key, err := payload.NewArmoredPayloadKey("TODO", nFlag)
	cmd.ReportErrorAndExit(err)

	encPayload, err := key.Encrypt(plaPayload)
	cmd.ReportErrorAndExit(err)
	
	if lFlag!="" {
		err = fs.SaveFile(encPayload, lFlag)
		cmd.ReportErrorAndExit(err)
	}

	shares, err := key.Split(int(tFlag), int(sFlag))
	cmd.ReportErrorAndExit(err)
	for i:=0; i<len(shares); i++{
		fName := fmt.Sprintf("temp%d.shab", i)
		fs.SaveFile(shares[i], fName)
	}

	if dFlag{
		err = fs.DeleteFile(pFlag)
		cmd.ReportErrorAndExit(err)
	}
}

