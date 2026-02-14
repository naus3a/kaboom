package remote_test

import (
	"time"
	"testing"
	"github.com/naus3a/kaboom/remote"
)

func TestChannelName(t *testing.T){
	salt := "cippa"
	now := time.Now().UTC()
	// dont test after 22:00 UTC pleeze ;)
	// (or add the relevant test if you want)
	future := now.Add(2*time.Hour)
	chanNow := remote.MakeChannelNameNow(salt)
	chanFuture := remote.MakeChannelName(future, salt)
	if chanNow != chanFuture {
		t.Errorf("FAIL: channel name ia different for different times on the same day")
	}
}
