package sign

import (
	"time"
	"bytes"
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

func NewHeartBeat(allGood bool, privKey []byte)(*HeartBeat, error){
	now := time.Now()
	epoch := now.Unix()
	hb := &HeartBeat{
		Epoch: epoch,
		AllGood: allGood,
	}
	err := hb.sign(privKey)
	if err != nil {
		return nil, err
	}
	return hb, nil
}

func (h *HeartBeat)sign(privKey []byte)( error){
	data, err := h.GetData()
	if err != nil{
		return err
	}
	h.Signature = base64.RawURLEncoding.EncodeToString(SignData(data, privKey))
	return nil
}
