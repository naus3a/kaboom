package remote

import(
	"fmt"
	"context"
	"sync"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/naus3a/kaboom/cmd"
)

type PubSubComms struct {
	chanName string
	theCtx context.Context
	theHost host.Host
	pubSub *pubsub.PubSub
	topic *pubsub.Topic
	sub *pubsub.Subscription
	theDht *dht.IpfsDHT
	route *drouting.RoutingDiscovery
	hasConnectedPeers bool

	OnPeerConnected func()
}

func NewPubSubComms(channelName string, ctx context.Context)(*PubSubComms, error){
	h, err := makeHost()
	if err!= nil{
		return nil, err
	}
	cmd.ColorPrintln("Host ready.", cmd.Green)
	kDht, err := initDht(ctx, h)
	if err!=nil{
		return nil, err
	}
	cmd.ColorPrintln("DHT ready.", cmd.Green)

	routeDiscovery := drouting.NewRoutingDiscovery(kDht)
	dutil.Advertise(ctx, routeDiscovery, channelName)
	cmd.ColorPrintln("Routing ready.", cmd.Green)

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil{
		return nil, err
	}
	t, err := ps.Join(channelName)
	if err != nil{
		return nil, err
	}
	cmd.ColorPrintln("Channel ready.", cmd.Green)

	return &PubSubComms{
		chanName: channelName,
		theCtx: ctx,
		theHost: h,
		pubSub: ps,
		topic: t,
		sub: nil,
		theDht: kDht,
		route: routeDiscovery,
		hasConnectedPeers: false,
		OnPeerConnected: nil,

	}, nil
} 

func (c *PubSubComms) DiscoverPeers(){
	if c.theCtx==nil || c.theHost==nil || c.theDht==nil || c.route==nil{
		cmd.ReportErrorAndExit(fmt.Errorf("PubSubComms is initialized"))
	}
	if c.hasConnectedPeers{
		cmd.ColorPrintln("Peers already connected.", cmd.Yellow)
		return
	}
	for !c.hasConnectedPeers{
		cmd.ColorPrintln("Peer Discovery started", cmd.Green)
		peerChan, err := c.route.FindPeers(c.theCtx, c.chanName)
		if err !=nil{
			cmd.ReportErrorAndExit(err)
		}
		for peer := range peerChan{
			if peer.ID == c.theHost.ID(){
				continue
			}
			err := c.theHost.Connect(c.theCtx, peer)
			if err==nil{
				txt := fmt.Sprintf("Connected to %s", peer.ID)
				cmd.ColorPrintln(txt, cmd.Green)
				if c.OnPeerConnected != nil{
					c.OnPeerConnected()
				}
				c.hasConnectedPeers=true
			}
		}
	}
	cmd.ColorPrintln("Discovery complete", cmd.Green)
}

func (c *PubSubComms) Send(data []byte) error{
	if !c.hasConnectedPeers {
		return fmt.Errorf("Caannot send data, because not connected to any peer")
	}
	return c.topic.Publish(c.theCtx, data)
}

func makeHost()(host.Host, error){
	return libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
}

func initDht(ctx context.Context, h host.Host)(*dht.IpfsDHT, error){
	kDht, err := dht.New(ctx, h)
	if err!=nil{
		return nil, err
	}
	err = kDht.Bootstrap(ctx)
	if err !=nil{
		return nil, err
	}

	var wg sync.WaitGroup
	for _, peerAddr := range dht.DefaultBootstrapPeers{
		peerInfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func(){
			defer wg.Done()
			if err := h.Connect(ctx, *peerInfo); err != nil {
				fmt.Println("Bootstrap warning:", err)
			}
		}()
	}
	wg.Wait()
	
	return kDht,nil
}


