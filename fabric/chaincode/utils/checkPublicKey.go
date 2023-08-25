package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/goledgerdev/cc-tools/errors"
)

func CheckPublicKey(pubKeyB64 string) (*rsa.PublicKey, error) {
	pubKey, err := base64.StdEncoding.DecodeString(pubKeyB64)
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "error decoding base64", 400)
	}

	block, _ := pem.Decode([]byte(pubKey))
	if block == nil {
		return nil, errors.NewCCError("could not decode public key", 400)
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "could not parse PKIX public key", 400)
	}

	switch pub.(type) {
	case *rsa.PublicKey:
		fmt.Println("pub is of type RSA, continuing...")
	default:
		return nil, errors.NewCCError("public key is not RSA", 400)
	}

	return pub.(*rsa.PublicKey), nil
}
