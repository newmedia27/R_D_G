package ws

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/contrib/v3/websocket"
	"project/internal/logger"
)

const (
	wait   = time.Second * 60
	period = (wait * 9) / 10

	maxMsgSize     = 1024 * 512
	sendBufferSize = 256
)

type Client struct {
	conn     *websocket.Conn
	UserID   string
	tokenExp time.Time
	send     chan []byte
	done     chan struct{}
}

func NewClient(c *websocket.Conn, userID string, tokenExp time.Time) *Client {
	return &Client{
		conn:     c,
		UserID:   userID,
		tokenExp: tokenExp,
		send:     make(chan []byte, sendBufferSize),
		done:     make(chan struct{}),
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	tokenTimer := time.NewTimer(time.Until(c.tokenExp) - time.Minute)
	defer tokenTimer.Stop()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logger.LogErrorContext(context.Background(), err)
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				logger.LogErrorContext(context.Background(), err)
				return
			}
		case <-tokenTimer.C:
			event, err := json.Marshal(OutgoingEvent{
				Type: EventAuthExpiringSoon,
			})
			if err != nil {
				logger.LogErrorContext(context.Background(), err)
				return
			}
			if err = c.conn.WriteMessage(websocket.TextMessage, event); err != nil {
				logger.LogErrorContext(context.Background(), err)
				return
			}
		case <-c.done:
			return
		}
	}
}
func (c *Client) readPump(h *Hub, onMessage func(*Client, []byte)) {
	defer func() {
		close(c.done)
		h.Unregister(c)
		err := c.conn.Close()
		if err != nil {
			logger.LogErrorContext(context.Background(), err)
		}
	}()
	c.conn.SetReadLimit(maxMsgSize)
	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(wait)); err != nil {
			logger.LogErrorContext(context.Background(), err)
			return err
		}
		return nil
	})
	if err := c.conn.SetReadDeadline(time.Now().Add(wait)); err != nil {
		logger.LogErrorContext(context.Background(), err)
	}
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		onMessage(c, msg)
	}
}
