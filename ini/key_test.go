package ini_test

import (
	"github.com/jeffutter/eini/ini"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeyName(t *testing.T) {
	cfg, _ := ini.Load([]byte(FILE))
	sections := cfg.GetSections()
	assert.Equal(t, sections[0].GetKeys()[0].Name(), "_public_key")
}

func TestKeyValue(t *testing.T) {
	cfg, _ := ini.Load([]byte(FILE))
	sections := cfg.GetSections()
	assert.Equal(t, sections[0].GetKeys()[0].Value(), "12345")
}

func TestGetComment(t *testing.T) {
	cfg, _ := ini.Load([]byte(FILE))
	sections := cfg.GetSections()
	assert.Equal(t, sections[0].GetKeys()[1].Comment(), "# I am a comment")
}
