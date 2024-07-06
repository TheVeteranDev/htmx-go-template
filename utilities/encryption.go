package utilities

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func Encrypt(s string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	return string(bytes), err
}

func EncodeBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func DecodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func EncryptRSA(s string) ([]byte, error) {
	pubKey, err := loadRSAPublicKey()
	if err != nil {
		return nil, err
	}
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, []byte(s), nil)
}

func DecryptRSA(bytes []byte) ([]byte, error) {
	privKey, err := loadRSAPrivateKey()
	if err != nil {
		return nil, err
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, bytes, nil)
}

func CompareHash(hash string, s string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
}

func CreateHMAC(s string) string {
	h := hmac.New(sha256.New, []byte(GetEnv("HMAC_KEY", "test-secret")))
	h.Write([]byte(s))
	return EncodeBase64(h.Sum(nil))
}

func loadRSAPublicKey() (*rsa.PublicKey, error) {
	data, err := os.ReadFile(GetEnv("RSA_PUBLIC_KEY_PATH", "/etc/ssl/certs/public_test_key.pem"))
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("Failed to decode RSA PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

func loadRSAPrivateKey() (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(GetEnv("RSA_PRIVATE_KEY_PATH", "/etc/ssl/certs/private_test_key.pem"))
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("Failed to decode RSA PEM block containing private key")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv.(*rsa.PrivateKey), nil
}
