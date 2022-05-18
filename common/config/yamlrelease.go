package config

import "io/ioutil"

var (
	yaml = "test"
)

//Release yaml if no yaml available.
func ReleaseYaml(path string) error {
	confByte := []byte(yaml)
	return ioutil.WriteFile(path, confByte, 0666)
}
