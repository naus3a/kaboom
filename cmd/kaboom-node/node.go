package main

import (
	"context"
	"fmt"
	"github.com/naus3a/kaboom/cmd"
	"github.com/naus3a/kaboom/remote"

	"bufio"
	"os"
	"sync"

	//"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

var topicNameFlag string

func main() {
	topicNameFlag = "cippa"

	ctx := context.Background()

	comms, err := remote.NewPubSubComms(topicNameFlag, ctx)
	cmd.ReportErrorAndExit(err)
	cmd.ColorPrintln("Comms ready.", cmd.Green)

	go discoverPeers(comms)

	go streamConsoleTo(comms.TheCtx, comms.Topic)

	err = comms.Listen()
	cmd.ReportErrorAndExit(err)
	cmd.ColorPrintln("Listening.", cmd.Green)

	printMessagesFrom(comms.TheCtx, comms.Sub)
}

func initDHT(ctx context.Context, h host.Host) *dht.IpfsDHT {
	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	kademliaDHT, err := dht.New(ctx, h)
	if err != nil {
		panic(err)
	}
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, peerAddr := range dht.DefaultBootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h.Connect(ctx, *peerinfo); err != nil {
				//fmt.Println("Bootstrap warning:", err)
			}
		}()
	}
	wg.Wait()

	return kademliaDHT
}

func discoverPeers(c *remote.PubSubComms) {
	//kademliaDHT := initDHT(ctx, h)
	err := c.InitDHT()
	cmd.ReportErrorAndExit(err)

	routingDiscovery := drouting.NewRoutingDiscovery(c.TheDht)
	dutil.Advertise(c.TheCtx, routingDiscovery, c.ChanName)

	// Look for others who have announced and attempt to connect to them
	anyConnected := false
	for !anyConnected {
		fmt.Println("Searching for peers...")
		peerChan, err := routingDiscovery.FindPeers(c.TheCtx, c.ChanName)
		if err != nil {
			panic(err)
		}
		for peer := range peerChan {
			if peer.ID == c.TheHost.ID() {
				continue // No self connection
			}
			err := c.TheHost.Connect(c.TheCtx, peer)
			if err != nil {
				//fmt.Printf("Failed connecting to %s, error: %s\n", peer.ID, err)
			} else {
				fmt.Println("Connected to:", peer.ID)
				anyConnected = true
			}
		}
	}
	fmt.Println("Peer discovery complete")
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
