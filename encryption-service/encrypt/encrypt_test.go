package encrypt

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestEncryptDecrypt(t *testing.T) {
	key := []byte("kz5h9VmyVjC1yFVDeHnsaoCqIcTfVbpU")
	c := AES256EncryptionClient{key}
	plaintext := "I'm a message"
	enc, err := c.Encrypt(&plaintext)
	assert.Equal(t, nil, err)

	dec, err := c.Decrypt(enc)
	assert.Equal(t, nil, err)
	assert.Equal(t, plaintext, *dec)
}

func TestEncryptDecryptBadKey(t *testing.T) {
	key := []byte("xxxyyy")
	c := AES256EncryptionClient{key}
	plaintext := "I'm a message"

	enc, err := c.Encrypt(&plaintext)
	assert.Equal(t, nil, enc)
	assert.NotEqual(t, nil, err)

	dec, err := c.Decrypt(&plaintext)
	assert.Equal(t, nil, dec)
	assert.NotEqual(t, nil, err)
}
