package main

import (
	"fmt"
	"flag"
	"context"
	"os"
	"github.com/naus3a/kaboom/fs"
	"github.com/naus3a/kaboom/cmd"
	"github.com/naus3a/kaboom/sign"
	"github.com/naus3a/kaboom/remote"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

const usage = `Usage:
kaboom-switch -s signing.keys

Options:
	-h, --help	this help screen
	-v, --version	prints version
	-s, --sign	path to the signing key
	-r, --release	if specified IT WILL RELEASE PAYLOAD (default: false)
`

func main() {
	var hFlag bool
	var vFlag bool
	var rFlag bool
	var sFlag string

	cmd.InitCli(usage)
	cmd.AddArg(&hFlag, false, "h", "help")
	cmd.AddArg(&vFlag, false, "v", "version")
	cmd.AddArg(&rFlag, false, "r", "release")
	cmd.AddArg(&sFlag, "", "s", "sign")

	flag.Parse()

	if hFlag {
		flag.Usage()
		os.Exit(0)
	}

	if vFlag {
		fmt.Println(cmd.Version)
		os.Exit(0)
	}

	allGood := !rFlag

	if sFlag==""{
		fmt.Println("You need to specify signing keys\n")
		flag.Usage()
		os.Exit(1)
	}
	signKeysJson, err := fs.LoadFile(sFlag)
	cmd.ReportErrorAndExit(err)
	signKeys, err:= sign.DeserializeSigningKeys(signKeysJson)
	cmd.ReportErrorAndExit(err)
	
	if !allGood{
		cmd.ColorPrintln("YOU ACTIVATED THE RELEASE PROTOCOL. SHIT GONNA HIT THE FAN.", cmd.Red)
	}

	hb, err := sign.NewHeartBeat(allGood, signKeys)
	cmd.ReportErrorAndExit(err)	

	ctx := context.Background()

	comms, err := remote.NewPubSubComms("cippa", ctx)
	cmd.ReportErrorAndExit(err)
	comms.OnPeerConnected = func() {
		comms.Send([]byte(hb.Signature))
		cmd.ColorPrintln("Heartbeat delivered.", cmd.Green)
	}
	comms.OnMessageParsed = func(m *pubsub.Message ){
		myId, err := comms.GetMyId()
		if err == nil{
			if myId == m.ReceivedFrom{
				os.Exit(0)
			}
		}
	}
	cmd.ColorPrintln("Comms ready.", cmd.Green)

	go comms.DiscoverPeers()

	err = comms.Listen()
	cmd.ReportErrorAndExit(err)
	cmd.ColorPrintln("Listening.", cmd.Green)
	comms.ParseMessages()
}
