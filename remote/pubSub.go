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
)

type PubSubCommsSender struct {
	theHost host.Host
	pubSub *pubsub.PubSub
	topic *pubsub.Topic
}

func NewPubSubCommsSender(channelName string, ctx context.Context)(*PubSubCommsSender, error){
	fmt.Print("Creating host...")
	h, err := makeHost()
	if err != nil{
		fmt.Print("failed :(\n")
		return nil, err
	}
	fmt.Print("done!\n")
	fmt.Print("Discovering peers...")
	err = discoverPeers(ctx, h, channelName)
	if err != nil{
		fmt.Print("failed :(\n")
		return nil, err
	}
	fmt.Print("done!\n")
	fmt.Print("Setting up pubsub...")
	ps, topic, err := makePubSub(ctx, h, channelName)
	if err != nil{
		fmt.Print("failed :(\n")
		return nil, err
	}
	fmt.Print("done!\n")
	return &PubSubCommsSender{
		theHost: h,
		pubSub: ps,
		topic: topic,
	}, nil
} 

func makeHost()(host.Host, error){
	return libp2p.New(
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/0",
			"/ip6/::/tcp/0",
		),
	)
}

func makePubSub(ctx context.Context, h host.Host, channelName string)(*pubsub.PubSub, *pubsub.Topic, error){
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		return nil, nil, err
	}
	topic, err := ps.Join(channelName)
	if err != nil {
		return nil, nil, err
	}
	return ps, topic, nil	
}

func makeDHT(ctx context.Context, h host.Host) (*dht.IpfsDHT, error){
	kDht, err := dht.New(ctx, h)
	if err!=nil{
		return nil, err
	}
	err = kDht.Bootstrap(ctx)
	if err != nil{
		return nil, err
	}
	var wg sync.WaitGroup
	for _, peerAddr := range dht.DefaultBootstrapPeers{
		peerInfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func(){
			defer wg.Done()
			err := h.Connect(ctx, *peerInfo)
			if err != nil{
				fmt.Println("Bootstrap issue: ", err)
			}
		}()
	}
	wg.Wait()
	return kDht, nil
}

func discoverPeers(ctx context.Context, h host.Host,  channelName string)error{
	kDht, err := makeDHT(ctx, h)
	if err != nil{
		return err
	}
	routeDiscovery := drouting.NewRoutingDiscovery(kDht)
	dutil.Advertise(ctx, routeDiscovery, channelName)
	anyConnected := false
	for !anyConnected{
		fmt.Println("Searching peers...")
		peerChan, err := routeDiscovery.FindPeers(ctx, channelName)
		if err != nil {
			return err
		}
		for peer := range peerChan{
			if peer.ID == h.ID(){
				//dont connect to self
				continue
			}
			err := h.Connect(ctx, peer)
			if err != nil{
				return err
			}
			fmt.Println("Connected to ", peer.ID)
			anyConnected = true
		}
	}
	fmt.Println("Peer discovery complete.")
	return nil
}
