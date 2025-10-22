package sss

import (
	"fmt"
)

func MakeSSSRecipient(threshold int, publicKeys ...[]byte) ([]byte, error) {
	if len(publicKeys) < threshold {
		return nil, fmt.Errorf("the number of shared must be bigger than the threshold")
	}
	return nil, nil
}
