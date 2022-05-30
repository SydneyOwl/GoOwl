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
func GetCwd() string {
	path, err := os.Getwd()
	if err != nil {
		panic("Cannot get currdir!")
	}
	return path
}
func CreateFile(filepath string) error {
	if !isExist(filepath) {
		f, err := os.Create(filepath)
		f.Close()
		return err
	}
	return nil
}
func CreateDirOnPwd() {
	CreateDir(GetCwd())
}

//调⽤os.MkdirAll递归创建⽂件夹
func CreateDir(filePath string) error {
	if !isExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

// 判断所给路径⽂件/⽂件夹是否存在(返回true是存在)
func isExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取⽂件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
