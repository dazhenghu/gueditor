package storage

import (
    "io"
    qiniu_storage "github.com/qiniu/api.v7/storage"
    "github.com/qiniu/api.v7/auth/qbox"
    "fmt"
    "context"
    "time"
    "github.com/dazhenghu/util/encryptutil"
    "bytes"
)

type QiniuConfig struct {
    AccessKey     string
    SecretKey     string
    Bucket        string
    PolicyExpires uint32              // 上传凭证的有效时间，单位秒
    Zone          *qiniu_storage.Zone //空间所在的机房
    UseHTTPS      bool                //是否使用https域名
    UseCdnDomains bool                //是否使用cdn加速域名
    CentralRsHost string              //中心机房的RsHost，用于list bucket
}

type Qiniu struct {
    BaseInterface
    upToken string
}

type QiniuRet struct {
    Key    string
    Hash   string
    Fsize  int
    Bucket string
}

func NewQiniu(config *QiniuConfig) (*Qiniu) {
    qiniuObj := &Qiniu{}

    putPolicy := qiniu_storage.PutPolicy{
        Scope: config.Bucket,
        ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)"}`,
    }

    if config.PolicyExpires != 0 {
        putPolicy.Expires = config.PolicyExpires
    }

    mac := qbox.NewMac(config.AccessKey, config.SecretKey)
    qiniuObj.upToken = putPolicy.UploadToken(mac)
    return qiniuObj
}

func (qn *Qiniu) SaveFileFromLocalPath(srcPath string, dstAbsPath, dstRelativePath string) (url string, err error) {
    cfg := qiniu_storage.Config{}
    formUploader := qiniu_storage.NewFormUploader(&cfg)

    ret := QiniuRet{}
    putExtra := qiniu_storage.PutExtra{}

    key := encryptutil.Md5(fmt.Sprintf("%d_%s", time.Now().Unix(), dstRelativePath))

    err = formUploader.PutFile(context.Background(), &ret, qn.upToken, key, srcPath, &putExtra)
    if err != nil {
        return
    }

    return
}

func (qn *Qiniu) SaveFile(srcFile io.Reader, srcFileSize int64, dstAbsPath, dstRelativePath string) (url string, err error) {
    cfg := qiniu_storage.Config{}
    formUploader := qiniu_storage.NewFormUploader(&cfg)

    ret := QiniuRet{}
    putExtra := qiniu_storage.PutExtra{}

    err = formUploader.Put(context.Background(), &ret, qn.upToken, dstRelativePath, srcFile, srcFileSize, &putExtra)

    url = ret.Key

    return
}

func (qn *Qiniu) SaveData(data []byte, dstAbsPath, dstRelativePath string) (url string, err error) {
    cfg := qiniu_storage.Config{}
    formUploader := qiniu_storage.NewFormUploader(&cfg)

    ret := QiniuRet{}
    putExtra := qiniu_storage.PutExtra{}

    dataBuffer := bytes.NewBuffer(data)

    err = formUploader.Put(context.Background(), ret, qn.upToken, dstRelativePath, dataBuffer, int64(dataBuffer.Len()), &putExtra)
    url = ret.Key
    return
}
