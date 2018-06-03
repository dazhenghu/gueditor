package gueditor

import "net/http"

type service struct {
    uploader UploaderInterface
}

func NewService(uploaderObj UploaderInterface, configFile string) (serv *service, err error) {
    if uploaderObj == nil {
        // 没有注入uploader接口，则使用默认的方法
        uploaderObj = &Uploader{}
    }

    serv = &service{
        uploader:uploaderObj,
    }

    if configFile != "" {
        // 加载默认配置
        err = loadConfigFromFile(configFile)
    }

    return
}

/**
上传图片
 */
func (serv *service) Uploadimage(r *http.Request) (err error) {
    params := &UploaderParams{
        PathFormat:GloabConfig.ImagePathFormat,
        MaxSize:GloabConfig.ImageMaxSize,
        AllowFiles:GloabConfig.ImageAllowFiles,
    }

    fieldName := GloabConfig.ImageFieldName
    serv.uploadFile(r, fieldName, params)
    return
}

/**
上传涂鸦
 */
func (serv *service) UploadScrawl(r *http.Request) (err error)  {
    params := &UploaderParams{
        PathFormat:GloabConfig.ScrawlPathFormat,
        MaxSize:GloabConfig.ScrawlMaxSize,
        AllowFiles:GloabConfig.ImageAllowFiles,
        OriName:"scrawl.png",
    }

    fieldName := GloabConfig.ScrawlFieldName
    data := r.PostFormValue(fieldName)
    serv.uploader.SetParams(params)

    err = serv.uploader.UpBase64(params.OriName, data)
    return
}

/**
上传视频
 */
func (serv *service) UploadVideo(r *http.Request) (err error) {
    params := &UploaderParams{
        PathFormat:GloabConfig.VideoPathFormat,
        MaxSize:GloabConfig.VideoMaxSize,
        AllowFiles:GloabConfig.VideoAllowFiles,
    }

    fieldName := GloabConfig.VideoFieldName
    err = serv.uploadFile(r, fieldName, params)
    return
}

/**
上传文件
 */
func (serv *service) UploadFile(r *http.Request) (err error) {
    params := &UploaderParams{
        PathFormat:GloabConfig.FilePathFormat,
        MaxSize:GloabConfig.FileMaxSize,
        AllowFiles:GloabConfig.FileAllowFiles,
    }

    fieldName := GloabConfig.FileFieldName
    err = serv.uploadFile(r, fieldName, params)
    return
}



func (serv *service) uploadFile(r *http.Request, fieldName string, params *UploaderParams) (err error) {
    file, fileHeader, err := r.FormFile(fieldName)
    if err != nil {
        return
    }

    serv.uploader.SetParams(params)

    err = serv.uploader.UpFile(file, fileHeader)
    return
}
