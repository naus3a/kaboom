package sign

import (
	"time"
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"encoding/base64"
)

type HeartBeat struct{
	Epoch int64
	AllGood bool
	Signature string
}

// GetData returns a byte buffer, ready to be signed
func (h *HeartBeat) GetData()([]byte, error){
	buf := new(bytes.Buffer)
	order := binary.BigEndian
	err := binary.Write(buf, order, h.Epoch)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, order, h.AllGood)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func NewHeartBeat(allGood bool, key *SigningKeys)(*HeartBeat, error){
	now := time.Now()
	epoch := now.Unix()
	hb := &HeartBeat{
		Epoch: epoch,
		AllGood: allGood,
	}
	err := hb.sign(key)
	if err != nil {
		return nil, err
	}
	return hb, nil
}

func (h *HeartBeat)sign(k *SigningKeys)(error){
	data, err := h.GetData()
	if err != nil{
		return err
	}
	h.Signature = base64.RawURLEncoding.EncodeToString(k.SignData(data))
	return nil
}

func VerifyHeartBeat(s *ArmoredShare, h *HeartBeat)(bool, error){
	data, err := h.GetData()
	if err != nil{
		return false, err
	}
	kData, err := base64.RawURLEncoding.DecodeString(s.AuthKey)
	if err != nil {
		return false, err
	}
	k := ed25519.PublicKey(kData)

	sig, err := base64.RawURLEncoding.DecodeString(h.Signature)
	if err != nil {
		return false, err
	}

	isGood := ed25519.Verify(k, data, sig)
	return isGood, nil
}
