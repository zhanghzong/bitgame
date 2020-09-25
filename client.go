// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitgame

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zhanghuizong/bitgame/structs"
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

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 2048
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the gws connection and the ClientManager.
type Client struct {
	// 客户端ID
	Id string

	// 用户ID
	Uid string

	// CommonKey 加密认证 key
	CommonKey string

	// 协议默认参数
	RequestJwt structs.RequestJwt

	// 管理
	Hub *ClientManager

	// The gws connection.
	Conn *websocket.Conn

	// 当前上下文
	Context *gin.Context

	// Buffered channel of outbound messages.
	Send chan []byte
}

// 接受消息
func (c *Client) read() {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("接受消息异常-read：", err, string(debug.Stack()))
		}
	}()

	defer func() {
		c.Hub.Unregister <- c
		cErr := c.Conn.Close()
		if cErr != nil {
			log.Println("接受消息 websocket 断开连接异常", cErr)
		}
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		msgType, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("读取消息异常：", msgType, message, err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v\n", err)
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
			log.Println("发送消息异常-write：", err, string(debug.Stack()))
		}
	}()

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.Hub.Unregister <- c
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The ClientManager closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current gws message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
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

	c.Send <- jsonByte
}
