package socket

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/openimsdk/open-im-server/v3/framework/clean/pkg/common"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func StartServer(ctx context.Context, addr string, loggerPrefix string, broker *common.Broker) error {
	logger := common.NewLogger(loggerPrefix)
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Printf("upgrade error: %v", err)
			return
		}
		defer conn.Close()

		sub := broker.Subscribe(64)
		defer broker.Unsubscribe(sub)

		errCh := make(chan error, 1)
		go func() {
			for msg := range sub {
				if writeErr := conn.WriteMessage(websocket.TextMessage, []byte("broadcast:"+msg)); writeErr != nil {
					errCh <- writeErr
					return
				}
			}
		}()

		for {
			_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			mt, message, readErr := conn.ReadMessage()
			if readErr != nil {
				logger.Printf("socket read error: %v", readErr)
				return
			}
			if writeErr := conn.WriteMessage(mt, message); writeErr != nil {
				logger.Printf("socket write error: %v", writeErr)
				return
			}
			select {
			case bgErr := <-errCh:
				logger.Printf("socket broadcast error: %v", bgErr)
				return
			default:
			}
		}
	})

	srv := &http.Server{Addr: addr, Handler: mux}
	logger.Printf("socket server listening at %s", addr)

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()
	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()

	err := <-errCh
	if err == http.ErrServerClosed {
		logger.Printf("socket server stopped")
		return nil
	}
	return err
}
