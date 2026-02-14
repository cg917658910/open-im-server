# Socket 工程框架（WebSocket 网关）

## 1. 对标本仓库的参考实现

- 网关主流程：`internal/msggateway/ws_server.go`
- 客户端/连接对象：`internal/msggateway/client.go`、`client_conn.go`
- 消息处理抽象：`internal/msggateway/message_handler.go`
- 网关启动入口：`cmd/openim-msggateway/main.go` + `pkg/common/cmd/msg_gateway.go`

## 2. 工程目录（建议）

```text
cmd/my-gateway/main.go
internal/msggateway/
  ws_server.go
  client.go
  client_conn.go
  dispatcher.go
  auth.go
  online.go
  subscription.go
pkg/common/cmd/my_gateway.go
config/my-gateway.yml
```

## 3. 生产级设计约束（强烈建议）

1. **状态变更单线程化**：连接注册/注销/踢下线走 channel 事件循环。
2. **连接读写分离**：read pump / write pump，避免并发写 websocket。
3. **慢消费者保护**：发送队列满后要丢弃或断连。
4. **多端在线策略可配置**：同平台互踢、全端并存、后台在线。
5. **节点间协同**：多机部署通过 RPC 通知实现跨节点踢线与在线同步。

## 4. 连接生命周期（推荐）

1) HTTP Upgrade  
2) token 校验 + userID/platform 绑定  
3) register 事件  
4) read loop（解码 + 校验 + 分发）  
5) write loop（串行发送）  
6) 心跳超时检测  
7) unregister（在线状态落库 + webhook）

## 5. 可复制模板（核心骨架）

```go
package msggateway

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

type Server struct {
	upgrader   websocket.Upgrader
	register   chan *Client
	unregister chan *Client
	kick       chan *Client
	users      map[string]map[*Client]struct{}
}

func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		register: make(chan *Client, 1024),
		unregister: make(chan *Client, 1024),
		kick: make(chan *Client, 1024),
		users: make(map[string]map[*Client]struct{}),
	}
}

func (s *Server) Run() {
	for {
		select {
		case c := <-s.register:
			if _, ok := s.users[c.UserID]; !ok {
				s.users[c.UserID] = map[*Client]struct{}{}
			}
			s.users[c.UserID][c] = struct{}{}
		case c := <-s.unregister:
			if m, ok := s.users[c.UserID]; ok {
				delete(m, c)
				if len(m) == 0 {
					delete(s.users, c.UserID)
				}
			}
			_ = c.Conn.Close()
		case c := <-s.kick:
			_ = c.Conn.Close()
		}
	}
}
```

## 6. tools / infra 复用点

- 鉴权、用户状态、消息处理都通过 RPC client（而不是网关内硬编码业务）。
- webhook 与日志统一走已有工具库；错误统一包装后落日志。
- 观测必须包含：在线人数、连接数、踢线数、发送队列堆积、消息处理耗时。

## 7. 故障场景压测清单

- [ ] 单机 5w+ 长连接稳定性（fd、内存、goroutine）。
- [ ] Redis/MQ 抖动时的退化策略（仅在线推送/降级丢弃）。
- [ ] 节点重启时连接回收和重连风暴控制。
- [ ] 大群广播时慢消费者隔离。


## 8. 脚手架文件（可直接复制）

已提供：

- `framework/socket/scaffold/cmd/my-gateway/main.go.tpl`
- `framework/socket/scaffold/internal/msggateway/ws_server.go.tpl`
- `framework/socket/scaffold/config/my-gateway.yml.tpl`
