package encryptor

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestWeakEncryptor_EncryptAndDecryptWhenSameKey(t *testing.T) {
	t.Parallel()

	key := []byte("some_password")
	startBytes := []byte("some_bytes")

	e := &weakEncryptor{}
	encrypted := e.Encrypt(startBytes, key)

	assert.Equal(t, startBytes, e.Decrypt(encrypted, key))
}

func TestWeakEncryptor_EncryptAndDecryptWhenDifferentKeys(t *testing.T) {
	t.Parallel()

	startBytes := []byte("some_bytes")

	e := &weakEncryptor{}
	encrypted := e.Encrypt(startBytes, []byte("some_password"))

	assert.NotEqual(t, startBytes, e.Decrypt(encrypted, []byte("another_password")))
}
