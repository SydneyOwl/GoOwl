package config

import "io/ioutil"

var (
	yaml = "test"
)

// Releaseyaml if no yaml is available; Deprecated in v0.1.3+.
func ReleaseYaml(path string) error {
	confByte := []byte(yaml)
	return ioutil.WriteFile(path, confByte, 0666)
}
