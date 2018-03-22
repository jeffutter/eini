package ini

import (
	"fmt"
	"github.com/jeffutter/eini/crypto"
	"io"
	"regexp"
	"strings"
)

var ignoreKeyRegex = regexp.MustCompile("^_.*")
var encryptedRegex = regexp.MustCompile("^EJ\\[.*\\]")
var decryptedRegex = regexp.MustCompile("(?i)decrypted")
var envTemplate = "if [ -z ${%s+x} ]; then\n    declare -x \"%s\"=\"%s\"\nfi\n"

func (cfg File) Encrypt(pubkey string, output io.Writer) error {
	encrypter, err := crypto.PrepareEncrypter(pubkey)
	if err != nil {
		return fmt.Errorf("Error setting up crypto: %s", err)
	}

	for _, sec := range cfg.GetSections() {
		for _, key := range sec.GetKeys() {
			if shouldEncrypt(key) {
				encrypted, err := encrypter.Encrypt(key.Value())
				if err != nil {
					return fmt.Errorf("Failed encrypting key %s: %s", key.Name(), err)
				}

				key.SetValue(encrypted)
			}
		}
	}

	_, err = cfg.WriteTo(output)
	if err != nil {
		return fmt.Errorf("Failed to write config: %s", err)

	}
	return nil
}

func (cfg File) Decrypt(pubkey string, privkey string) ([]string, error) {
	var lines []string

	decrypter, err := crypto.PrepareDecrypter(pubkey, privkey)
	if err != nil {
		return lines, fmt.Errorf("Error setting up crypto: %s", err)
	}

	for _, sec := range cfg.GetSections() {
		for _, key := range sec.GetKeys() {
			if !ignoreKeyRegex.MatchString(key.Name()) {
				decrypted, err := decrypter.Decrypt(key.Value())
				if err != nil {
					return lines, fmt.Errorf("Failed decrypting key %s: %s", sec.Name(), err)
				}

				var keyName string
				if sec.Name() == "DEFAULT" {
					keyName = key.Name()
				} else {
					keyName = fmt.Sprintf("%s_%s", strings.ToUpper(sec.Name()), key.Name())
				}
				lines = append(lines, fmt.Sprintf(envTemplate, keyName, keyName, decrypted))
			}
		}
	}
	return lines, nil
}

func shouldEncrypt(key Key) bool {
	return !ignoreKeyRegex.MatchString(key.Name()) && !decryptedRegex.MatchString(key.Comment()) && !encryptedRegex.MatchString(key.Value())
}
