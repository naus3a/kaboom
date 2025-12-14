package main

import (
	"fmt"
	"flag"
	"context"
	"os"
	"github.com/naus3a/kaboom/cmd"
	"github.com/naus3a/kaboom/remote"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

const usage = `Usage:
kaboom-switch -s signing.keys

Options:
	-h, --help	this help screen
	-v, --version	prints version
	-s, --sign	path to the signing key
`

func main() {
	var hFlag bool
	var vFlag bool
	var sFlag string

	cmd.InitCli(usage)
	cmd.AddArg(&hFlag, false, "h", "help")
	cmd.AddArg(&vFlag, false, "v", "version")
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

	

	ctx := context.Background()

	comms, err := remote.NewPubSubComms("cippa", ctx)
	cmd.ReportErrorAndExit(err)
	comms.OnPeerConnected = func() {
		comms.Send([]byte("cippa"))
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
