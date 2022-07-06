# go-rustdesk-server

[<a href="README-English.md">English</a>] | [<a href="README.md">中文</a>]

rustdesk远程桌面软件，服务端golang实现。

参考[官方实现](https://github.com/rustdesk/rustdesk-server)

更多功能开发中~

# 已实现功能

- 走中继的连接
- 局域网的连接
- 安全连接
- 中继的安全连接

# 配置详解
server - id注册服务

relay - 在无法穿透情况下使用的中继服务

- whiteList 是否启用白名单模式，false为黑名单
- ipList ip名单列表，黑名单模式下在内的ip无法连接
- debug 开发模式，为true会输出debug日志
- reg_server relay注册时服务端地址填写公网地址 仅relay配置
- relay_name relay名称，不为空时会启动relay服务 仅relay配置
- server_port 服务端启动端口 仅server配置
- reg_port 服务端启动的relay注册监听端口 仅server配置
- must_key 必须带key才能连接
# docker-compose安装

下载仓库中的`docker-compose.yml`、`config.json`
修改`config.json`

执行`docker-compose up -d`即可

请开放对应端口且最好使用默认端口。

使用`docker-compose logs`查看生成的key用于加密连接

## 如需只启动relay

在`docker-compose.yml`修改启动参数

`command: /app/go_rustdesk_server -server=false`

并确保relay的配置有值

## 如需只启动server

去掉`config.json`中relay的配置的值

# 更多功能 后续开发

- webapi？
- websocket连接

