package ini

import (
	goIni "gopkg.in/ini.v1"
)

type Section struct {
	section *goIni.Section
}

func (section Section) GetKeys() []Key {
	ks := section.section.Keys()

	keys := make([]Key, len(ks))
	for i, v := range ks {
		keys[i] = Key{v}
	}
	return keys
}

func (section Section) Name() string {
	return section.section.Name()
}
