package config

import (
	"errors"

	"github.com/spf13/viper"
)

var (
	ConfigFormatNotSupportedError = errors.New("unknown setting format")
	ConfigContentError            = errors.New("yaml Error-Check your yaml content")
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
			return ConfigFormatNotSupportedError
		} else {
			return ConfigContentError
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
