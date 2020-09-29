// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitgame

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/zhanghuizong/bitgame/app/structs"
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

	// CommonKey 加密认证 key
	CommonKey string

	// 协议默认参数
	ParamJwt structs.ParamJwt

	// 管理
	Hub *ClientManager

	// websocket 连接资源
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func closeClient(c *Client) {
	if c == nil {
		return
	}

	c.Hub.Unregister <- c
	c.conn.Close()
}

// 接受消息
func (c *Client) read() {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("接受消息异常：", err, string(debug.Stack()))
		}
	}()

	defer func() {
		closeClient(c)

		log.Println("...............接受消息..............")
	}()

	pongWaitErr := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if pongWaitErr != nil {
		log.Println("设置 SetReadDeadline 异常", pongWaitErr)
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
		log.Println("............设置 websocket 离线处理....................", code, text)

		value, ok := getHandlers("offline")
		if ok {
			value(c, nil)
		}

		// TODO
		// 删除 redis 登录记录

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
			log.Println("发送消息异常", err, string(debug.Stack()))
		}
	}()

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		closeClient(c)
		ticker.Stop()

		log.Println("............发送消息............")
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
				log.Println("websocket 发送消息异常", wErr, message)
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
func (c *Client) SendMsg(data interface{}) {
	jsonByte, err := json.Marshal(data)
	if err != nil {
		log.Fatalln("SendMsg, JSON 编码异常", err, string(debug.Stack()))
		return
	}

	// 启用加密传输
	if utils.IsAuth() {
		encodeRes := aes.Encode(jsonByte, []byte(c.CommonKey))
		jsonByte = []byte("0" + base64.StdEncoding.EncodeToString(encodeRes))
	}

	c.send <- jsonByte
}

// 系统错误消息推送
func (c *Client) pushError(data map[string]interface{}) {
	pushError(c, data)
}

// 系统错误消息推送
func pushError(c *Client, res map[string]interface{}) {
	data := map[string]interface{}{
		"cmd": "conn::error",
		"res": res,
	}

	c.SendMsg(data)
}
