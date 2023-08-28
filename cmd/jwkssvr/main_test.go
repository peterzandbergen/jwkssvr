package main

import (
	"io"
	"testing"
)

func TestGetBytes(t *testing.T) {
	r, err := getBody(JWKSPing)
	if err != nil {
		t.Errorf("error getting bytes: %s", err)
	}
	b, _ := io.ReadAll(r)
	t.Logf("received %d bytes", len(b))
}

func TestGetAllJWKS(t *testing.T) {
	ks, err := getJWKS(JWKSPing)
	if err != nil {
		t.Errorf("error getting bytes: %s", err)
	}
	t.Logf("received %d keys", len(ks.Keys))
}

func TestGetAlgJWKS(t *testing.T) {
	ks, err := getJWKS(JWKSPing, filterWithAlg)
	if err != nil {
		t.Errorf("error getting bytes: %s", err)
	}
	t.Logf("received %d keys", len(ks.Keys))
}
