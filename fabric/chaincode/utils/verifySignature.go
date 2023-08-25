package utils

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"

	"github.com/goledgerdev/cc-tools/errors"
)

func verifySignature(pubKey *rsa.PublicKey, message, signature string) error {
	hash := sha256.Sum256([]byte(message))
	hashBytes := hash[:]

	sign, err := hex.DecodeString(signature)
	if err != nil {
		return errors.WrapErrorWithStatus(err, "error decoding hex", 400)
	}

	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashBytes, []byte(sign))
	if err != nil {
		return errors.WrapErrorWithStatus(err, "RSA signature failed", 400)
	}
	return nil
}
