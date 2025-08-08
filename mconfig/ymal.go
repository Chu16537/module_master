package mconfig

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// 讀取
func LoadYaml(path string, conf interface{}) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, conf)
	if err != nil {
		return err
	}

	return nil
}
