# 老版本不支持多个

`Docker`17.05以后才支付, 支持多阶段构建

# 老版本为什么不支持多个`FROM`指令

*   `Docker`镜像并非只是一个文件，而是由一堆文件组成，最主要的文件是 层。
*   `Dockerfile`中,大多数指令都会生成一个层

```bash

# 示例一，foo 镜像的Dockerfile
# 基础镜像中已经存在若干个层了
FROM ubuntu:16.04
 
# RUN指令会增加一层，在这一层中，安装了 git 软件
RUN apt-get update \
  && apt-get install -y --no-install-recommends git \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*
```

*   `Docker`镜像每一层只记录文件变更, `Docker`会将镜像的各个层进行计算，最后生成一个文件系统，这个被称为 **联合挂载**。
*   `Docker`各个层是有相关性的, 联合挂载过程中,要求一个`Docker`镜像只能有一个起始层, 只能有一个根
*   所以`Dockerfile`中只允许一个`FROM`指令,因为多个`FROM`指令会造成多根,则是无法实现的

# 多个`FROM`指定意义

多阶段构建, 每条`FROM`指令都是一个构建阶段, 最后生成的镜像只会是最后一个阶段的结果

**能够将前置阶段中文件拷贝到后边阶段中,这就是多阶段的意义**

例子: 运用golang镜像构建出exe,再使用空镜像执行exe

```bash
# 编译阶段
FROM golang:1.10.3
 
COPY server.go /build/
 
WORKDIR /build
 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOARM=6 go build -ldflags '-w -s' -o server
 
# 运行阶段
FROM scratch
 
# 从编译阶段的中拷贝编译结果到当前镜像中
# --from=0 参数, 0代表第一个阶段
COPY --from=0 /build/server /
 
ENTRYPOINT ["/server"]

```

或者

```bash
# 编译阶段 命名为 builder
FROM golang:1.10.3 as builder
 
# ... 省略
 
# 运行阶段
FROM scratch
 
# 从编译阶段的中拷贝编译结果到当前镜像中
COPY --from=builder /build/server /
```

`copy --from `还可以直接从一个已存在的镜像中拷贝

```bash
FROM ubuntu:16.04
 
COPY --from=quay.io/coreos/etcd:v3.3.9 /usr/local/bin/etcd /usr/local/bin/
```

