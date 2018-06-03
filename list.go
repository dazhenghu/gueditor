package gueditor

import (
    "os"
    "path/filepath"
    "strings"
)

type ListFileItem struct {
    State string      `json:"state"`
    List  []*FileItem `json:"list"`
    Start int         `json:"start"`
    Total int         `json:"total"`
}

type FileItem struct {
    Url string `json:"url"` // 文件url
    Mtime int64 `json:"mtime"` // 最后编辑时间
}

type ListInterface interface {
    GetFileList(start int, size int) (fileList []*FileItem, err error) // 获取文件夹下的文件列表
    SetParams(params *ListParams) error // 设置参数
    SetRootPath(path string) error // 设置根目录
}

type ListParams struct {
    AllowFiles []string // 允许的文件类型
    ListSize int // 列表分页大小
    Path string  // 资源路径
}

type List struct {
    RootPath string // 项目根目录
    Params *ListParams // 参数
}

func (l *List) SetParams(params *ListParams) error {
    l.Params = params
    return nil
}

func (l *List) SetRootPath(path string) error {
    l.RootPath = path
    return nil
}

/**
获取资源列表
 */
func (l *List) GetFileList(start int, size int) (fileList []*FileItem, err error) {

    files := make([]*FileItem, 0)

    path := filepath.Join(l.RootPath, l.Params.Path)

    l.getFiles(path, l.Params.AllowFiles, &files)

    if size == 0 {
        size = l.Params.ListSize
    }

    end := start + size

    i := end
    listLen := len(files)
    if i > listLen {
        i = listLen
    }

    fileList = make([]*FileItem, 0)
    for i := i - 1; i < listLen && i >=0 && i >= start; i--  {
        fileList = append(fileList, files[i])
    }

    return
}

/**
递归获取文件列表
 */
func (l *List) getFiles(path string, allowFiles []string, files *[]*FileItem) (err error)  {
    fileInfo, err := os.Stat(path)
    if err == nil {
        if !fileInfo.IsDir() {
            return
        }

        filepath.Walk(path, func(fileName string, info os.FileInfo, err error) error {
            if fileName != "." && fileName != ".." && fileName != "walk:.DS_Store" && fileName != path {
                path2 := fileName
                if info.IsDir() {
                    l.getFiles(path2, allowFiles, files)
                } else {
                    ext := filepath.Ext(fileName)
                    for _, allowItem := range allowFiles {
                        urlPath := path2[len(l.RootPath):]
                        if strings.ToLower(ext) == allowItem {
                            *files = append(*files, &FileItem{
                                Url: urlPath,
                                Mtime: info.ModTime().Unix(),
                            })
                            break
                        }
                    }
                }
            }

            return nil
        })
        return
    }

    return
}