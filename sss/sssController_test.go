package sss_test

import (
	"github.com/naus3a/kaboom/sss"
	"testing"
)

func TestSSSRecipientMaking(t *testing.T) {
	publicKeys := [][]byte{
		[]byte("onmvpenr"),
		[]byte("jnacer"),
		[]byte("ljanc;vjnr"),
	}
	_, err := sss.MakeSSSRecipient(4, publicKeys...)
	if err == nil {
		t.Errorf("FAIL: should throw a threshold error")
	}

	_, err = sss.MakeSSSRecipient(2, publicKeys...)
	if err != nil {
		t.Errorf("FAIL: could not make SSS recipient: %v", err)
	}
}
