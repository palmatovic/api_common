package api_common

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

// InitServiceConfiguration initializes configuration reading by config
// file with filepath set in ENV_CONFIGFILEPATH env variable
func InitServiceConfiguration(envConfigFilePath string) (MicroserviceConfiguration, error) {
	log.Traceln("initializing service configuration")
	var ServiceConfiguration MicroserviceConfiguration
	configFilepath := os.Getenv(envConfigFilePath)
	if len(configFilepath) == 0 {
		return MicroserviceConfiguration{}, fmt.Errorf("missing value for environment variable %s", envConfigFilePath)
	}
	log.Tracef("reading configuration file %s", configFilepath)
	configFile, errReadFile := os.ReadFile(configFilepath)
	if errReadFile != nil {
		return MicroserviceConfiguration{}, fmt.Errorf("cannot initialize configuration: file %s not found", configFilepath)
	}
	log.Tracef("unmarshaling configuration file %s", configFilepath)
	errUnmarshal := yaml.Unmarshal(configFile, &ServiceConfiguration)
	if errUnmarshal != nil {
		return MicroserviceConfiguration{}, fmt.Errorf("cannot initialize configuration, check the syntax of the file")
	}
	log.Traceln("service configuration initialized")
	return ServiceConfiguration, nil
}

// GetSecretString returns the content of the secret file
// into a string and trims the content
func GetSecretString(filepath string) (string, error) {
	log.Tracef("getting secret %s", filepath)
	secret, errReadFile := os.ReadFile(filepath)
	if errReadFile != nil {
		log.WithField("error", errReadFile.Error()).Errorf("cannot get secret %s", filepath)
		return "", fmt.Errorf("cannot get secret %s", filepath)
	}
	secretString := string(secret)
	secretString = strings.TrimSpace(secretString)
	return secretString, nil
}

// GetSecretPrivateKey returns the RSA private key
// available on filepath
func GetSecretPrivateKey(filepath string) (*rsa.PrivateKey, error) {
	log.Tracef("getting secret private key %s", filepath)
	secret, errReadFile := os.ReadFile(filepath)
	if errReadFile != nil {
		log.WithField("error", errReadFile.Error()).Errorf("cannot get secret private key %s", filepath)
		return nil, fmt.Errorf("cannot get secret private key %s", filepath)
	}
	pemBlock, _ := pem.Decode(secret)
	if pemBlock == nil {
		log.Errorf("cannot parse pem block secret %s", filepath)
		return nil, fmt.Errorf("cannot parse pem block secret %s", filepath)
	}
	key, errParse := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if errParse != nil {
		log.WithField("error", errParse.Error()).Errorf("cannot parse private key %s", filepath)
		return nil, fmt.Errorf("cannot parse private key %s", filepath)
	}
	return key, nil
}

// GetSecretPublicKey returns the RSA public key
// available on filepath
func GetSecretPublicKey(filepath string) (*rsa.PublicKey, error) {
	log.Tracef("getting secret public key %s", filepath)
	secret, errReadFile := os.ReadFile(filepath)
	if errReadFile != nil {
		log.WithField("error", errReadFile.Error()).Errorf("cannot get public private key %s", filepath)
		return nil, fmt.Errorf("cannot get public private key %s", filepath)
	}
	pemBlock, _ := pem.Decode(secret)
	if pemBlock == nil {
		log.Errorf("cannot parse pem block secret %s", filepath)
		return nil, fmt.Errorf("cannot parse pem block secret %s", filepath)
	}
	key, errParse := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if errParse != nil {
		log.WithField("error", errParse.Error()).Errorf("cannot parse public key %s", filepath)
		return nil, fmt.Errorf("cannot parse public key %s", filepath)
	}
	switch key := key.(type) {
	case *rsa.PublicKey:
		return key, nil
	default:
		break // fall through
	}
	return nil, fmt.Errorf("cannot parse public key %s: key type is not RSA", filepath)
}

// GetStringArrayFromFile opens a file and reads line by line
// returning the string array built with those lines or an error
func GetStringArrayFromFile(filepath string) ([]string, error) {
	var array []string
	array = []string{}
	file, errOpen := os.Open(filepath)
	if errOpen != nil {
		return []string{}, fmt.Errorf("cannot open file %s: %s", filepath, errOpen.Error())
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		array = append(array, fileScanner.Text())
	}
	errScan := fileScanner.Err()
	if errScan != nil {
		return []string{}, fmt.Errorf("cannot scan file %s: %s", filepath, errScan.Error())
	}
	return array, nil
}
