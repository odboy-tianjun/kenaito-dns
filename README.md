# kenaito-dns

DNS服务器，可通过Web管理界面随意设置灵活的解析规则。为了纯血自研devops平台而生。

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