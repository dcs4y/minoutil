package configclient

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

var vipers = make(map[string]*viper.Viper)

// LoadViper 读取自定义配置文件
// 显示调用Set设置值 > 命令行参数（flag）> 环境变量 > 配置文件 > key/value存储 > 默认值
// 目前Viper配置的键（Key）是大小写不敏感的。
// 没有扩展名的配置文件名称，具体扩展名，会自动以这个顺序去尝试：json>toml>yaml>yml>properties>props>prop>hcl>dotenv>env>ini。
// 所以，如果相同路径有 my.json 又有 my.toml，会读入 my.json。
func LoadViper(name string, path string) (*viper.Viper, error) {
	v := viper.GetViper()
	if name != "" {
		v = viper.New()
	}
	v.SetConfigFile(path) // 指定配置文件路径
	//viper.SetConfigName("spider-dev") // 配置文件名称(无扩展名)
	//viper.SetConfigType("yaml") // 如果配置文件的名称中没有扩展名，则需要配置此项
	//viper.AddConfigPath(path)   // 查找配置文件所在的路径
	//viper.AddConfigPath(".")    // 多次调用以添加多个搜索路径。还可以在工作目录中查找配置。
	err := v.ReadInConfig()
	if err != nil {
		fmt.Println("读取配置文件失败：" + err.Error())
	}
	vipers[name] = v
	return v, err
}

// GetConfig 解析主配置文件信息
func GetConfig(key string, value interface{}) error {
	return GetConfigByName("", key, value)
}

// GetConfigByName 解析自定义配置文件信息
func GetConfigByName(name string, key string, value interface{}) error {
	v, ok := vipers[name]
	if !ok {
		return errors.New("not init config " + name)
	}
	return v.UnmarshalKey(key, value)
}
