- [asynq](#asynq)
  - [Features](#features)
- [快速开始](#快速开始)
- [Tasks的生命周期](#tasks的生命周期)
- [自带配置的webUI 和Cli](#自带配置的webui-和cli)
- [个人理解](#个人理解)

# asynq
源码地址: https://github.com/hibiken/asynq

Highlevel overview of how Asynq works:

- Client puts tasks on a queue 
- Server pulls tasks off queues and starts a worker goroutine for each task 
- Tasks are processed concurrently by multiple workers

## Features
- Guaranteed [at least one execution](https://www.cloudcomputingpatterns.org/at_least_once_delivery/) of a task
- Scheduling of tasks
- [Retries](https://github.com/hibiken/asynq/wiki/Task-Retry) of failed tasks
- Automatic recovery of tasks in the event of a worker crash
- [Weighted priority queues](https://github.com/hibiken/asynq/wiki/Queue-Priority#weighted-priority)
- [Strict priority queues](https://github.com/hibiken/asynq/wiki/Queue-Priority#strict-priority)
- Low latency to add a task since writes are fast in Redis
- De-duplication of tasks using [unique option](https://github.com/hibiken/asynq/wiki/Unique-Tasks)
- Allow [timeout and deadline per task](https://github.com/hibiken/asynq/wiki/Task-Timeout-and-Cancelation)
- Allow [aggregating group of tasks](https://github.com/hibiken/asynq/wiki/Task-aggregation) to batch multiple successive operations
- [Flexible handler interface with support for middlewares](https://github.com/hibiken/asynq/wiki/Handler-Deep-Dive)
- [Ability to pause queue](/tools/asynq/README.md#pause) to stop processing tasks from the queue
- [Periodic Tasks](https://github.com/hibiken/asynq/wiki/Periodic-Tasks)
- [Support Redis Cluster](https://github.com/hibiken/asynq/wiki/Redis-Cluster) for automatic sharding and high availability
- [Support Redis Sentinels](https://github.com/hibiken/asynq/wiki/Automatic-Failover) for high availability
- Integration with [Prometheus](https://prometheus.io/) to collect and visualize queue metrics
- [Web UI](#web-ui) to inspect and remote-control queues and tasks
- [CLI](#command-line-tool) to inspect and remote-control queues and tasks
- 
# 快速开始
参考: https://github.com/hibiken/asynq#quickstart

# Tasks的生命周期
参考: https://github.com/hibiken/asynq/wiki/Life-of-a-Task
```golang
// Task 1 : Scheduled to be processed 24 hours later.
client.Enqueue(task1, asynq.ProcessIn(24*time.Hour))

// Task 2 : Enqueued to be processed immediately.
client.Enqueue(task2)

// Task 3: Enqueued with a Retention option.
client.Enqueue(task3, asynq.Retention(2*time.Hour))
```
`task1` 的任务会在开始在`Scheduled`(redis中的`scheduled`名称的`zset`)中, 到期执行到`Pending` 

`task2`直接执行,处于`Pending`状态

`task3`任务执行完,会保留在`Completed`两个小时, 如果不设定Retention,则会直接删除

```bash
+-------------+            +--------------+          +--------------+           +-------------+
|             |            |              |          |              | Success   |             |
|  Scheduled  |----------->|   Pending    |--------->|    Active    |---------> |  Completed  |
|  (Optional) |            |              |          |              |           |  (Optional) |
+-------------+            +--------------+          +--------------+           +-------------+
                                  ^                       |                            |
                                  |                       |                            | Deletion
                                  |                       | Failed                     |
                                  |                       |                            V
                                  |                       |
                                  |                       |
                           +------+-------+               |        +--------------+
                           |              |               |        |              |
                           |    Retry     |<--------------+------->|   Archived   |
                           |              |                        |              |
                           +--------------+                        +--------------+


```
# 自带配置的webUI 和Cli
Cli: https://github.com/hibiken/asynq/tree/master/tools/asynq

[Asynqmon](https://github.com/hibiken/asynqmon) is a web based tool for monitoring and administrating Asynq queues and tasks.


# 个人理解
初始化启动`Server`,会在redis中添加key: `asynq:server`,`asynq:workers`, `asynq:queues`, 保存配置的元信息 

`Server`主要执行原理是从Redis不断获取`Queue`names, 高级别优先, 然后从队列中`zset`或者`set`类型,获取任务`id`, 再从Redis `get`这个任务id的payload执行

`Client`可以通过配置各种任务的`options`, 然后进行`enqueue`操作, 内容传递到redis中,最后交给`Server`调度

其中更多的功能,可以查看其`rdb`源码的实现, 利用执行一系列`redis_script`redis脚本


