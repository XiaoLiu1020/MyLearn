# 背景

`docker` 容器之间都是互相隔离的，不能互相访问，下面介绍三种方法解决容器互访问题：

# 方式一　直接使用虚拟ip访问

安装`Docker`时，`docker`会默认创建一个内部桥接网络`docker0`,　每创建一个容器分配一个虚拟网卡，容器直接可以跟据`ip`互相访问

```bash
ifconfig     # 直接查看

# 在容器内部ping的别的容器的虚拟网卡
ping
```

# 方式二　\`run --link\`\`参数

`--link`格式： `--link <name or id>:alias`

其中，`name`和`id`是源容器的`name`和`id`，　`alias`是源容器在`link`下的别名

比如：

**源容器**

```bash
docker run -d --name selenium_hub selenium/hub

创建并启动名为selenium_hub的容器
```

**接收容器**

```bash
docker run -d --name node --link selenium_hub:hub selenium/node-chrome-debug

创建并启动名为node的容器，并把该容器和名为selenium_hub的容器连接起来

--link selenium_hub:hub
```

> 　　站在node容器角度，　selenium\_hub 和　hub 都是　selenium\_hub.image\_id容器的名称，　并且作为容器的hostname, node中都可以用这名字与之通信，　docker通过DNS自动解析

```bash
docker exec -it node /bin/bash

ping selenium_hub

ping hub
```

## `--link`下容器间的通信

源容器和接收容器之间传递数据通过以下两种方式：

*   设置环境变量
*   更新`/etc/hosts`文件

### 设置环境变量

`docker`会在接收容器中设置名为　`<alias>_NAME`的环境变量，该环境变量的值为: `<alias>_NAME=/接收容器名/源容器alias`

另外，`docker`还会接收容器中创建关于　**源容器暴露的端口号的环境变量**，这些环境变量都有统一的前缀名称：

`<name>PORT<port>_<protocol>`

其中：
`<name>`表示链接的源容器`alias`
`<port>`是源容器暴露的端口号
`<protocol>`是通信协议：　`TCP` or `UDP`

docker用上面定义的前缀定义3个环境变量：

```bash
<name>PORT<port>_<protocol>ADDR
<name>PORT<port><protocol>PORT
<name>PORT<port><protocol>_PROTO

# 查看： env |grep -i <name>PORT<port>_<protocol>_
```

## 更新/etc/hosts

`docker`会将源容器`host`更新到目标容器中的`/etc/hosts/`中

```bash
172.17.0.2  hub 1cbbf6f07804 selenium_hub   # 源容器的alias名称等
172.17.0.3  c4cc05d832e0        # node容器的ip
```

如果重启了源容器，接收容器的`/etc/hosts`会自动更新源容器的新的`ip`

# 方式三， 创建bridge网络

创建`bridge`网络命令：　`docker network create testnet`

运行容器连接到`testnet`网络

```bash
docker run -it --name <容器名> --network <bridge> --network-alias <网络别名> <镜像名>

docker run -it --name centos-1 --network testnet --network-alias centos-1 docker.io/centos:lastest
```

若访问容器中服务,可以使用这种方式访问：　`<网络别名>:<服务端口号>`

推荐这种方法，自定义网络，不用顾虑`ip`是否变动
