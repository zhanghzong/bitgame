# bitgame

包名称：`github.com/zhanghuizong/bitgame`

## 实现功能

- websocket
    - 路由分发
        - websocket hook (online|offline)
    - 异地登录
    - 通讯协议加解密
- MySQL
- Redis
- 配置文件

# 配置文件

> 注意：文件名称必须为 `config.properties`

- 应用配置文件
> 支持热更新

-  阿波罗
> 支持热更新

> `注意：如果以上两种配置文件都存在，将优先采用应用配置文件`


# 运行模式

bitgame 模块不支持单独运行，必须以模块安装。在业务逻辑中启动服务

在 `go.mod` 中增加如下：

```
require (
    github.com/zhanghuizong/bitgame v1.0.1
)
```

# 业务逻辑层次结构

> 注意：后续会提供脚手架，一键生成。统一采用

```
├── README.md
├── app
│     ├── app.go
│     ├── constants
│     │     ├── errConst
│     │     ├── redisConst
│     │     ├── roomConst
│     │     └── userConst
│     ├── controllers
│     │     ├── hall.go
│     │     ├── offline.go
│     │     ├── online.go
│     │     ├── room.go
│     │     └── user.go
│     ├── logics
│     │     ├── GameLogic
│     │     ├── HallLogic
│     │     ├── RoomLogic
│     │     └── UserLogic
│     ├── models
│     │     ├── RoomModel
│     │     ├── UserModel
│     │     └── model.go
│     └── utils
│         └── utils.go
├── app.go
├── config
│     ├── app.yaml
│     └── app.yaml.example
├── go.mod
├── go.sum
├── go_build_game.exe
├── main.go
├── routes
│     ├── web_routes.go
│     └── ws_routes.go
└── test
    ├── crypto-js.js
    ├── jsencrypt.js
    ├── test.html
    └── utils.js

19 directories, 22 files
zhanghuizong@DESKTOP-FD0NVQI:/mnt/d/workspaces/go_game_server$ tree -L 4
.
├── README.md
├── app
│     ├── app.go
│     ├── constants
│     │     ├── errConst
│     │     ├── redisConst
│     │     ├── roomConst
│     │     └── userConst
│     ├── controllers
│     │     ├── hall.go
│     │     ├── offline.go
│     │     ├── online.go
│     │     ├── room.go
│     │     └── user.go
│     ├── logics
│     │     ├── GameLogic
│     │     ├── HallLogic
│     │     ├── RoomLogic
│     │     └── UserLogic
│     ├── models
│     │     ├── RoomModel
│     │     ├── UserModel
│     │     └── model.go
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

# 安装

```
go get -u github.com/zhanghuizong/bitgame 
```

> 在实际开发中，可以将代码复制 `gopath` 目录对应位置，可避免反复安装过程

## 说明
1. 提交代码需要去 github 打最新版本 `tag`，否则无法生效
2. 执行完 go get 命令后，会将代码打包上传到 `go.dev` 网站

## 后续
- 计划搭建内网私有源