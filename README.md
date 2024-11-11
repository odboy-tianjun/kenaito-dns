# kenaito-dns

## 背景

Bind9不能直接支持API的方式添加解析记录， 通过脚本修改Bind服务器配置这种事情实在是太冒险了，而且没有发现有开源的、性能嘎嘎好的DNS服务器项目。

## 简介

一个轻量级 DNS 服务器，让变更解析记录简单、优雅！为了纯血自研devops平台而生。

## 环境依赖

- gcc
- go version >= 1.20

## 接口文档

[在线阅读](https://oss.odboy.cn/blog/files/onlinedoc/kenaito-dns/index.html)

## 项目结构

- constant 常量
- controller api接口
- core dns解析
- dao 数据库交互
- domain 各种领域模型
- util 工具函数

## 项目耗时

``
more than 7 hours
``

## 主要特性

- 纳秒级、毫秒级( <= 5 )响应时间
- 支持API变更解析记录
- 支持解析记录回滚
- 支持A、AAAA、MX、TXT、CNAME记录解析

## 待办清单

- 支持回滚 2024-11-08 [ok]
- 添加缓存 2024-11-09 [ok]

## 运行配置

#### Go代理地址配置

[去看看](https://blog.odboy.cn/go%E5%85%A8%E5%B1%80%E9%85%8D%E7%BD%AE%E5%9B%BD%E5%86%85%E6%BA%90-by-odboy/)

#### window安装gcc(记得配置环境变量哦, 记得重启电脑哦)

- [去看看](https://github.com/niXman/mingw-builds-binaries/releases)
- [去下载](https://github.com/niXman/mingw-builds-binaries/releases/download/14.2.0-rt_v12-rev0/x86_64-14.2.0-release-posix-seh-msvcrt-rt_v12-rev0.7z)
- [国内下载](https://oss.odboy.cn/blog/files/windows-gcc/x86_64-14.2.0-release-posix-seh-msvcrt-rt_v12-rev0.7z)

#### window验证gcc

```shell
gcc -v
```

#### 编译

- GOOS代表程序构建环境的目标操作系统，其值可以是liunx，windows，freebsd，darwin
- GORACH代表程序构建环境的目标计算架构，其值可以是386，amd64或arm

```text
首先，配置代理，并在项目根目录执行 go mod tidy 安装依赖
```

```shell
# Windows
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
go build -o ./bin/kenaito-dns_windows_amd64 main.go
```

```shell
# Linux
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o ./bin/kenaito-dns_linux_amd64 main.go
```

```shell
# Mac
set GOOS=darwin
set GOARCH=amd64
set CGO_ENABLED=0
go build -o ./bin/kenaito-dns_darwin_amd64 main.go
```

#### 运行

```shell
# 例子：在 Mac 平台上
# 将编译产出的 kenaito-dns_darwin_amd64 与 dns.sqlite3 文件放在同一目录下， 执行以下命令运行即可
./kenaito-dns_darwin_amd64
```

![jietu1](https://oss.odboy.cn/blog/files/onlinedoc/kenaito-dns/jietu1.png)

#### 测试

```text
DNS服务所在的ip地址为 192.168.43.130

所需工具：brew install watch

测试命令：watch -n 2 nslookup demo2024.odboy.cn 192.168.43.130
```

- 新增A解析记录 [演示视频](https://oss.odboy.cn/blog/files/onlinedoc/kenaito-dns/AddRR.mp4)
- 删除A解析记录 [演示视频](https://oss.odboy.cn/blog/files/onlinedoc/kenaito-dns/RemoveRR.mp4)
- 修改A解析记录 [演示视频](https://oss.odboy.cn/blog/files/onlinedoc/kenaito-dns/ModifyRR.mp4)
- 回滚解析记录 [演示视频](https://oss.odboy.cn/blog/files/onlinedoc/kenaito-dns/RollbackRR.mp4)

## 常见问题

#### nslookup 命令不存在解决

```shell
yum install bind-utils -y
```

#### nslookup指定dns服务器查询

```shell
# 这里dns服务器为 192.168.1.103
nslookup example.com 192.168.1.103
```

## 特别鸣谢

- [数据库操作 - xorm](http://xorm.topgoer.com/)
- [DNS解析 - miekg/dns](https://github.com/miekg/dns)
- [Web - gin](https://gin-gonic.com/zh-cn/docs/quickstart/)

## 代码托管（以私人仓库Gitea为准）

- Gitea: [https://gitea.odboy.cn/odboy/kenaito-dns](https://gitea.odboy.cn/odboy/kenaito-dns)
- Github: [https://github.com/odboy-tianjun/kenaito-dns](https://github.com/odboy-tianjun/kenaito-dns)
- Gitee(已关闭，单纯的不想放在gitee): [https://gitee.com/odboy/kenaito-dns](https://gitee.com/odboy/kenaito-dns)

## 微信交流群

![wxcode](https://oss.odboy.cn/blog/files/userinfo/MyWxCode.png)

(扫码添加微信，备注：kenaito-dns，邀您加入群聊)

加入群聊的好处：

- 第一时间收到项目更新通知。
- 第一时间收到项目 bug 通知。
- 第一时间收到新增开源案例通知。
- 和众多大佬一起互相 (huá shuǐ) 交流 (mō yú)。