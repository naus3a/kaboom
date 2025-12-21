package remote_test

import(
	"context"
	"testing"
	"github.com/naus3a/kaboom/remote"
)

func TestPubSubSender(t *testing.T){
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(
		func(){
			cancel()
		},
	)

	comms, err := remote.NewPubSubComms("cippa", ctx)
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
