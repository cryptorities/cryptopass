package util

import (
	"github.com/cryptorities/cryptopass/pkg/app"
	"os"
)

var PrivateKey, _ = app.Encoding.DecodeString(os.Getenv("CRYPTOPASS_PRIVATE_KEY"))
var PublicKey, _ = app.Encoding.DecodeString(os.Getenv("CRYPTOPASS_PUBLIC_KEY"))


func PromptPrivateKey() ([]byte, error) {

	if len(PrivateKey)> 0 {
		return PrivateKey, nil
	}

	privateKeyHex := PromptPassword("Enter Private Key: ")

	privateKey, err := app.Encoding.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func PromptPublicKey() ([]byte, error) {

	if len(PublicKey)> 0 {
		return PublicKey, nil
	}

	publicKeyHex := PromptPassword("Enter Public Key: ")

	publicKey, err := app.Encoding.DecodeString(publicKeyHex)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}