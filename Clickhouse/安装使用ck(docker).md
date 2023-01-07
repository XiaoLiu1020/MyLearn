- [使用`docker`安装`ck`](#使用docker安装ck)
- [启动\&连接](#启动连接)
- [数据库基本操作](#数据库基本操作)

# 使用`docker`安装`ck`
```bash
# server
docker pull yandex/clickhouse-server
# client, 也可以直接进server容器里 执行clickhouse-client
docker pull yandex/clickhouse-client
```

介绍文档: https://hub.docker.com/r/yandex/clickhouse-server/

# 启动&连接
```bash
docker run -d --name some-clickhouse-server --ulimit nofile=262144:262144 yandex/clickhouse-server

# 允许连接
docker run -it --rm --link some-clickhouse-server:clickhouse-server yandex/clickhouse-client --host clickhouse-server
```
`client` 使用: 可以查看文档: https://clickhouse.com/docs/en/interfaces/cli/

大概:
```bash
clickhouse-client --host HOSTNAME.clickhouse.cloud \
  --secure \
  --port 9440 \
  --user default \
  --password PASSWORD \
  --query "INSERT INTO cell_towers FORMAT CSVWithNames" \
  < cell_towers.csv
```

# 数据库基本操作
跟Mysql基本一致语法, 每个类型可以去查看具体Ck官方文档

```sql
# 创建数据库
create table user(id UInt8,name String,address String)engine=MergeTree order by id

# 插入
 insert into user (id, name, address) values (1, '你要相信', '地址')

# 增加列
 alter table user add column age Int8;

# 修改列
alter table user modify column age String;

# 查看
select * from user;
```