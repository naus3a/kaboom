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

	_, err := remote.NewPubSubCommsSender("cippa", ctx)
	if err!= nil {
		t.Errorf("FAIL: cannot create a pubsub sender: %v", err)
	}
}
