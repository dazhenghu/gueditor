package gueditor

import (
    "runtime"
    "errors"
    "path/filepath"
)

var currDir string

/**
获取当前组件的绝对路径
 */
func getCurrAbsDir() (path string, err error) {

    if currDir == "" {
        _, file, _, ok := runtime.Caller(0)
        if !ok {
            err = errors.New("runtime get caller err")
            return
        }

        currDir, err = filepath.Abs(filepath.Dir(file))
        if err != nil {
            return
        }

        path = currDir
        return
    }

    path = currDir

    return
}

func getDefaultConfigFile() (cnfFilePath string, err error)  {
    dirPath, err := getCurrAbsDir()
    if err != nil {
        return
    }

    cnfFilePath = filepath.Join(dirPath, "config.json")
    return
}