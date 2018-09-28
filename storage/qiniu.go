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
    Zone          *qiniu_storage.Zone // 空间所在的机房
    UseHTTPS      bool                // 是否使用https域名
    UseCdnDomains bool                // 是否使用cdn加速域名
    CentralRsHost string              // 中心机房的RsHost，用于list bucket
    Domain        string              // 外链域名
}

type Qiniu struct {
    BaseInterface
    putPolicy *qiniu_storage.PutPolicy // 上传策略
    config *QiniuConfig // 配置信息
}

type QiniuRet struct {
    Key    string
    Hash   string
    Fsize  int
    Bucket string
}

func NewQiniu(config *QiniuConfig) (*Qiniu) {
    qiniuObj := &Qiniu{}

    putPolicy := &qiniu_storage.PutPolicy{
        Scope: config.Bucket,
        ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)"}`,
    }

    if config.PolicyExpires != 0 {
        putPolicy.Expires = config.PolicyExpires
    }

    qiniuObj.putPolicy = putPolicy
    qiniuObj.config = config
    return qiniuObj
}

/**
从本地文件保存
 */
func (qn *Qiniu) SaveFileFromLocalPath(srcPath string, dstAbsPath, dstRelativePath string) (url string, err error) {
    cfg := qiniu_storage.Config{}
    formUploader := qiniu_storage.NewFormUploader(&cfg)

    ret := QiniuRet{}
    putExtra := qiniu_storage.PutExtra{}

    key := encryptutil.Md5(fmt.Sprintf("%d_%s", time.Now().Unix(), dstRelativePath))

    // 按照官方文档建议，每次上传都重新请求一次上传token
    err = formUploader.PutFile(context.Background(), &ret, qn.getUpToken(), key, srcPath, &putExtra)
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

    err = formUploader.Put(context.Background(), &ret, qn.getUpToken(), dstRelativePath, srcFile, srcFileSize, &putExtra)

    // 生成外链
    url = qiniu_storage.MakePublicURL(qn.config.Domain, ret.Key)
    return
}

func (qn *Qiniu) SaveData(data []byte, dstAbsPath, dstRelativePath string) (url string, err error) {
    cfg := qiniu_storage.Config{}
    formUploader := qiniu_storage.NewFormUploader(&cfg)

    ret := QiniuRet{}
    putExtra := qiniu_storage.PutExtra{}

    dataBuffer := bytes.NewBuffer(data)

    err = formUploader.Put(context.Background(), ret, qn.getUpToken(), dstRelativePath, dataBuffer, int64(dataBuffer.Len()), &putExtra)

    url = qiniu_storage.MakePublicURL(qn.config.Domain, ret.Key)
    return
}

/**
获取上传请求token
 */
func (qn *Qiniu) getUpToken() string {
    mac := qbox.NewMac(qn.config.AccessKey, qn.config.SecretKey)
    return qn.putPolicy.UploadToken(mac)
}
