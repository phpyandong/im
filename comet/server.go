package comet

import (
	"github.com/phpyandong/im/pkg/websocket"
	"net/http"
	"time"
	"os"
	"os/signal"
	"syscall"
	"context"
	stdlog "github.com/go-kratos/kratos/v2/log"

	"sync"
	"fmt"
	"strconv"
)

type Server struct {
	log *stdlog.Helper
	groups []*Session
	sessions sync.Map/*Session*/
	sessionLock sync.RWMutex
}

func NewServer() *Server {
	logger := stdlog.NewStdLogger(os.Stderr)
	log := stdlog.NewHelper("im_server",logger)
	return &Server{
		log: log,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
//自定义handler
type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
func longquery(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("开始"))

	time.Sleep(30 * time.Second)
	w.Write([]byte("结束"))
}

func (s *Server) Run() {
	mux := http.NewServeMux()
	mux.Handle("/", &myHandler{})
	mux.HandleFunc("/long", longquery)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		uidstr := query.Get("uid")
		uid, err := strconv.ParseInt(uidstr, 10, 64)
		if err != nil{
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.log.Errorf("Server Run conn New err :%v\n", err)

		}
		session := NewSession(uid,conn)
		fmt.Println("session:before:",s.sessions)
		s.addSession(session)
		fmt.Println("session:add:",s.sessions)

		session.run()
		session.mesDTO <- []byte("100")
		go func() {
			for   {
				if value,ok := s.sessions.Load(int64(1));ok{
					fmt.Println("hahha")
					sess1 := value.(*Session)
					sess1.mesDTO <- []byte("111")
				}
			}

		}()
		go func() {
			for {
				if value, ok := s.sessions.Load(int64(2)); ok {
					sess1 := value.(*Session)
					sess1.mesDTO <- []byte("222")
				}
			}
		}()

		//s.conn = conn
		//go s.Recv()
		//go s.Send()
	})

	server := &http.Server{
		Addr:         ":8888",
		WriteTimeout: time.Second * 3, //设置3秒的写超时
		Handler:      mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Errorf("server run err: %+v", err)
		}
	}()
	// 一个通知退出的chan
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case mes := <-quit:
		s.log.Infof("Server shutdown by quit mes :%v\n", mes)
		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			s.log.Errorf("Server forced to shutdown:%v", err)
		}

	}
	s.log.Infof("Server Out exist :%v\n", time.Now())

}
func (s *Server) addSession(sess *Session){
	s.sessions.LoadOrStore(sess.id,sess)
	fmt.Println(s.sessions.Load(sess.id))
}
//const (
//	// Time allowed to write a message to the peer.
//	writeWait = 10 * time.Second
//
//	// Time allowed to read the next pong message from the peer.
//	pongWait = 6 * time.Second
//
//	// Send pings to peer with this period. Must be less than pongWait.
//	pingPeriod = (pongWait * 9) / 10
//
//	// Maximum message size allowed from peer.
//	maxMessageSize = 51200
//)
//
//var (
//	newline = []byte{'\n'}
//	space   = []byte{' '}
//)
//
//func (c *Server) Recv() {
//	c.log.Info("Server Recv Init ...")
//
//	defer func() {
//		//c.hub.unregister <- c
//		c.conn.Close()
//	}()
//	c.conn.SetReadLimit(maxMessageSize)
//	c.conn.SetReadDeadline(time.Now().Add(pongWait))
//	c.conn.SetPongHandler(
//		func(string) error {
//			c.conn.SetReadDeadline(time.Now().Add(pongWait));
//			c.log.Infof("pong %v \n", time.Now())
//			return nil
//		})
//	for {
//		size, message, err := c.conn.ReadMessage()
//		c.log.Infof("Sever Recv message: %v size: %v err:%+v\n", message, size, err)
//		if err != nil {
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
//				c.log.Errorf("error: %v", err)
//			}
//			break
//		}
//		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
//		c.message <- message
//	}
//}

//func (c *Server) Send() {
//	c.log.Info("Server Send Init ...")
//
//	ticker := time.NewTicker(pingPeriod)
//	defer func() {
//		ticker.Stop()
//		c.conn.Close()
//	}()
//	for {
//		select {
//		case message, ok := <-c.message:
//			c.log.Infof("Server Send mess:%v", message)
//			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
//			if !ok {
//				// The hub closed the channel.
//				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
//				return
//			}
//
//			w, err := c.conn.NextWriter(websocket.TextMessage)
//			if err != nil {
//				return
//			}
//			w.Write(message)
//
//			// Add queued chat messages to the current websocket message.
//			n := len(c.message)
//			for i := 0; i < n; i++ {
//				w.Write(newline)
//				w.Write(<-c.message)
//				w.Write(<-c.message)
//			}
//			if err := w.Close(); err != nil {
//				return
//			}
//		case <-ticker.C:
//			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
//			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
//				c.log.Infof("ping mess %v\n", time.Now())
//				return
//			}
//		}
//	}
//}
