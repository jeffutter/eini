package crypto_test

import (
	"github.com/jeffutter/eini/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateKeypair(t *testing.T) {
	pubkey, privkey, err := crypto.GenerateKeypair()
	assert.NoError(t, err)

	assert.Equal(t, len(pubkey), 64)
	assert.Equal(t, len(privkey), 64)
}

func TestPrepareDecrypter(t *testing.T) {
	pubkey, privkey, err := crypto.GenerateKeypair()
	assert.NoError(t, err)

	decrypter, err := crypto.PrepareDecrypter(pubkey, privkey)
	assert.NoError(t, err)

	assert.IsType(t, crypto.Decrypter{}, decrypter)
}

func TestPrepareEncrypter(t *testing.T) {
	pubkey, _, err := crypto.GenerateKeypair()
	assert.NoError(t, err)

	encrypter, err := crypto.PrepareEncrypter(pubkey)
	assert.NoError(t, err)

	assert.IsType(t, crypto.Encrypter{}, encrypter)
}

func TestEncryptDecryptCycle(t *testing.T) {
	pubkey, privkey, err := crypto.GenerateKeypair()
	assert.NoError(t, err)

	encrypter, err := crypto.PrepareEncrypter(pubkey)
	assert.NoError(t, err)
	decrypter, err := crypto.PrepareDecrypter(pubkey, privkey)
	assert.NoError(t, err)

	message := "Encrypt Me"

	encrypted, err := encrypter.Encrypt(message)
	assert.NoError(t, err)
	decrypted, err := decrypter.Decrypt(encrypted)
	assert.NoError(t, err)

	assert.Equal(t, message, decrypted)
}
