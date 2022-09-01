package file

import "os"

// CheckYamlReadable check if file is readable;
func CheckYamlReadable(addr *string) (bool, error) {
	_, err := os.ReadFile(*addr)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CheckPathExists check if file exists.
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

// GetCwd returns current path.
func GetCwd() string {
	path, err := os.Getwd()
	if err != nil {
		panic("Cannot get current directory! GoOwl Stop.")
	}
	return path
}

// CreateFile creates file on specified path.
func CreateFile(filepath string) error {
	if !isExist(filepath) {
		f, err := os.Create(filepath)
		f.Close()
		return err
	}
	return nil
}

// CreateDir create dir recursively.
func CreateDir(filePath string) error {
	if !isExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

// isExist check if pathexists.
func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// CalcSize calculates size of specified file/dir. Returns byte.
func CalcSize(path string) (int64, error) {
	if file, err := os.Stat(path); err != nil {
		return 0, err
	} else {
		return file.Size(), nil
	}
}
