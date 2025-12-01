# GolangLab5

基于 Go 语言的微服务短视频平台后端系统，采用 Kitex RPC 框架构建。

## 技术栈

| 组件     | 技术                                        |
| -------- | ------------------------------------------- |
| Web 框架 | [Hertz](https://github.com/cloudwego/hertz) |
| RPC 框架 | [Kitex](https://github.com/cloudwego/kitex) |
| 数据库   | MySQL                                       |
| 缓存     | Redis                                       |
| 消息队列 | Kafka                                       |
| 服务注册 | etcd                                        |

## 快速开始

### 安装运行

```bash
# 克隆项目
git clone https://github.com/nnieie/GolangLab5.git
cd GolangLab5

# 修改配置
vim config/config.yaml

# 初始化数据库
mysql -u root -p < config/sql/init.sql

# 启动服务
./start.sh
```

## 项目结构

```
.
├── cmd                     # 各个可执行服务的入口
│   ├── api                 # API 网关 / HTTP 服务
│   │   ├── biz             # API 层业务代码与路由
|   |   |   ├── handler
│   │   │   ├── model
│   │   │   └── router
│   │   ├── pack            # 响应/请求的打包、转换工具
│   │   ├── rpc             # 对外/内部 RPC 客户端封装
│   │   ├── script          # 启动/部署相关脚本
│   │   └── ws              # WebSocket 相关处理
│   ├── chat                # 聊天服务子项目
│   │   ├── dal             # 数据访问层
│   │   ├── pack            # 消息封装与转换
│   │   ├── rpc             # 对外 RPC 封装
│   │   ├── script          # 启动脚本
│   │   └── service         # 聊天业务实现
│   ├── interaction         # 点赞/评论等交互服务
│   ├── social              # 社交相关服务
│   ├── user                # 用户服务
│   └── video               # 视频服务
├── config                  # 配置文件和类型定义
├── docs                    # 项目文档与架构说明
├── idl                     # Thrift 接口定义
├── kitex_gen               # Kitex 生成的代码
└── pkg                     # 公共库/工具包
    ├── constants           # 常量定义
    ├── errno               # 错误码定义与处理
    ├── kafka               # Kafka 生产/消费封装
    ├── logger              # 日志封装
    ├── oss                 # 对象存储工具
    └── utils               # 通用工具
```

## 文档

[架构设计](./docs/architecture.md)

[缓存流程](./docs/cache-flowcharts.md)

[消息队列](./docs/kafka-flowcharts.md)
