package config

import (
	"errors"

	"github.com/spf13/viper"
	"github.com/sydneyowl/GoOwl/common/file"
	"github.com/sydneyowl/GoOwl/common/logger"
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
func InitConfig(yamlAddr *string){
	if readable, err := file.CheckYamlReadable(yamlAddr); !readable {
		logger.Fatal(err.Error(),"GoOwl-MainLog")
		return
	}
	rawConfig, err := LoadConfigFromYaml(*yamlAddr) //returns raw viper obj
	if err := CheckViperErr(err); err != nil {
		logger.Error(err.Error(),"GoOwl-MainLog")
		return
	}
	if err := rawConfig.Unmarshal(YamlConfig); err != nil {
		logger.Warning("Unknown Error occurred!","GoOwl-MainLog")
		return
	}
}