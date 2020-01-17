package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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
	Padding  string `json:"padding"`
	Client   bool
}

// API 相关配置
type API struct {
	Port   int  `json:"port"`
	Enable bool `json:"enable"`
}
type ipMask struct {
	ipaddr uint32
	mask   uint8
}

// Server 服务端配置
type Server struct {
	// Port 默认 0 ，采用随机端口，客户端需要自行更新 port
	Port      int      `json:"port"`
	API       API      `json:"api"`
	WhiteList []string `json:"whiteList"`
	whiteList []ipMask
	BlackList []string `json:"blackList"`
	blackList []ipMask
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

// Init 日志初始化
func Init() {
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
	for _, v := range C.Server.WhiteList {
		if vv := parseIPv4(v); vv != nil {
			C.Server.whiteList = append(C.Server.whiteList, vv.(ipMask))
		} else {
			fmt.Fprintf(os.Stderr, "%s %s\n", v, "parse ipv4 failed,only support for x.x.x.x or x.x.x.x/xx")
		}
	}
	for _, v := range C.Server.BlackList {
		if vv := parseIPv4(v); vv != nil {
			C.Server.blackList = append(C.Server.blackList, vv.(ipMask))
		} else {
			fmt.Fprintf(os.Stderr, "%s %s\n", v, "parse ipv4 failed,only support for x.x.x.x or x.x.x.x/xx")
		}
	}
}
func parseIPv4(str string) interface{} {
	str2 := strings.Split(str, "/")
	result := ipMask{}
	if len(str2) == 1 {
		result.mask = 0
	} else if len(str2) == 2 {
		mask, _ := strconv.Atoi(str2[1])
		if mask < 0 || mask > 32 {
			return nil
		}
		result.mask = uint8(mask)
	} else {
		return nil
	}
	ipdata := strings.Split(str2[0], ".")
	if len(ipdata) != 4 {
		return nil
	}
	for _, v := range ipdata {
		eachdata, _ := strconv.Atoi(v)
		result.ipaddr = result.ipaddr<<8 | uint32(eachdata)
	}
	if result.ipaddr == 0 && result.mask == 0 {
		for i := 0; i < 4; i++ {
			result.ipaddr |= 255 << 8
		}
	}
	return result
}

// CheckIPAllow str must be ip:port in ipv4
// such as x.x.x.x:22
func CheckIPAllow(str string) bool {
	str2 := strings.Split(str, ":")
	if len(str2) != 2 {
		return false
	}
	var data uint32
	ipdata := strings.Split(str2[0], ".")
	if len(ipdata) != 4 {
		return false
	}
	for _, v := range ipdata {
		eachdata, _ := strconv.Atoi(v)
		data = data<<8 | uint32(eachdata)
	}
	for _, v := range C.Server.whiteList {
		if data&v.ipaddr == v.ipaddr {
			return true
		}
	}
	return false
}

// CheckIPNotAllow str must be ip:port in ipv4
// such as x.x.x.x:22
func CheckIPNotAllow(str string) bool {
	str2 := strings.Split(str, ":")
	if len(str2) != 2 {
		return true
	}
	var data uint32
	ipdata := strings.Split(str2[0], ".")
	if len(ipdata) != 4 {
		return true
	}
	for _, v := range ipdata {
		eachdata, _ := strconv.Atoi(v)
		data = data<<8 | uint32(eachdata)
	}
	for _, v := range C.Server.blackList {
		if data&v.ipaddr == v.ipaddr { // now allowed
			return true
		}
	}
	return false
}
