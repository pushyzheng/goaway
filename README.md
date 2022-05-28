## goaway

ga~~te~~way -> g**o**away

## Introduction

Goaway 是一个轻量级的网关，可以用来保护你的网站。

功能：
1. 无需任何的存储引擎，只需要通过 yaml 文件配置
2. 支持 Web 和 File 两种 Server
3. 身份验证和权限
4. Prometheus 监控

## Quick start

首先在根目录下创建 conf.yaml，指定 Goaway 监听的端口以及你的域名：

```yaml
server:
  port: 3000
  domain: example.com

accounts:
  admin:
    enable: true
    is-admin: true
    password: admin

applications:
  example:
    enable: true
    server-type: web
    port: 5000
```

运行 goaway server：

```shell
$ go build .
$ ./goaway -env prod
```

运行测试的服务器：

```shell
$ go build example-server/web/app.go
$ ./app -p 5000
```

配置 Nginx：

```nginx
server {
    listen 80;
    server_name localhost;

    location / {
        proxy_set_header APPLICATION_NAME 'example';
        proxy_pass http://127.0.0.1:9000;
    }
}
```

访问： http://localhost:80