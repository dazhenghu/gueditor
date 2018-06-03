package gueditor

import (
    "io"
    "mime/multipart"
    "os"
    "time"
    "strings"
    "strconv"
    "path/filepath"
    "encoding/base64"
    "io/ioutil"
    "net/url"
    "errors"
    "net/http"
)

const (
    BIGGER_THAN_UPLOAD_MAX_FILESIZE string = "文件大小超出 upload_max_filesize 限制"
    BIGGER_THAN_MAX_FILE_SIZE       string = "文件大小超出 MAX_FILE_SIZE 限制"
    FILE_NOT_COMPLETE               string = "文件未被完整上传"
    NO_FILE_UPLOAD                  string = "没有文件被上传"
    UPLOAD_FILE_IS_EMPTY            string = "上传文件为空"
    ERROR_TMP_FILE                  string = "临时文件错误"
    ERROR_TMP_FILE_NOT_FOUND        string = "找不到临时文件"
    ERROR_SIZE_EXCEED               string = "文件大小超出网站限制"
    ERROR_TYPE_NOT_ALLOWED          string = "文件类型不允许"
    ERROR_CREATE_DIR                string = "目录创建失败"
    ERROR_DIR_NOT_WRITEABLE         string = "目录没有写权限"
    ERROR_FILE_MOVE                 string = "文件保存时出错"
    ERROR_FILE_NOT_FOUND            string = "找不到上传文件"
    ERROR_WRITE_CONTENT             string = "写入文件内容错误"
    ERROR_UNKNOWN                   string = "未知错误"
    ERROR_DEAD_LINK                 string = "链接不可用"
    ERROR_HTTP_LINK                 string = "链接不是http链接"
    ERROR_HTTP_CONTENTTYPE          string = "链接contentType不正确"
    INVALID_URL                     string = "非法 URL"
    INVALID_IP                      string = "非法 IP"
    ERROR_BASE64_DATA               string = "base64图片解码错误"
    ERROR_FILE_STATE                string = "文件系统错误"
    ERRPR_READ_REMOTE_DATA          string = "读取远程链接出错"
)

type ResFileInfo struct {
    URL      string `json:"url"`
    Title    string `json:"title"`
    Original string `json:"original"`
    Type     string `json:"type"`
    Size     int64  `json:"size"`
}

type ResFileInfoWithState struct {
    State string `json:"state"`
    ResFileInfo
}

type ResFilesInfoWithStates struct {
    State string `json:"state"`
    List []*ResFileInfo `json:"list"`
}

func (resState *ResFileInfoWithState) fromResFileInfo(res *ResFileInfo)  {
    resState.URL = res.URL
    resState.Title = res.Title
    resState.Original = res.Original
    resState.Type = res.Type
    resState.Size = res.Size
}

// 上传文件的参数
type UploaderParams struct {
    PathFormat string   /* 上传保存路径,可以自定义保存路径和文件名格式 */
    MaxSize    int      /* 上传大小限制，单位B */
    AllowFiles []string /* 上传格式限制 */
    OriName    string   /* 原始文件名 */
}

type UploaderInterface interface {
    UpFile(file multipart.File, handle *multipart.FileHeader) (fileInfo *ResFileInfo, err error) //上传文件的方法
    UpBase64(fileName, base64data string) (fileInfo *ResFileInfo, err error)                     //处理base64编码的图片上传
    SaveRemote(remoteUrl string) (fileInfo *ResFileInfo, err error)                              // 拉取远程图片
    SetParams(params *UploaderParams) error                         // 设置参数信息
    SetRootPath(path string) error // 设置根目录
}

type Uploader struct {
    RootPath string // 项目根目录
    params *UploaderParams
}

/**
新建uploader
 */
func NewUploader(upParams *UploaderParams) *Uploader {
    uploaderObj := &Uploader{}
    uploaderObj.SetParams(upParams)
    return uploaderObj
}

func (up *Uploader) SetParams(params *UploaderParams) (err error)  {
    up.params = params
    return
}

func (up *Uploader) SetRootPath(path string) error {
    up.RootPath = path
    return nil
}

/**
上传文件
 */
func (up *Uploader) UpFile(file multipart.File, fileHeader *multipart.FileHeader) (fileInfo *ResFileInfo, err error)  {
    if file == nil || fileHeader == nil {
        // 上传文件为空
        err = errors.New(UPLOAD_FILE_IS_EMPTY)
        return
    }

    // 校验文件大小
    err = up.checkSize(fileHeader.Size)
    if err != nil {
        return
    }

    // 校验文件类型
    ext := filepath.Ext(fileHeader.Filename)
    err = up.checkType(ext)
    if err != nil {
        return
    }

    fullName := up.getFullName(fileHeader.Filename)
    fileAbsPath := up.getFilePath(fullName)
    fileDir  := filepath.Dir(fileAbsPath)
    exists, err := pathExists(fileDir)
    if err != nil {
        err = errors.New(ERROR_FILE_STATE)
        return
    }

    if !exists {
        // 文件夹不存在，创建
        if err = os.MkdirAll(fileDir, 0666); err != nil {
            err = errors.New(ERROR_CREATE_DIR)
            return
        }
    }

    dstFile, err := os.OpenFile(fileAbsPath, os.O_WRONLY | os.O_RDONLY, 0666)
    if err != nil {
        err = errors.New(ERROR_DIR_NOT_WRITEABLE)
        return
    }
    defer func() {
        dstFile.Close()
    }()

    _, err = io.Copy(dstFile, file)
    if err != nil {
        err = errors.New(ERROR_WRITE_CONTENT)
        return
    }

    fileInfo.Size = fileHeader.Size
    fileInfo.Type = ext
    fileInfo.Title = filepath.Base(fileAbsPath)
    fileInfo.Original = fileHeader.Filename
    fileInfo.URL = fullName
    return
}

/**
删除base64数据文件
 */
func (up *Uploader) UpBase64(fileName, base64data string) (fileInfo *ResFileInfo, err error)  {
    imgData, err := base64.StdEncoding.DecodeString(base64data)
    if err != nil {
        err = errors.New(ERROR_BASE64_DATA)
        return
    }

    fileSize := len(imgData)
    // 校验文件大小
    err = up.checkSize(int64(fileSize))
    if err != nil {
        return
    }

    ext := filepath.Ext(fileName)
    err = up.checkType(ext)
    if err != nil {
        return
    }

    fullName := up.getFullName(fileName)
    fileAbsPath := up.getFilePath(fullName)
    fileDir  := filepath.Dir(fileAbsPath)
    exists, err := pathExists(fileDir)
    if err != nil {
        err = errors.New(ERROR_FILE_STATE)
        return
    }

    if !exists {
        // 文件夹不存在，创建
        if err = os.MkdirAll(fileDir, 0666); err != nil {
            err = errors.New(ERROR_CREATE_DIR)
            return
        }
    }

    err = ioutil.WriteFile(fileAbsPath, imgData, 0666)
    if err != nil {
        err = errors.New(ERROR_WRITE_CONTENT)
        return
    }

    fileInfo.Size = int64(fileSize)
    fileInfo.Type = ext
    fileInfo.Title = filepath.Base(fileAbsPath)
    fileInfo.Original = up.params.OriName
    fileInfo.URL = fullName

    return
}

/**
拉取远程文件并保存
 */
func (up *Uploader) SaveRemote(remoteUrl string) (fileInfo *ResFileInfo, err error) {
    urlObj, err := url.Parse(remoteUrl)
    if err != nil {
        err = errors.New(INVALID_URL)
        return
    }

    scheme := strings.ToLower(urlObj.Scheme)
    if scheme != "http" && scheme != "https" {
        err = errors.New(ERROR_HTTP_LINK)
        return
    }

    // 校验文件类型
    ext := filepath.Ext(urlObj.Path)
    err = up.checkType(ext)
    if err != nil {
        return
    }

    fileName := filepath.Base(urlObj.Path)
    fullName := up.getFullName(fileName)
    fileAbsPath := up.getFilePath(fullName)
    fileDir  := filepath.Dir(fileAbsPath)
    exists, err := pathExists(fileDir)
    if err != nil {
        err = errors.New(ERROR_FILE_STATE)
        return
    }

    if !exists {
        // 文件夹不存在，创建
        if err = os.MkdirAll(fileDir, 0666); err != nil {
            err = errors.New(ERROR_CREATE_DIR)
            return
        }
    }

    client := http.Client{Timeout: 5 * time.Second}
    // 校验是否是可用的链接
    headResp, err := client.Head(remoteUrl)
    if err == nil {
        defer func() {
            headResp.Body.Close()
        }()
    }
    if err != nil || headResp.StatusCode != http.StatusOK {
        err = errors.New(ERROR_DEAD_LINK)
        return
    }
    // 校验content-type
    contentType := headResp.Header.Get("Content-Type")
    if !strings.Contains(strings.ToLower(contentType), "image") {
        err = errors.New(ERROR_HTTP_CONTENTTYPE)
        return
    }

    resp, err := client.Get(remoteUrl)
    if err == nil {
        defer func() {
            resp.Body.Close()
        }()
    }
    if err != nil || resp.StatusCode != http.StatusOK {
        err = errors.New(ERROR_DEAD_LINK)
        return
    }

    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        errors.New(ERRPR_READ_REMOTE_DATA)
        return
    }

    err = ioutil.WriteFile(fileAbsPath, data, 0666)
    if err != nil {
        err = errors.New(ERROR_WRITE_CONTENT)
        return
    }

    fileInfo.Size = int64(len(data))
    fileInfo.Type = ext
    fileInfo.Title = filepath.Base(fileAbsPath)
    fileInfo.Original = fileName
    fileInfo.URL = fullName

    return
}

/**
根据原始文件名生成新文件名
 */
func (up *Uploader) getFullName(oriName string) string {
    timeNow := time.Now()
    timeNowFormat := time.Now().Format("2006_01_02_15_04_05")
    timeArr := strings.Split(timeNowFormat, "_")

    format := up.params.PathFormat

    format = strings.Replace(format, "{yyyy}", timeArr[0], 1)
    format = strings.Replace(format, "{mm}", timeArr[1], 1)
    format = strings.Replace(format, "{dd}", timeArr[2], 1)
    format = strings.Replace(format, "{hh}", timeArr[3], 1)
    format = strings.Replace(format, "{ii}", timeArr[4], 1)
    format = strings.Replace(format, "{ss}", timeArr[5], 1)

    timestamp := strconv.FormatInt(timeNow.UnixNano(), 10)
    format = strings.Replace(format, "{time}", string(timestamp), 1)

    ext := filepath.Ext(oriName)

    return format + ext
}

func (up *Uploader) getFilePath(fullName string) string  {
    return filepath.Join(up.RootPath, fullName)
}

/**
校验文件大小
 */
func (up *Uploader) checkSize(fileSize int64) (err error) {
    if fileSize > int64(up.params.MaxSize) {
        err = errors.New(ERROR_SIZE_EXCEED)
        return
    }
    return
}

/**
校验文件类型
 */
func (up *Uploader) checkType(fileType string) (err error)  {
    isvalid := false
    for _, fileTypeValid := range up.params.AllowFiles {
        if strings.ToLower(fileType) == fileTypeValid {
            isvalid = true
            break
        }
    }

    if !isvalid {
        err = errors.New(ERROR_TYPE_NOT_ALLOWED)
        return
    }

    return
}

