package ini_test

import (
	"bytes"
	"github.com/jeffutter/eini/ini"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	FILE = `_public_key = 12345
# I am a comment
foo         = bar

[section1]
baz = bang

[section2]

`
)

func TestLoadIniFromString(t *testing.T) {
	_, err := ini.Load([]byte(FILE))

	assert.NoError(t, err)
}

func TestReadsThePubKey(t *testing.T) {
	cfg, err := ini.Load([]byte(FILE))
	pubkey, err := cfg.PubKey()

	assert.NoError(t, err)
	assert.Equal(t, pubkey, "12345")
}

func TestListsSections(t *testing.T) {
	cfg, _ := ini.Load([]byte(FILE))

	sections := cfg.GetSections()
	assert.Len(t, sections, 3)
}

func TestWritesBackToAnIni(t *testing.T) {
	cfg, _ := ini.Load([]byte(FILE))

	buf := new(bytes.Buffer)

	_, err := cfg.WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, buf.String(), FILE)
}
