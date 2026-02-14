# clean-framework: 干净的 HTTP + Socket + RPC 最小框架

这是从业务代码中抽离出的**独立开发骨架**，目标是：

- 不依赖 OpenIM 现有业务域模型
- 保留工程能力：配置、日志、错误、中间件、分层
- 内置最小 demo 业务（Hello/Echo/Publish）

## 架构

- HTTP: Gin，提供 REST 接口
- Socket: Gorilla WebSocket，提供实时连接
- RPC: `net/rpc`，提供服务间调用（可替换为 gRPC）
- Broker: 内存发布订阅，演示消息中间件角色
- Store: 内存存储，演示数据库抽象层

## 目录

```text
framework/clean/
  cmd/app/main.go                # 一键启动 http + socket + rpc
  internal/httpapi/              # HTTP 层 + store
  internal/socket/               # Socket 层
  internal/rpc/                  # RPC 层
  pkg/common/                    # 配置、日志、错误、broker
```

## 快速启动

```bash
go run ./framework/clean/cmd/app
```

默认端口：

- HTTP: `:18080`
- Socket: `:18081`（ws 路径 `/ws`）
- RPC: `:18082`

## Demo 业务

1. HTTP 调 RPC：
   - `GET /api/v1/hello?name=codex`
2. HTTP Echo：
   - `POST /api/v1/echo`，body: `{"message":"hi"}`
3. HTTP Publish（入库 + 广播）：
   - `POST /api/v1/publish`，body: `{"message":"hello-room"}`
   - `GET /api/v1/messages` 查看已保存消息
4. Socket Echo + 广播：
   - 连接 `ws://127.0.0.1:18081/ws`
   - 发送任意文本会原样返回（echo）
   - 调用 publish 后，socket 会收到 `broadcast:<message>`

## 说明

- 这是一套“干净框架 + 最小业务闭环”的起点。
- 后续把 store/broker 替换成 MySQL/Redis/Kafka 即可演进到生产实现。
