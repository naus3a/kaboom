package sign_test

import (
	"testing"
	"github.com/naus3a/kaboom/sign"
	"time"
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
		t.Errorf("FAIL: didnt handle correctly an outdated hearrbeat")
	}

	key2, _ := sign.NewSigningKeys()
	id2 := key2.ToSerializedKeys()
	hb2, _ := sign.NewHeartBeat(true, key2)
	hbl.LogHeartBeat(id2.Public, hb2)
	if hbl.GetNumIds()!=2{
		t.Errorf("FAIL: didnt add a new identity")
	}
}

func TestHeartBeatExpiration(t *testing.T){
	key, _ := sign.NewSigningKeys()
	hb, _ := sign.NewHeartBeat(true, key)

	now := time.Now().Unix()
	const aDay int64 = 24

	if hb.IsExpired(aDay, now) {
		t.Errorf("FAIL: heartbeat expires too soon")
	}
	
	yesterday := now - (25*60*60)
	if !hb.IsExpired(aDay, yesterday){
		t.Errorf("FAIL: heartbeat did not expire")
	}
}

func TestHeartBeatExpirationWithShare(t *testing.T){
	key, _ := sign.NewSigningKeys()
	key2, _ := sign.NewSigningKeys()
	id := key.ToSerializedKeys()
	fakeShare := []byte("111111111")
	fakeShare2 := []byte("222222222")
	s := sign.NewArmoredShare(fakeShare, 24, key)
	s2:= sign.NewArmoredShare(fakeShare2, 24, key2)
	hb, _ := sign.NewHeartBeat(true, key)
	hbl := sign.NewHeartBeatLog()
	hbl.LogHeartBeat(id.Public, hb)
	now := time.Now().Unix()
	yesterday := now -(25*60*60)
	if hbl.IsExpired(s, now) {
		t.Errorf("FAIL: good share shows up expired")
	}
	if !hbl.IsExpired(s, yesterday){
		t.Errorf("FAIL: share did not expire on an old heartbeat")
	}
	if hbl.IsExpired(s2, yesterday){
		t.Errorf("FAIL: found share which wasnt logged")
	}
}

func forgeTestHeartBeat(epoch int64)(*sign.HeartBeat){
	return &sign.HeartBeat{
		Epoch: epoch,
		AllGood: true,
		Signature: "",
	}
}
