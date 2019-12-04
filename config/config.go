package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

// C 全局配置
var C *Config

// Config 配置
type Config struct {
	Server  *Server   `toml:"server"`
	Routers []*Router `toml:"routers"`
}

// Router 代理路由
type Router struct {
	Upstream   string `toml:"upstream"`
	Path       string `toml:"path"`
	CopyStream string `toml:"copy_stream"`
	Strip      int    `toml:"strip"`
}

// Server 服务端配置
type Server struct {
	Host   string `toml:"host"`
	Port   string `toml:"port"`
	Path   string `toml:"path"`
	Scheme string `toml:"scheme"`
}

func init() {
	configPath := flag.String("c", "./conf.toml", "config path")
	flag.Parse()
	if C == nil {
		fdata, openError := ioutil.ReadFile(*configPath)
		if openError != nil {
			fmt.Fprintf(os.Stderr, "%s", openError.Error())
		}
		C = &Config{}
		toml.Decode(string(fdata), C)
	}
}
