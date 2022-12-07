- [MyLearn](#mylearn)
- [编写错误或者引用问题](#编写错误或者引用问题)
- [目录](#目录)

# MyLearn

从2019年开始,我就偶尔记录自己的学习、工作和摘抄的一些笔记，以前都是存在了有道云文档，现在陆续把它放上来，刚好也给自己一个笔记进行整理复习吧

# 编写错误或者引用问题

如有编写错误,或者某些文档引用没写清著作人的,请联系我更改或者删除,提出issues就行,十分抱歉

# 目录
├─Docker
│      ![Dockerfile使用.pdf](Dockerfile使用.pdf)
│      Docker基础版&使用Docker安装FastDFS.pdf
│      Docker教程.pdf
│      Docker语句.pdf
│      
├─Golang
│  ├─etcd
│  │      etcd.md
│  │      
│  ├─golang常用包
│  │  │  bytes.md
│  │  │  errors 包.md
│  │  │  flag-行参数解析.md
│  │  │  fmt.Printf格式化输出.md
│  │  │  go-grpc使用.md
│  │  │  gob和msgpack.md
│  │  │  gopsutil获取系统性能数据.md
│  │  │  json数据格式转换.md
│  │  │  net.http--服务端和客户端.md
│  │  │  protobuf.md
│  │  │  runtime.md
│  │  │  template.md
│  │  │  上下文环境Context包.md
│  │  │  数据类型转换(strconv包).md
│  │  │  文件操作-file-bufio-os.md
│  │  │  日志库-Uber-go的Zap Logger.md
│  │  │  标准库log.md
│  │  │  
│  │  └─GoMicro
│  │      │  README.md
│  │      │  框架&Cli.md
│  │                      
│  ├─基础
│  │      01-Go基础.md
│  │      02-函数、变量.md
│  │      03-数据类型与操作.md
│  │      04-流程控制.md
│  │      05-数组.md
│  │      06-切片slice.md
│  │      07-映射map.md
│  │      08-函数func.md
│  │      09-指针.md
│  │      10-结构体struct.md
│  │      11-包package.md
│  │      12-接口interface.md
│  │      13-并发goroutine.md
│  │      13-并发实例--Runner.md
│  │      13-并发示例--连接池pool.md
│  │      14-并发sync控制.md
│  │      15-反射reflect.md
│  │      16-网络编程.md
│  │      17-单元测试go test.md
│  │      
│  ├─项目常用
│  │      golang redis分布式锁.md
│  │      Golang静态检查.pdf
│  │      gopay聚合支付包.md
│  │      gorm源码.md
│  │      Viper读取配置文件.pdf
│  │      Worker协程池--gpools.pdf
│  │      Zap记录日志,分等级,日志切割.pdf
│  │      利用反射自定义Verify.pdf
│  │      提高邮件送达率.md
│  │      自己以前的socket项目.md
│  │      
│  └─高级优化
│      │  30+张图讲解：Golang调度器GMP原理与调度全分析.md
│      │  for-loop 与 json.Unmarshal 性能分析概要.md
│      │  for-range排坑指南.md
│      │  golang Channel用法总结.pdf
│      │  Golang CSP并发模型.md
│      │  Go之pprof性能分析&实践.md
│      │  Go装饰器.md
│      │  Go语言中的内存对齐.pdf
│      │  控制goroutine并发数量.md
│      │  深入Golang深入Golang调度器之GMP模型.pdf
│      │  逃逸分析.md
│      │  面向并发的内存模型讲解--goroutine.md
│      │  
│      └─原理解析
│              array数组原理.md
│              channel原理解析.md
│              Context原理解析.md
│              iota原理解析.md
│              map原理解析.pdf
│              Mutex锁原理分析.md
│              range源码解析.md
│              string原理解析.md
│              struct原理解析.md
│              sync.RWMutex原理解析.md
│              
├─Kafka
│      Kafka架构.md
│      
├─Linux
│      curl 使用.pdf
│      Linux-VIM高级快捷键使用.pdf
│      Linux_ps进程管理命令.pdf
│      Linux信号量signal.md
│      Linux原理--各个文件目录.pdf
│      Linux基础---创建文件-目录.pdf
│      linux如何查看端口被哪个进程占用？.pdf
│      Linux常用20条命令.pdf
│      linux服务器性能查看.md
│      Linux高级--查找,权限,压缩,登录.pdf
│      Ubuntu设置根据情况设置shell启动运行代码--如export变量.pdf
│      vmstat虚拟内存统计.pdf
│      实战linux命令大全.pdf
│      磁盘管理.pdf
│      防火墙-iptables&Firewall.pdf
│      
├─Mysql
│  │  Explain 详解.pdf
│  │  mysql如何保证数据一致性.pdf
│  │  Mysql所有命令.pdf
│  │  Mysql锁.pdf
│  │  SQLAlchemy.pdf
│  │  分布式事务解决方案实战.pdf
│  │  基础总结.md
│  │  
│  ├─PostgreSQL
│  │      binlog介绍.pdf
│  │      mysql与pg的主键索引说明.pdf
│  │      MySQL数据库优化.pdf
│  │      pgSQL命令速查表.pdf
│  │      pg与mysql的MVCC对比.pdf
│  │      pg常用工具.pdf
│  │      pg索引类型基本原理和应用场景.pdf
│  │      事务与索引.pdf
│  │      外键约束.pdf
│  │      
│  └─高级语法
│          datetime日期类型处理总结.md
│          decode()函数.md
│          Docker-Compose.md
│          Dockerfile ENV和ARG应用.md
│          DockerFile中多个From&copy--from.md
│          Docker语句二十条.pdf
│          Docker连接网络的三种方式.md
│          Group by 与Order by 顺序问题.md
│          substring截取内容分组.md
│          实战 SQL窗口函数.pdf
│          将一张表中数据批量导入另一张表.md
│          
├─Nginx
│      Nginx部署基础.pdf
│      Nginx项目部署.md
│      使用Nginx转发Fastdfs及解决同步问题.pdf
│      
├─Python
│      README.md
│      
├─Redis
│      Redis 高级知识.pdf
│      Redis-让值过期.md
│      redis基础知识.pdf
│      redis的五大数据类型实现原理.pdf
│      redis的底层数据结构.pdf
│      Redis集群.pdf
│      
├─selenium和爬虫相关
│      Selenium(PhantomJS&Chromedriver安装).pdf
│      selenium使用.pdf
│      selenium扩展--headless无头模式.pdf
│      selenium等待.pdf
│      selenium输入验证码.pdf
│      实例-踩的坑.pdf
│      
├─数据结构与算法
│      B-树.pdf
│      二分查找.pdf
│      二叉树(1).pdf
│      剑指offer.pdf
│      总结：.pdf
│      排序算法.pdf
│      散列表.pdf
│      栈实现.pdf
│      算法基础.pdf
│      红黑树.pdf
│      链表.pdf
│      队列.pdf
│      顺序表.pdf
│      
└─科学计算
        Matplotlib基础.pdf
        Numpy.pdf
        Pandas.pdf
