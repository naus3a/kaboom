package sign

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"encoding/binary"
)

type ArmoredShare struct {
	TTL	uint64
	Share     string
	AuthKey   string
	Signature string
	ShortId   string
	TtlSignature string
}

func NewArmoredShare(share []byte, ttl uint64, key *SigningKeys) *ArmoredShare {
	sig := key.SignData(share)
	sid := makeShortId(sig)
	bufTtl := make([]byte, 8)
	binary.BigEndian.PutUint64(bufTtl, ttl)
	sit := key.SignData(bufTtl)
	return &ArmoredShare{
		Share:     base64.RawURLEncoding.EncodeToString(share),
		AuthKey:   base64.RawURLEncoding.EncodeToString(key.Public),
		Signature: base64.RawURLEncoding.EncodeToString(sig),
		ShortId:   sid,
		TTL: ttl,
		TtlSignature: base64.RawURLEncoding.EncodeToString(sit),
	}
}

func (s *ArmoredShare) verify(data []byte, signed []byte) (bool, error){
	kData, err := base64.RawURLEncoding.DecodeString(s.AuthKey)
	if err!= nil{
		return false, err
	}
	k := ed25519.PublicKey(kData)
	isGood := ed25519.Verify(k, data, signed)
	return isGood, nil
}

func (s *ArmoredShare) VerifyShare(other *ArmoredShare) (bool, error) {
	if other.AuthKey != s.AuthKey {
		return false, nil
	}
	sig, err := base64.RawURLEncoding.DecodeString(other.Signature)
	if err != nil {
		return false, nil
	}
	data, err := base64.RawURLEncoding.DecodeString(other.Share)
	if err != nil {
		return false, nil
	}
	return s.verify(data, sig)
}

func (s *ArmoredShare) VerifyTtl()(bool, error){
	bufTtl := make([]byte, 8)
	binary.BigEndian.PutUint64(bufTtl, s.TTL)
	sig, err := base64.RawURLEncoding.DecodeString(s.TtlSignature)
	if err != nil{
		return false, err
	}
	return s.verify(bufTtl, sig)
}

// Serialize returns a json version of the signed share
func (s *ArmoredShare) Serialize() ([]byte, error) {
	return json.Marshal(s)
}

func (s *ArmoredShare) GetData() ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s.Share)
}

func makeShortId(signature []byte) string {
	const lenId = 12
	hasher := sha256.New()
	hasher.Write(signature)
	fullHash := hasher.Sum(nil)
	shortHash := fullHash[:lenId]
	return base64.RawURLEncoding.EncodeToString(shortHash)
}

// DeserializeShare deserializes from json
func DeserializeShare(data []byte) (*ArmoredShare, error) {
	s := ArmoredShare{}
	err := json.Unmarshal(data, &s)
	return &s, err
}

type SigningKeys struct {
	Private []byte
	Public  []byte
}

// MakeSigningKeys return a nrew pair of signing/autheticating keys
func NewSigningKeys() (*SigningKeys, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &SigningKeys{
		Private: ed25519.PrivateKey(priv),
		Public:  ed25519.PublicKey(pub),
	}, nil
}

// SignData signs arbitrary data
func (k *SigningKeys) SignData(data []byte) []byte {
	return ed25519.Sign(k.Private, data)
}

// SignShares returns an array of signed shares
func (k *SigningKeys) SignShares(shares [][]byte, ttl uint64) []*ArmoredShare {
	signed := make([]*ArmoredShare, len(shares))
	for i := 0; i < len(shares); i++ {
		signed[i] = NewArmoredShare(shares[i], ttl, k)
	}
	return signed
}

type SerializedSigningKeys struct {
	Public  string
	Private string
}

func (k *SigningKeys) ToSerializedKeys()*SerializedSigningKeys{
	return &SerializedSigningKeys{
		Public: base64.RawURLEncoding.EncodeToString(k.Public),
		Private: base64.RawURLEncoding.EncodeToString(k.Private), 
	}
}

func (k *SigningKeys) Serialize() ([]byte, error) {
	serialized := k.ToSerializedKeys()
	return json.Marshal(serialized)
}

func DeserializeSigningKeys(jsonData []byte) (*SigningKeys, error) {
	serialized := SerializedSigningKeys{}
	err := json.Unmarshal(jsonData, &serialized)
	if err != nil {
		return nil, err
	}
	pubData, err := base64.RawURLEncoding.DecodeString(serialized.Public)
	if err != nil {
		return nil, err
	}
	privData, err := base64.RawURLEncoding.DecodeString(serialized.Private)
	if err != nil {
		return nil, err
	}
	return &SigningKeys{
		Public:  pubData,
		Private: privData,
	}, nil
}
