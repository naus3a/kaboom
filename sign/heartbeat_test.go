package sign_test

import(
	"testing"
	"github.com/naus3a/kaboom/sign"
)

func TestHeartBeat(t *testing.T){
	pubKey, privKey, _ := sign.MakeSigningKeys()
	hb, err := sign.NewHeartBeat(true, privKey)
	if err != nil {
		t.Errorf("FAIL: coulndt create a heartbeat: %v", err)
	}

	fakeShare := []byte("123456")
	share := sign.NewArmoredShare(fakeShare, pubKey, privKey)

	b, err := sign.VerifyHeartBeat(share, hb)
	if err != nil {
		t.Errorf("FAIL: error verifying heartbeat: %v", err)
	}
	if !b {
		t.Errorf("FAIL: could not verify a good heartbeat")
	}
	
	hb.Epoch = hb.Epoch + 10
	b, err = sign.VerifyHeartBeat(share, hb)
	if err !=nil{
		t.Errorf("FAIL: error verifying a tampered heartbeat: %v", err)
	}
	if b {
		t.Errorf("FAIL: did not catch a tampered heartbeat")
	}
}
