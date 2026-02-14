package main

import (
	"os"
	"fmt"
	"flag"
	"time"
	"sync"
	"github.com/naus3a/kaboom/fs"
	"github.com/naus3a/kaboom/cmd"
	"github.com/naus3a/kaboom/sign"
)

const usage = `Usage:
kaboom-node -s a.shab,b.shab

Options:
	-h, --help	this help screen
	-v, --version	prints the version
	-s, --shares	a list  of csv share paths
	-l, --log	path to the log file (default: log.logb)
`

const intervalExpiredHb = 1*time.Hour
const intervalRotChan = 5*time.Minute

var muShares sync.RWMutex
var shares []*sign.ArmoredShare
var log *sign.HeartBeatLog
var comms []*CommsListener
var lFlag string

func main() {
	//
	// arg parsing
	//
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
	
	//
	// file loading
	//
	err := loadShares(sFlag)
	cmd.ReportErrorAndExit(err)

	err = loadLog(lFlag)
	cmd.ReportErrorAndExit(err)

	checkExpiredHeartBeats()

	//
	// heartbeat listening
	//
	createCommsListeners()

	tExpiredHb := time.NewTicker(intervalExpiredHb)
	tRotChan := time.NewTicker(intervalRotChan)

	for{
		select{
			case <-tExpiredHb.C:
				checkExpiredHeartBeats()
			case <-tRotChan.C:
				updateChannelNameIfNeeded()

		}
	}
}

///
/// comms
///

func createCommsListeners(){
	var err error = nil
	comms = make([]*CommsListener, len(shares))
	for i:=0;i<len(shares);i++{
		comms[i], err = NewCommsListener(shares[i], logHeartBeat, handleTamperedHeartBeat) 
		cmd.ReportErrorAndExit(err)
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

///
/// channel name rotation
//

func updateChannelNameIfNeeded(){
	var err error = nil
	for i:=0;i<len(comms);i++{
		err = comms[i].UpdateChannelNameIfNeeded()
		cmd.ReportErrorAndExit(err)
	}
}

///
/// fs
///

func loadShares(csv string) error{
	pthShares, err := cmd.UnpackCsvArg(&csv)
	if err!=nil{
		return err
	}
	if len(pthShares)<1{
		return fmt.Errorf("you need to specify at least a share file")
	}
	muShares.Lock()
	defer muShares.Unlock()
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

//
// check expired heartbeats
//

func checkExpiredHeartBeats(){
	now := time.Now().Unix()
	muShares.RLock()
	for _,s := range shares{
		if log.IsExpired(s, now){
			startReleaseProtocol(s)
		}	
	}
	muShares.RUnlock()
}

func startTimedExpiredHeartBeatsCheck(){
	const interval = 2
	ticker := time.NewTicker(interval *time.Hour)
	defer ticker.Stop()
	for {
		select{
			case <- ticker.C:
				checkExpiredHeartBeats()	
		}
	}
}

// release

func startReleaseProtocol(s *sign.ArmoredShare){
	cmd.ColorPrintln(fmt.Sprintf("RELEASING PROTOCOL STARTED FOR %s", s.AuthKey), cmd.Red)
	//TODO
}


