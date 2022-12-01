# `my_runner`
```golang
package my_runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

//一个执行者，可以执行任何任务，但是这些任务是限制完成的
//该执行者，可以通过发送终止信号终止它
type Runner struct {
	tasks	[]func(int)		//执行任务
	complete chan error		//用于通知任务全部完成
	timeout	<-chan time.Time	//这些任务多久内完成
	interrupt chan os.Signal	//可以强制终止的信号
}

//工厂函数
func New(tm time.Duration) *Runner {
	return &Runner{
		tasks:     nil,
		complete:  make(chan error),	//实例通道, 同步通道
		timeout:   time.After(tm),		//返回 <-chan Time
		interrupt: make(chan os.Signal, 1),	 // 实例通道，容量为1, 缓冲通道，发送信号时不会被阻塞
	}
}

// 添加任务方法, tasks 为一个切片， 可以添加一个，甚至同时多个
// 相当于 打包，解包
func (r *Runner) Add(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

// 用于返回错误类型
var ErrTimeOut = errors.New("执行者执行超时")
var ErrInterrupt = errors.New("执行者被中断")

//执行任务, 中断会返回中断错误
func (r *Runner) run() error{
	for id, task := range r.tasks {
		//tasks为切片，返回Index 和func(int)
		if r.isInterrupt() {
			return ErrInterrupt
		}
		task(id)
	}
	return nil
}

//判断是否中断了,获取 chan 传来的os.Signal
func (r *Runner) isInterrupt() bool {
	select {
	case <- r.interrupt:	//os.Signal发生
		signal.Stop(r.interrupt)
		return true
	default:
		return false
	}
}

// 开始执行所有任务，并且监视通道事件
func (r *Runner) Start() error {
	//希望接受哪些系统信息,注册监听事件
	signal.Notify(r.interrupt, os.Interrupt)
	go func() {
		r.complete <- r.run()
	}()

	select {
	case err := <-r.complete:
		return err
	case <-r.timeout:
		return ErrTimeOut
	}
}
```

# `main.go`
```golang
package main

import (
	runner "liukaitao.com/m/my_runner"
	"log"
	"os"
	"time"
)

func main() {
	log.Println("...Tasks Starting ...")

	timeout := 3 * time.Second

	//使用工厂函数生产出*runner
	r := runner.New(timeout)

	//增加任务
	r.Add(createTask(), createTask(), createTask())

	if err := r.Start(); err != nil {
		switch err {
		case runner.ErrTimeOut:
			log.Println(err)
			os.Exit(1)
		case runner.ErrInterrupt:
			log.Println(err)
			os.Exit(2)
		}
	}
	log.Println("======Tasks finish ======")
	log.Println("======Tasks end ======")
}

func createTask() func(int) {
	return func(id int) {
		log.Printf("the task is doing, id is %d", id)
		time.Sleep(time.Duration(id) * time.Second)
	}
}
```