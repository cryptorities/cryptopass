package crypto

import (
	"crypto/rand"
	"github.com/cryptorities/cryptopass/pkg/app"
	"github.com/cryptorities/cryptopass/pkg/util"
	"github.com/pkg/errors"
	"crypto/ed25519"
	"golang.org/x/crypto/nacl/box"
	"io"
	"io/ioutil"
	"os"
)

/**
	Alex Shvid
*/

const (
	// file format
	NonceSize = 24
	// content
	PublicKeySize = 32
)

var (
	zero32 [32]byte
  	zero64 [64]byte
)

func EncryptFile(inputFile, outputFile string, publicKeyProv PublicKeyProvider) (int, error) {

	fileContent, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return 0, err
	}

	publicKeyEncoded, err := publicKeyProv()
	if err != nil {
		return 0, err
	}

	publicKey, err := app.Encoding.DecodeString(publicKeyEncoded)
	if err != nil {
		return 0, err
	}

	if len(publicKey) != ed25519.PublicKeySize {
		return 0, errors.Errorf("invalid ed25519 public key len %d", len(publicKey))
	}

	var edwardsPublicKey [32]byte
	copy(edwardsPublicKey[:], publicKey)

	var curvePublicKey [32]byte
	if !util.PublicKeyToCurve25519(&curvePublicKey, &edwardsPublicKey) {
		return 0, errors.New("can not convert ed25519 public key to curve25519 public key")
	}

	content, err := encrypt(&curvePublicKey, fileContent)
	if err != nil {
		return 0, err
	}

	return len(content), ioutil.WriteFile(outputFile, content, os.FileMode(0660))
}

func DecryptFile(inputFile, outputFile string, privateKeyProv PrivateKeyProvider) (int, error) {

	fileContent, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return 0, err
	}

	privateKeyEncoded, err := privateKeyProv()
	if err != nil {
		return 0, err
	}

	privateKey, err := app.Encoding.DecodeString(privateKeyEncoded)
	if err != nil {
		return 0, err
	}

	if len(privateKey) != ed25519.PrivateKeySize {
		return 0, errors.Errorf("invalid ed25519 private key len %d", len(privateKey))
	}

	var edwardsPrivateKey [64]byte
	copy(edwardsPrivateKey[:], privateKey)

	var curvePrivateKey [32]byte
	util.PrivateKeyToCurve25519(&curvePrivateKey, &edwardsPrivateKey)

	// clean
	copy(edwardsPrivateKey[:], zero64[:])

	content, err := decrypt(&curvePrivateKey, fileContent)
	if err != nil {
		return 0, err
	}

	// clean
	copy(curvePrivateKey[:], zero32[:])

	return len(content), ioutil.WriteFile(outputFile, content, os.FileMode(0660))
}

func encrypt(recipient *[32]byte, content []byte) ([]byte, error) {

	boxPublicKey, boxPrivateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	var nonce [24]byte
	_, err = io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, err
	}

	ciphertext := box.Seal(nonce[:], content, &nonce, recipient, boxPrivateKey)
	return append(ciphertext, boxPublicKey[:]...), nil

}

func decrypt(privateKey *[32]byte, content []byte) ([]byte, error) {

	if len(content) < PublicKeySize + NonceSize {
		return nil, errors.Errorf("insufficient file size %d", len(content))
	}

	var decryptNonce [24]byte
	copy(decryptNonce[:], content[:24])
	content = content[24:]

	var peerPubKey [32]byte
	copy(peerPubKey[:], content[len(content)-32:])
	content = content[:len(content)-32]

	plaintext, ok := box.Open(nil, content, &decryptNonce, &peerPubKey, privateKey)
	if !ok {
		return nil, errors.New("unseal nacl box error")
	}
	return plaintext, nil

}
