# Redis 流程

这个项目里 Redis 主要做三件事：点赞缓存、聊天消息暂存、MFA 临时数据。

## 点赞

```text
点赞 / 取消点赞请求
        |
        v
检查 Redis 里的点赞关系
        |
        v
更新点赞关系和计数
        |
        v
发送 Kafka 事件
        |
        v
返回结果
```

- 点赞关系用 `like:video:{userID}:{videoID}` 和 `like:comment:{userID}:{commentID}`。
- 点赞计数用 `like_count:video:{videoID}` 和 `like_count:comment:{commentID}`。
- MySQL 不在主链路里，数据库更新交给 Kafka 消费端处理。

## 聊天

```text
发送消息
   |
   v
写消息详情到 Redis
   |
   +-- 更新消息列表 ZSet
   |
   +-- 推入待同步队列
   |
   v
SyncWorker 每 5 秒批量刷 MySQL
```

- 私聊消息详情前缀是 `msg:private:`，群聊消息详情前缀是 `msg:group:`。
- 私聊同步队列是 `msg:queue:private`，群聊同步队列是 `msg:queue:group`。
- 消息详情和消息列表都会设置 7 天过期时间，消息列表最多保留 1000 条。
- `SyncWorker` 每批最多同步 200 条，私聊和群聊分别处理。

历史消息读取流程比较直接：

```text
读历史消息
   |
   v
先查 Redis
   |
   +-- 命中 -> 直接返回
   |
   +-- 未命中 -> 查 MySQL
```

## MFA

```text
申请绑定 MFA
   |
   v
生成 secret 和二维码
   |
   v
加密后写 Redis（15 分钟）

提交验证码
   |
   v
从 Redis 取 secret
   |
   +-- 不存在 -> 重新申请二维码
   |
   +-- 存在 -> 校验验证码 -> 绑定成功
```

- 临时 secret 的 key 前缀是 `totp_secret:{userID}`。
- 过期时间是 15 分钟，过期后需要重新走一遍绑定流程。
