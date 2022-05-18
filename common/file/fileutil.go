package file

import "os"

//check file is readable;modify addr if addr "=="
func CheckYamlReadable(addr *string) (bool, error) {
	if *addr == "" { //Using dafaultaddr
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		addrtmp := path + "/config/settings.yaml"
		*addr = addrtmp
	}
	_, err := os.ReadFile(*addr)
	if err != nil {
		return false, err
	}
	return true, nil
}

//check file is exist
func CheckPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
