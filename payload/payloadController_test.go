package payload_test

import (
	"testing"
	"github.com/naus3a/kaboom/payload"
)

func TestNewArmoredPayloadKey(t *testing.T){
	key, err := payload.NewArmoredPayloadKey("testAddress", "testNote")
	if err != nil {
		t.Fatalf("NewArmoredKey threw error: %v", err)
	}
	// 32 bites -> 256 bits/6
	const expectedL = 43
	if len(key.Key) != expectedL{
		t.Errorf("FAIL: wrong key length: expected %d, found %d", expectedL, len(key.Key))
	}
}

func TestPayloadSerialization(t * testing.T){
	key, err := payload.NewArmoredPayloadKey("testAddress", "testNote")
	if err != nil {
		t.Fatalf("NewArmoredKey threw error: %v", err)
	}

	var data []byte
	data, err = key.Serialize()
	if err!= nil {
		t.Fatalf("Could not serialize key")
	}
	var key1 *payload.ArmoredPayloadKey
	key1, err = payload.Deserialize(data)
	if err != nil{
		t.Errorf("FAIL: cannot serialize key")
	}
	if key.Key != key1.Key {
		t.Errorf("FAIL: deserializing serialized key corrupts data")
	}
}

func TestPayloadEncryptDecrypt(t *testing.T){
	key, err := payload.NewArmoredPayloadKey("testAddress", "testNote")
	if err != nil {
		t.Fatalf("NewArmoredKey threw error: %v", err)
	}
	testPayload := "this is a test"
	cipherText, err := key.Encrypt([]byte(testPayload))
	if err!= nil{
		t.Errorf("FAIL: could not encrypt")
	}
	plainText, err := key.Decrypt(cipherText)
	if err != nil {
		t.Errorf("FAIL: could not decrypt")
	}
	decrypted := string(plainText)
	if decrypted != testPayload{
		t.Errorf("Exoected: %s, but got %s", testPayload, decrypted)
	}
}
