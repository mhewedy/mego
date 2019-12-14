package user

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_encrypt_decrypt(t *testing.T) {

	enc, err := encrypt("hello world")
	if err != nil {
		log.Fatal(err)
	}

	plain, err := decrypt(enc)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, "hello world", plain)
}

func Test_encrypt_decrypt_negative(t *testing.T) {

	enc, err := encrypt("hello world123")
	if err != nil {
		log.Fatal(err)
	}

	plain, err := decrypt(enc)
	if err != nil {
		log.Fatal(err)
	}

	assert.NotEqual(t, "hello world", plain)
}
