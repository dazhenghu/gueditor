package gueditor

import (
    "io/ioutil"
    "encoding/json"
    "errors"
    "fmt"
)

func init()  {
    loadDefaultConfig()
}

var GloabConfig *Config

type Config struct {
    ImageActionName         string   `json:"imageActionName"`
    ImageFieldName          string   `json:"imageFieldName"`
    ImageMaxSize            int      `json:"imageMaxSize"`
    ImageAllowFiles         []string `json:"imageAllowFiles"`
    ImageCompressEnable     bool     `json:"imageCompressEnable"`
    ImageCompressBorder     int      `json:"imageCompressBorder"`
    ImageInsertAlign        string   `json:"imageInsertAlign"`
    ImageURLPrefix          string   `json:"imageUrlPrefix"`
    ImagePathFormat         string   `json:"imagePathFormat"`
    ScrawlActionName        string   `json:"scrawlActionName"`
    ScrawlFieldName         string   `json:"scrawlFieldName"`
    ScrawlPathFormat        string   `json:"scrawlPathFormat"`
    ScrawlMaxSize           int      `json:"scrawlMaxSize"`
    ScrawlURLPrefix         string   `json:"scrawlUrlPrefix"`
    ScrawlInsertAlign       string   `json:"scrawlInsertAlign"`
    SnapscreenActionName    string   `json:"snapscreenActionName"`
    SnapscreenPathFormat    string   `json:"snapscreenPathFormat"`
    SnapscreenURLPrefix     string   `json:"snapscreenUrlPrefix"`
    SnapscreenInsertAlign   string   `json:"snapscreenInsertAlign"`
    CatcherLocalDomain      []string `json:"catcherLocalDomain"`
    CatcherActionName       string   `json:"catcherActionName"`
    CatcherFieldName        string   `json:"catcherFieldName"`
    CatcherPathFormat       string   `json:"catcherPathFormat"`
    CatcherURLPrefix        string   `json:"catcherUrlPrefix"`
    CatcherMaxSize          int      `json:"catcherMaxSize"`
    CatcherAllowFiles       []string `json:"catcherAllowFiles"`
    VideoActionName         string   `json:"videoActionName"`
    VideoFieldName          string   `json:"videoFieldName"`
    VideoPathFormat         string   `json:"videoPathFormat"`
    VideoURLPrefix          string   `json:"videoUrlPrefix"`
    VideoMaxSize            int      `json:"videoMaxSize"`
    VideoAllowFiles         []string `json:"videoAllowFiles"`
    FileActionName          string   `json:"fileActionName"`
    FileFieldName           string   `json:"fileFieldName"`
    FilePathFormat          string   `json:"filePathFormat"`
    FileURLPrefix           string   `json:"fileUrlPrefix"`
    FileMaxSize             int      `json:"fileMaxSize"`
    FileAllowFiles          []string `json:"fileAllowFiles"`
    ImageManagerActionName  string   `json:"imageManagerActionName"`
    ImageManagerListPath    string   `json:"imageManagerListPath"`
    ImageManagerListSize    int      `json:"imageManagerListSize"`
    ImageManagerURLPrefix   string   `json:"imageManagerUrlPrefix"`
    ImageManagerInsertAlign string   `json:"imageManagerInsertAlign"`
    ImageManagerAllowFiles  []string `json:"imageManagerAllowFiles"`
    FileManagerActionName   string   `json:"fileManagerActionName"`
    FileManagerListPath     string   `json:"fileManagerListPath"`
    FileManagerURLPrefix    string   `json:"fileManagerUrlPrefix"`
    FileManagerListSize     int      `json:"fileManagerListSize"`
    FileManagerAllowFiles   []string `json:"fileManagerAllowFiles"`
}

/**
加载默认配置
 */
func loadDefaultConfig() (err error) {
    filePath, err := getDefaultConfigFile()
    if err != nil {
        return
    }
    cnfJson, err := ioutil.ReadFile(filePath)
    if err != nil {
        return
    }
    err = json.Unmarshal(cnfJson, GloabConfig)
    if err != nil {
        return
    }

    return
}

/**
加载配置
 */
func loadConfigFromFile(filePath string) (err error) {
    exists, err := pathExists(filePath)
    if !exists {
        err = errors.New(fmt.Sprintf("config file not exists:%s", filePath))
        return
    }

    cnfJson, err := ioutil.ReadFile(filePath)
    if err != nil {
        return
    }
    err = json.Unmarshal(cnfJson, GloabConfig)
    if err != nil {
        return
    }

    return
}
