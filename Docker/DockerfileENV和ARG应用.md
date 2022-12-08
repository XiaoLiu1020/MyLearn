`ENV`,`ARG`,`.env`都是关于变量应用,很容易混淆

# `.env`文件

和`docker-compose.yml`配合使用, `docker-compose.yml`默认使用`.env`文件里面的变量

# `env_file`

在`Dockefile`中使用,环境变量较多时候,可使用该参数,**指定对应变量文件**

# `ARG`

在`Dockerfile`中使用, 仅仅在`build docker image`过程中(包括`CMD`和`ENTRYPOINT`)有效,既构建镜像时候, `image`镜像被创建后,无效

*   如果在`Dockerfile`中使用了`ARG`但未给定初始值, 运行`docker build`时候也没指定该`ARG`变量, 则会失败

设置ARG和使用ARG编译image, 实例如下：

```bash
# In the Dockerfile
ARG some_variable_name
# or with a hard-coded default:
#ARG some_variable_name=default_value
 
RUN echo "Oh dang look at that $some_variable_name"
 
# In the shell command
docker build --build-arg some_variable_name=a_value
 
# Then you will get
Oh dang look at that a_value
```

# `ENV`

在`Dockerfile`中使用,在`image`被创建和`container`启动后作为环境变量依然有效,并且可以重写覆盖, `printenv`可查看值

设置ENV和使用env，实例如下

    # no default value
    ENV hey
    # a default value
    ENV foo /bar
    # or ENV foo=/bar
     
    # ENV values can be used during the build
    ADD . $foo
    # or ADD . ${foo}
    # translates to: ADD . /bar
     
    # Use the following docker commands to set env
     
    docker run -e "env_var_name=another_value" alpine env
    docker run -e env_var_name alpine env
    docker run --env-file=env_file_name alpine env

有时候，`ARG`和`ENV`一起使用，实例如下图：

    # expect a build-time variable
    ARG A_VARIABLE
    # use the value to set the ENV var default
    ENV an_env_var=$A_VARIABLE
    # if not overridden, that value of an_env_var will be available to your containers!

