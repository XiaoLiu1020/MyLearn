- [`mypool`](#mypool)
	- [`main.go`](#maingo)
- [修改地方：](#修改地方)


# `mypool`
```golang
package my_pool

import (
	"errors"
	"io"
	"log"
	"sync"
)

//一个安全的资源池，被管理的资源必须都实现io.Close接口
type Pool struct {
	M       sync.Mutex                //锁
	res     chan io.Closer            //缓冲通道，保存共享资源, 类型为io.Closer, 实现了io.Closer接口类型都可以作为资源，交给资源池管理
	factory func() (io.Closer, error) //生成新资源, 返回是io.Closer func自己定义
	closed  bool                      //表示是否关闭，关闭后，再访问报错
}

//定义访问错误
var ErrPoolClosed = errors.New("资源池已经关闭了")

//工厂函数,返回资源池 , 接受参数 生成资源函数， size 大小
func New(fn func() (io.Closer, error), size uint) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("size的值不能少于0")
	}
	return &Pool{
		res:     make(chan io.Closer, size),	//最多存下两个连接
		factory: fn,
		closed:  false,
	}, nil
}

// 从资源池获取一个资源
func (p *Pool) Acquire() (io.Closer, error) {
	for {
		log.Println("寻找连接池资源中")
		select {
		case r, ok := <-p.res:
			log.Println("Acquire: 共享资源")
			if !ok {
				return nil, ErrPoolClosed
			}
			return r, nil
		default:
		    // 这样写，并发较低，只能使用有限资源，容易造成死锁
			if len(p.res) < 2 {
				log.Println("Acquire: 新生成资源")
				log.Println("当前连接池Pool长度：", len(p.res))
				return p.factory()
			}
		}
	}
}

//释放资源方法
func (p *Pool) Release(r io.Closer){
	//保证该操作和Close方法的操作是安全的
	p.M.Lock()
	defer p.M.Unlock()

	//如果资源池已经关闭了，也把这个资源关闭了
	if p.closed {
		_ = r.Close()
		return
	}
	select {
	case p.res <- r:
		log.Println("资源释放到池子里了")
	// 问题： 并发低
/*	default:
		log.Println("资源池满了，释放这个资源吧")
		_ = r.Close()*/
	}
}

//关闭资源池
func (p *Pool) Close() {
	p.M.Lock()
	defer p.M.Unlock()

	if p.closed{
		return
	}
	p.closed = true

	// 关闭通道，不让写入了
	close(p.res)

	//关闭通道里所有资源
	for r := range p.res {
		_ = r.Close()
	}
}

```
## `main.go`
```golang
package main

import (
	"io"
	pool "liukaitao.com/m/my_pool"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	//模拟最大的goroutine
	maxGoroutine 	= 5
	//资源池大小
	poolRes 		= 2
)

func main() {
	var wg sync.WaitGroup
	wg.Add(maxGoroutine)

	// 实例，参数1： 传入创建资源方法，方法需要传入io.Closer, error方法 ，
	// 参数2 传入size, 池子的大小
	p, err := pool.New(createConnection, poolRes)
	if err != nil {
		log.Println(err)
		return
	}
	// 模拟好几个goroutine同时使用资源池查询数据
	for query := 0; query < maxGoroutine; query++ {
		go func(q int) {
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
			dbQuery(q, p)
			wg.Done()
		}(query)	// 这里传入query就是 q
	}
	time.Sleep(5 * time.Second)

	wg.Wait()
	log.Println("关闭连接池")
	p.Close()
}

func dbQuery(query int, p *pool.Pool) {
	conn, err := p.Acquire()
	if err != nil {
		log.Println(err)
		return
	}
	defer p.Release(conn)		// 这里io.Closer 就是一个单独的 &dbConnection{id}

	// 模拟查询
	time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
	log.Printf("第%d个查询，使用的是Id为%d的数据库连接", query, conn.(*dbConnection).ID)
}

type dbConnection struct {
	ID int32	//连接标记
}

//io.Closer接口类型，需要实现Close()方法
func (d *dbConnection) Close() error {
	log.Println("关闭连接: ", d.ID)
	return nil
}

var idCounter int32

//生成数据库连接方法，以供连接池使用
func createConnection() (io.Closer,error) {
	//并发安全，给数据库连接生成唯一标记
	id := atomic.AddInt32(&idCounter, 1)
	log.Println("生成连接池，id为 ,", id)
	return &dbConnection{id}, nil
}
```

# 修改地方：
```golang
func (p *Pool) Acquire() (io.Closer, error) {
        ...
        default:
				log.Println("Acquire: 新生成资源")
				return p.factory()
				}
	    ...
	    
func (p *Pool) Release(r io.Closer){
        ...
        default:
    		log.Println("资源池满了，释放这个资源吧")
    		_ = r.Close()
		...
```