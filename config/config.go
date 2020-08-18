package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/guomio/go-template/tools"
)

// NewConfigOption config 配置
type NewConfigOption struct {
	File string
}

// NewConfig 新建配置
func NewConfig(opt *NewConfigOption, config interface{}) error {
	data, err := ioutil.ReadFile(opt.File)

	if err != nil {
		return err
	}

	return json.Unmarshal(data, config)
}

// 初始化部分

const (
	// EnvConfigPath 配置文件目录
	envConfigPath = "ENV_CONFIG_PATH"

	// DefaultPort 默认端口
	defaultPort = "8888"
)

// Config 配置项
type Config struct {
	Log   string
	Port  string
	Mysql string
}

// C Config 实例
var C Config

// Init 初始化
func Init() {
	cwd, err := os.Getwd()

	if err != nil {
		log.Fatalln("os.Getwd fail:", err)
	}

	configPath := tools.EmptyToString(os.Getenv(envConfigPath), tools.Join(cwd, "config.json"))

	log.Println(configPath)

	err = NewConfig(&NewConfigOption{File: configPath}, &C)

	if err != nil {
		log.Fatalln("NewConfig fail:", err)
	}

	C.Port = tools.EmptyToString(C.Port, defaultPort)
	C.Log = tools.EmptyToString(C.Log, tools.Join(cwd, "logs", "guomio.log"))
}

// GetConfig 获取配置项
func GetConfig() Config {
	return C
}
