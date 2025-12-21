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
	chanName          string
	theCtx            context.Context
	theHost           host.Host
	pubSub            *pubsub.PubSub
	topic             *pubsub.Topic
	sub               *pubsub.Subscription
	theDht            *dht.IpfsDHT
	route             *drouting.RoutingDiscovery
	hasConnectedPeers bool

	OnPeerConnected func()
	OnMessageParsed func(*pubsub.Message)
}

func NewPubSubComms(channelName string, ctx context.Context) (*PubSubComms, error) {
	c := new(PubSubComms)
	var err error
	c.chanName = channelName
	c.theCtx = ctx
	c.theHost, err = libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))

	if err != nil {
		return nil, err
	}

	err = c.initPubSub()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *PubSubComms) GetMyId()(peer.ID	, error){
	if c.theHost==nil {
		return peer.ID(""), fmt.Errorf("host not ready")
	}
	return c.theHost.ID(), nil
}

func (c *PubSubComms) initPubSub() error {
	if c.theCtx == nil || c.theHost == nil {
		return fmt.Errorf("Cannot init PubSub: not initialized")
	}
	var err error
	c.pubSub, err = pubsub.NewGossipSub(c.theCtx, c.theHost)
	if err != nil {
		return err
	}
	c.topic, err = c.pubSub.Join(c.chanName)
	if err != nil {
		return err
	}
	return nil
}

func (c *PubSubComms) Listen() error {
	if c.topic == nil {
		return fmt.Errorf("Cannot listen: not initialized")
	}
	var err error
	c.sub, err = c.topic.Subscribe()
	return err
}

func (c* PubSubComms) ParseMessages(){
	for{
		m, err := c.sub.Next(c.theCtx)
		if err!= nil{
			continue
		}
		cmd.ColorPrintln(fmt.Sprintf("%s: %s", m.ReceivedFrom, string(m.Message.Data)), cmd.Green)
		if c.OnMessageParsed != nil{
			c.OnMessageParsed(m)
		}
	}
}

func (c *PubSubComms) initDHT() error {
	if c.theCtx == nil || c.theHost == nil {
		return fmt.Errorf("Cannot init DHT: not initialized")
	}
	var err error
	c.theDht, err = dht.New(c.theCtx, c.theHost)
	if err != nil {
		return err
	}
	err = c.theDht.Bootstrap(c.theCtx)
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
			err = c.theHost.Connect(c.theCtx, *peerInfo)
			if err != nil {

			}
		}()
	}
	wg.Wait()
	return nil
}

func (c *PubSubComms) Send(data []byte) error {
	return c.topic.Publish(c.theCtx, data)
}

func (c *PubSubComms) DiscoverPeers() {
	err := c.initDHT()
	cmd.ReportErrorAndExit(err)
	routingDiscovery := drouting.NewRoutingDiscovery(c.theDht)
	dutil.Advertise(c.theCtx, routingDiscovery, c.chanName)
	anyConnected := false
	for !anyConnected {
		cmd.ColorPrintln("Searching for peers...", cmd.Yellow)
		peerChan, err := routingDiscovery.FindPeers(c.theCtx, c.chanName)
		if err != nil {
			panic(err)
		}
		for peer := range peerChan {
			if peer.ID == c.theHost.ID() {
				continue // No self connection
			}
			err := c.theHost.Connect(c.theCtx, peer)
			if err != nil {
			} else {
				cmd.ColorPrintln(fmt.Sprintf("Connected to %s", peer.ID), cmd.Green)
				anyConnected = true
			}
		}
	}
	cmd.ColorPrintln("Peer discovery complete.", cmd.Green)
	if c.OnPeerConnected != nil {
		c.OnPeerConnected()
	}
}
