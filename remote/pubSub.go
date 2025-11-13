package remote

import (
	"context"
	// "fmt"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	// "github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	// dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	// "github.com/naus3a/kaboom/cmd"
	// "sync"
)

type PubSubComms struct {
	chanName          string
	TheCtx            context.Context
	TheHost           host.Host
	pubSub            *pubsub.PubSub
	Topic             *pubsub.Topic
	Sub               *pubsub.Subscription
	theDht            *dht.IpfsDHT
	route             *drouting.RoutingDiscovery
	hasConnectedPeers bool

	OnPeerConnected func()
}

func NewPubSubComms(channelName string, ctx context.Context) (*PubSubComms, error) {
	c := new(PubSubComms)
	var err error
	c.TheCtx = ctx
	c.TheHost, err = libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))

	if err != nil {
		return nil, err
	}

	err = c.InitPubSub()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *PubSubComms) InitPubSub() error {
	var err error
	c.pubSub, err = pubsub.NewGossipSub(c.TheCtx, c.TheHost)
	if err != nil {
		return err
	}
	c.Topic, err = c.pubSub.Join(c.chanName)
	if err != nil {
		return err
	}
	return nil
}

func (c *PubSubComms) Listen() error {
	var err error
	c.Sub, err = c.Topic.Subscribe()
	return err
}

/*
func NewPubSubComms(channelName string, ctx context.Context) (*PubSubComms, error) {
	comms := new(PubSubComms)
	var err error
	comms.chanName = channelName
	comms.theCtx = ctx
	comms.theHost, err = libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		return nil, err
	}
	go discoverPeers(comms.theCtx, comms.theHost, comms.chanName)
	comms.pubSub, err = pubsub.NewGossipSub(comms.theCtx, comms.theHost)
	if err != nil {
		return nil, err
	}
	comms.topic, err = comms.pubSub.Join(comms.chanName)
	if err != nil {
		return nil, err
	}

	return comms, nil
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
				fmt.Println("Bootstrap warning:", err)
			}
		}()
	}
	wg.Wait()

	return kademliaDHT
}

func discoverPeers(ctx context.Context, h host.Host, chanName string) {
	kademliaDHT := initDHT(ctx, h)
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, chanName)

	// Look for others who have announced and attempt to connect to them
	anyConnected := false
	for !anyConnected {
		fmt.Println("Searching for peers...")
		peerChan, err := routingDiscovery.FindPeers(ctx, chanName)
		if err != nil {
			panic(err)
		}
		for peer := range peerChan {
			if peer.ID == h.ID() {
				continue // No self connection
			}
			err := h.Connect(ctx, peer)
			if err != nil {
				fmt.Printf("Failed connecting to %s, error: %s\n", peer.ID, err)
			} else {
				fmt.Println("Connected to:", peer.ID)
				anyConnected = true
			}
		}
	}
	fmt.Println("Peer discovery complete")
}
*/
