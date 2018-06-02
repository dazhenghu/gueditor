package gueditor


type UploaderInterface interface {
    UpFile() error //上传文件的方法
    UpBase64() error //处理base64编码的图片上传
    UpRemote() error // 拉取远程图片
}

type Uploader struct {

}

func (up *Uploader) UpFile() error  {
    return nil
}

func (up *Uploader) UpBase64() error  {
    return nil
}

func (up *Uploader) UpRemote() error  {
    return nil
}