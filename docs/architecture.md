# 系统架构

这个项目是一个拆开的短视频后端。客户端统一从 `api` 进入，HTTP 请求走 Hertz，服务之间走 Kitex RPC，聊天连接走 WebSocket。

```text
Client (HTTP / WebSocket)
          |
          v
API Gateway (Hertz :8888)
  - JWT
  - Router / Handler
  - RPC Client
  - WebSocket Hub
          |
          v
+-------------+  +-------------+  +-------------+
| user        |  | video       |  | social      |
+-------------+  +-------------+  +-------------+
       \              |               /
        \             |              /
         +------------+-------------+
                      |
          +-----------+-----------+
          |                       |
          v                       v
   +-------------+         +-------------+
   | interaction |         | chat        |
   +-------------+         +-------------+
          \                       /
           \                     /
            +-------------------+
                    |
                    v
MySQL / Redis / Kafka / etcd / Cloudflare R2
```

## 服务分工

- `api`：统一入口，负责鉴权、路由、RPC 转发和 WebSocket 连接管理。
- `user`：注册、登录、MFA、用户资料、头像上传。
- `video`：视频投稿、查询、列表、点赞数聚合。
- `social`：关注、粉丝、好友关系。
- `interaction`：点赞和评论，点赞事件会异步写 Kafka。
- `chat`：私聊、群聊，消息先写 Redis，再批量落 MySQL。

## 基础组件

- `MySQL`：存业务数据。
- `Redis`：做缓存，也临时保存聊天消息和 MFA secret。
- `Kafka`：只在点赞这条链路上用来异步落库。
- `etcd`：做服务注册和发现。
- `Cloudflare R2`：存头像和视频文件。
