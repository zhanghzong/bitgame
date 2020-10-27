// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/zhanghuizong/bitgame/app/definition"
	"github.com/zhanghuizong/bitgame/app/models"
	"github.com/zhanghuizong/bitgame/utils"
	"github.com/zhanghuizong/bitgame/utils/aes"
	"log"
	"runtime/debug"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 2048
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	// 客户端ID
	SocketId string

	// 用户ID
	Uid string

	// 协议默认参数
	Jwt definition.ParamJwt

	// 请求数据
	Msg *definition.RequestMsg

	// 管理
	Hub *ClientManager

	// websocket 连接资源
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// commonKey 加密认证 key
	commonKey string

	*logrus.Entry
}

func closeClient(c *Client) {
	if c == nil {
		return
	}

	c.Hub.unregister <- c
	c.conn.Close()
}

// 接受消息
func (c *Client) read() {
	defer func() {
		err := recover()
		if err != nil {
			c.Warnf("接受消息异常. err:%s, stack:%s", err, string(debug.Stack()))
		}
	}()

	defer func() {
		closeClient(c)
	}()

	pongWaitErr := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if pongWaitErr != nil {
		c.Warnf("设置 SetReadDeadline 异常. err:%s", pongWaitErr)
		return
	}

	// 设置 读取消息体大小
	c.conn.SetReadLimit(maxMessageSize)

	// 设置 pong 方法
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// 设置 websocket 离线处理
	c.conn.SetCloseHandler(func(code int, text string) error {
		model := new(models.LoginModel)
		uid := c.Uid
		connSocketId := model.GetSocketId(uid)
		if connSocketId != c.SocketId {
			return nil
		}

		c.Infof("客户端离线, 错误码：%d, 错误：%s", code, text)

		// offline
		value, ok := getHandlers("offline")
		if ok {
			value(c)
		}

		// 删除 redis 登录记录
		model.DelSocketId(uid)

		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("websocket IsUnexpectedCloseError", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// 解析数据格式
		parseMsg(c, message)
	}
}

// 发送消息
func (c *Client) write() {
	defer func() {
		err := recover()
		if err != nil {
			c.Warnf("发送消息异常, err:%s, stack:%s", err, string(debug.Stack()))
		}
	}()

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		closeClient(c)
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			_, wErr := w.Write(message)
			if wErr != nil {
				c.Warnf("websocket 发送消息异常. err:%s, msg:%s", wErr, message)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 消息单播
func (c *Client) sendMsg(data interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			c.Errorf("发送消息异常. err:%s", err)
		}
	}()

	if c == nil {
		closeClient(c)
		return
	}

	jsonByte, err := json.Marshal(data)
	if err != nil {
		c.Warnf("sendMsg, JSON 编码异常. err:%s, stack:%s", err, string(debug.Stack()))
		return
	}

	c.Infof("消息推送:%s", jsonByte)

	// 启用加密传输
	if utils.IsAuth() {
		encodeRes := aes.Encode(jsonByte, []byte(c.commonKey))
		jsonByte = []byte("0" + base64.StdEncoding.EncodeToString(encodeRes))
	}

	c.send <- jsonByte
}

// 统一消息推送格式
// 正确消息单播
func (c *Client) Success(cmd string, data interface{}) {
	if c == nil {
		return
	}

	pushClient(c, pushSuccess(cmd, data))
}

// 统一消息推送格式
// 错误消息单播
func (c *Client) Error(cmd string, row definition.ErrMsgStruct) {
	if c == nil {
		return
	}

	pushClient(c, pushError(cmd, row))
}

// 统一消息推送格式
// uid 模式消息推送
func (c *Client) Broadcast(users []string, cmd string, data interface{}) {
	broadcast(users, cmd, data)
}

// websocket 系统内部错（兼容历史数据格式）
// 系统错误消息推送
func (c *Client) insidePushError(data map[string]interface{}) {
	insidePushError(c, data)
}
