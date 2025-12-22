package sign_test

import (
	"testing"
	"github.com/naus3a/kaboom/sign"
)

func TestLogHeartBeat(t *testing.T){
	key, _ := sign.NewSigningKeys()
	id := key.ToSerializedKeys()
	hb, _ := sign.NewHeartBeat(true, key)
	hbl := sign.NewHeartBeatLog()
	hbl.LogHeartBeat(id.Public, hb)

	h, found := hbl.GetLastHeartBeat(id.Public)
	if !found {
		t.Errorf("FAIL: didnt add heartbeat")
	}
	if !h.Equals(hb){
		t.Errorf("FAIL: added a corrupted heartbeat")
	}

	newerHb := forgeTestHeartBeat(hb.Epoch + 100)
	hbl.LogHeartBeat(id.Public, newerHb)
	h, found = hbl.GetLastHeartBeat(id.Public)
	if !h.Equals(newerHb){
		t.Errorf("FAIL: didnt update heartbeat")
	}

	olderHb := forgeTestHeartBeat(hb.Epoch - 100)
	hbl.LogHeartBeat(id.Public, olderHb)
	h, found = hbl.GetLastHeartBeat(id.Public)
	if !h.Equals(newerHb){
		t.Errorf("FAIL: didnt handle correctly an ourdated hearrbeat")
	}

	key2, _ := sign.NewSigningKeys()
	id2 := key2.ToSerializedKeys()
	hb2, _ := sign.NewHeartBeat(true, key2)
	hbl.LogHeartBeat(id2.Public, hb2)
	if hbl.GetNumIds()!=2{
		t.Errorf("FAIL: didnt add a new identity")
	}
}

func forgeTestHeartBeat(epoch int64)(*sign.HeartBeat){
	return &sign.HeartBeat{
		Epoch: epoch,
		AllGood: true,
		Signature: "",
	}
}
