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

	for i := 0; i < len(signed); i++ {
		if signed[i].Share != base64.RawURLEncoding.EncodeToString(fakeShares[i]) {
			t.Errorf("FAIL: corrupted share")
		}
		if signed[i].AuthKey != base64.RawURLEncoding.EncodeToString(pubKey) {
			t.Errorf("FAIL: corrupted authentication key")
		}
	}
}
