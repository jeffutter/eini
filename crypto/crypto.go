package crypto

import (
	"encoding/hex"
	"errors"
	"fmt"
	ejsonCrypto "github.com/Shopify/ejson/crypto"
	"regexp"
)

type Decrypter = *ejsonCrypto.Decrypter
type Encrypter = *ejsonCrypto.Encrypter

func PrepareDecrypter(pubKey string, privKey string) (Decrypter, error) {
	var priv [32]byte
	var pub [32]byte

	pubkey, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error Decoding Public Key: %s", err))
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
	return myKP.Decrypter(), nil
}

func PrepareEncrypter(pubKey string) (Encrypter, error) {
	var pub [32]byte

	pubkey, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}
	copy(pub[:], pubkey)

	var myKP ejsonCrypto.Keypair
	if err := myKP.Generate(); err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to generate Keypair: %s", err))
	}
	return myKP.Encrypter(pub), nil
}

func Decrypt(decrypter Decrypter, s string) (string, error) {
	encryptedRegex, _ := regexp.Compile("^EJ\\[.*\\]")

	if encryptedRegex.MatchString(s) {
		decrypted, err := decrypter.Decrypt([]byte(s))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s", decrypted), nil
	} else {
		return s, nil
	}
}