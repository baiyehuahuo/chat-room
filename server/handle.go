package server

import (
	"net/http"
	"os"
	"path/filepath"
)

var rootDir string

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
	rootDir = infer(cwd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func RegisterHandle() {
	http.HandleFunc("/", nil)
	http.HandleFunc("/ws", nil)
}
