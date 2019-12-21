FROM golang:1.13.5 AS go-builder

# 设置golang环境变量和禁用CGO,开启go mod机制
ENV  GO111MODULE=on CGO_ENABLED=0 \
    GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct

WORKDIR /mygo
COPY .  /mygo

RUN go build -o go-demo


FROM alpine:3.10

#解决docker时区问题和中文乱码问题
ENV TZ=Asia/Shanghai LANG="zh_CN.UTF-8" 

# 解决http x509证书问题，需要安装证书
RUN echo "export LC_ALL=$LANG"  >>  /etc/profile \
    && echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.10/main/" > /etc/apk/repositories \
    && apk update \
    && apk upgrade \
    && apk --no-cache add tzdata ca-certificates bash vim bash-doc bash-completion curl \
    && ln -snf  /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone \
    && rm -rf /var/cache/apk/* /tmp/* /var/tmp/* $HOME/.cache

WORKDIR /mygo

EXPOSE 1338 2338

COPY --from=go-builder /mygo/go-demo .

CMD ["./go-demo"]
