package global

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
)

/*
	用于推断当前项目根路径
*/

var (
	RootDir string // 项目根路径
	once    = new(sync.Once)
)

func exist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || errors.Is(err, os.ErrExist)
}

// 计算项目路径
func inferRootDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var infer func(path string) string
	infer = func(path string) string {
		if exist(path + "/config") {
			return path
		}
		return infer(filepath.Dir(path))
	}
	return infer(cwd)
}

func init() {
	once.Do(func() {
		RootDir = inferRootDir()
	})
}
