package main

import(
	"os"
	"fmt"
	"flag"
	"github.com/naus3a/kaboom/cmd"
)

const usage = `Usage:
	kaboom-decrypt [-s a.shab,b.shab]

	Options:
		-h, --help	this help screen
		-v, --version	prints version
		-s, shares	a list of csv share paths
`

func main(){
	var hFlag bool
	var vFlag bool
	var sFlag string
	
	hasAtLeast1Task := false

	cmd.InitCli(usage)
	cmd.AddArg(&hFlag, false, "h", "help")
	cmd.AddArg(&vFlag, false, "v", "version")
	cmd.AddArg(&sFlag, "", "s", "shares")
	flag.Parse()

	if hFlag{
		flag.Usage()
		os.Exit(0)
	}

	if vFlag{
		fmt.Println(cmd.Version)
		os.Exit(0)
	}

	if sFlag!=""{
		pthShares, err := cmd.UnpackCsvArg(&sFlag)
		cmd.ReportErrorAndExit(err)
	}

	if !hasAtLeast1Task {
		flag.Usage()
		os.Exit(1)
	}
}
