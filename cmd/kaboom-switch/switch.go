package main

import(
	"fmt"
	"context"
	"github.com/naus3a/kaboom/cmd"
	"github.com/naus3a/kaboom/remote"
)

func main(){
	ctx := context.Background()
	comms, err := remote.NewPubSubComms("cippa", ctx)
	cmd.ReportErrorAndExit(err)
	
	comms.OnPeerConnected = func(){
		fmt.Print(cmd.ColorString("Sending Message...", cmd.Yellow))
		
		err = comms.Send([]byte("cippa"))
		if err!= nil{
			fmt.Println("\r"+cmd.GetAnsiCode(cmd.ClearLine)+cmd.ColorString("Message sent.", cmd.Green))
		}else{
			fmt.Println("\r"+cmd.GetAnsiCode(cmd.ClearLine)+cmd.ColorString("Message failed.", cmd.Red))
		}
	}
	comms.DiscoverPeers()
}
