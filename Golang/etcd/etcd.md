- [1.`ETCD`简介与使用](#1etcd简介与使用)
- [3. 架构](#3-架构)
  - [3.1 `etcd`概念词汇](#31-etcd概念词汇)
- [4.`docker` 启动`etcd`单机](#4docker-启动etcd单机)
- [5. 集群化应用实践](#5-集群化应用实践)
  - [5.1 集群启动](#51-集群启动)
  - [5.2 静态配置](#52-静态配置)
  - [5.3 `etcd`自发现模式](#53-etcd自发现模式)
  - [5.4 `DNS`自发现模式](#54-dns自发现模式)
  - [5.5 关键部分源码解析](#55-关键部分源码解析)
  - [5.6 `peer-urls` `client-urls`](#56-peer-urls-client-urls)
  - [5.7 运行时节点变更](#57-运行时节点变更)
- [6 `Proxy`模式](#6-proxy模式)
- [7 数据存储](#7-数据存储)
  - [7.1 `WAL`体系](#71-wal体系)
  - [7.2 预写式日志`（WAL）`](#72-预写式日志wal)
  - [7.3 `WAL`与`snapshot`命名规则](#73-wal与snapshot命名规则)
  - [7.4 源码解析--`read`和`append`模式](#74-源码解析--read和append模式)
- [8. `Raft`](#8-raft)
  - [8.1 常见问题](#81-常见问题)
  - [8.2 关键部分源码](#82-关键部分源码)
- [9. `Store`--提供`API`支持](#9-store--提供api支持)
  - [自动化`CAS`操作](#自动化cas操作)
  - [总结：](#总结)
- [10. `go`操作`etcd`](#10-go操作etcd)
  - [安装](#安装)
    - [1. 修改依赖为v1.26.0](#1-修改依赖为v1260)
  - [`API`` V2`与`V3`区别](#api-v2与v3区别)
  - [`put`和`get`操作](#put和get操作)
  - [`watch`操作](#watch操作)
  - [`lease`租约--定时](#lease租约--定时)
  - [`keepAlive`--续租](#keepalive--续租)
  - [基于`etcd`实现分布式锁](#基于etcd实现分布式锁)
  - [事务\&超时](#事务超时)
  - [服务注册](#服务注册)
  - [服务发现](#服务发现)

# 1.`ETCD`简介与使用

`TODO` 复制在家的那份--包含简介与应用场景

# 3. 架构

<img src="https://upload-images.jianshu.io/upload_images/5525735-078bf1d0626f4a82.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/1200/format/webp" style />

*   `HTTP Server`: 处理用户发送的`API`请求以及其他`etcd`节点的同步与心跳信息请求。
*   `Store`: 处理`etcd`支持的各类功能事务， 是`etcd`对用户提供的大多数`API`功能具体实现
*   `Raft`： 一致性算法实现
*   `WAL`： `Write Ahead Log`预写式日志， 数据储存方式，`WAL`中，所有的数据提交前都会事先记录日志。`Snapshot`是为了防止数据过多而进行的状态快照；`Entry`表示存储的具体日志内容。

**通常，一个用户的请求发送过来，会经由`HTTP Server`转发给Store进行具体的事务处理，如果涉及到节点的修改，则交给`Raft`模块进行状态的变更、日志的记录，然后再同步给别的`etcd`节点以确认数据提交，最后进行数据的提交，再次同步。**

## 3.1 `etcd`概念词汇

*   `Raft`：`etcd`所采用的保证分布式系统强一致性的算法。
*   `Node`：一个`Raft`状态机实例。
*   `Membe`r： 一个`etcd`实例。它管理着一个`Node`，并且可以为客户端请求提供服务。
*   `Cluster`：由多个`Member`构成可以协同工作的`etcd`集群。
*   `Peer`：对同一个`etcd`集群中另外一个`Member`的称呼。
*   `Client`： 向`etcd`集群发送HTTP请求的客户端。
*   `WAL`：预写式日志，`etcd`用于持久化存储的日志格式。
*   `snapshot`：`etcd`防止`WAL`文件过多而设置的快照，存储`etcd`数据状态。
*   `Proxy`：`etcd`的一种模式，为`etcd`集群提供反向代理服务。
*   `Leader`：`Raft`算法中通过竞选而产生的处理所有数据提交的节点。
*   `Follower`：竞选失败的节点作为`Raft`中的从属节点，为算法提供强一致性保证。
*   `Candidate`：当`Follower`超过一定时间接收不到`Leader`的心跳时转变为`Candidate`开始竞选。
*   `Term`：某个节点成为`Leader`到下一次竞选时间，称为一个`Term`。
*   `Index`：数据项编号。`Raft`中通过`Term`和`Index`来定位数据。

# 4.`docker` 启动`etcd`单机

`etcd`官网使用` gcr.io/etcd-development/etcd:`为外国资源， 采用了别人的`zhaowenlei/etcd-development:v3.3.18`

```bash
rm -rf /tmp/etcd-data.tmp && mkdir -p /tmp/etcd-data.tmp && \
  docker rmi zhaowenlei/etcd-development:lastest || true && \
  docker run \
  --rm \
  -p 2379:2379 \
  -p 2380:2380 \
  --mount type=bind,source=/tmp/etcd-data.tmp,destination=/etcd-data \
  -e ETCDCTL_API=3 \		# 设置命令行版本
  --name etcd \
  zhaowenlei/etcd-development:lastest \
  /usr/local/bin/etcd \
  --name s1 \
  --data-dir /etcd-data \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379 \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-advertise-peer-urls http://0.0.0.0:2380 \
  --initial-cluster s1=http://0.0.0.0:2380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new \
  --debug 'info' \
  --logger 'zap' \     			 		# 似乎暂时没有
  --log-output 'stderr'			# 
```

启动成功， 输入测试

`etcdctl`默认使用`ETCDCTL_API=2`

设置环境变量`ETCDCTL_API=3`

```bash
docker exec etcd /bin/sh -c "/usr/local/bin/etcd --version"

docker exec etcd /bin/sh -c "/usr/local/bin/etcdctl version"

docker exec etcd /bin/sh -c "/usr/local/bin/etcdctl  endpoint health"

docker exec etcd /bin/sh -c "/usr/local/bin/etcdctl put foo bar"

docker exec etcd /bin/sh -c "/usr/local/bin/etcdctl get foo"
```

# 5. 集群化应用实践

`etcd`作为一个高可用键值存储系统，天生就是为集群化而设计的。由于`Raft`算法在做决策时需要多数节点的投票，所以`etcd`一般部署集群推荐奇数个节点，推荐的数量为3、5或者7个节点构成一个集群。

## 5.1 集群启动

方式：

*   静态配置启动
*   `etcd`自身服务发现
*   通过`DNS`进行服务发现。

**它摒弃了使用配置文件进行参数配置的做法，转而使用命令行参数或者环境变量的做法来配置参数。**

## 5.2 静态配置

适用于离线环境

启动集群之前就预先清楚要配置的集群大小，以及节点信息

通过`initial-cluster`参数进行集群启动

```bash
ETCD_INITIAL_CLUSTER="infra0=http://10.0.1.10:2380,infra1=http://10.0.1.11:2380,infra2=http://10.0.1.12:2380"
ETCD_INITIAL_CLUSTER_STATE=new
```

值得注意:

`-initial-cluster`参数中配置的`url`地址必须与各个节点启动时设置的`initial-advertise-peer-urls`参数相同。（`initial-advertise-peer-urls`参数表示节点监听其他节点同步信号的地址）

为避免命名意外发生:

`-initial-cluster-token`参数为每个集群单独配置一个token认证。这样就可以确保每个集群和集群的成员都拥有独特的ID。

三个`etcd`集群，分别使用命令

```bash
# 10
$ etcd -name infra0 -initial-advertise-peer-urls http://10.0.1.10:2380 \
  -listen-peer-urls http://10.0.1.10:2380 \
  -initial-cluster-token etcd-cluster-1 \
  -initial-cluster infra0=http://10.0.1.10:2380,infra1=http://10.0.1.11:2380,infra2=http://10.0.1.12:2380 \
  -initial-cluster-state new

# 11
$ etcd -name infra1 -initial-advertise-peer-urls http://10.0.1.11:2380 \
  -listen-peer-urls http://10.0.1.11:2380 \
  -initial-cluster-token etcd-cluster-1 \
  -initial-cluster infra0=http://10.0.1.10:2380,infra1=http://10.0.1.11:2380,infra2=http://10.0.1.12:2380 \
  -initial-cluster-state new

# 12
$ etcd -name infra2 -initial-advertise-peer-urls http://10.0.1.12:2380 \
  -listen-peer-urls http://10.0.1.12:2380 \
  -initial-cluster-token etcd-cluster-1 \
  -initial-cluster infra0=http://10.0.1.10:2380,infra1=http://10.0.1.11:2380,infra2=http://10.0.1.12:2380 \
  -initial-cluster-state new
```

在初始化完成后，`etcd`还提供动态增、删、改`etcd`集群节点的功能，这个需要用到`etcdctl`命令进行操作。

## 5.3 `etcd`自发现模式

通过自发现的方式启动`etcd`集群需要事先准备一个`etcd`集群。

假设有一个三个节点的`etcd`集群

```bash
$ curl -X PUT http://myetcd.local/v2/keys/discovery/6c007a14875d53d9bf0ef5a6fc0257c817f0fb83

/_config/size -d value=3

```

需要使用 ` http://myetcd.local/v2/keys/discovery/6c007a14875d53d9bf0ef5a6fc0257c817f0fb83` 作为`-discovery`参数启动`etcd`, 节点会自动使用此`url`目录进行`etcd`的注册和发现服务

某个启动上最终启动`etcd`命令

```bash
 etcd -name infra0 -initial-advertise-peer-urls http://10.0.1.10:2380 \
  -listen-peer-urls http://10.0.1.10:2380 \
  -discovery http://myetcd.local/v2/keys/discovery/6c007a14875d53d9bf0ef5a6fc0257c817f0fb83
```

同样的，当你完成了集群的初始化后，这些信息就失去了作用。当你需要增加节点时，需要使用`etcdctl`来进行操作。

## 5.4 `DNS`自发现模式

`etcd`还支持使用`DNS ` `SRV`记录进行启动。

另找资料

## 5.5 关键部分源码解析

在`etcdmain/etcd.go`中的`setupCluster`函数可以看到，根据不同`etcd`的参数，启动集群的方法略有不同，但是最终需要的就是一个`IP`与端口构成的字符串。

在静态配置的启动方式中，集群的所有信息都已经在给出，所以直接解析用逗号隔开的集群`url`信息就好了。

`DNS`发现的方式类似，会预先发送一个`tcp`的`SRV`请求，先查看`etcd-server-ssl._tcp.example.com`下是否有集群的域名信息，如果没有找到，则去查看`etcd-server._tcp.example.com`。根据找到的域名，解析出对应的`IP`和端口，即集群的`url`信息。

较为复杂是`etcd`式的自发现启动。首先就用自身单个的`url`构成一个集群，然后在启动的过程中根据参数进入`discovery/discovery.go`源码的`JoinCluster`函数。因为我们事先是知道启动时使用的`etcd`的`token`地址的，里面包含了集群大小`(size)`信息。在这个过程其实是个不断监测与等待的过程。启动的第一步就是在这个`etcd`的`token`目录下注册自身的信息，然后再监测`token`目录下所有节点的数量，如果数量没有达标，则循环等待。当数量达到要求时，才结束，进入正常的启动过程。

## 5.6 `peer-urls` `client-urls`

配置`etcd`过程中通常要用到两种`url`地址容易混淆，一种用于`etcd`集群同步信息并保持连接，通常称为`peer-urls`；另外一种用于接收用户端发来的HTTP请求，通常称为`client-urls`。

*   `peer-urls`：通常监听的端口为`2380`（老版本使用的端口为`7001`），包括所有已经在集群中正常工作的所有节点的地址。
*   `client-urls`：通常监听的端口为`2379`（老版本使用的端口为`4001`），为适应复杂的网络环境，新版`etcd`监听客户端请求的`url`从原来的1个变为现在可配置的多个。这样`etcd`可以配合多块网卡同时监听不同网络下的请求。

## 5.7 运行时节点变更

`etcd`集群启动完毕后，可以在运行的过程中对集群进行重构，包括核心节点的增加、删除、迁移、替换等。**运行时重构使得etcd集群无须重启即可改变集群的配置，这也是新版`etcd`区别于旧版包含的新特性。**

<https://blog.csdn.net/bbwangj/article/details/82584988>

# 6 `Proxy`模式

`etcd`作为一个反向代理把客户的请求转发给可用的`etcd`集群

![](https://res.infoq.com/articles/etcd-interpretation-application-scenario-implement-principle/zh/resources/0129011.jpg)

新版`etcd`中，只会在最初启动`etcd`集群时，发现核心节点的数量已经满足要求时，自动启用`Proxy`模式，反之则并未实现。主要原因如下。

*   `etcd`是用来保证高可用的组件
*   `etcd`集群是支持高可用的
*   自动转换使得`etcd`集群变得复杂

基于上述原因，目前`Proxy`模式有转发代理功能，而不会进行角色转换。

# 7 数据存储

分为**内存存储** 和 \*\*持久化（硬盘）\*\*存储

内存中的存储除了`顺序化的记录下所有用户对节点数据变更的记录外，还会对用户数据进行索引、建堆等方便查询的操作`。

持久化则使用`预写式日志（WAL：Write Ahead Log）`进行记录存储。

## 7.1 `WAL`体系

在`WAL`的体系中，所有的数据在提交之前都会进行日志记录。

在持久化存储目录中，有两个子目录：

*   一个是`WAL`，存储着所有事务的变化记录；
*   另一个则是`snapshot`，用于存储某一个时刻`etcd`所有目录的数据。

既然有了`WAL`实时存储了所有的变更，为什么还需要`snapshot`呢？

> 随着使用量的增加，`WAL`存储的数据会暴增，为了防止磁盘很快就爆满，`etcd`默认每10000条记录做一次`snapshot`，经过`snapshot`以后的`WAL`文件就可以删除。而通过`API`可以查询的历史`etcd`操作默认为1000条。

用户需要避免`etcd`从一个过期的数据目录中重新启动

> 因为使用过期的数据目录启动的节点会与集群中的其他节点产生不一致（如：之前已经记录并同意`Leader`节点存储某个信息，重启后又向`Leader`节点申请这个信息）。所以，为了最大化集群的安全性，一旦有任何数据损坏或丢失的可能性，你就应该把这个节点从集群中移除，然后加入一个不带数据目录的新节点。

## 7.2 预写式日志`（WAL）`

`WAL（Write Ahead Log）`最大的作用是记录了整个数据变化的全部历程。

在`etcd`中，所有数据的修改在提交前，都要先写入到`WAL`中

好处： 拥有两个重要功能：

*   **故障快速恢复**
*   **数据回滚（undo）/重做（redo）**： 因为所有的修改操作都被记录在`WAL`中，需要回滚或重做，`只需要反向或正向执行日志中的操作即可。`

## 7.3 `WAL`与`snapshot`命名规则

`WAL`文件以`$seq-$index.wal`的格式存储。最初始的`WAL`文件是`0000000000000000-0000000000000000.wal`。

表示是所有`WAL`文件中的第0个，每次切分后自增一， 初始的`Raft`状态编号为0，是根据实际存储的Raft起始状态来定。

`snapshot`的存储命名则比较容易理解，以`$term-$index.wal`格式进行命名存储。

`term`和`index`就表示存储`snapshot`时数据所在的raft节点状态，当前的任期编号以及数据项位置信息。

## 7.4 源码解析--`read`和`append`模式

`WAL`有两种模式，读模式`（read）`和数据添加`（append）`模式，两种模式不能同时成立。

一个新创建的`WAL`文件处于`append模式`，并且不会进入到`read模式`。

一个本来存在的`WAL`文件被打开的时候必然是`read模式`，并且只有在所有记录都被读完的时候，才能进入`append模式`，进入`append模式`后也不会再进入`read模式`。

> *   集群在进入到`etcdserver/server.go`的`NewServer`函数准备启动一个`etcd`节点时，会检测是否存在以前的遗留`WAL`数据。
> *   检测的第一步是查看snapshot文件夹下是否有符合规范的文件，从snapshot中获得集群的配置信息，包括token、其他节点的信息等等，然后载入`WAL`目录的内容，从小到大进行排序。
> *   根据snapshot中得到的term和index，找到`WAL`紧接着snapshot下一条的记录，然后向后更新，直到所有`WAL`包的entry都已经遍历完毕，Entry记录到`ents`变量中存储在内存里。此时`WAL`就进入`append`模式，为数据项添加进行准备。

# 8. `Raft`

`raft`包就是对`Raft一致性算法`的具体实现。

## 8.1 常见问题

*   **Raft中一个`Term（任期）`是什么意思？**。从时间上，一个任期讲即从一次竞选开始到下一次竞选开始。

    *   如果`Follower`接收不到`Leader`节点的心跳信息，就会结束当前任期，变为`Candidate`发起竞选，有助于`Leader`节点故障时集群的恢复。发起竞选投票时，任期值小的节点不会竞选成功
    *   投票出现冲突也有可能直接进入下一任再次竞选。
    *   如果集群不出现故障，那么一个任期将无限延续下去。

*   **Raft状态机怎么切换？**

    1.  Raft刚开始运行时，节点`默认进入Follower状态`，等待`Leader`发来心跳信息。
    2.  等待超时，则状态由`Follower`切换到`Candidate`进入下一轮`term`发起竞选，`等到收到集群多数节点的投票时，该节点转变为Leader`。
    3.  `Leader`节点有可能出现网络等故障，导致别的节点发起投票成为新`term`的`Leader`，此时原先的老`Leader`节点会切换为`Follower`。
    4.  `Candidate`在等待其它节点投票的过程中如果发现别的节点已经竞选成功成为`Leader`了，也会切换为`Follower`节点。

    ![](https://res.infoq.com/articles/etcd-interpretation-application-scenario-implement-principle/zh/resources/0129013.jpg)

*   **如何保证最短时间内竞选出Leader，防止竞选冲突？**

    *   在`Candidate`状态下， 有一个`times out`，这里的`times out时间是个随机值`，也就是说，每个机器成为`Candidate`以后，超时发起新一轮竞选的时间是各不相同的，这就会出现一个时间差。
    *   在时间差内，如果`Candidate1`收到的竞选信息比自己发起的竞选信息`term`值大（即对方为`新一轮term`）
    *   并且新一轮想要成为Leader的`Candidate2`包含了所有提交的数据，那么`Candidate1`就会投票给`Candidate2`。

*   **如何防止别的Candidate在遗漏部分数据的情况下发起投票成为Leader？**

*   如果发起竞选的节点在`上一个term`中保存的已提交数据不完整，节点就会拒绝投票给它。通过这种机制就可以防止遗漏数据的节点成为Leader。

*   **Raft某个节点宕机后会如何？**

    *   如果是`Follower`节点宕机，如果剩余可用节点数量超过半数，集群可以几乎没有影响的正常工作。

    *   如果是`Leader`节点宕机，那么`Follower`就收不到心跳而超时，发起竞选获得投票，成为新一轮term的Leader，继续为集群提供服务。

    > **需要注意的是；`etcd`目前没有任何机制会自动去变化整个集群总共的节点数量**， 即如果没有人为的调用`API`，`etcd`宕机后的节点仍然被计算为总节点数中，任何请求被确认需要获得的投票数都是这个总数的半数以上。

<img src="https://res.infoq.com/articles/etcd-interpretation-application-scenario-implement-principle/zh/resources/0129014.jpg" style />

*   **用户从集群中哪个节点读写数据？**
    *   `Raft`为了保证数据的强一致性，所有的数据流向都是一个方向，从Leader流向Follower，也就是所有Follower的数据必须与`Leader`保持一致，如果不一致会被覆盖。
    *   每个节点都有`Raft`已提交数据准确的备份（最坏的情况也只是已提交数据还未完全同步），所以读的请求任意一个节点都可以处理。

## 8.2 关键部分源码

在`etcd`中，对Raft算法的调用如下，你可以在`etcdserver/raft.go`中的`startNode`找到：

```go
storage := raft.NewMemoryStorage()
n := raft.StartNode(0x01, []int64{0x02, 0x03}, 3, 1, storage)
```

首先，你需要把从集群的其他机器上收到的信息推送到Node节点，你可以在`etcdserver/server.go`中的`Process`函数看到。

```go
func (s *EtcdServer) Process(ctx context.Context, m raftpb.Message) error {
    if m.Type == raftpb.MsgApp {
        s.stats.RecvAppendReq(types.ID(m.From).String(), m.Size())
    }
    return s.node.Step(ctx, m)
}
```

其次，你需要把日志项存储起来，在你的应用中执行提交的日志项，然后把完成信号发送给集群中的其它节点，再通过`node.Ready()`监听等待下一次任务执行。

最后，你需要保持一个心跳信号`Tick()`。

> Raft有两个很重要的地方用到超时机制：心跳保持和Leader竞选。需要用户在其raft的Node节点上周期性的调用Tick()函数，以便为超时机制服务。

综上所述，整个raft节点的状态机循环类似如下所示：

```go
for {
    select {
    case &lt;-s.Ticker:
        n.Tick()
    case rd := &lt;-s.Node.Ready():
        saveToStorage(rd.State, rd.Entries)
        send(rd.Messages)
        process(rd.CommittedEntries)
        s.Node.Advance()
    case &lt;-s.done:
        return
    }
}
```

# 9. `Store`--提供`API`支持

为用户提供五花八门的`API`支持，处理用户的各项请求。要理解`Store`，只需要从`etcd`的`API`入手即可。

[etcd\_api\_doc文档](https://etcd.io/docs/v3.4.0/)

`API`中提到的`目录（Directory）和键（Key）`，上文中也可能称为`etcd节点（Node）`。

*
    ## 为`etcd`存储的键赋值

```bash
curl http://127.0.0.1:2379/v2/keys/message -X PUT -d value="Hello world"
# 反馈的内容
{
    "action": "set",
    "node": {
        "createdIndex": 2,			#节点每次有变化时都会自增一个，　除用户请求外，内部运行也会引起变化
        "key": "/message",
        "modifiedIndex": 2,
        "value": "Hello world"
    },
    # 如果前面有set操作过
    "prevNode":{"key":"/message","value":"hello world2","modifiedIndex":7,"createdIndex":7}
}
```

反馈的内容含义如下：

*   `action`: 刚刚进行的动作名称。

*   `node.key`: 请求的HTTP路径。`etcd`使用一个类似文件系统的方式来反映键值存储的内容。

*   `node.value`: 刚刚请求的键所存储的内容。

*   `node.createdIndex`:` etcd`节点每次有变化时都会自增的一个值，除了用户请求外，`etcd`内部运行（如启动、集群信息变化等）也会对节点有变动而引起这个值的变化。

*   `node.modifiedIndex`: 类似`node.createdIndex`，能引起`modifiedIndex`变化的操作包括`set, delete, update, create, compareAndSwap and compareAndDelete`。

*

    ## 查看`etcd`键存储的值

    ```bash
    curl http://127.0.0.1:2379/v2/keys/message -X GET
    # 反馈
    {"action":"get",
    	"node":{
    		"key":"/message",
    		"value":"hello world2",
    		"modifiedIndex":8,
    		"createdIndex":8
    		}
     }
    ```

*
    ## 其他`API`

*   修改键值：与创建新值几乎相同，但是反馈时会有一个`prevNode`值反应了修改前存储的内容。

    ```bash
    curl http://127.0.0.1:2379/v2/keys/message -XPUT -d value="Hello etcd"
    ```

*   删除一个值

```bash
curl http://127.0.0.1:2379/v2/keys/message -XDELETE
```

*   对一个键进行定时删除：`etcd`中对键进行定时删除，设定一个`TTL`值，当这个值到期时键就会被删除。反馈的内容会给出`expiration`项告知超时时间，`ttl`项告知设定的时长。

```bash
curl http://127.0.0.1:2379/v2/keys/foo -XPUT -d value=bar -d ttl=5
# 反馈类似
"expiration":"2020-04-27T04:30:06.183849063Z"
```

*   取消定时删除任务

    ```bash
    curl http://127.0.0.1:2379/v2/keys/foo -XPUT -d value=bar -d ttl= -d prevExist=true
    ```

*   对键值修改进行监控：`etcd`提供的这个`API`让用户可以监控一个值或者递归式的监控一个目录及其子目录的值，当目录或值发生变化时，`etcd`会主动通知。

```bash
curl http://127.0.0.1:2379/v2/keys/foo?wait=true
# 会阻塞等待
# 使用另外client 
# curl http://localhost:2379/v2/keys/message -X PUT -d value=Hello1
# 改变后，收到信息
{"action":"set",
	"node":{
		"key":"/message",
		"value":"Hello1",
		"modifiedIndex":19,
		"createdIndex":19
		},
	"prevNode":{
	"key":"/message",
	"value":"",
	"modifiedIndex":18,
	"createdIndex":18
	}
}
```

*   对过去的键值操作进行查询：类似上面提到的监控，只不过监控时加上了过去某次修改的索引编号，就可以查询历史操作。默认可查询的历史记录为1000条。

```bash
curl 'http://127.0.0.1:2379/v2/keys/foo?wait=true&waitIndex=7'
```

*   自动在目录下创建有序键。在对创建的目录使用`POST`参数，会自动在该目录下创建一个以`createdIndex`值为键的值，这样就相当于以创建时间先后严格排序了。这个`API`对分布式队列这类场景非常有用。

```bash
curl http://127.0.0.1:2379/v2/keys/queue -XPOST -d value=Job1
{
    "action": "create",
    "node": {
        "createdIndex": 6,
        "key": "/queue/6",
        "modifiedIndex": 6,
        "value": "Job1"
    }
}
```

*   按顺序列出所有创建的有序键。

```bash
curl -s 'http://127.0.0.1:2379/v2/keys/queue?recursive=true&sorted=true'
```

*   创建定时删除的目录：就跟定时删除某个键类似。如果目录因为超时被删除了，其下的所有内容也自动超时删除。

```bash
curl http://127.0.0.1:2379/v2/keys/dir -XPUT -d ttl=30 -d dir=true
```

*   刷新超时时间。

    ```bash
        curl http://127.0.0.1:2379/v2/keys/dir -XPUT -d ttl=30 -d dir=true -d prevExist=true
    ```

## 自动化`CAS`操作

自动化`CAS（Compare-and-Swap）`操作：`etcd`强一致性最直观的表现就是这个`API`，通过设定条件，阻止节点二次创建或修改。即用户的指令被执行当且仅当`CAS`的条件成立。条件有以下几个。

*   `prevValue` 先前节点的值，如果值与提供的值相同才允许操作。
*   `prevIndex `先前节点的编号，编号与提供的校验编号相同才允许操作。
*   `prevExist `先前节点是否存在。如果存在则不允许操作。这个常常被用于分布式锁的唯一获取。

假设先进行了如下操作：设定了`foo`的值。

```bash
curl http://127.0.0.1:2379/v2/keys/foo -XPUT -d value=one
```

然后再进行操作：

```bash
curl http://127.0.0.1:2379/v2/keys/foo?prevExist=false -XPUT -d value=three
```

就会返回创建失败的错误。

*   条件删除（Compare-and-Delete）：与`CAS`类似，条件成立后才能删除。

*   创建目录

    ```bash
    curl http://127.0.0.1:2379/v2/keys/dir -XPUT -d dir=true
    ```

*   列出目录下所有的节点信息，最后以`/`结尾。还可以通过`recursive`参数递归列出所有子目录信息。

    ```bash
    curl http://127.0.0.1:2379/v2/keys/
    ```

*   删除目录：默认情况下只允许删除空目录，如果要删除有内容的目录需要加上`recursive=true`参数。

    ```bash
    curl 'http://127.0.0.1:2379/v2/keys/foo_dir?dir=true' -XDELETE
    ```

*   创建一个隐藏节点：命名时名字以下划线`_`开头默认就是隐藏键。

    ```bash
    curl http://127.0.0.1:2379/v2/keys/_message -XPUT -d value="Hello hidden world"
    ```

## 总结：

*   `Store` 对`etcd`下存储的数据进行加工，创建出如文件系统般的树状结构供用户快速查询。
*   它有一个`Watcher`用于节点变更的实时反馈，还需要维护一个`WatcherHub`对所有`Watcher`订阅者进行通知的推送。
*   它还维护了一个由定时键构成的小顶堆，快速返回下一个要超时的键。
*   所有这些`API`的请求都以事件的形式存储在事件队列中等待处理。

# 10. `go`操作`etcd`

文档：　<https://etcd.io/docs/v3.4.0/integrations/>

## 安装

```bash
go get go.etcd.io/etcd/clientv3
```

会出现报错

解决方法:

将grpc版本替换成`v1.26.0版本`

### 1. 修改依赖为v1.26.0

```bash
go mod edit -require=google.golang.org/grpc@v1.26.0
```

有时候还有错

```bash
replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
```

## `API`` V2`与`V3`区别

*   事务：`ETCD V3`提供了多键条件事务（`multi-key conditional transactions`），应用各种需要使用事务代替原来的`Compare-And-Swap`操作。
*   平键空间（`Flat key space`）：`ETCD V3`不再使用目录结构，只保留键。例如：”/a/b/c/“是一个键，而不是目录。`V3`中提供了前缀查询，来获取符合前缀条件的所有键值，这变向实现了`V2`中查询一个目录下所有子目录和节点的功能。
*   简洁的响应：像`DELETE`这类操作成功后将不再返回操作前的值。如果希望获得删除前的值，可以使用事务，来实现一个原子操作，先获取键值，然后再删除。
*   租约：租约代替了`V2`中的`TTL`实现，`TTL`绑定到一个租约上，键再附加到这个租约上。当TTL过期时，租约将被销毁，同时附加到这个租约上的键也被删除。

***

以下操作来自官方文档

## `put`和`get`操作

`put`命令用来设置键值对数据，`get`命令用来根据key获取值。

```go
package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// etcd client put/get demo
// use etcd/clientv3

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
    fmt.Println("connect to etcd success")
	defer cli.Close()
	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "message", "dsb")
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "message")
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
}
```

## `watch`操作

`watch`用来获取未来更改的通知。

```go
package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// watch demo

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()
	// watch key:message change
	rch := cli.Watch(context.Background(), "message") // <-chan WatchResponse
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}

// 或者
	// see https://github.com/coreos/etcd/blob/master/clientv3/example_watch_test.go
	log.Println("监视")
	rch := cli.Watch(context.Background(), "", clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
```

会阻塞不断接收到通知

## `lease`租约--定时

```go
package main

import (
	"fmt"
	"time"
)

// etcd lease

import (
	"context"
	"log"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connect to etcd success.")
	defer cli.Close()

	// 创建一个5秒的租约
	resp, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
	}

	// 5秒钟之后, /lease/ 这个key就会被移除
	_, err = cli.Put(context.TODO(), "/lease/", "dsb", clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}
}
```

## `keepAlive`--续租

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// etcd keepAlive

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connect to etcd success.")
	defer cli.Close()
	
    // 建立租约
	resp, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
	}
	
    // 定时修改
	_, err = cli.Put(context.TODO(), "/lease/", "dsb", clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}

	// the key 'foo' will be kept forever
	ch, kaerr := cli.KeepAlive(context.TODO(), resp.ID)
	if kaerr != nil {
		log.Fatal(kaerr)
	}
	for {
		ka := <-ch
		fmt.Println("ttl:", ka.TTL)
	}
}
```

## 基于`etcd`实现分布式锁

`go.etcd.io/etcd/clientv3/concurrency`在`etcd`之上实现并发操作，如分布式锁、屏障和选举。

```go
import "go.etcd.io/etcd/clientv3/concurrency"
```

```go
cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
if err != nil {
    log.Fatal(err)
}
defer cli.Close()

// 创建两个单独的会话用来演示锁竞争
s1, err := concurrency.NewSession(cli)
if err != nil {
    log.Fatal(err)
}
defer s1.Close()
m1 := concurrency.NewMutex(s1, "/my-lock/")

s2, err := concurrency.NewSession(cli)
if err != nil {
    log.Fatal(err)
}
defer s2.Close()
m2 := concurrency.NewMutex(s2, "/my-lock/")

// 会话s1获取锁
if err := m1.Lock(context.TODO()); err != nil {
    log.Fatal(err)
}
fmt.Println("acquired lock for s1")

m2Locked := make(chan struct{})
go func() {
    defer close(m2Locked)
    // 阻塞等待直到会话s1释放了/my-lock/的锁
    if err := m2.Lock(context.TODO()); err != nil {
        log.Fatal(err)
    }
}()

if err := m1.Unlock(context.TODO()); err != nil {
    log.Fatal(err)
}
fmt.Println("released lock for s1")

// 关闭了m2Locked就会继续进行
<-m2Locked
fmt.Println("acquired lock for s2")
```

输出

```bash
acquired lock for s1
released lock for s1
acquired lock for s2
```

## 事务&超时

```go
log.Println("事务&超时")
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	// 流式处理
	_, err = cli.Txn(ctx).
		If(clientv3.Compare(clientv3.Value("key"), ">", "abc")). // txn value comparisons are lexical
		Then(clientv3.OpPut("key", "XYZ")).                      // this runs, since 'xyz' > 'abc'
		Else(clientv3.OpPut("key", "ABC")).
		Commit()
	cancel()
	if err != nil {
		log.Fatal(err)
	}

```

## 服务注册

```go
package main

import (
    "context"
    "fmt"
    "go.etcd.io/etcd/clientv3"
    "time"
)

//创建租约注册服务
type ServiceReg struct {
    client        *clientv3.Client
    lease         clientv3.Lease
    leaseResp     *clientv3.LeaseGrantResponse
    canclefunc    func()
    keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
    key           string
}

func NewServiceReg(addr []string, timeNum int64) (*ServiceReg, error) {
    conf := clientv3.Config{
        Endpoints:   addr,
        DialTimeout: 5 * time.Second,
    }

    var (
        client *clientv3.Client
    )

    if clientTem, err := clientv3.New(conf); err == nil {
        client = clientTem
    } else {
        return nil, err
    }

    ser := &ServiceReg{
        client: client,
    }

    if err := ser.setLease(timeNum); err != nil {
        return nil, err
    }
    go ser.ListenLeaseRespChan()
    return ser, nil
}

//设置租约
func (this *ServiceReg) setLease(timeNum int64) error {
    lease := clientv3.NewLease(this.client)

    //设置租约时间
    leaseResp, err := lease.Grant(context.TODO(), timeNum)
    if err != nil {
        return err
    }

    //设置续租
    ctx, cancelFunc := context.WithCancel(context.TODO())
    leaseRespChan, err := lease.KeepAlive(ctx, leaseResp.ID)

    if err != nil {
        return err
    }

    this.lease = lease
    this.leaseResp = leaseResp
    this.canclefunc = cancelFunc
    this.keepAliveChan = leaseRespChan
    return nil
}

//监听 续租情况
func (this *ServiceReg) ListenLeaseRespChan() {
    for {
        select {
        case leaseKeepResp := <-this.keepAliveChan:
            // 续租通道关闭会返回nil, 关闭此协程
            if leaseKeepResp == nil {
                fmt.Printf("已经关闭续租功能\n")
                return
            } else {
                fmt.Printf("续租成功\n")
            }
        }
    }
}

//通过租约 注册服务
func (this *ServiceReg) PutService(key, val string) error {
    kv := clientv3.NewKV(this.client)
    _, err := kv.Put(context.TODO(), key, val, clientv3.WithLease(this.leaseResp.ID))
    return err
}


//撤销租约
func (this *ServiceReg) RevokeLease() error {
    this.canclefunc()
    time.Sleep(2 * time.Second)
    _, err := this.lease.Revoke(context.TODO(), this.leaseResp.ID)
    return err
}

func main() {
    ser,_ := NewServiceReg([]string{"127.0.0.1:2379"},5)
    ser.PutService("/node/111","heiheihei")
    select{}
}
```

## 服务发现

达到动态修改`ServiceList`目的

```go
import (
    "go.etcd.io/etcd/clientv3"
    "time"
    "context"
    "go.etcd.io/etcd/mvcc/mvccpb"
    "sync"
    "log"
)

type ClientDis struct {
    client        *clientv3.Client
    serverList    map[string]string
    lock          sync.Mutex	// 同步锁
}

func NewClientDis (addr []string)( *ClientDis, error){
    conf := clientv3.Config{
        Endpoints:   addr,
        DialTimeout: 5 * time.Second,
    }
    if client, err := clientv3.New(conf); err == nil {
        return &ClientDis{
            client:client,
            serverList:make(map[string]string),
        }, nil
    } else {
        return nil ,err
    }
}


func (this * ClientDis) GetService(prefix string) ([]string ,error){
    // resp 获取到的是pre(node) 节点下的所有服务地址
    resp, err := this.client.Get(context.Background(), prefix, clientv3.WithPrefix())
    if err != nil {
        return nil, err
    }
    // 初始化
    addrs := this.extractAddrs(resp)

    go this.watcher(prefix)
    return addrs ,nil
}

//　监听节点服务变化
func (this *ClientDis) watcher(prefix string) {
    rch := this.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
    for wresp := range rch {
        for _, ev := range wresp.Events {
            // 根据类型　变化ServiceList
            switch ev.Type {
            case mvccpb.PUT:	// 修改
                this.SetServiceList(string(ev.Kv.Key),string(ev.Kv.Value))
            case mvccpb.DELETE:	// 节点删除
                this.DelServiceList(string(ev.Kv.Key))
            }
        }
    }
}

// 
func (this *ClientDis) extractAddrs(resp *clientv3.GetResponse) []string {
    addrs := make([]string,0)
    // 如果节点没有地址，返回空字符串切片，　Kvs: []map[string]string 数组嵌套字典
    if resp == nil || resp.Kvs == nil {
        return addrs
    }
    for i := range resp.Kvs {
        if v := resp.Kvs[i].Value; v != nil {
            this.SetServiceList(string(resp.Kvs[i].Key),string(resp.Kvs[i].Value))
            addrs = append(addrs, string(v))
        }
    }
    return addrs
}

// 设置serverList发现的服务列表
func (this *ClientDis) SetServiceList(key,val string) {
    this.lock.Lock()
    defer this.lock.Unlock()
    this.serverList[key] = string(val)
    log.Println("set data key :",key,"val:",val)
}

func (this *ClientDis) DelServiceList(key string) {
    this.lock.Lock()
    defer this.lock.Unlock()
    delete(this.serverList,key)
    log.Println("del data key:", key)
}


func (this *ClientDis) SerList2Array()[]string {
    this.lock.Lock()
    defer this.lock.Unlock()
    addrs := make([]string,0)

    for _, v := range this.serverList {
        addrs = append(addrs,v)
    }
    return addrs
}

func main () {
    cli,_ := NewClientDis([]string{"127.0.0.1:2379"})
    cli.GetService("/node")
    select {}
}
```




