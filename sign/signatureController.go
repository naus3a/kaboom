package sign

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
)

type ArmoredShare struct {
	Share     string
	AuthKey   string
	Signature string
	ShortId	 string
}

func NewArmoredShare(share []byte, pubKey []byte, privKey []byte) *ArmoredShare {
	sig := signData(share, privKey)
	sid := makeShortId(sig)
	return &ArmoredShare{
		Share:     base64.RawURLEncoding.EncodeToString(share),
		AuthKey:   base64.RawURLEncoding.EncodeToString(pubKey),
		Signature: base64.RawURLEncoding.EncodeToString(sig),
		ShortId: sid,
	}
}

func (s *ArmoredShare) VerifyShare(other *ArmoredShare) (bool, error) {
	if other.AuthKey != s.AuthKey {
		return false, nil
	}
	kData, err := base64.RawURLEncoding.DecodeString(s.AuthKey)
	if err != nil {
		return false, err
	}
	k := ed25519.PublicKey(kData)
	sig, err := base64.RawURLEncoding.DecodeString(other.Signature)
	if err != nil {
		return false, nil
	}
	data, err := base64.RawURLEncoding.DecodeString(other.Share)
	if err != nil {
		return false, nil
	}
	isGood := ed25519.Verify(k, data, sig)
	return isGood, nil
}

// Serialize returns a json version of the signed share
func (s * ArmoredShare) Serialize() ([]byte, error){
	return json.Marshal(s)
}

func makeShortId(signature []byte)string{
	const lenId = 12
	hasher := sha256.New()
	hasher.Write(signature)
	fullHash := hasher.Sum(nil)
	shortHash := fullHash[:lenId]
	return base64.RawURLEncoding.EncodeToString(shortHash)
}

// DeserializeShare deserializes from json
func DeserializeShare(data []byte)(*ArmoredShare, error){
	s := ArmoredShare{}
	err := json.Unmarshal(data, &s)
	return &s, err
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
