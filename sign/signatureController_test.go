package sign_test

import (
	"encoding/base64"
	"github.com/naus3a/kaboom/sign"
	"testing"
)

func TestShareSigning(t *testing.T) {
	fakeShares := make([][]byte, 3)
	fakeShares[0] = []byte("000000000")
	fakeShares[1] = []byte("111111111")
	fakeShares[2] = []byte("222222222")

	pubKey, privKey, err := sign.MakeSigningKeys()
	if err != nil {
		t.Errorf("FAIL: cannot generate keys")
	}

	signed := sign.SignShares(pubKey, privKey, fakeShares)

	//test the share making
	for i := 0; i < len(signed); i++ {
		if signed[i].Share != base64.RawURLEncoding.EncodeToString(fakeShares[i]) {
			t.Errorf("FAIL: corrupted share")
		}
		if signed[i].AuthKey != base64.RawURLEncoding.EncodeToString(pubKey) {
			t.Errorf("FAIL: corrupted authentication key")
		}
	}

	b, err := signed[0].VerifyShare(signed[1])
	if err != nil {
		t.Errorf("FAIL: verifying signature throws error: %v", err)
	}
	if !b {
		t.Errorf("FAIL: verifying good signature failed")
	}

	wrongPub, wrongPriv, _ := sign.MakeSigningKeys()
	wrongShares := sign.SignShares(wrongPub, wrongPriv, fakeShares)
	b, _ = signed[0].VerifyShare(wrongShares[0])
	if b {
		t.Errorf("FAIL: failed to identify different authentication key")
	}

	tampered := &sign.ArmoredShare{
		Share:     "badstuffhere",
		Signature: signed[1].Signature,
		AuthKey:   signed[1].AuthKey,
	}
	b, _ = signed[0].VerifyShare(tampered)
	if b {
		t.Errorf("FAIL: failed to identify tampered share")
	}
}
