# Web 工程框架（Gin + RPC 聚合 + 统一中间件）

## 1. 对标本仓库的参考实现

- 路由与中间件装配：`internal/api/router.go`
- 进程启动入口：`cmd/openim-api/main.go` + `pkg/common/cmd/api.go`

## 2. 工程目录（建议）

```text
cmd/my-api/main.go
internal/api/
  init.go
  router.go
  auth.go user.go group.go msg.go
  middleware/
    auth.go
    trace.go
    ratelimit.go
pkg/common/cmd/my_api.go
config/my-api.yml
```

## 3. 必接能力（生产级）

### 3.1 中间件顺序（建议固定）

1) access log  
2) prommetrics  
3) panic recover  
4) CORS  
5) operationID / traceID 注入  
6) token 解析与鉴权  
7) admin 身份注入  
8) handler

### 3.2 错误返回规范

- handler 内业务错误统一映射到 `servererrs` code。
- 不向客户端暴露底层错误栈（DB/MQ/内部地址等）。
- 请求参数错误全部归一到 `ArgsError` 风格码。

### 3.3 RPC client 注入规范

- 不要在 handler 里 new client。
- 在 router 初始化阶段集中构造并注入（auth/user/group/friend/msg/...）。
- handler 只关心调用接口，不关心连接管理。

## 4. 可复制模板（cmd + router）

### 4.1 `cmd/my-api/main.go`

```go
package main

import (
	"github.com/openimsdk/open-im-server/v3/pkg/common/cmd"
	"github.com/openimsdk/tools/system/program"
)

func main() {
	if err := cmd.NewApiCmd().Exec(); err != nil {
		program.ExitWithError(err)
	}
}
```

### 4.2 `internal/api/router.go`（结构化模板）

```go
package api

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Clients struct {
	User any
	Msg  any
}

func NewRouter(ctx context.Context, c *Clients) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// 建议固定顺序：log -> metrics -> recover -> cors -> trace -> auth -> admin
	// r.Use(...)

	user := r.Group("/user")
	{
		user.POST("/get_users_info", func(g *gin.Context) {
			g.JSON(200, gin.H{"errCode": 0, "data": gin.H{"ok": true}})
		})
	}

	return r, nil
}
```

## 5. 与 tools 库的复用建议

- 日志：统一用 `tools/log`，避免 `fmt.Println`。
- 错误：统一用 `tools/errs` 包装内部错误，再映射 `servererrs`。
- 通用中间件：优先复用 `tools/mw` 体系。

## 6. 落地清单（上线前）

- [ ] 所有接口具备 operationID。
- [ ] 所有接口有 errCode/errMsg 统一结构。
- [ ] 核心接口打 metrics（QPS、P99、错误率）。
- [ ] 鉴权中间件支持白名单接口（登录/注册等）。
- [ ] 限流配置可动态调整（配置中心或热更新）。


## 7. 脚手架文件（可直接复制）

已提供：

- `framework/web/scaffold/cmd/my-api/main.go.tpl`
- `framework/web/scaffold/internal/api/router.go.tpl`
- `framework/web/scaffold/config/my-api.yml.tpl`
