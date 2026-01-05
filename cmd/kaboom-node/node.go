package main

import (
	"os"
	"fmt"
	"flag"
	"time"
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
	-l, --log	path to the log file (default: log.logb)
`

var shares []*sign.ArmoredShare
var log *sign.HeartBeatLog
var lFlag string

func main() {
	var hFlag bool
	var vFlag bool
	var sFlag string

	cmd.InitCli(usage)
	cmd.AddArg(&hFlag, false, "h", "help")
	cmd.AddArg(&vFlag, false, "v", "version")
	cmd.AddArg(&sFlag, "", "s", "shares")
	cmd.AddArg(&lFlag, "log.logb", "l", "log")
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

	err = loadLog(lFlag)
	cmd.ReportErrorAndExit(err)

	checkExpiredHeartBeats()

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

func loadLog(pth string) error{
	data, err := fs.LoadFile(pth)
	if err != nil {
		fmt.Printf("No log file at %s; creating one\n", pth)
		log = sign.NewHeartBeatLog()
		return nil
	}else{
		log, err = sign.DeserializeHeartBeatLog(data)
		return err
	}


}

func saveLog(pth string) error{
	data, err := log.Serialize()
	if err!= nil{
		return err
	}
	return fs.SaveFile(data, pth)
}

func checkExpiredHeartBeats(){
	now := time.Now().Unix()
	for _,s := range shares{
		if log.IsExpired(s, now){
			startReleaseProtocol(s)
		}	
	}
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
	log.LogHeartBeat(s.AuthKey, hb)
	err := saveLog(lFlag)
	if err != nil {
		cmd.ColorPrintln("Cannot save log", cmd.Red)
	}
}

func handleTamperedHeartBeat(s * sign.ArmoredShare){
	cmd.ColorPrintln(fmt.Sprintf("Tampered heartbeat for %s", s.AuthKey), cmd.Red)
}

func startReleaseProtocol(s *sign.ArmoredShare){
	cmd.ColorPrintln(fmt.Sprintf("RELEASING PROTOCOL STARTED FOR %s", s.AuthKey), cmd.Red)
	//TODO
}
