package remote

import (
	"context"
	"fmt"
	"sync"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/naus3a/kaboom/cmd"
)

type PubSubComms struct {
	chanName          string
	theCtx            context.Context
	cancel		  context.CancelFunc
	wg		  sync.WaitGroup
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

func NewPubSubComms(channelName string) (*PubSubComms, error) {
	c := new(PubSubComms)
	c.theCtx, c.cancel = context.WithCancel(context.Background())
	var err error
	c.chanName = channelName
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

func (c* PubSubComms) GetChannelName()string{
	return c.chanName
}

func (c* PubSubComms) Stop() error{
	c.cancel()
	c.wg.Wait()

	err := c.theHost.Close()
	if err!=nil{
		return err
	}
	fmt.Println("%s: Comms Closed.", c.chanName)	
	return nil
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
	c.wg.Add(1)
	defer c.wg.Done()
	for{
		if c.theCtx.Err()!=nil{
			fmt.Println("Message parsing stopped")
			return
		}
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
	c.wg.Add(1)
	defer c.wg.Done()

	err := c.initDHT()
	cmd.ReportErrorAndExit(err)
	routingDiscovery := drouting.NewRoutingDiscovery(c.theDht)
	dutil.Advertise(c.theCtx, routingDiscovery, c.chanName)
	anyConnected := false
	for !anyConnected {
		if c.theCtx.Err()!=nil{
			fmt.Println("%s: Discovery stopped.", c.chanName)
			return
		}
		cmd.ColorPrintln(fmt.Sprintf("%s: Searching for peers...", c.chanName), cmd.Yellow)
		peerChan, err := routingDiscovery.FindPeers(c.theCtx, c.chanName)
		if err != nil {
			panic(err)
		}
		for peer := range peerChan {
			if c.theCtx.Err()!=nil{
				fmt.Println("%s: Discovery stopped.", c.chanName)
			return
		}
			if peer.ID == c.theHost.ID() {
				continue // No self connection
			}
			err := c.theHost.Connect(c.theCtx, peer)
			if err != nil {
			} else {
				cmd.ColorPrintln(fmt.Sprintf("%s: Connected to %s", c.chanName, peer.ID), cmd.Green)
				anyConnected = true
			}
		}
	}
	cmd.ColorPrintln(fmt.Sprintf("%s: Peer discovery complete.", c.chanName), cmd.Green)
	if c.OnPeerConnected != nil {
		c.OnPeerConnected()
	}
}
