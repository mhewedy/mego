package user

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

var encKey []byte

func init() {
	encKey = make([]byte, 32)
	_, err := rand.Read(encKey)
	if err != nil {
		log.Fatal(err)
	}
}

//https://tutorialedge.net/golang/go-encrypt-decrypt-aes-tutorial/
func encrypt(plain string) (string, error) {

	c, err := aes.NewCipher(encKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	seal := gcm.Seal(nonce, nonce, []byte(plain), nil)

	return string(seal), nil
}

func decrypt(enc string) (string, error) {

	c, err := aes.NewCipher(encKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	encBytes := []byte(enc)

	nonceSize := gcm.NonceSize()
	if len(encBytes) < nonceSize {
		return "", err
	}

	nonce, encBytes := encBytes[:nonceSize], encBytes[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, encBytes, nil)
	if err != nil {
		return "", nil
	}

	return string(plaintext), nil
}
