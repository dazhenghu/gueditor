# 百度ueditor的go(golang)语言后台服务程序

# 说明：

>百度提供的富文本框插件ueditor，因官方没有提供go版本后台，此项目旨在提供一个go的后台接口

# 注意:
>因个人精力有限，代码中还有一些不足之处，希望使用该库的同学能够发扬开源精神一起完善该库

# 接口：
>主方法在service.go文件中

# 示例：

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


