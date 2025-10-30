package sign

import (
	"time"
	"bytes"
	"encoding/base64"
	"encoding/binary"
)

type HeartBeat struct{
	Epoch int64
	Signature string
}

func NewHeartBeat(privKey []byte)(*HeartBeat, error){
	now := time.Now()
	epoch := now.Unix()
	sig, err := SignTime(epoch, privKey)
	if err != nil {
		return nil, err
	}
	return &HeartBeat{
		Epoch: epoch,	
		Signature: base64.RawURLEncoding.EncodeToString(sig),
	}, nil
}

func SignTime(epoch int64, privKey []byte)([]byte, error){
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, epoch)
	if err != nil{
		return nil, err
	}
	data := buf.Bytes()
	return SignData(data, privKey), nil
}
