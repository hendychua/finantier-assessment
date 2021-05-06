package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

type IEncryptionClient interface {
	Encrypt(plaintext *string) (*string, error)
	Decrypt(encrypted *string) (*string, error)
}

type AES256EncryptionClient struct {
	Key []byte
}

func (c AES256EncryptionClient) Encrypt(plaintext *string) (*string, error) {
	plaintextBytes := []byte(*plaintext)

	aesGCM, err := c.getAESGCM()
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	//Add the nonce as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := hex.EncodeToString(aesGCM.Seal(nonce, nonce, plaintextBytes, nil))
	return &ciphertext, nil
}

func (c AES256EncryptionClient) Decrypt(encrypted *string) (*string, error) {
	encryptedBytes, err := hex.DecodeString(*encrypted)
	if err != nil {
		return nil, err
	}

	aesGCM, err := c.getAESGCM()
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := encryptedBytes[:nonceSize], encryptedBytes[nonceSize:]

	plaintextBytes, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	plaintext := string(plaintextBytes)

	return &plaintext, nil
}

func (c AES256EncryptionClient) getAESGCM() (cipher.AEAD, error) {
	block, err := aes.NewCipher(c.Key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}
