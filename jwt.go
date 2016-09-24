package apnsapi

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/big"
	"time"
)

type jwtHeader struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
}

type jwtClaim struct {
	Iss string `json:"iss"`
	Iat int64  `json:"iat"`
}

type ecdsaSignature struct {
	R, S *big.Int
}

func CreateToken(key *ecdsa.PrivateKey, kid string, teamID string) (string, error) {
	header := jwtHeader{
		Alg: "ES256",
		Kid: kid,
	}
	claim := jwtClaim{
		Iss: teamID,
		Iat: time.Now().Unix(),
	}
	var b bytes.Buffer

	headerJson, err := json.Marshal(&header)
	if err != nil {
		return "", err
	}
	if err := base64Encode(&b, headerJson); err != nil {
		return "", err
	}
	b.WriteString(".")

	claimJson, err := json.Marshal(&claim)
	if err != nil {
		return "", err
	}
	if err := base64Encode(&b, claimJson); err != nil {
		return "", err
	}

	h := crypto.SHA256.New()
	h.Write(b.Bytes())
	msg := h.Sum(nil)

	r, s, err := ecdsa.Sign(rand.Reader, key, msg)
	if err != nil {
		return "", err
	}

	sig, err := asn1.Marshal(ecdsaSignature{r, s})
	if err != nil {
		return "", err
	}

	b.WriteString(".")
	if err := base64Encode(&b, sig); err != nil {
		return "", err
	}

	return b.String(), nil
}

func base64Encode(w io.Writer, byt []byte) error {
	enc := base64.NewEncoder(base64.RawURLEncoding, w)
	enc.Write(byt)
	return enc.Close()
}
