package ini_test

import (
	"bytes"
	"fmt"
	"github.com/jeffutter/eini/crypto"
	"github.com/jeffutter/eini/ini"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
)

func makeEini(pubkey string) []byte {
	contents := fmt.Sprintf(`
_public_key = %s
# I am a comment
foo         = bar

[section1]
baz = bang

[section2]
#decrypted
decrypted = decrypted
`, pubkey)
	return []byte(contents)
}

func TestEncrypt(t *testing.T) {
	pub, _, err := crypto.GenerateKeypair()
	assert.NoError(t, err)

	file := makeEini(pub)

	cfg, err := ini.Load(file)
	assert.NoError(t, err)

	var buffer bytes.Buffer
	err = cfg.Encrypt(pub, &buffer)
	assert.NoError(t, err)

	output := buffer.String()

	assert.Regexp(t, regexp.MustCompile(fmt.Sprintf("_public_key.*=.*%s", pub)), output)
	assert.NotRegexp(t, regexp.MustCompile("foo.*=.*bar"), output)
	assert.Regexp(t, regexp.MustCompile("foo.*=\\s+\\w+"), output)
	assert.Regexp(t, regexp.MustCompile("decrypted.*=.*decrypted"), output)
}

func TestDecrypt(t *testing.T) {
	pub, priv, err := crypto.GenerateKeypair()
	assert.NoError(t, err)

	file := makeEini(pub)

	cfg, err := ini.Load(file)
	assert.NoError(t, err)

	var buffer bytes.Buffer
	err = cfg.Encrypt(pub, &buffer)
	assert.NoError(t, err)

	encryptedOutput := buffer.String()

	cfg, err = ini.Load([]byte(encryptedOutput))
	assert.NoError(t, err)

	lines, err := cfg.Decrypt(pub, priv)
	assert.NoError(t, err)

	output := strings.Join(lines, "\n")

	assert.NotRegexp(t, regexp.MustCompile(pub), output)
	assert.Regexp(t, regexp.MustCompile("foo.*=.*bar"), output)
	assert.Regexp(t, regexp.MustCompile("decrypted.*=.*decrypted"), output)
}
