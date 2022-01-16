package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"

	"github.com/sirupsen/logrus"
)

// function to encrypt message to be sent
func Encrypt(msg string, key rsa.PublicKey) string {

	label := []byte("OAEP Encrypted")
	rng := rand.Reader

	// * using OAEP algorithm to make it more secure
	// * using sha256
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &key, []byte(msg), label)
	// check for errors
	if err != nil {
		logrus.WithError(err).Fatal("unable to encrypt")

	}

	return base64.StdEncoding.EncodeToString(ciphertext)
}

// function to decrypt message to be received
func Decrypt(cipherText string, key rsa.PrivateKey) string {

	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")
	rng := rand.Reader

	// decrypting based on same parameters as encryption
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &key, ct, label)
	// check for errors
	if err != nil {
		logrus.WithError(err).Fatal("Decrypt proccess fail")
	}
	return string(plaintext)
}
