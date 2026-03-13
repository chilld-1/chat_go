# GoChat - 实时聊天系统

GoChat 是一个基于 Go 语言开发的实时聊天系统，支持 WebSocket 和 TCP 协议，使用 RabbitMQ 作为消息队列，提供高性能、可靠的实时通信能力。

## 功能特性

- **多协议支持**：同时支持 WebSocket 和 TCP 协议，满足不同场景的需求
- **实时通信**：基于消息队列的异步处理，确保消息实时传递
- **用户认证**：使用 JWT 进行用户认证，确保系统安全
- **消息互传**：WebSocket 和 TCP 客户端之间可以互相发送消息
- **广播机制**：支持消息广播给所有客户端
- **Redis 应用**：使用 Redis 进行用户在线状态管理、消息未读计数等
- **可扩展性**：模块化设计，易于扩展和维护

## 技术栈

- **后端语言**：Go 1.25+
- **数据库**：MySQL
- **缓存**：Redis
- **消息队列**：RabbitMQ
- **Web 框架**：Gin
- **认证**：JWT
- **网络协议**：WebSocket、TCP

## 项目结构

```
gochat/
├── api/              # API 层
│   ├── handler/      # 请求处理器
│   └── router/       # 路由配置
├── config/           # 配置管理
├── connect/          # 连接管理
│   ├── channel.go    # 通道管理
│   ├── websocket.go  # WebSocket 处理
│   └── tcpSever.go   # TCP 服务器
├── db/               # 数据库连接
├── logic/            # 业务逻辑
│   └── dao/          # 数据访问对象
├── model/            # 数据模型
├── tools/            # 工具类
│   ├── jwt.go        # JWT 认证
│   ├── rabbitmq.go   # RabbitMQ 工具
│   └── redis.go      # Redis 工具
├── main.go           # 主文件
├── go.mod            # 依赖管理
└── README.md         # 项目说明
```

## 安装与运行

### 前置条件

- Go 1.25+
- MySQL
- Redis
- RabbitMQ

### 安装步骤

1. **克隆代码**

   ```bash
   git clone <repository-url>
   cd gochat
   ```

2. **安装依赖**

   ```bash
   go mod tidy
   ```

3. **配置文件**

   编辑 `config/config.toml` 文件，设置数据库、Redis 和 RabbitMQ 的连接信息。

4. **运行服务**

   - 启动 API 服务
     ```bash
     go run main.go -module api
     ```

   - 启动 WebSocket 服务
     ```bash
     go run main.go -module connect_websocket
     ```

   - 启动 TCP 服务
     ```bash
     go run main.go -module connect_tcp
     ```

## API 接口

### 用户注册

- **URL**：`/v1/user/register`
- **方法**：POST
- **请求体**：
  ```json
  {
    "username": "testuser",
    "password": "123456"
  }
  ```
- **响应**：
  ```json
  {
    "code": 0,
    "msg": "success",
    "data": null
  }
  ```

### 用户登录

- **URL**：`/v1/user/login`
- **方法**：POST
- **请求体**：
  ```json
  {
    "username": "testuser",
    "password": "123456"
  }
  ```
- **响应**：
  ```json
  {
    "code": 0,
    "msg": "success",
    "data": {
      "token": "<jwt-token>",
      "username": "testuser"
    }
  }
  ```

## WebSocket 连接

- **URL**：`ws://localhost:7000/ws?token=<jwt-token>`
- **消息格式**：
  ```json
  {
    "type": "chat",
    "from": "user1",
    "to": "user2",
    "content": "Hello",
    "timestamp": 1620000000
  }
  ```

## TCP 连接

- **地址**：`localhost:7001`
- **连接流程**：
  1. 建立 TCP 连接
  2. 发送 token 进行认证
  3. 发送和接收消息

## 消息格式

所有消息使用 JSON 格式，包含以下字段：

- `type`：消息类型（chat、notification 等）
- `from`：发送者
- `to`：接收者（为空时广播给所有客户端）
- `content`：消息内容
- `timestamp`：时间戳

## 测试

### WebSocket 测试

使用 `ws_client.go` 进行 WebSocket 客户端测试：

```bash
go run ws_client.go <jwt-token>
```

### TCP 测试

使用 `tcp_client.go` 进行 TCP 客户端测试：

```bash
go run tcp_client.go
```

