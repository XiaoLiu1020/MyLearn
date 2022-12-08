# &#x20;

# 文档

介绍类 \:Redis分布式锁-golang实现  <http://turbock79.cn/?p=3911>

<https://github.com/bsm/redislock/blob/main/redislock.go>

# &#x20;使用

```go
import (
  "fmt"
  "time"

  "github.com/bsm/redislock"
  "github.com/go-redis/redis/v9"
)

func main() {
	// Connect to redis.
	client := redis.NewClient(&redis.Options{
		Network:	"tcp",
		Addr:		"127.0.0.1:6379",
	})
	defer client.Close()

	// Create a new lock client.
	locker := redislock.New(client)

	ctx := context.Background()

	// Try to obtain lock.
	lock, err := locker.Obtain(ctx, "my-key", 100*time.Millisecond, nil)
	 // lock 里面包含redisClient, key, value
	if err == redislock.ErrNotObtained {
		fmt.Println("Could not obtain lock!")
	} else if err != nil {
		log.Fatalln(err)
	}

	// Don't forget to defer Release.
	defer lock.Release(ctx)
	fmt.Println("I have a lock!")

	// Sleep and check the remaining TTL.
	time.Sleep(50 * time.Millisecond)
	if ttl, err := lock.TTL(ctx); err != nil {
		log.Fatalln(err)
	} else if ttl > 0 {
		fmt.Println("Yay, I still have my lock!")
	}

	// Extend my lock.
	if err := lock.Refresh(ctx, 100*time.Millisecond, nil); err != nil {
		log.Fatalln(err)
	}

	// Sleep a little longer, then check.
	time.Sleep(100 * time.Millisecond)
	if ttl, err := lock.TTL(ctx); err != nil {
		log.Fatalln(err)
	} else if ttl == 0 {
		fmt.Println("Now, my lock has expired!")
	}

}
```

# 源码解析

<https://github.com/bsm/redislock/blob/main/redislock.go#L51>&#x20;

主要看locker.Obtain方法, 可以添加locker.Options

原理:&#x20;

1- context.WithDealined 设置获取不到锁,返回退出的时间

2- 在有限时间内,for循环获取redis lock, 即通过redis.SetNx(ctx, key, value, ttl)

3- 判断locker是否有Options, 在总的有限Dealiend里,间隔多久去请求一次锁, 这里使用timer.NewTicker(backoff)实现

其他方法: Release,Refresh,TTL 都使用了redis的LUA脚本的原子性
