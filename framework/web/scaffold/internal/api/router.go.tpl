package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Clients struct {
	User any
	Msg  any
}

func NewRouter(ctx context.Context, c *Clients) (*gin.Engine, error) {
	_ = ctx
	_ = c

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// middleware order: logger -> metrics -> recover -> cors -> trace -> auth -> admin
	// r.Use(...)

	r.GET("/healthz", func(g *gin.Context) {
		g.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": "ok"})
	})

	user := r.Group("/user")
	{
		user.POST("/get_users_info", func(g *gin.Context) {
			g.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": gin.H{"ok": true}})
		})
	}
	return r, nil
}
