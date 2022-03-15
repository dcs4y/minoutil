package common

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	WS            *workSpaceModel
	RootPath      string
	DataPath      string
	ResourcesPath string
)

// 读取配置文件信息
func init() {
	RootPath, _ = os.Getwd()
	if strings.HasSuffix(os.Args[0], ".test.exe") {
		// 测试时使用固定路径
		RootPath = "D:\\ideaWorkspace\\game\\"
	}
	DataPath = filepath.Join(RootPath, "data")
	ResourcesPath = filepath.Join(RootPath, "resources")
	configPath := filepath.Join(ResourcesPath, "config")
	// 读取运行环境
	readConfig(filepath.Join(configPath, "application.yml"), &WS)
	// 根据运行环境重新读取配置文件信息
	if WS.ServerConfig.Active != "" {
		readConfig(filepath.Join(configPath, "application-"+WS.ServerConfig.Active+".yml"), &WS)
		// 设置默认路径
		if WS.ServerConfig.DataPath == "" {
			WS.ServerConfig.DataPath = DataPath
		} else {
			DataPath = WS.ServerConfig.DataPath
		}
	}
}

func readConfig(configPath string, target interface{}) {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println("读取配置文件信息失败！", err)
		panic("读取配置文件信息失败！")
	}
	err = yaml.Unmarshal(configData, target)
	if err != nil {
		log.Println("读取配置文件信息失败！", err)
		panic("读取配置文件信息失败！")
	}
}

// 对应配置文件内容config/application*.yml
type workSpaceModel struct {
	ServerConfig *ServerModel   `yaml:"server"`
	WebConfig    *WebModel      `yaml:"web"`
	DBConfig     *DatabaseModel `yaml:"database"`
	RedisConfig  *RedisModel    `yaml:"redis"`
	EmailConfig  *EmailModel    `yaml:"email"`
}

type ServerModel struct {
	Active   string `yaml:"active"` // 环境变量：dev,test,prod
	DataPath string `yaml:"dataPath"`
}

type DatabaseModel struct {
	Url     string `yaml:"url"`
	BbType  string `yaml:"dbType"`
	ShowSql bool   `yaml:"showSql"`
}

type RedisModel struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database int    `yaml:"database"`
	Password string `yaml:"password"`
}

type WebModel struct {
	Enable bool `yaml:"enable"`
	Port   int  `yaml:"port"`
}

type EmailModel struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
