package storage

import (
    "io"
)

const (
    SUCCESS = "success"
    FAIL    = "fail"
)

type BaseInterface interface {

    /**
    从本地路径读取文件并保存
    srcPath 本地文件路径
    dstAbsPath 目的位置绝对路径
    dstRelativePath 目的位置相对路径
     */
    SaveFileFromLocalPath(srcPath string, dstAbsPath, dstRelativePath string) (url string, err error)

    /**
    将本地文件保存到目标位置
    srcFile 源文件读取接口
    dstAbsPath 目的位置绝对路径
    dstRelativePath 目的位置相对路径
     */
    SaveFile(srcFile io.Reader, srcFileSize int64, dstAbsPath, dstRelativePath string) (url string, err error)

    /**
    将数据保存到目标位置
    data 数据
    dstAbsPath 目的位置绝对路径
    dstRelativePath 目的位置相对路径
     */
    SaveData(data []byte, dstAbsPath, dstRelativePath string) (url string, err error)
}
