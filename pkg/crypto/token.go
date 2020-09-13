package crypto

import (
	"github.com/cryptorities/cryptopass/pkg/app"
	"github.com/cryptorities/cryptopass/pkg/util"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
	"time"
	"fmt"
	"crypto/ed25519"
)

/**
	Alex Shvid
 */

type PublicKeyProvider func() ([]byte, error)
type PrivateKeyProvider func() ([]byte, error)

func Issue(username, expirationDate string, privateKeyProv PrivateKeyProvider) (string, error) {
	return Sign(username, expirationDate, app.IssueSep, privateKeyProv)
}

func Revoke(username, expirationDate string, privateKeyProv PrivateKeyProvider) (string, error) {
	return Sign(username, expirationDate, app.RevokeSep, privateKeyProv)
}

func Sign(username, date, sep string, privateKeyProv PrivateKeyProvider) (string, error) {

	rig := username
	wal := ""

	rigIdx := strings.IndexByte(username, '.')
	if rigIdx != -1 {
		rig = username[rigIdx+1:]
		wal = username[:rigIdx]
	}

	if wal != "" {
		fmt.Fprintf(os.Stderr,"Warn: wal %s is not using for generating token\n", wal)
	}

	datetime, err := time.Parse(app.DateFormatISO, date)
	if err != nil {
		return "", errors.Errorf("invalid expiration date '%s', %v", date, err)
	}

	if datetime.Before(time.Now()) {
		return "", errors.Errorf("expiration date is before than now '%s'", datetime.String())
	}

	days := util.DaysOffset(&datetime)

	signingString := fmt.Sprintf("%s%s%d", rig, sep, days)

	//fmt.Printf("Signing rig='%s' days=%d, signingString='%s'\n", rig, days, signingString)

	privateKey, err := privateKeyProv()
	if err != nil {
		return "", err
	}

	if len(privateKey) != ed25519.PrivateKeySize {
		return "", errors.Errorf("invalid ed25519 private key len %d", len(privateKey))
	}

	sig := ed25519.Sign(ed25519.PrivateKey(privateKey), []byte(signingString))

	sigBase := app.Encoding.EncodeToString(sig)

	return fmt.Sprintf("%s.%d", sigBase, days), nil

}

func VerifyIssued(username string, token string, publicKeyProv PublicKeyProvider) (bool, string, error) {
	return Verify(username, token, app.IssueSep, publicKeyProv)
}

func VerifyRevoked(username string, token string, publicKeyProv PublicKeyProvider) (bool, string, error) {
	return Verify(username, token, app.RevokeSep, publicKeyProv)
}

func Verify(username string, token string, sep string, publicKeyProv PublicKeyProvider) (bool, string, error) {

	rig := username
	wal := ""

	rigIdx := strings.IndexByte(username, '.')
	if rigIdx != -1 {
		rig = username[rigIdx+1:]
		wal = username[:rigIdx]
	}

	if wal != "" {
		fmt.Fprintf(os.Stderr,"Warn: wal %s is not using for signing token\n", wal)
	}

	daysIdx := strings.LastIndexByte(token, '.')
	if daysIdx == -1 {
		return false, "", errors.New("days separator not found in token")
	}

	signStr := token[:daysIdx]
	daysStr := token[daysIdx+1:]

	days, err := strconv.Atoi(daysStr)
	if err != nil {
		return false, "", errors.Errorf("invalid days '%s' in token, %v", daysStr, err)
	}

	expiration := util.ParseDaysOffset(days)
	expirationFmt := expiration.Format(app.DateFormatISO)

	if sep == app.IssueSep && expiration.Before(time.Now()) {
		return false, expirationFmt, app.ErrExpired
	}

	sign, err := app.Encoding.DecodeString(signStr)
	if err != nil {
		return false, expirationFmt, errors.Errorf("invalid sign '%s' in token, %v", signStr, err)
	}

	signingString := fmt.Sprintf("%s%s%d", rig, sep, days)

	//fmt.Printf("Verifying rig='%s' days=%d, signingString='%s'\n", rig, days, signingString)

	publicKey, err := publicKeyProv()
	if err != nil {
		return false, expirationFmt, err
	}

	if len(publicKey) != ed25519.PublicKeySize {
		return false, expirationFmt, errors.Errorf("invalid ed25519 public key len %d", len(publicKey))
	}

	return ed25519.Verify(ed25519.PublicKey(publicKey), []byte(signingString), sign), expirationFmt, nil

}

