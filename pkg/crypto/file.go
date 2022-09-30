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

	content, err := EncryptData(fileContent, publicKeyProv)
	if err != nil {
		return 0, err
	}

	return len(content), ioutil.WriteFile(outputFile, content, os.FileMode(0660))
}

func EncryptStream(inputStream io.ReadCloser, outputStream io.WriteCloser, publicKeyProv PublicKeyProvider) (int, error) {

	fileContent, err := ioutil.ReadAll(inputStream)
	if err != nil {
		return 0, err
	}

	defer inputStream.Close()

	content, err := EncryptData(fileContent, publicKeyProv)
	if err != nil {
		return 0, err
	}

	nw, err := outputStream.Write(content)
	defer outputStream.Close()

	if nw < len(content) && err == nil {
		err = errors.New("invalid write")
	}

	return nw, err

}

func EncryptData(fileContent []byte, publicKeyProv PublicKeyProvider) ([]byte, error) {

	publicKeyEncoded, err := publicKeyProv()
	if err != nil {
		return nil, err
	}

	publicKey, err := app.Encoding.DecodeString(publicKeyEncoded)
	if err != nil {
		return nil, err
	}

	if len(publicKey) != ed25519.PublicKeySize {
		return nil, errors.Errorf("invalid ed25519 public key len %d", len(publicKey))
	}

	var edwardsPublicKey [32]byte
	copy(edwardsPublicKey[:], publicKey)

	var curvePublicKey [32]byte
	if !util.PublicKeyToCurve25519(&curvePublicKey, &edwardsPublicKey) {
		return nil, errors.New("can not convert ed25519 public key to curve25519 public key")
	}

	return Encrypt(&curvePublicKey, fileContent)

}

func DecryptFile(inputFile, outputFile string, privateKeyProv PrivateKeyProvider) (int, error) {

	fileContent, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return 0, err
	}

	content, err := DecryptData(fileContent, privateKeyProv)
	if err != nil {
		return 0, err
	}

	return len(content), ioutil.WriteFile(outputFile, content, os.FileMode(0660))
}

func DecryptStream(inputStream io.ReadCloser, outputStream io.WriteCloser, privateKeyProv PrivateKeyProvider) (int, error) {

	fileContent, err := ioutil.ReadAll(inputStream)
	if err != nil {
		return 0, err
	}

	defer inputStream.Close()

	content, err := DecryptData(fileContent, privateKeyProv)
	if err != nil {
		return 0, err
	}

	nw, err := outputStream.Write(content)
	defer outputStream.Close()

	if nw < len(content) && err == nil {
		err = errors.New("invalid write")
	}

	return nw, err
}

func DecryptData(fileContent []byte, privateKeyProv PrivateKeyProvider) ([]byte, error) {

	privateKeyEncoded, err := privateKeyProv()
	if err != nil {
		return nil, err
	}

	privateKey, err := app.Encoding.DecodeString(privateKeyEncoded)
	if err != nil {
		return nil, err
	}

	if len(privateKey) != ed25519.PrivateKeySize {
		return nil, errors.Errorf("invalid ed25519 private key len %d", len(privateKey))
	}

	var edwardsPrivateKey [64]byte
	copy(edwardsPrivateKey[:], privateKey)

	var curvePrivateKey [32]byte
	util.PrivateKeyToCurve25519(&curvePrivateKey, &edwardsPrivateKey)

	// clean
	copy(edwardsPrivateKey[:], zero64[:])

	content, err := Decrypt(&curvePrivateKey, fileContent)
	if err != nil {
		return nil, err
	}

	// clean
	copy(curvePrivateKey[:], zero32[:])

	return content, err
}

func Encrypt(recipient *[32]byte, content []byte) ([]byte, error) {

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

func Decrypt(privateKey *[32]byte, content []byte) ([]byte, error) {

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
