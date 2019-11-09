# 构建镜像

    $ docker build -t go-demo:v1 .

# 运行

    $ docker run -it go-demo:v1
    2019/11/09 14:28:38 111
    2019/11/09 14:28:38 hello
    2019/11/09 14:28:38 111
    2019/11/09 14:28:39 111

# 验证 docker 中的时区

    $ docker run -it go-demo:v1 /bin/bash
    bash-5.0# date
    Sat Nov  9 13:59:24 CST 2019
    bash-5.0#

# 指定 docker 运行的 name

    $ docker run -it --name=go-demo go-demo:v1

    查看进程
    $ docker exec -it go-demo  /bin/bash
    $ top
