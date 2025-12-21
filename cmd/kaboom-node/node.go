package main

import (
	"os"
	"fmt"
	"flag"
	"context"
	"github.com/naus3a/kaboom/fs"
	"github.com/naus3a/kaboom/cmd"
	"github.com/naus3a/kaboom/sign"
	"github.com/naus3a/kaboom/remote"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

const usage = `Usage:
kaboom-node -s a.shab,b.shab

Options:
	-h, --help	this help screen
	-v, --version	prints the version
	-s, --shares	a list  of csv share paths
`

var shares []*sign.ArmoredShare

func main() {
	var hFlag bool
	var vFlag bool
	var sFlag string

	cmd.InitCli(usage)
	cmd.AddArg(&hFlag, false, "h", "help")
	cmd.AddArg(&vFlag, false, "v", "version")
	cmd.AddArg(&sFlag, "", "s", "shares")
	flag.Parse()

	if hFlag{
		flag.Usage()
		os.Exit(0)
	}

	if vFlag{
		fmt.Println(cmd.Version)
		os.Exit(0)
	}

	if sFlag==""{
		fmt.Println("You need to specify at least a share file")
		flag.Usage()
		os.Exit(1)
	}

	err := loadShares(sFlag)
	cmd.ReportErrorAndExit(err)

	ctx := context.Background()

	comms, err := remote.NewPubSubComms("cippa", ctx)
	cmd.ReportErrorAndExit(err)
	cmd.ColorPrintln("Comms ready", cmd.Green)
	
	comms.OnMessageParsed = handleMessage

	go comms.DiscoverPeers()

	err = comms.Listen()
	cmd.ReportErrorAndExit(err)
	cmd.ColorPrintln("Listening.", cmd.Green)
	
	comms.ParseMessages()
}

func loadShares(csv string) error{
	pthShares, err := cmd.UnpackCsvArg(&csv)
	if err!=nil{
		return err
	}
	if len(pthShares)<1{
		return fmt.Errorf("you need to specify at least a share file")
	}
	shares = make([]*sign.ArmoredShare, len(pthShares))
	for i:=0; i<len(pthShares);i++{
		jsonData, err := fs.LoadFile(pthShares[i])
		if err!=nil{
			return err
		}
		shares[i], err = sign.DeserializeShare(jsonData)
		if err!=nil{
			return err
		}
	}
	return nil
}

func handleMessage(m *pubsub.Message){
	hb, err := sign.DecodeBinaryHeartBeat(m.Message.Data)
	if err != nil{
		return
	}
	for i:=0; i<len(shares); i++{
		good, err := sign.VerifyHeartBeat(shares[i], &hb)
		if err == nil {
			if good {
				logHeartBeat(shares[i], &hb)
			}else{
				handleTamperedHeartBeat(shares[i])
			}
			return
		}
	} 
}

func logHeartBeat(s * sign.ArmoredShare, hb *sign.HeartBeat){
	if !hb.AllGood{
		startReleaseProtocol(s)
		return
	}
	cmd.ColorPrintln(fmt.Sprintf("Good heartbeat from %s", s.AuthKey), cmd.Green)
	//TODO
}

func handleTamperedHeartBeat(s * sign.ArmoredShare){
	cmd.ColorPrintln(fmt.Sprintf("Tampered heartbeat for %s", s.AuthKey), cmd.Red)
}

func startReleaseProtocol(s *sign.ArmoredShare){
	cmd.ColorPrintln(fmt.Sprintf("RELEASING PROTOCOL STARTED FOR %s", s.AuthKey), cmd.Red)
	//TODO
}
