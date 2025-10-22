package payload

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"filippo.io/age"
	"fmt"
	"github.com/hashicorp/vault/shamir"
	"io"
)

type ArmoredPayloadKey struct {
	Key         string
	IPFSAddress string
	Notes       string
}

// NewArmoredPayloadKey creates an object with a random key, the payload address and any additional note
func NewArmoredPayloadKey(ipfsAddr string, notes string) (*ArmoredPayloadKey, error) {
	k, err := MakeSymmetricKey()
	if err != nil {
		return nil, err
	}
	return &ArmoredPayloadKey{
		Key:         k,
		IPFSAddress: ipfsAddr,
		Notes:       notes,
	}, nil
}

// Serialize serialize the armored key to json
func (k *ArmoredPayloadKey) Serialize() ([]byte, error) {
	return json.Marshal(k)
}

//Encrypt encrypts plaintext using symmetric key
func (k *ArmoredPayloadKey) Encrypt(plainText []byte) ([]byte, error) {
	rec, err := age.NewScryptRecipient(k.Key)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt payload: %w", err)
	}
	var buf bytes.Buffer
	w, err := age.Encrypt(&buf, rec)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt payload: %w", err)
	}
	_, err = io.Copy(w, bytes.NewReader(plainText))
	if err != nil {
		w.Close()
		return nil, fmt.Errorf("could not encrypt payload: %w", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("could not encrypt payload: %w", err)
	}
	return buf.Bytes(), nil

}

// Decrypt decrypta ciphertezt using armored key
func (k *ArmoredPayloadKey) Decrypt(cipherText []byte) ([]byte, error) {
	identity, err := age.NewScryptIdentity(k.Key)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt payload: %w", err)
	}
	r, err := age.Decrypt(bytes.NewReader(cipherText), identity)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt payload: %w", err)
	}
	plainText, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt payload: %w", err)
	}
	return plainText, nil
}

// Split an armored key in shares using Shamir Shared Secret
func (k *ArmoredPayloadKey) Split(threshold int, numParts int) ([][]byte, error) {
	if threshold > numParts {
		return nil, fmt.Errorf("number of parts should be bigger than threshold")
	}
	jsonData, err := k.Serialize()
	if err != nil {
		return nil, err
	}
	return shamir.Split(jsonData, numParts, threshold)
}

// CombineSharesInArmoredPayloadKey recombines a number of shares into the original ArmoredPayloadKey
func CombineSharesInArmoredPayloadKey(shares [][]byte) (*ArmoredPayloadKey, error) {
	jsonData, err := shamir.Combine(shares)
	if err != nil {
		return nil, err
	}
	return Deserialize(jsonData)
}

// Deserialize deserialize json to armored key
func Deserialize(data []byte) (*ArmoredPayloadKey, error) {
	k := ArmoredPayloadKey{}
	err := json.Unmarshal(data, &k)
	return &k, err
}

// MakeRandomData creates a random byte array of a given size
func MakeRandomData(sz int) ([]byte, error) {
	d := make([]byte, sz)
	_, err := rand.Read(d)
	return d, err
}

// MakeSymmetricKey creates a secure symmetric kwy
func MakeSymmetricKey() (string, error) {
	d, err := MakeRandomData(32)
	if err != nil {
		fmt.Println("could not generate random.data")
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(d), nil
}
