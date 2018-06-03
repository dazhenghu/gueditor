# 百度ueditor的go(golang)语言后台服务程序

## 说明：

>百度提供的富文本框插件ueditor，因官方没有提供go版本后台，此项目旨在提供一个go的后台接口

## 注意:
>因个人精力有限，代码中还有一些不足之处，希望使用该库的同学能够发扬开源精神一起完善该库

## 接口：
>主方法在service.go文件中

## 示例一、：

```
// 设置umask主要是防止创建文件夹和文件的时候能够正常设置其权限
syscall.Umask(0)
// rootPath为项目的根目录，根据情况自行设定，以image举例，rootPath取值原则是 rootPath + imagePathFormat 为图片保存路径和
rootPath, _    := fileutil.GetCurrentDirectory()
// configFilePath为ueditor的相关配置，库里加载了默认配置，此处配置会与默认配置合并。配置说明参见ueditor官网
configFilePath := filepath.Join(rootPath, "config/ueditor.json") 
// 生成ueditor的服务对象
uedService, _ = gueditor.NewService(nil, nil, rootPath, configFilePath)
```
## 示例二、：

>基于gin框架的样例（https://github.com/dazhenghu/ginCms 中的admin模块）

```go
package controller

import (
    "github.com/dazhenghu/ginApp/controller"
    "github.com/gin-gonic/gin"
    "github.com/dazhenghu/gueditor"
    "github.com/dazhenghu/util/fileutil"
    "net/http"
    "path/filepath"
    "syscall"
)

type ueditorController struct {
    controller.Controller
}

var ueditorInstance *ueditorController
var uedService *gueditor.Service

func init() {
    ueditorInstance = &ueditorController{}
    ueditorInstance.Init(ueditorInstance)

    syscall.Umask(0)
    rootPath, _    := fileutil.GetCurrentDirectory()
    configFilePath := filepath.Join(rootPath, "config/ueditor.json") // 设置自定义配置文件路径

    rootPath      = filepath.Join(rootPath, "../") // 设置项目根目录
    uedService, _ = gueditor.NewService(nil, nil, rootPath, configFilePath)

    ueditorInstance.PostAndGet("/ueditor", ueditorInstance.index)
}

func (ued *ueditorController) index(context *gin.Context) {
    action := context.Query("action")

    switch action {
    case "config":
        // config接口
        ued.config(context)
    case "uploadimage":
        // 上传图片
        ued.uploadImage(context)
    case "uploadscrawl":
        // 上传涂鸦
        ued.uploadScrawl(context)
    case "uploadvideo":
        // 上传视频
        ued.uploadVideo(context)
    case "uploadfile":
        // 上传附件
        ued.uploadfile(context)
    case "listfile":
        // 查询上传的文件列表
        ued.listFile(context)
    case "listimage":
        // 查询上传的图片列表
        ued.listImage(context)
    }

}

func (ued *ueditorController) config(context *gin.Context) {
    cnf := uedService.Config()
    context.JSON(http.StatusOK, cnf)
}

func (ued *ueditorController) uploadImage(context *gin.Context) {
    res, _ := uedService.Uploadimage(context.Request)
    context.JSON(http.StatusOK, res)
}

func (ued *ueditorController) uploadScrawl(context *gin.Context)  {
    res, _ := uedService.UploadScrawl(context.Request)
    context.JSON(http.StatusOK, res)
}

func (ued *ueditorController) uploadVideo(context *gin.Context)  {
    res, _ := uedService.UploadVideo(context.Request)
    context.JSON(http.StatusOK, res)
}

func (ued *ueditorController) uploadfile(context *gin.Context)  {
    res, _ := uedService.UploadFile(context.Request)
    context.JSON(http.StatusOK, res)
}

func (ued *ueditorController) listFile(context *gin.Context) {
    listFileItem := &gueditor.ListFileItem{}
    uedService.Listfile(listFileItem, 0, 10)
    context.JSON(http.StatusOK, listFileItem)
}

func (ued *ueditorController) listImage(context *gin.Context)  {
    listFileItem := &gueditor.ListFileItem{}
    uedService.ListImage(listFileItem, 0, 10)
    context.JSON(http.StatusOK, listFileItem)
}
```

