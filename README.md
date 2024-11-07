# kenaito-dns

DNS服务器，可通过Web管理界面随意设置灵活的解析规则。为了纯血自研devops平台而生。

# 环境依赖

- gcc
- go version >= 1.20

# Go代理地址配置

[去看看](https://blog.odboy.cn/go%E5%85%A8%E5%B1%80%E9%85%8D%E7%BD%AE%E5%9B%BD%E5%86%85%E6%BA%90-by-odboy/)

# nslookup 命令不存在解决

```shell
yum install bind-utils -y
```

# nslookup指定dns服务器查询

```shell
# 这里dns服务器为 192.168.1.103
nslookup example.com 192.168.1.103
```

# 本程序所用依赖，感谢开源者的无私奉献

- [数据库操作](http://xorm.topgoer.com/)
- [DNS解析](https://github.com/miekg/dns)
- [web](https://gin-gonic.com/zh-cn/docs/quickstart/)

# window安装gcc(记得配置环境变量哦, 记得重启电脑哦)

- [去看看](https://github.com/niXman/mingw-builds-binaries/releases)
- [去下载](https://github.com/niXman/mingw-builds-binaries/releases/download/14.2.0-rt_v12-rev0/x86_64-14.2.0-release-posix-seh-msvcrt-rt_v12-rev0.7z)
- [国内下载](https://oss.odboy.cn/blog/files/windows-gcc/x86_64-14.2.0-release-posix-seh-msvcrt-rt_v12-rev0.7z)

# window验证gcc

```shell
gcc -v
```

# sql转各种在线工具

[去看看](https://gotool.top/)