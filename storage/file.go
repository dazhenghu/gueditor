package storage

import (
    "os"
    "io"
    "errors"
    "github.com/dazhenghu/gueditor/common"
    "path/filepath"
    "io/ioutil"
    "github.com/dazhenghu/util/fileutil"
)

type LocalFile struct {
    BaseInterface
}

func (lf *LocalFile) SaveFileFromLocalPath(srcPath string, dstAbsPath, dstRelativePath string) (url string, err error) {
    exists, err := fileutil.PathExists(srcPath)
    if err != nil {
        err = errors.New(common.ERROR_FILE_STATE)
        return
    }
    if !exists {
        err = errors.New(common.ERROR_FILE_NOT_FOUND)
        return
    }

    return
}

/**
保存文件到本地
 */
func (lf *LocalFile) SaveFile(srcFile io.Reader, srcFileSize int64, dstAbsPath, dstRelativePath string) (url string, err error) {

    fileDir  := filepath.Dir(dstAbsPath)
    exists, err := pathExists(fileDir)
    if err != nil {
        err = errors.New(common.ERROR_FILE_STATE)
        return
    }

    if !exists {
        // 文件夹不存在，创建
        if err = os.MkdirAll(fileDir, 0766); err != nil {
            err = errors.New(common.ERROR_CREATE_DIR)
            return
        }
    }

    dstFile, err := os.OpenFile(dstAbsPath, os.O_WRONLY | os.O_RDONLY | os.O_CREATE, 0666)
    if err != nil {
        err = errors.New(common.ERROR_DIR_NOT_WRITEABLE)
        return
    }
    defer func() {
        dstFile.Close()
    }()
    _, err = io.Copy(dstFile, srcFile)
    if err != nil {
        err = errors.New(common.ERROR_WRITE_CONTENT)
        return
    }

    url = dstRelativePath
    return
}

/**
保存数据到本地
 */
func (lf *LocalFile) SaveData(data []byte, dstAbsPath, dstRelativePath string) (url string, err error)  {
    fileDir  := filepath.Dir(dstAbsPath)
    exists, err := pathExists(fileDir)
    if err != nil {
        err = errors.New(common.ERROR_FILE_STATE)
        return
    }

    if !exists {
        // 文件夹不存在，创建
        if err = os.MkdirAll(fileDir, 0766); err != nil {
            err = errors.New(common.ERROR_CREATE_DIR)
            return
        }
    }

    err = ioutil.WriteFile(dstAbsPath, data, 0666)
    if err != nil {
        err = errors.New(common.ERROR_WRITE_CONTENT)
        return
    }
    url = dstRelativePath
    return
}

/**
判断对应路径是否存在
 */
func pathExists(path string) (bool, error)  {
    // 获取path的信息
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }

    if os.IsNotExist(err) {
        return false, nil
    }

    return false, err
}
