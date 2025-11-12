package main

import(
//	"fmt"
	"context"
	"github.com/naus3a/kaboom/cmd"
	"github.com/naus3a/kaboom/remote"
)

func main(){
	ctx := context.Background()
	comms, err := remote.NewPubSubComms("cippa", ctx)
	cmd.ReportErrorAndExit(err)
	comms.DiscoverPeers()
}
