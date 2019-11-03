FROM golang AS go-builder

ENV GO111MODULE=auto CGO_ENABLED=0 GOPROXY=https://goproxy.cn

WORKDIR /mygo
COPY .  /mygo

RUN go build -o go-demo


FROM alpine:latest
WORKDIR /mygo

COPY --from=go-builder /mygo/go-demo .

CMD ["./go-demo"]
