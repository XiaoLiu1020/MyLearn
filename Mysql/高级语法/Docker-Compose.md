# `docker-compose`简介
`docker-compose`是 `docker`官方的容器编排工具, 可以让用户编写一个简单的`yml`文件达到部署容器的应用集群, 实现快速编排

`DockerFile`仅仅只是达到开启一个服务镜像而且, `compose`可以使用多个容器配合完成任务

它允许用户通过一个单独的`docker-compose.yml`模板文件`（YAML格式）`来定义一组相关联的应用容器作为一个项目`（project）`

`Compose`中有两个重要概念:
* 服务`(service)`: 一个应用容器, 实际上可以包含若干运行相同镜像的容器实例
* 项目`(project)`: 由一组关联的应用容器组成的一个完成业务单元, 在`docker-compose.yml`文件中定义。

# `docker-compose`安装
`Compose`项目使用`python`编写, 可以通过`python`的`pip`安装

```bash
[root@localhost ~]# yum install -y python-pip
[root@localhost ~]# pip install -U docker-compose
[root@localhost ~]# docker-compose version
```

`win`安装略

# `Compose`命令
```bash
Usage:
  docker-compose [-f <arg>...] [options] [COMMAND] [ARGS...]
  docker-compose -h|--help

Options:
  -f, --file FILE             指定使用的Compose模板文件，默认为docker-compose.yml,可多次指定；                        
  -p, --project-name NAME     指定项目名称，默认将使用所在目录名称作为项目名 ；                           
  --verbose                   输出更多调试信息；

  -v, --version               打印版本信息；

Commands:
  build              构建项目中的服务容器
  help               获得一个命令的帮助
  images             列出所有镜像
  kill               通过发送SIGKILL信号来强制停止服务容器
  logs               查看服务器容器的输出
  pause              暂停一个服务容器
  port               打印某个容器的端口所映射的公共端口
  ps                 列出项目中目前的所有容器
  pull               拉取服务依赖的镜像
  push               推送服务依赖的镜像
  restart            重启项目中的服务
  rm                 删除所有的服务器容器（停止状态中的）
  run                在指定服务上执行一个命令
  scale              设置指定服务运行的容器个数
  start              启动已经存在的服务容器
  stop               停止已经处于运行状态的容器，但不删除它
  top                展示运行的进程
  unpause            恢复处于暂停状态中的服务
  up                 自动完成包括构建镜像、创建服务、启动服务并关联服务相关容器的一系列操作
  version            打印docker-compose的版本信息 
```

# `Compose`模板文件
默认模板文件名称: `docker-compose.yml`, 格式为 `yaml`

举例:
```yaml
version: "2"
service:
   webapp: 
      image: examplses/web
      ports:
        - "80:80"
      volumes:
        - "/data"
        
```
每个服务都必须通过`image`制定镜像或者`build`命令(需要`Dockerfile`), 如果使用了`build`, 在`Dockerfile`中设置的选项(例如`CMD, EXPOSE, ENV`等)将自动获取,不必再`docker-compse`再次设置

# 以下为模板的主要指令和功能
## `build`指令
指定`Dockerfile`所在文件夹路径(绝对或者相对`docker-compose.yml`文件路径)

```bash
build: /path/to/build/dir
```
## `command`
覆盖容器启动后默认执行命令

## `container_name`
指定容器名称, 默认使用 **项目名称_服务名称_序号** 格式

## `device`
指定设备映射关系
```bash
devices:
    - "/dev/ttyUSB1:/dev/ttyUSB0"
```
## `env_file`
从文件中获取环境变量, 可以单独文件路径或者列表
```yaml
env_file: .env
env_file:
    - ./common.env
    - ./apps/web.env
    - ./opt/secrets.env
环境变量文件中每一行都必须符合格式，支持#开头的注释行
```
如果有变量名称和`environment`指令冲突, 按照惯例,后者为准

## `environment`
设置环境变量, 可以使用数组或者字典两种格式
```bash
environment:
    RACK_ENV: development
    SESSION_SECRET
或者：
environment:
    - RACK_ENV=development
    - SESSION_SECRET
```
## `expose`
暴露端口, 但不映射到宿主机, 仅可以指定内部端口为参数
```bash
expose:
   - "3000"
   - "8000"
```

## `extends`
基于其他模板文件进行扩展, 例如已经有`webapp`服务, 定义一个基础模板文件为`common.yml`, 相当于继承
```yaml
# common.yml
webapp:
    build: ./webapp
    environment:
        - DEBUG=false
        - SEND_EMAILS=false
```
使用`common.yml`进行扩展:
```yaml
#development.yml
web:
    extends:
        file: common.yml
        service: webapp
    ports:
        - "8000:8000"
    links:
        - db
    environment:
        - DEBUG=true
db:
    image: postgres
```
* `development.yml`会自动继承`common.yml`中`webapp`服务及环境变量定义
* 应该避免出现循环依赖
* `extends`不会继承`links`和`volume_from`中定义的容器和数据卷资源

## `links`
链接到其他服务中的容器, 使用服务名称(同时作为别名)
```yaml
links:
    - db
    - db:database
    - redis
使用的别名会将自动在服务容器中的/etc/hosts里创建。例如：
172.17.2.186 db
172.17.2.186 database
172.17.2.187 redis
所连接容器中相应的环境变量也将创建
```


## `external_link`
链接到`docker-compose.yml`外部容器, 可以是非`compose`管理的外部容器, 参数格式和`links`参数类似
```yaml
external_links:
    - redis_1
    - project_db_1:mysql
    - project_db_1:postgresql
```
> 在使用Docker过程中，会有许多单独使用 docker run 启动的容器的情况，为了使 `Compose` 能够连接这些不在`docker-compose.yml` 配置文件中定义的容器，那么就需要一个特殊的标签，就是 `external_links`，它可以让`Compose` 项目里面的容器连接到那些项目配置外部的容器（前提是外部容器中必须至少有一个容器是连接到与项目内的服务的同一个网络里面）


## `volumes`
数据卷挂载路径设置
```yaml
volumes:
    - /var/lib/mysql
    - cache/:/tmp/cache
    - ~/configs:/etc/configs/:ro  # 只读
```

```yaml
version: "3.2"
services:
  web:
    image: nginx:alpine
    volumes:
      - type: volume
        source: mydata
        target: /data
        volume:
          nocopy: true
      - type: bind
        source: ./static
        target: /opt/app/static

  db:
    image: postgres:latest
    volumes:
      - "/var/run/postgres/postgres.sock:/var/run/postgres/postgres.sock"
      - "dbdata:/var/lib/postgresql/data"

volumes:
  mydata:
  dbdata:

```
### `LONG`语法
```yaml
type：安装类型，可以为 volume、bind 或 tmpfs
source：安装源，主机上用于绑定安装的路径或定义在顶级 volumes密钥中卷的名称 ,不适用于 tmpfs 类型安装。
target：卷安装在容器中的路径
read_only：标志将卷设置为只读
bind：配置额外的绑定选项
propagation：用于绑定的传播模式
volume：配置额外的音量选项
nocopy：创建卷时禁止从容器复制数据的标志
tmpfs：配置额外的 tmpfs 选项
size：tmpfs 的大小，以字节为单位

version: "3.2"
services:
  web:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - type: volume
        source: mydata
        target: /data
        volume:
          nocopy: true
      - type: bind
        source: ./static
        target: /opt/app/static

networks:
  webnet:

volumes:
  mydata:

```


## `ports`
暴露端口信息
```yaml
ports:
    - "3000"
    - "8000:8000"
    - "49100:22"
    - "127.0.0.1:8081:8081"
```


## `net`
设置网络模式, 参数类似`docker client`的`-net`参数
```yaml
net: "bridge"
net: "none"
net: "container:[name or id]"
net: "host"
```

## `deploy` 指定与部署和运行服务相关配置
```yaml
version: '3'
services:
  redis:
    image: redis:alpine
    deploy:
      replicas: 6
      update_config:
        parallelism: 2
        delay: 10s
      restart_policy:
        condition: on-failure
```
其中有很多子项   ,需要百度了

## `entrypoint`
用于指定接入点, 容器开启后一定会执行的命令
```bash
entrypoint 也可以是一个列表，方法类似于 dockerfile

entrypoint:
    - php
    - -d
    - zend_extension=/usr/local/lib/php/extensions/no-debug-non-zts-20100525/xdebug.so
    - -d
    - memory_limit=-1
    - vendor/bin/phpunit
```

## `healthcheck`
用于检查测试服务使用的容器是否正常
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost"]
  interval: 1m30s
  timeout: 10s
  retries: 3
  start_period: 40s
```