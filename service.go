package gueditor

import (
    "net/http"
    "encoding/json"
    "errors"
)

type service struct {
    uploader UploaderInterface
}

func NewService(uploaderObj UploaderInterface, listObj ListInterface, rootPath string, configFilePath string) (serv *service, err error) {
    if uploaderObj == nil {
        // 没有注入uploader接口，则使用默认的方法
        uploaderObj = &Uploader{}
    }
    uploaderObj.SetRootPath(rootPath)

    if listObj == nil {
        listObj = &List{}
    }
    listObj.SetRootPath(rootPath)

    serv = &service{
        uploader:uploaderObj,
    }

    if configFilePath != "" {
        // 加载默认配置
        err = loadConfigFromFile(configFilePath)
    }

    return
}

/**
上传图片
 */
func (serv *service) Uploadimage(r *http.Request) (res *ResFileInfoWithState, err error) {
    params := &UploaderParams{
        PathFormat:GloabConfig.ImagePathFormat,
        MaxSize:GloabConfig.ImageMaxSize,
        AllowFiles:GloabConfig.ImageAllowFiles,
    }

    fieldName := GloabConfig.ImageFieldName
    res, err =  serv.uploadFile(r, fieldName, params)
    return
}

/**
上传涂鸦
 */
func (serv *service) UploadScrawl(r *http.Request) (res *ResFileInfoWithState, err error)  {
    params := &UploaderParams{
        PathFormat:GloabConfig.ScrawlPathFormat,
        MaxSize:GloabConfig.ScrawlMaxSize,
        AllowFiles:GloabConfig.ImageAllowFiles,
        OriName:"scrawl.png",
    }

    fieldName := GloabConfig.ScrawlFieldName
    data := r.PostFormValue(fieldName)
    serv.uploader.SetParams(params)

    fileInfo, err := serv.uploader.UpBase64(params.OriName, data)
    if err == nil {
        res.fromResFileInfo(fileInfo)
        res.State = SUCCESS
        return
    }
    res.State = err.Error()
    return
}

/**
上传视频
 */
func (serv *service) UploadVideo(r *http.Request) (res *ResFileInfoWithState, err error) {
    params := &UploaderParams{
        PathFormat:GloabConfig.VideoPathFormat,
        MaxSize:GloabConfig.VideoMaxSize,
        AllowFiles:GloabConfig.VideoAllowFiles,
    }

    fieldName := GloabConfig.VideoFieldName
    res, err = serv.uploadFile(r, fieldName, params)
    return
}

/**
上传文件
 */
func (serv *service) UploadFile(r *http.Request) (res *ResFileInfoWithState, err error) {
    params := &UploaderParams{
        PathFormat:GloabConfig.FilePathFormat,
        MaxSize:GloabConfig.FileMaxSize,
        AllowFiles:GloabConfig.FileAllowFiles,
    }

    fieldName := GloabConfig.FileFieldName
    res, err = serv.uploadFile(r, fieldName, params)
    return
}



func (serv *service) uploadFile(r *http.Request, fieldName string, params *UploaderParams) (fileInfo *ResFileInfoWithState, err error) {
    file, fileHeader, err := r.FormFile(fieldName)
    if err != nil {
        return
    }

    serv.uploader.SetParams(params)

    resFileInfo, err := serv.uploader.UpFile(file, fileHeader)
    if err == nil {
        fileInfo.fromResFileInfo(resFileInfo)
        fileInfo.State = SUCCESS
        return
    }
    fileInfo.State = err.Error()
    return
}

/**
读取配置信息
 */
func (serv *service) Config() ([]byte, error) {
    return json.Marshal(GloabConfig)
}

/**
获取图片列表
 */
func (serv *service) ListImage(rootPath string, fileList []*FileItem, start int, size int)  {
    listParams := &ListParams{
        AllowFiles: GloabConfig.ImageManagerAllowFiles,
        ListSize: GloabConfig.ImageManagerListSize,
        Path: GloabConfig.ImageManagerListPath,
    }

    list := &List{
        RootPath: rootPath,
        Params: listParams,
    }

    list.GetFileList(fileList, start, size)
}

/**
获取文件列表
 */
func (serv *service) Listfile(rootPath string, fileList []*FileItem, start int, size int)  {
    listParams := &ListParams{
        AllowFiles: GloabConfig.FileManagerAllowFiles,
        ListSize: GloabConfig.FileManagerListSize,
        Path: GloabConfig.FileManagerListPath,
    }

    list := &List{
        RootPath: rootPath,
        Params: listParams,
    }

    list.GetFileList(fileList, start, size)
}

/**
从远程拉取图片
 */
func (serv *service) CatchImage(r *http.Request) (listRes *ResFilesInfoWithStates, err error) {
    params := &UploaderParams{
        PathFormat:GloabConfig.CatcherPathFormat,
        MaxSize:GloabConfig.CatcherMaxSize,
        AllowFiles:GloabConfig.CatcherAllowFiles,
        OriName:"remote.png",
    }

    serv.uploader.SetParams(params)

    fieldName := GloabConfig.CatcherFieldName

    err = r.ParseForm()
    if err != nil {
        err = errors.New("form parse error")
        return
    }

    source, _ := r.PostForm[fieldName + "[]"]

    filesInfos := make([]*ResFileInfo, 0)
    for _, imgurl := range source {
        fileInfo, _ := serv.uploader.SaveRemote(imgurl)
        filesInfos = append(filesInfos, fileInfo)
    }

    if len(filesInfos) > 0 {
        listRes.State = SUCCESS
        listRes.List = filesInfos
        return
    }
    listRes.State = ERROR
    return
}