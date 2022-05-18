package config

import (
	"errors"

	"github.com/spf13/viper"
)

var (
	ConfigFormatNotSupportedError = errors.New("unknown setting format")
	ConfigContentError            = errors.New("yaml Error-Check your yaml content")
)

//Returns raw viper stu-could read directly without sp.to struct
func LoadConfigFromYaml(configFile string) (*viper.Viper, error) {
	yamlConfig := viper.New()
	//设置配置文件的名字
	yamlConfig.SetConfigFile(configFile)
	if err := yamlConfig.ReadInConfig(); err != nil {
		return nil, err
	}
	return yamlConfig, nil
}

//This give explaination of error and a bool of whether error exists.
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

//Check if element is in slice.
func CheckInSlice(compareTo []string, from string) bool {
	for _, v := range compareTo {
		if v == from {
			return true
		}
	}
	return false
}
