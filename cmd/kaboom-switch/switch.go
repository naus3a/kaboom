package main

import (
	"context"
	"fmt"
	"github.com/naus3a/kaboom/cmd"
	"github.com/naus3a/kaboom/remote"
	"os"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var topicNameFlag string

func main() {
	topicNameFlag = "cippa"

	ctx := context.Background()

	comms, err := remote.NewPubSubComms(topicNameFlag, ctx)
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

	printMessagesFrom(comms.TheCtx, comms.Sub)
}

func printMessagesFrom(ctx context.Context, sub *pubsub.Subscription) {
	for {
		m, err := sub.Next(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println(m.ReceivedFrom, ": ", string(m.Message.Data))
		os.Exit(0)
	}
}
