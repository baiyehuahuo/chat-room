package global

import (
	"os"
	"path/filepath"
)

func init() {
	inferRootDir()
	initConfig()
}

var RootDir string

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func inferRootDir() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var infer func(d string) string
	infer = func(d string) string {
		// 递归在路径下寻找包含 template 的父目录
		if exists(d + "/template") {
			return d
		}
		return infer(filepath.Dir(d))
	}
	RootDir = infer(cwd)
}
