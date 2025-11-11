package remote

import(
//	"fmt"
	"context"
//	"sync"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
//	"github.com/libp2p/go-libp2p/core/peer"
//	dht "github.com/libp2p/go-libp2p-kad-dht"
//	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
//	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type PubSubComms struct {
	chanName string
	theCtx context.Context
	theHost host.Host
	pubSub *pubsub.PubSub
	topic *pubsub.Topic
	sub *pubsub.Subscription
}

func NewPubSubComms(channelName string, ctx context.Context)(*PubSubComms, error){
	h, err := makeHost()
	if err!= nil{
		return nil, err
	}

	return &PubSubComms{
		chanName: channelName,
		theCtx: ctx,
		theHost: h,
		pubSub: nil,
		topic: nil,
		sub: nil,
	}, nil
} 

func makeHost()(host.Host, error){
	return libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
}
