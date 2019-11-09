# 构建镜像

    $ docker build -t go-demo:v1 .

# 运行

    $ docker run -it -v /etc/localtime:/etc/localtime :ro go-demo:v1
