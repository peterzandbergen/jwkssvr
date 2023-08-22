package main

import "testing"

func TestGetBytes(t *testing.T) {
	b, err := getJWKSBytes(JWKSPing)
	if err != nil {
		t.Errorf("error getting bytes: %s", err)
	}
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

