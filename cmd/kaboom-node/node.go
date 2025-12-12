package main

import (
	"context"
	"fmt"
	"github.com/naus3a/kaboom/cmd"
	"github.com/naus3a/kaboom/remote"

	"bufio"
	"os"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var topicNameFlag string

func main() {
	topicNameFlag = "cippa"

	ctx := context.Background()

	comms, err := remote.NewPubSubComms(topicNameFlag, ctx)
	cmd.ReportErrorAndExit(err)
	cmd.ColorPrintln("Comms ready.", cmd.Green)

	go comms.DiscoverPeers()

	go streamConsoleTo(comms.TheCtx, comms.Topic)

	err = comms.Listen()
	cmd.ReportErrorAndExit(err)
	cmd.ColorPrintln("Listening.", cmd.Green)

	printMessagesFrom(comms.TheCtx, comms.Sub)
}

func streamConsoleTo(ctx context.Context, topic *pubsub.Topic) {
	reader := bufio.NewReader(os.Stdin)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if err := topic.Publish(ctx, []byte(s)); err != nil {
			fmt.Println("### Publish error:", err)
		}
	}
}

func printMessagesFrom(ctx context.Context, sub *pubsub.Subscription) {
	for {
		m, err := sub.Next(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println(m.ReceivedFrom, ": ", string(m.Message.Data))
	}
}
