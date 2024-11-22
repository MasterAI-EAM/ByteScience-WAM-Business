# ByteScience-WAM-Business

## 项目背景
`ByteScience-WAM-Business` 是一个用于支持业务系统管理的服务，旨在提供全面的业务功能。该项目采用 Go 语言开发，并与其他服务协作，构建了高效的服务体系。其技术栈主要依赖 MySQL 和 Redis 数据库。
## 安装依赖
本项目依赖以下服务：
- **MySQL**: 用于存储系统的持久化数据。
- **Redis**: 用于缓存数据和管理会话等。

### 环境要求
- Go 1.18 及以上
- MySQL 8.0 或更高版本
- Redis 7.2.4 或更高版本

### 安装 MySQL 和 Redis
确保你已经安装并配置了 MySQL 和 Redis。如果没有安装，可以参考以下链接进行安装：
- [MySQL 安装教程](https://dev.mysql.com/doc/refman/8.0/en/installing.html)
- [Redis 安装教程](https://redis.io/docs/getting-started/)


## 服务启动
* 安装依赖
```azure
    go get -u
```
* 启动服务
```azure
    go run main.go
```
