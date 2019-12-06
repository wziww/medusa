package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

// C 全局配置
var C *Config

// Config 配置
type Config struct {
	Base   Base   `json:"base"`
	Server Server `json:"server"`
	Client Client `json:"client"`
	Log    Log    `json:"log"`
}

// Base 基础配置
type Base struct {
	Password string `json:"password"`
	Crypto   string `json:"crypto"`
	Client   bool
}

// API 相关配置
type API struct {
	Port   int  `json:"port"`
	Enable bool `json:"enable"`
}

// Server 服务端配置
type Server struct {
	// Port 默认 0 ，采用随机端口，客户端需要自行更新 port
	Port int `json:"port"`
	API  API `json:"api"`
}

// Log 日志相关配置
type Log struct {
	LogLevel []string `json:"logLevel"`
	LogPath  string   `json:"logPath"`
}

// Client 客户端配置
type Client struct {
	RemoteAddress string `json:"remoteAddress"`
	// Port 默认 0 ，采用随机端口，客户端需要自行更新 port
	Port int `json:"port"`
	API  API `json:"api"`
}

func init() {
	/**
	 * 文件加载流程 => 使用指定文件
	 * 默认文件：config 文件夹下 conf.json 或者项目根目录下 conf.json
	 */
	configPath := flag.String("c", "../conf.json", "config path")
	flag.Parse()
	if C == nil {
		fdata, openError := ioutil.ReadFile(*configPath)
		if openError != nil {
			fmt.Fprintf(os.Stderr, "%s\n", openError.Error())
			openError = nil
			fdata, openError = ioutil.ReadFile("./conf.json")
			if openError != nil {
				fmt.Fprintf(os.Stderr, "%s", openError.Error())
				os.Exit(0)
			}
		}
		C = &Config{}
		json.Unmarshal(fdata, C)
	}
}
