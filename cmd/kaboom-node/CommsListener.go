package main

import(
	"fmt"
	"github.com/naus3a/kaboom/cmd"
	"github.com/naus3a/kaboom/sign"
	"github.com/naus3a/kaboom/remote"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type CommsListener struct{
	comms *remote.PubSubComms
	share *sign.ArmoredShare
	onGoodHeartBeat func(*sign.ArmoredShare, *sign.HeartBeat)
	onBadHeartBeat func(*sign.ArmoredShare)
}

func NewCommsListener(s *sign.ArmoredShare, onGoodHb func(*sign.ArmoredShare, *sign.HeartBeat), onBadHb func(*sign.ArmoredShare)) (*CommsListener, error){
	chanName := remote.MakeChannelNameNow(s.AuthKey)
	l := new(CommsListener)
	l.comms = nil
	l.share = s
	l.onGoodHeartBeat = onGoodHb
	l.onBadHeartBeat = onBadHb
	err := l.startCommsOnChannel(chanName)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (l *CommsListener)startCommsOnChannel(chanName string) error{
	l.logSuccess(fmt.Sprintf("Channel name: %s", chanName))

	if l.comms != nil {
		return nil
	}
	var err error = nil
	l.comms, err = remote.NewPubSubComms(chanName)
	if err != nil {
		return err
	}
	l.logSuccess("Comms Ready!")

	l.comms.OnMessageParsed = l.handleMessage
	go l.comms.DiscoverPeers()
	
	err = l.comms.Listen()
	if err!=nil{
		return err
	}
	l.logSuccess("Listening")

	go l.comms.ParseMessages()

	return nil
}

func (l *CommsListener)stop()error{
	if l.comms==nil{
		return nil
	}
	err := l.comms.Stop()
	if err!=nil{
		return err
	}
	l.comms = nil
	return nil
}

func (l *CommsListener)UpdateChannelNameIfNeeded()error{
	chanName := remote.MakeChannelNameNow(l.share.AuthKey)
	if chanName == l.comms.GetChannelName(){
		return nil
	}
	fmt.Println("Rotating channel")
	err := l.stop()
	if err != nil{
		return err
	}
	return l.startCommsOnChannel(chanName)
}

func (l *CommsListener)handleMessage(m *pubsub.Message){
	hb, err := sign.DecodeBinaryHeartBeat(m.Message.Data)
	if err != nil {
		return
	}

	good, err := sign.VerifyHeartBeat(l.share, &hb)
	if err!= nil {
		return
	}
	if good {
		l.onGoodHeartBeat(l.share, &hb)
	}else{
		l.onBadHeartBeat(l.share)
	}


}

func (l *CommsListener)log(msg string)string{
	return fmt.Sprintf("%s: %s", l.share.ShortId, msg)
}

func (l *CommsListener)logColored(msg string, c cmd.AnsiCode){
	cmd.ColorPrintln(l.log(msg), c)
}

func (l *CommsListener)logSuccess(msg string){
	l.logColored(msg, cmd.Green)
}
