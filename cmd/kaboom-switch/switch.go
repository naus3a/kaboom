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
	
	comms.OnPeerConnected = func(){
		cmd.ColorPrintln("Sending msg...", cmd.Green)
		err = comms.Send([]byte("cippa"))
		if err!= nil{
			cmd.ColorPrintln("...done", cmd.Green)
		}
	}
	comms.DiscoverPeers()
}
