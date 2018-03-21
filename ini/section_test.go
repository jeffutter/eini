package ini_test

import (
	"github.com/jeffutter/eini/ini"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetsKeys(t *testing.T) {
	cfg, _ := ini.Load([]byte(FILE))

	sections := cfg.GetSections()

	keys := sections[0].GetKeys()
	assert.Len(t, keys, 2)

	keys = sections[1].GetKeys()
	assert.Len(t, keys, 1)
}

func TestName(t *testing.T) {
	cfg, _ := ini.Load([]byte(FILE))

	sections := cfg.GetSections()

	assert.Equal(t, sections[0].Name(), "DEFAULT")
	assert.Equal(t, sections[1].Name(), "section1")
}
