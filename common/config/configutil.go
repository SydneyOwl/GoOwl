package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"github.com/sydneyowl/GoOwl/common/file"
)

// LoadConfigFromYaml Returns raw viper object that could be read directly.
func LoadConfigFromYaml(configFile string) (*viper.Viper, error) {
	yamlConfig := viper.New()
	yamlConfig.SetConfigFile(configFile)
	if err := yamlConfig.ReadInConfig(); err != nil {
		return nil, err
	}
	return yamlConfig, nil
}

// CheckViperErr give explanation of error and a bool of whether error exists.
func CheckViperErr(err error) error {
	if err != nil {
		if _, ok := err.(viper.UnsupportedConfigError); ok {
			return errors.New("unknown setting format")
			// ConfigContentError            =
		} else {
			return errors.New("yaml Error-Check your yaml content")
		}
	}
	return nil
}

// CheckInSlice check if element is in slice.
func CheckInSlice(compareTo []string, from string) bool {
	for _, v := range compareTo {
		if v == from {
			return true
		}
	}
	return false
}

// InitConfig loads config from yaml.
func InitConfig(yamlAddr *string) {
	if readable, err := file.CheckYamlReadable(yamlAddr); !readable {
		fmt.Println(err.Error())
		return
	}
	rawConfig, err := LoadConfigFromYaml(*yamlAddr) //returns raw viper obj
	if err := CheckViperErr(err); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := rawConfig.Unmarshal(YamlConfig); err != nil {
		fmt.Println("Unknown Error occurred!", "GoOwl-MainLog")
		return
	}
}
