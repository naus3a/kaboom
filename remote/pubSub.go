package remote

import(
	"fmt"
	"context"
	"sync"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

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
