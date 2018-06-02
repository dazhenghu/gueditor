package gueditor



type service struct {
    uploader UploaderInterface
}

func NewService(uploader UploaderInterface, configFile string) (serv *service, err error) {
    if uploader == nil {
        // 没有注入uploader接口，则使用默认的方法
        uploader = &Uploader{}
    }

    serv = &service{
        uploader:uploader,
    }

    if configFile != "" {
        // 加载默认配置
        err = loadDefaultConfig()
    }

    return
}

func (serv *service) Uploadimage()  {
    serv.uploader.UpFile()
}
