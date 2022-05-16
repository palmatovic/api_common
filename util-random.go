package api_common

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// RandomGenerateUuid generates a random string identifier
// including (or not) hypens
func RandomGenerateUuid(withHyphen bool) string {
	id := uuid.New()
	if !withHyphen {
		return strings.Replace(id.String(), "-", "", -1)
	}
	return id.String()
}

// RandomGenerateUuidWithLength generates a random string identifier
// including (or not) hypens, with max length
func RandomGenerateUuidWithLength(withHyphen bool, length int) string {
	out := RandomGenerateUuid(withHyphen)
	return out[:length]
}

// RandomGenerateToken generates a random string of the given length
func RandomGenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// RandomGenerateNumeric generates a random string of the given length
// with only digits
func RandomGenerateNumeric(length int) string {
	pin := ""
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(9))
		if err != nil {
			log.WithError(err).Errorln("BIG PROBLEMS IN RANDOM STUFF")
			n = big.NewInt(1)
		}
		pin += n.String()
	}
	return pin
}
