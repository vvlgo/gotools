package yamlconf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func Read(yamlPath string, v interface{}) error {
	// Read config file
	buf, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		panic(err)
		return err
	}
	err = yaml.Unmarshal(buf, v)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}
