FROM golang:1.18-alpine AS build-env
MAINTAINER DanBai

#修改镜像源为国内
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && apk update
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go env -w GOPATH="/go"
#安装所需工具
RUN apk add gcc g++ make upx git
#配置时区为中国
RUN apk add tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

#拉取代码
RUN mkdir /build
ADD ./ /build
#构建
WORKDIR /build
RUN go build -ldflags '-w -s' -o go_rustdesk_server
RUN upx go_rustdesk_server


FROM alpine:latest
#运行环境
LABEL maintainer="danbai@88.com"
LABEL description="go-rustdesk-server build image file"
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && apk update
RUN apk --no-cache add tzdata ca-certificates libc6-compat libgcc libstdc++ apache2-utils vim
#时区
ENV TZ=Asia/Shanghai

#配置时区为中国
RUN apk add tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

RUN mkdir /app
WORKDIR /app
COPY --from=build-env /build/search_trace /app/go_rustdesk_server
RUN chmod +x /app/go_rustdesk_server
CMD ["/app/go_rustdesk_server"]