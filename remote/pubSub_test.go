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

	sender, err := remote.NewPubSubCommsSender("cippa", ctx)
	if err!= nil {
		t.Errorf("FAIL: cannot create a pubsub sender: %v", err)
	}
	err = sender.Send([]byte("cippa"))
	if err != nil {
		t.Errorf("FAIL: cannot send message: %v", err)
	}
}
