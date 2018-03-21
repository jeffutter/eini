package ini

import (
	"fmt"
	goIni "gopkg.in/ini.v1"
	"io"
)

type File struct {
	cfg    *goIni.File
	source interface{}
}

func Load(source interface{}) (File, error) {
	cfg, err := goIni.Load(source)
	return File{cfg, source}, err
}

func (file File) PubKey() (string, error) {
	pubkey, err := file.cfg.Section("").GetKey("_public_key")
	if err != nil {
		return "", fmt.Errorf("Couldn't read public key from ini - %v: %s", file.source, err)
	}
	return pubkey.Value(), nil
}

func (file File) GetSections() []Section {
	secs := file.cfg.Sections()

	sections := make([]Section, len(secs))
	for i, v := range secs {
		sections[i] = Section{v}
	}
	return sections
}

func (file File) WriteTo(output io.Writer) (int64, error) {
	bytes, err := file.cfg.WriteTo(output)
	return bytes, err
}
