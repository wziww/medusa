### 配置示例
```javascript
{
  "base": {
    // 32 位密码
    "password": "AES256Key-32Characters1234567890",
    // 加密方式
    "crypto": "aes-128-gcm",
    //  填充方式
    "padding": ""
  },
  // 服务端配置
  "server": {
    "port": 8080,
    "api": {
      "port": 8082, // 监听地址
      "enable": true // 是否启用
    }
  },
  // 客户端配置
  "client": {
    "port": 8081,
    // 服务端地址
    "remoteAddress": "0.0.0.0:8080",
    "api": {
      "port": 8083, // 监听地址
      "enable": true // 是否启用
    }
  },
  // 日志配置，多级别
  "log": {
    "logLevel": [
      "LOGINFO",
      "LOGERROR",
      "LOGEDEBUG"
    ],
    // stdout 或者指定目录
    "logPath": "stdout"
  }
}
```