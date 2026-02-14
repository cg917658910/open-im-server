package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	frpc "github.com/openimsdk/open-im-server/v3/framework/clean/internal/rpc"
	"github.com/openimsdk/open-im-server/v3/framework/clean/pkg/common"
)

type EchoReq struct {
	Message string `json:"message" binding:"required"`
}

type PublishReq struct {
	Message string `json:"message" binding:"required"`
}

func StartServer(ctx context.Context, addr, rpcAddr, loggerPrefix string, broker *common.Broker, store MessageStore) error {
	logger := common.NewLogger(loggerPrefix)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), requestID(), requestLogger(logger))

	var rpcClient *frpc.Client
	var err error
	for i := 0; i < 20; i++ {
		rpcClient, err = frpc.NewClient(rpcAddr)
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		return err
	}

	api := r.Group("/api/v1")
	api.GET("/hello", func(c *gin.Context) {
		name := c.Query("name")
		msg, callErr := rpcClient.SayHello(name)
		if callErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": callErr.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": msg})
	})

	api.POST("/echo", func(c *gin.Context) {
		var req EchoReq
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": common.ErrBadRequest.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": req.Message})
	})

	api.POST("/publish", func(c *gin.Context) {
		var req PublishReq
		if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": common.ErrBadRequest.Error()})
			return
		}
		store.Save(req.Message)
		broker.Publish(req.Message)
		c.JSON(http.StatusOK, gin.H{"published": true})
	})

	api.GET("/messages", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"messages": store.List()})
	})

	srv := &http.Server{Addr: addr, Handler: r}
	logger.Printf("http server listening at %s", addr)

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()
	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()

	err = <-errCh
	if err == http.ErrServerClosed {
		logger.Printf("http server stopped")
		return nil
	}
	return err
}

func requestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		c.Set("requestID", id)
		c.Writer.Header().Set("X-Request-ID", id)
		c.Next()
	}
}

func requestLogger(logger interface{ Printf(string, ...any) }) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		requestIDVal, _ := c.Get("requestID")
		logger.Printf("rid=%v %s %s %d", requestIDVal, c.Request.Method, c.Request.URL.Path, c.Writer.Status())
	}
}
