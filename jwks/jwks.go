package jwks

type JWK struct {
	KeyType      string `json:"kty,omitempty"`
	KeyID        string `json:"kid,omitempty"`
	PublicKeyUse string `json:"use,omitempty"`
	Algorithm    string `json:"alg,omitempty"`
	// RSA
	N string `json:"n,omitempty"`
	E string `json:"e,omitempty"`
	// Elliptic
	X     string `json:"x,omitempty"`
	Y     string `json:"y,omitempty"`
	Curve string `json:"crv,omitempty"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}
