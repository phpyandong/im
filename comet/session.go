package comet

import (
	"github.com/phpyandong/im/pkg/websocket"
	"time"
	"os"
	stdlog "github.com/go-kratos/kratos/v2/log"

	"sync"
)

type Session struct {
	id int64
	conn    *websocket.Conn
	//mesDTO  chan *model.MsgDTO
	mesDTO  chan []byte
	log 	*stdlog.Helper
	rw 		sync.RWMutex
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 6 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 51200
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (c *Session) Recv() {
	c.log.Info("Server Recv Init ...")

	defer func() {
		//c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(
		func(string) error {
			c.conn.SetReadDeadline(time.Now().Add(pongWait));
			//c.log.Infof("pong %v \n", time.Now())
			return nil
		})
	for {
		size, message, err := c.conn.ReadMessage()
		c.log.Infof("Sever Recv message: %v size: %v err:%+v\n", message, size, err)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.log.Errorf("error: %v", err)
			}
			break
		}
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.mesDTO <- message
	}
}
func (c *Session) Send() {
	c.log.Info("Server Send Init ...")

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.mesDTO:
			c.log.Infof("Server Send mess:%v", message)
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.mesDTO)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.mesDTO)
				w.Write(<-c.mesDTO)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				//c.log.Infof("ping mess %v\n", time.Now())
				return
			}
		}
	}
}
func NewSession(uid int64,conn *websocket.Conn) *Session{
	logger := stdlog.NewStdLogger(os.Stderr)
	logs := stdlog.NewHelper("session",logger)
	return &Session{
		id:uid,
		conn :conn,
		mesDTO:make(chan []byte),
		log :logs,
	}
}
func (session *Session) run(){

	go session.Recv()
	go session.Send()
}

