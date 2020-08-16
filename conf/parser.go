package conf

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// ConfigParser  parser yaml config file into config struct
func ConfParser(file string, in interface{}) error {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		msg := fmt.Sprintf("failed to read '%s' with err: %s", file, err.Error())
		return errors.New(msg)
	}
	err = yaml.UnmarshalStrict(yamlFile, in)
	if err != nil {
		msg := fmt.Sprintf("failed to unmarshal '%s' with err: %s", file, err.Error())
		return errors.New(msg)
	}
	return nil
}
