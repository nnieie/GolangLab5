package chat

import (
	"encoding/json"
	"time"

	"github.com/hertz-contrib/websocket"

	"github.com/nnieie/golanglab5/cmd/api/biz/model/common"
	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/nnieie/golanglab5/pkg/utils"
)

const (
	maxMessageChannelSize = 256
)

type Client struct {
	userID int64
	conn   *websocket.Conn
	send   chan []byte
	svc    *ChatService
}

func (c *Client) readPump() {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("panic in readPump: %v", r)
		}
		c.conn.Close()
		hub.unregister <- c
	}()
	c.conn.SetReadLimit(constants.MaxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(constants.PongWait)); err != nil {
		logger.Errorf("set read deadline error: %v", err)
		return
	}
	c.conn.SetPongHandler(func(appData string) error {
		return c.conn.SetReadDeadline(time.Now().Add(constants.PongWait))
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			logger.Debugf("read message error for user %v: %v", c.userID, err)
			break
		}

		var msg = &common.Message{}
		if err := json.Unmarshal(message, msg); err != nil {
			c.sendStruct(utils.BuildBaseResp(err))
			continue
		}

		// 消息统一处理
		if err := c.handleMessage(msg); err != nil {
			c.sendStruct(utils.BuildBaseResp(err))
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(constants.PingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					logger.Errorf("write close message error: %v", err)
				}
				return
			}
			err := c.conn.SetWriteDeadline(time.Now().Add(constants.WriteWait))
			if err != nil {
				logger.Errorf("set write deadline error: %v", err)
				return
			}
			writer, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.Errorf("failed to get next writer: %v", err)
				return
			}
			_, err = writer.Write(message)
			if err != nil {
				logger.Errorf("write message error: %v", err)
				return
			}
			n := len(c.send)
			for i := 0; i < n; i++ {
				if _, err := writer.Write([]byte{'\n'}); err != nil {
					logger.Errorf("write newline error: %v", err)
					return
				}
				if _, err = writer.Write(<-c.send); err != nil {
					logger.Errorf("write additional message error: %v", err)
					return
				}
			}
			if err := writer.Close(); err != nil {
				return
			}
			logger.Debugf("%v send %v", c.userID, string(message))
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(constants.WriteWait)); err != nil {
				logger.Errorf("set write deadline error: %v", err)
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Debugf("%v dead", c.userID)
				return
			}
		}
	}
}

func ServeWs(svc *ChatService, conn *websocket.Conn) {
	if svc == nil {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(`{"code":401,"msg":"unauthorized"}`)); err != nil {
			logger.Errorf("write message error: %v", err)
		}
		if err := conn.Close(); err != nil {
			logger.Errorf("close connection error: %v", err)
		}
		return
	}

	if conn == nil {
		logger.Errorf("websocket connection is nil")
		return
	}

	client := &Client{
		userID: svc.userID,
		conn:   conn,
		send:   make(chan []byte, maxMessageChannelSize),
		svc:    svc,
	}

	logger.Debugf("new websocket client connected: user %v", client.userID)
	hub.register <- client

	go client.writePump()
	client.readPump()

	logger.Debugf("websocket client disconnected: user %v", client.userID)
}

func (c *Client) sendStruct(s any) {
	msg, err := json.Marshal(s)
	if err != nil {
		c.send <- []byte(err.Error())
	}
	c.send <- msg
}
