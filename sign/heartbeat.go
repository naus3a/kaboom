package sign

import (
	"time"
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"encoding/base64"
	"encoding/gob"
)

type HeartBeat struct{
	Epoch int64
	AllGood bool
	Signature string
}

// GetData returns a byte buffer, ready to be signed
func (h *HeartBeat) getData()([]byte, error){
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

func (h *HeartBeat)Equals(other *HeartBeat)bool{
	if other==nil {
		return false
	}
	return h.Epoch==other.Epoch && h.AllGood==other.AllGood && h.Signature==other.Signature
}

func (h *HeartBeat)Encode()([]byte, error){
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(h)
	return buf.Bytes(), err
}

func DecodeBinaryHeartBeat(data []byte)(HeartBeat, error){
	var hb HeartBeat
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(&hb)
	return hb, err
}

func (h *HeartBeat)sign(k *SigningKeys)(error){
	data, err := h.getData()
	if err != nil{
		return err
	}
	h.Signature = base64.RawURLEncoding.EncodeToString(k.SignData(data))
	return nil
}

func VerifyHeartBeat(s *ArmoredShare, h *HeartBeat)(bool, error){
	data, err := h.getData()
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
