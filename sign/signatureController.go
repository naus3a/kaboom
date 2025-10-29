package sign

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
)

type ArmoredShare struct {
	Share     string
	AuthKey   string
	Signature string
}

func NewArmoredShare(share []byte, pubKey []byte, privKey []byte) *ArmoredShare {
	sig := signData(share, privKey)
	return &ArmoredShare{
		Share:     base64.RawURLEncoding.EncodeToString(share),
		AuthKey:   base64.RawURLEncoding.EncodeToString(pubKey),
		Signature: base64.RawURLEncoding.EncodeToString(sig),
	}
}

// MakeSigningKeys returns a public and a private key (int this order) + an error if anything goes wrong
func MakeSigningKeys() ([]byte, []byte, error) {
	return ed25519.GenerateKey(rand.Reader)
}

func signData(data []byte, privKey []byte) []byte {
	pk := ed25519.PrivateKey(privKey)
	signature := ed25519.Sign(pk, data)
	return signature
}

// SignShares retyrns an array of signed shares
func SignShares(pubKey []byte, privKey []byte, shares [][]byte) []*ArmoredShare {
	signed := make([]*ArmoredShare, len(shares))
	for i := 0; i < len(shares); i++ {
		signed[i] = NewArmoredShare(shares[i], pubKey, privKey)
	}
	return signed
}
