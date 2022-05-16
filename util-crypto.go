package api_common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	"golang.org/x/crypto/bcrypt"
)

// CryptoHashText applies the predefined hashing algorithm
// to te given string and returns the hashed value or an
// error
func CryptoHashText(plainText string, bcryptCost int) (string, error) {
	now := time.Now()
	baHashedText, errGenerateFromPassword := bcrypt.GenerateFromPassword([]byte(plainText), bcryptCost)
	if errGenerateFromPassword != nil {
		return "", fmt.Errorf("cannot hash provided string: %s", errGenerateFromPassword.Error())
	}
	log.Tracef("completed CryptoHashText operation in %s", time.Since(now))
	return string(baHashedText), nil
}

// CryptoCompareHashedText compares the given strings, one
// hashed and the other in plain to see if they match
// It automatically reads the hashed string cost from the
// string itself
func CryptoCompareHashedText(hashedText string, plainText string) bool {
	now := time.Now()
	errCompareHashAndPassword := bcrypt.CompareHashAndPassword([]byte(hashedText), []byte(plainText))
	if errCompareHashAndPassword != nil {
		return false
	}
	log.Tracef("completed CryptoCompareHashedText operation in %s", time.Since(now))
	return true
}

// CryptoEncryptText executes encryption on the given string
// using the predefined encryption algorithm and returns the
// encrypted value or an error
func CryptoEncryptText(plainText string, encryptionKey string) (string, error) {
	now := time.Now()
	if len(encryptionKey) != 32 {
		return "", fmt.Errorf("cannot create aes-256 cipher, missing 32 bytes encryption key (passed a key of length %d)",
			len(encryptionKey))
	}
	aes256Cipher, errNewCipher := aes.NewCipher([]byte(encryptionKey))
	if errNewCipher != nil {
		return "", fmt.Errorf("cannot create aes-256 cipher: %s", errNewCipher)
	}
	gcm, errNewGcm := cipher.NewGCM(aes256Cipher)
	if errNewGcm != nil {
		return "", fmt.Errorf("cannot create gcm: %s", errNewGcm.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	_, errRandRead := rand.Read(nonce)
	if errRandRead != nil {
		return "", fmt.Errorf("cannot create random nonce: %s", errRandRead.Error())
	}
	encryptedText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	log.Tracef("completed CryptoEncryptText operation in %s", time.Since(now))
	return hex.EncodeToString(encryptedText), nil
}

// CryptoDecryptText executes decryption on the given encrypted
// string using the predefined encryption algorithm and returns
// the decrypted value or an error
func CryptoDecryptText(encryptedText string, encryptionKey string) (string, error) {
	now := time.Now()
	if len(encryptionKey) != 32 {
		return "", fmt.Errorf("cannot create aes-256 cipher, missing 32 bytes encryption key (passed a key of length %d)",
			len(encryptionKey))
	}
	aes256Cipher, errNewCipher := aes.NewCipher([]byte(encryptionKey))
	if errNewCipher != nil {
		return "", fmt.Errorf("cannot create aes-256 cipher: %s", errNewCipher)
	}
	gcm, errNewGcm := cipher.NewGCM(aes256Cipher)
	if errNewGcm != nil {
		return "", fmt.Errorf("cannot create gcm: %s", errNewGcm.Error())
	}
	baEncryptedText, errHexDecodeString := hex.DecodeString(encryptedText)
	if errHexDecodeString != nil {
		return "", fmt.Errorf("cannot hex decode string: %s", errHexDecodeString.Error())
	}
	if len(baEncryptedText) <= gcm.NonceSize() {
		return "", fmt.Errorf("cannot decrypt text with invalid length")
	}
	plainText, errGcmOpen := gcm.Open(nil, baEncryptedText[:gcm.NonceSize()], baEncryptedText[gcm.NonceSize():], nil)
	if errGcmOpen != nil {
		return "", fmt.Errorf("cannot decrypt text with gcm open: %s", errGcmOpen.Error())
	}
	log.Tracef("completed CryptoDecryptText operation in %s", time.Since(now))
	return string(plainText), nil
}

// CryptoSha256String returns the sha256 hash of the given input string
func CryptoSha256String(input string) string {
	h := sha256.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}
