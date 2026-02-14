# Framework

已按你的要求新增“干净版框架”：`framework/clean`。

> 目标：不依赖本项目既有业务域模型，提供完整可运行的 `http + socket + rpc` 工程骨架，并附最小 demo 业务代码。

## 入口

- `framework/clean/README.md`
- `framework/clean/cmd/app/main.go`

## 一键运行

```bash
go run ./framework/clean/cmd/app
```

然后可以直接验证：

- `GET /api/v1/hello?name=codex`
- `POST /api/v1/echo`
- `ws://127.0.0.1:18081/ws`
