package ini

import (
	goIni "gopkg.in/ini.v1"
)

type Key struct {
	key *goIni.Key
}

func (key Key) Name() string {
	return key.key.Name()
}

func (key Key) Value() string {
	return key.key.Value()
}

func (key Key) Comment() string {
	return key.key.Comment
}

func (key Key) SetValue(value string) {
	key.key.SetValue(value)
}
