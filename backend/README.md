# Grok2API Backend

Grok2API 的 Go 后端，负责上游账号调度、协议转换、额度管理、请求审计和管理 API，并可直接托管前端构建产物。

## 技术栈

- Go 1.26、Gin、GORM
- SQLite / PostgreSQL
- Memory / Redis
- Grok Build OAuth 与 Grok Web SSO Provider

## 本地运行

启动配置默认位于仓库根目录 `config.yaml`。若文件不存在，进程会在首次启动时自动生成默认配置。

凭据加密密钥推荐通过环境变量注入；`jwtSecret` 会由该密钥自动派生，无需单独配置：

```bash
export GROK2API_CREDENTIAL_ENCRYPTION_KEY="$(openssl rand -base64 32)"
cd backend
go run ./cmd/grok2api
```

也可先参考示例手动准备配置：

```bash
cp config.example.yaml config.yaml
export GROK2API_CREDENTIAL_ENCRYPTION_KEY="$(openssl rand -base64 32)"
cd backend
go run ./cmd/grok2api
```

服务默认监听 `http://127.0.0.1:8000`。也可以显式指定配置文件或监听地址：

```bash
go run ./cmd/grok2api --config /path/to/config.yaml --listen 0.0.0.0:8000
```

首次创建管理员时，若未在配置中指定 `bootstrapAdmin`，默认账号为 `admin` / `grok2api`。

## Docker 运行

容器内配置文件固定为 `/app/config.yaml`。若不存在会在启动时自动生成；数据库与媒体目录位于 `/app/data`。

```bash
export GROK2API_CREDENTIAL_ENCRYPTION_KEY="$(openssl rand -base64 32)"

docker run -d \
  --name grok2api \
  --restart unless-stopped \
  -p 8000:8000 \
  -e TZ=Asia/Shanghai \
  -e GROK2API_CREDENTIAL_ENCRYPTION_KEY \
  -v grok2api-data:/app/data \
  ghcr.io/jians1/grok2api:latest
```

Compose 与更多部署说明见仓库根目录 [`README.md`](../README.md)。

## 配置与存储

启动配置由根目录 `config.yaml` 管理，字段说明见 [`config.example.yaml`](../config.example.yaml)。Provider、服务容量、批量任务、路由、媒体、审计和客户端密钥默认限制由管理端设置页持久化；除页面明确标记“重启生效”的字段外均会热加载。

| 场景 | 数据库 | 运行态存储 |
| --- | --- | --- |
| 本地开发 / 单实例 | SQLite | Memory |
| 多实例部署 | PostgreSQL | Redis |

| 变量 / 字段 | 说明 |
| --- | --- |
| `GROK2API_CREDENTIAL_ENCRYPTION_KEY` | 凭据加密主密钥，推荐通过环境变量注入 |
| `secrets.credentialEncryptionKey` | 环境变量未设置时，可改由 YAML 提供 |
| `secrets.jwtSecret` | 由凭据加密密钥自动派生，无需手写 |
| `bootstrapAdmin` | 仅在数据库无管理员时用于创建初始账号 |

关系型数据库保存账号、凭据、模型、额度、客户端密钥、审计和媒体任务；Redis 仅承载限流、并发租约、粘滞路由、分布式锁和事件通知。敏感凭据使用 AES-256-GCM 加密，`credentialEncryptionKey` 必须长期保留且不得提交到版本库。

可热加载的 Provider、服务容量、批量任务并发、路由、审计、媒体和代理参数由管理端设置页维护；数据库驱动、监听地址、Redis 与加密密钥仍通过启动配置生效。

## 服务入口

- `/v1/*`：兼容 API
- `/api/admin/v1/*`：管理 API
- `/healthz`、`/readyz`：健康与就绪探针
- `/swagger/index.html`：公开 API Swagger，仅在 `server.swaggerEnabled: true` 时注册
- `frontend.staticPath`：前端静态目录，默认 `./frontend/dist`

详细协议说明见 [`docs`](./docs)。

修改公开接口注释后，在仓库根目录执行 `make swagger` 更新 `backend/docs/docs.go`、`swagger.json` 与 `swagger.yaml`。生产配置应保持 `server.swaggerEnabled: false`。

## 代码结构

```text
cmd/grok2api/       进程入口
internal/domain/    领域模型与规则
internal/application/ 应用服务与用例
internal/infra/     数据库、Provider、运行态与安全实现
internal/transport/ HTTP 路由、鉴权与协议适配
internal/repository/ 持久化接口
```

依赖方向保持为 Transport → Application → Domain，基础设施通过接口接入，不在领域层依赖 HTTP、数据库或具体 Provider。

## 验证

```bash
go test ./...
go test -race ./...
go vet ./...
go build ./cmd/grok2api
```
