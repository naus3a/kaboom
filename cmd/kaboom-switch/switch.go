package main

import (
	"context"

	"github.com/naus3a/kaboom/cmd"
	"github.com/naus3a/kaboom/remote"
)

func main() {
	ctx := context.Background()

	comms, err := remote.NewPubSubComms("cippa", ctx)
	cmd.ReportErrorAndExit(err)
	comms.OnPeerConnected = func() {
		comms.Send([]byte("cippa"))
		cmd.ColorPrintln("Heartbeat delivered.", cmd.Green)
	}
	cmd.ColorPrintln("Comms ready.", cmd.Green)

	go comms.DiscoverPeers()

	err = comms.Listen()
	cmd.ReportErrorAndExit(err)
	cmd.ColorPrintln("Listening.", cmd.Green)
	comms.ParseMessages()
}
