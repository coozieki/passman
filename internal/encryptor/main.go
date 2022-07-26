package encryptor

import (
	"encoding/base64"
	"log"
	"passman/internal/interfaces"
)

type weakEncryptor struct {
}

func (e *weakEncryptor) Encrypt(bytes []byte, key []byte) []byte {
	var keyIterator int

	bts := make([]byte, len(bytes))
	copy(bts, bytes)

	for i, b := range bts {
		encryptedByte := int(b) + int(key[keyIterator])
		if encryptedByte > 255 {
			encryptedByte -= 255
		}

		bts[i] = byte(encryptedByte)

		keyIterator++
		if keyIterator == len(key) {
			keyIterator = 0
		}
	}

	res := make([]byte, base64.StdEncoding.EncodedLen(len(bts)))
	base64.StdEncoding.Encode(res, bts)
	
	return res
}

func (e *weakEncryptor) Decrypt(bytes []byte, key []byte) []byte {
	bts := make([]byte, base64.StdEncoding.DecodedLen(len(bytes)))

	n, err := base64.StdEncoding.Decode(bts, bytes)
	if err != nil {
		log.Fatal("error while decoding base64: ", err)
	}

	bts = bts[:n]

	var keyIterator int

	for i, b := range bts {
		encryptedByte := int(b) - int(key[keyIterator])
		if encryptedByte < 0 {
			encryptedByte += 255
		}

		bts[i] = byte(encryptedByte)

		keyIterator++
		if keyIterator == len(key) {
			keyIterator = 0
		}
	}

	return bts
}

func NewWeakEncryptor() interfaces.Encryptor {
	return &weakEncryptor{}
}
