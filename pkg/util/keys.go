package util

import (
	"os"
)

var PrivateKey = os.Getenv("CRYPTOPASS_PRIVATE_KEY")
var PublicKey = os.Getenv("CRYPTOPASS_PUBLIC_KEY")


func PromptPrivateKey() (string, error) {

	if len(PrivateKey)> 0 {
		return PrivateKey, nil
	}

	PrivateKey = PromptPassword("Enter Private Key: ")

	return PrivateKey, nil
}

func PromptPublicKey() (string, error) {

	if len(PublicKey)> 0 {
		return PublicKey, nil
	}

	PublicKey = PromptPassword("Enter Public Key: ")

	return PublicKey, nil
}