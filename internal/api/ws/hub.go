package chat

import (
	"strconv"
	"time"

	"github.com/nnieie/golanglab5/internal/api/rpc"
	"github.com/nnieie/golanglab5/kitex_gen/user"
	"github.com/nnieie/golanglab5/pkg/logger"
)

type Hub struct {
	clients    map[string]*Client
	broadcast  chan *Broadcast
	register   chan *Client
	unregister chan *Client
}

type Broadcast struct {
	TargetUserIDs []int64
	Payload       []byte
}

var hub = newHub()

func InitChatHub() {
	go hub.run()
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Broadcast),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			if client != nil {
				h.clients[client.userID] = client
				logger.Infof("user %s connected", client.userID)
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client.userID]; ok {
				err := rpc.UpdateUserLastLogoutTime(client.svc.ctx, &user.UpdateLastLogoutTimeRequest{
					UserId:     client.userID,
					LogoutTime: time.Now().Unix(),
				})
				if err != nil {
					logger.Errorf("update user last logout time error: %v", err)
				}
				delete(h.clients, client.userID)
				close(client.send)
				if client.conn != nil {
					client.conn.Close()
				}
			}
		case msg := <-h.broadcast:
			go func() {
				for _, uid := range msg.TargetUserIDs {
					uidStr := strconv.FormatInt(uid, 10)
					logger.Debugf("broadcasting message to %d", uid)
					if client, ok := h.clients[uidStr]; ok {
						logger.Debugf("found connected client for %d", uid)
						select {
						case client.send <- msg.Payload:
							logger.Debugf("msg sent to %d", uid)
						default:
							h.unregister <- client
						}
					}
				}
			}()
		}
	}
}
