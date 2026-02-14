package remote_test

import(
	"testing"
	"github.com/naus3a/kaboom/remote"
)

func TestPubSubSender(t *testing.T){
	comms, err := remote.NewPubSubComms("cippa")
	if err!= nil {
		t.Errorf("FAIL: cannot create pubsub comms: %v", err)
	}
	err = comms.Send([]byte("cippa"))
	if err != nil {
		t.Errorf("FAIL: cannot send message: %v", err)
	}
	err = comms.Listen()
	if err != nil{
		t.Errorf("FAIL: cannot listen to topic: %v", err)
	}
}
