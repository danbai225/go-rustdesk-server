# go-rustdesk-server

[<a href="README-English.md">English</a>] | [<a href="README.md">中文</a>]

rustdesk remote desktop software, server-side golang implementation.

Reference [official implementation](https://github.com/rustdesk/rustdesk-server)

This project does not guarantee availability, please move to the official open source server.

# Implemented features

- Relay connection
- LAN connection
- Secure connection
- Secure connection for trunking

# Configuration details
server - id registration service

relay - relay service to be used in case of impenetration

- whiteList whether to enable whitelist mode, false for blacklist
- ipList ip list, blacklist mode, the inner ip cannot be connected
- debug development mode, true will output debug log
- reg_server relay server address when registering only relay configuration
- relay_name relay name, if it is not empty, the relay service will be started.
- server_port Server start port Only server configuration
- reg_port The port where the relay is registered to listen when the server is started.
- must_key key is required to Connect

# docker-compose installation

Download `docker-compose.yml`, `config.json` from the repository
Modify `config.json`

Execute `docker-compose up -d`.

Please open the corresponding port and preferably use the default port.

Use `docker-compose logs` to view the generated keys for encrypted connections

## If you want to start only relay

Change the start parameters in `docker-compose.yml`

`command: /app/go_rustdesk_server -server=false`

and make sure the relay configuration has the value

## If you want to start only the server

Remove the value of the relay configuration from `config.json`
