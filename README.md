# bitgame

## 实现功能

- websocket
    - 路由分发
    - 异地登录
- 通讯协议加解密
- MySQL
- Redis
- websocket hook (online|offline)

# 运行模式

bitgame 模块不支持单独运行，必须以模块安装。在业务逻辑中启动服务

在 `go.mod` 中增加如下：

```
require (
    github.com/zhanghuizong/bitgame v1.0.1
)
```

# 业务逻辑层次结构

> 注意：后续会提供脚手架，一键生成

```
├── README.md
├── app
│     ├── app.go
│     ├── constants
│     ├── controllers
│     ├── logics
│     ├── models
│     └── utils
├── app.go
├── config
│     ├── app.yaml
│     └── app.yaml.example
├── go.mod
├── go.sum
├── main.go
├── routes
│     ├── web_routes.go
│     └── ws_routes.go
└── test
    ├── crypto-js.js
    ├── jsencrypt.js
    ├── test.html
    └── utils.js
```