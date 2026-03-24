# Kafka 流程

当前 Kafka 只用在点赞事件这条链路上，目的很简单：接口先把缓存改掉，再慢慢把结果刷进 MySQL。

## 主流程

```text
点赞 / 取消点赞
        |
        v
interaction 更新 Redis
        |
        v
发送 LikeEvent 到 topic: interaction-events
        |
        v
consumer group: interaction_like_group
        |
        v
批量处理并折叠重复状态
        |
        v
写入 MySQL likes 表
```

## 生产端

- 入口在 `interaction` 服务。
- 先更新 Redis 里的点赞关系和计数，再发送 `LikeEvent`。
- 事件里会带 `userID`、`videoID` 或 `commentID`、`action`。
- 接口返回时不等 MySQL 写完，所以主链路会短一些。

## 消费端

- 消费组是 `interaction_like_group`。
- 批量阈值是 500 条，或者每 1 秒刷一次。
- 同一批里如果同一个用户对同一个目标反复操作，只保留最后一次状态。
- 点赞走 `BatchLikeAction`，取消点赞走 `BatchUnlikeAction`。

## 这条链路的作用

- 请求返回更快。
- MySQL 写入可以合并，压力更小。
- 读路径优先看 Redis，点赞状态和计数更新也更直接。
