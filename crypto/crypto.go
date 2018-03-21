package crypto

import (
	"encoding/hex"
	"fmt"
	ejsonCrypto "github.com/Shopify/ejson/crypto"
	"regexp"
)

type Decrypter struct {
	decrypter *ejsonCrypto.Decrypter
}
type Encrypter struct {
	encrypter *ejsonCrypto.Encrypter
}

func GenerateKeypair() (string, string, error) {
	var kp ejsonCrypto.Keypair
	if err := kp.Generate(); err != nil {
		return "", "", err
	}
	return kp.PublicString(), kp.PrivateString(), nil
}

func PrepareDecrypter(pubKey string, privKey string) (Decrypter, error) {
	var priv [32]byte
	var pub [32]byte

	pubkey, err := hex.DecodeString(pubKey)
	if err != nil {
		return Decrypter{}, fmt.Errorf("Error Decoding Public Key: %s", err)
	}

	privkey, _ := hex.DecodeString(privKey)
	// For some reason this returns an error but it still works
	// if err != nil {
	// 	return nil, errors.New(fmt.Sprintf("Error Decoding Private Key: %s", err))
	// }

	copy(pub[:], pubkey)
	copy(priv[:], privkey)

	myKP := ejsonCrypto.Keypair{
		Public:  pub,
		Private: priv,
	}
	return Decrypter{myKP.Decrypter()}, nil
}

func (encrypter Encrypter) Encrypt(s string) (string, error) {
	encrypted, err := encrypter.encrypter.Encrypt([]byte(s))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", encrypted), nil
}

func PrepareEncrypter(pubKey string) (Encrypter, error) {
	var pub [32]byte

	pubkey, err := hex.DecodeString(pubKey)
	if err != nil {
		return Encrypter{}, err
	}
	copy(pub[:], pubkey)

	var myKP ejsonCrypto.Keypair
	if err := myKP.Generate(); err != nil {
		return Encrypter{}, fmt.Errorf("Failed to generate Keypair: %s", err)
	}
	return Encrypter{myKP.Encrypter(pub)}, nil
}

func (decrypter Decrypter) Decrypt(s string) (string, error) {
	encryptedRegex, _ := regexp.Compile("^EJ\\[.*\\]")

	if encryptedRegex.MatchString(s) {
		decrypted, err := decrypter.decrypter.Decrypt([]byte(s))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s", decrypted), nil
	} else {
		return s, nil
	}
}
