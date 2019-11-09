FROM golang:1.13.4 AS go-builder

# 设置环境变量和解决中文乱码问题
# 禁用CGO,开启go mod机制
ENV LANG="zh_CN.UTF-8"  \ 
    GO111MODULE=on CGO_ENABLED=0 \
    GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct

WORKDIR /mygo
COPY .  /mygo

RUN go build -o go-demo


FROM alpine:3.10

# 解决时区和包依赖，http x509证书问题
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo 'Asia/Shanghai' >/etc/timezone \
    && echo "export LC_ALL=zh_CN.UTF-8"  >>  /etc/profile \
    && echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.10/main/" > /etc/apk/repositories \
    && apk update \
    && apk upgrade \
    && apk add --no-cache ca-certificates bash vim \
    bash-doc \
    bash-completion curl \
    && rm -rf /var/cache/apk/*

WORKDIR /mygo

COPY --from=go-builder /mygo/go-demo .

CMD ["./go-demo"]
