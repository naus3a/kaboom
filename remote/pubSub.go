package remote

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/naus3a/kaboom/cmd"
	"sync"
)

type PubSubComms struct {
	ChanName          string
	TheCtx            context.Context
	TheHost           host.Host
	pubSub            *pubsub.PubSub
	Topic             *pubsub.Topic
	Sub               *pubsub.Subscription
	TheDht            *dht.IpfsDHT
	route             *drouting.RoutingDiscovery
	hasConnectedPeers bool

	OnPeerConnected func()
}

func NewPubSubComms(channelName string, ctx context.Context) (*PubSubComms, error) {
	c := new(PubSubComms)
	var err error
	c.ChanName = channelName
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
	if c.TheCtx == nil || c.TheHost == nil {
		return fmt.Errorf("Cannot init PubSub: not initialized")
	}
	var err error
	c.pubSub, err = pubsub.NewGossipSub(c.TheCtx, c.TheHost)
	if err != nil {
		return err
	}
	c.Topic, err = c.pubSub.Join(c.ChanName)
	if err != nil {
		return err
	}
	return nil
}

func (c *PubSubComms) Listen() error {
	if c.Topic == nil {
		return fmt.Errorf("Cannot listen: not initialized")
	}
	var err error
	c.Sub, err = c.Topic.Subscribe()
	return err
}

func (c *PubSubComms) InitDHT() error {
	if c.TheCtx == nil || c.TheHost == nil {
		return fmt.Errorf("Cannot init DHT: not initialized")
	}
	var err error
	c.TheDht, err = dht.New(c.TheCtx, c.TheHost)
	if err != nil {
		return err
	}
	err = c.TheDht.Bootstrap(c.TheCtx)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, peerAddr := range dht.DefaultBootstrapPeers {
		peerInfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			err = c.TheHost.Connect(c.TheCtx, *peerInfo)
			if err != nil {

			}
		}()
	}
	wg.Wait()
	return nil
}

func (c *PubSubComms) DiscoverPeers() {
	err := c.InitDHT()
	cmd.ReportErrorAndExit(err)
	routingDiscovery := drouting.NewRoutingDiscovery(c.TheDht)
	dutil.Advertise(c.TheCtx, routingDiscovery, c.ChanName)
	anyConnected := false
	for !anyConnected {
		cmd.ColorPrintln("Searching for peers...", cmd.Yellow)
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
			} else {
				cmd.ColorPrintln(fmt.Sprintf("Connected to %s", peer.ID), cmd.Green)
				anyConnected = true
			}
		}
	}
	cmd.ColorPrintln("Peer discovery complete.", cmd.Green)
}
