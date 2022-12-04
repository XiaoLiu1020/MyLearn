- [`context`初始化](#context初始化)
- [`StartServer(ctx *context.Context)`](#startserverctx-contextcontext)
  - [开启`queueMessage`队列`ctx.Queue.Subscribe()`](#开启queuemessage队列ctxqueuesubscribe)
    - [`queue.NewRedis(ctx.RedisConfig, ctx.Logger)` 使用redis起一个列表](#queuenewredisctxredisconfig-ctxlogger-使用redis起一个列表)
    - [消息的结构定义`Message`](#消息的结构定义message)
    - [获取订阅消息的通道](#获取订阅消息的通道)
  - [把`RedisConfig`存入`ctx.Status`, 起`status`中带有一个客户端`client`](#把redisconfig存入ctxstatus-起status中带有一个客户端client)
  - [开启触发任务通道`ctx.Ticker`, 会从`Ticker.tasks`发送`task`出来](#开启触发任务通道ctxticker-会从tickertasks发送task出来)
  - [使用`socketconfig`配置监听`socket`](#使用socketconfig配置监听socket)
- [在开启`StartServer`中继续开启 `contextRoutine`](#在开启startserver中继续开启-contextroutine)
  - [`contextRoutine`,传入一个是接受`msg`的通道，一个是接收`Ticker.Tasks`中出来的`Task`的通道](#contextroutine传入一个是接受msg的通道一个是接收tickertasks中出来的task的通道)
  - [`closeConnection`关闭连接](#closeconnection关闭连接)
  - [`handleQueueMessage(ctx, m)`处理订阅消息](#handlequeuemessagectx-m处理订阅消息)
  - [`sendNormal` 把msg根据协议，打包发送到session的Data通道里，如果有retry，把还需要发送到\*ticker.Ticker.tasks列表中](#sendnormal-把msg根据协议打包发送到session的data通道里如果有retry把还需要发送到tickertickertasks列表中)
- [`server.Accept()`,接受socket连接](#serveraccept接受socket连接)
  - [`go receiveRoutine(ctx, s)` 接受socket信息](#go-receiveroutinectx-s-接受socket信息)
  - [解包协议`Unpack`](#解包协议unpack)
  - [`receive`连接时，执行`addSession`,`s.key`为空才会执行, 请求系统获取key，并添加到ctx.clients](#receive连接时执行addsessionskey为空才会执行-请求系统获取key并添加到ctxclients)
  - [处理好信息，添加到clients中后, 使用`handleClientMessage`](#处理好信息添加到clients中后-使用handleclientmessage)
- [`sendRoutine` 回复socket连接](#sendroutine-回复socket连接)


# `context`初始化

```golang

//上下文对象实例    --第一步
type Context struct {
	*Config
	*log.Logger
	clients map[*session.Session]bool // just for set struct
	Queue   queue.Queue
	Status  *status.Status
	Ticker  *ticker.Ticker
}

ctx := &Context{
    clients: make(map[*session.Session]bool),   //字典：session对象作为键，bool= 1 表示存在
    Config: config,                             //从json文件读取的配置
    Logger: logger,                             //从main传来的logger
    Ticker: ticker.New(logger)                             //Ticker指针类型,使用logger初始化
}

type Ticker struct {
    *log.Logger //日志配置
    sync.Mute   //用于锁
    closed bool // 用于循环中关闭
    tasks *list.List    //任务列表
}

func New(logger *log.Logger) *Ticker {
    return &Ticker{
        Logger: logger,             
        tasks: list.New(),          //Ticker主要存放着tasks列表
    }
}

type Session struct {
    Sn, Key, Operator, Charset string
    In, Out byte
    Conn    net.Conn        //代表着这个session的socket连接
    Data    chan []byte     // 通道，传输session数据, 
    Ack     chan byte       // Ack通道， 存byte类型
    Tasks   map[byte]*ticker.Task   //键为byte, 值存放Task指针类型
}

type Task struct {
    At          time.Time   // 任务发起时间
    Deleted     bool        // 是否删除了
    Interval    int         // 重复的时间间隔
    Retries     int         // 重试次数
    Payload     interface{} // 内容
}

```

# `StartServer(ctx *context.Context)`

## 开启`queueMessage`队列`ctx.Queue.Subscribe()`

### `queue.NewRedis(ctx.RedisConfig, ctx.Logger)` 使用redis起一个列表

```golang
type RedisQueue struct {
	*log.Logger
	*RedisConfig
	client *redis.Client
	closed bool     //用于关闭redis
}

func NewRedis(config *RedisConfig, logger *log.Logger) *RedisQueue{
    return &RedisQueue{
        Logger:         logger,
        ReidsConfig:    config,
        client: redis.NewClient(&redis.Options{     //redis 配置
            Addr:       config.URL,
            DB:         config.Db,
        })
    }
}

ctx.Queue = queue.NewRedis(ctx.RedisConfig, ctx.Logger)
```

### 消息的结构定义`Message`

```golang
type Message struct {
    Operator    string      `json:"operator,omitempty"`
    Sn          string      `json:"sn,omitempty"`
    Action      int         `json:"action"`
    Detail      interface{} `json:"detail"` 
}
```

### 获取订阅消息的通道

```golang
queueMessage := ctx.Queue.Subscribe()

//队列的三个方法
type Queue interface {
	Close()
	Publish(m *Message, verbose bool) error
	Subscribe() <-chan *Message
}

func (r *RedisQueue) Subscribe() <- chan *Message {     //返回的是*Message指针类型
    ch := make(chan *Message, 100)  开启通道
    go func() {
        for !r.closed{
            r.newRunner(ch)     //如果没有关闭队列,就不断循环runner
        }
        close(ch)   //关闭通道
    }()
    return ch       //即返回 newRunner中的data通道，其中内容为redis订阅的Message
}

//作用是不断接受redis订阅传来的消息
func (r *RedisQueue) newRunner(data chan *Message) {
    // sub为订阅的句柄
    sub := r.client.Subscribe(r.Inchannel)
    defer func() {
        sub.Close() //最终需要关闭句柄
    }
    offset := time.Now()        //计算偏移时间
    
    //队列没关闭情况下
    for !r.closed {
        now: time.Now()         //不断计算现在
        //如果现在已经大于 最初开始+连接市场的时间
        if now.After(offset.Add(pingInterval * time.Second)) {
            offset = now        //重新计算偏移，回归到now
            //如果连接不通，返回退出
            if err := sub.Ping(""); err != nil {
                return
            }
        }
        msg,err := sub.ReceiveMessage() //接受信息
        if err....
        
        // 反序列化获取msg.Payload信息载体
        m := new(Message)
        if err := json.Unmarshal([]byte(msg.Payload),m); err != nil{
            continue
        }
        
        //把消息发送到data 通道
        data <- m   
    }
}
```

## 把`RedisConfig`存入`ctx.Status`, 起`status`中带有一个客户端`client`

```golang
ctx.Status = status.New(ctx.RedisConfig)

func New(config *queue.RedisConfig) *Status {
    return &Status{
        client: redis.NewClient(&redis.Options{
            Addr: config.URL,
            DB:   config.Db,
        })
    }
}
```

## 开启触发任务通道`ctx.Ticker`, 会从`Ticker.tasks`发送`task`出来

```golang
tickerTask := ctx.Ticker.Fire()

//返回时一个chan，存放着Task指针类型
func (t *Ticker) Fire() <- chan *Task {
    ch := make(chan *Task, 100)
    go func() {
        for !t.closed {
            time.Sleep(time.Second)
            t.handler(ch)   //处理其中通道,go协程开启
        }
        close(ch)
    }()
    return ch
}

func (t *Ticker) handler(fire chan *Task) {
    now := time.Now()
    t.Lock()    //加锁,互斥锁   
    
    //计算任务列表开始长度
    begin := t.tasks.Len()  //初始为0
    e := t.tasks.Front()    //取出列表最前面那个
    
    //第一个任务不为空时开始处理
    for e != nil {          
        //  把列表第一个的值取出，并且实例化为task类型
        task := e.Value.(*Task)
        
        //如果任务已经删除了
        if task.Deleted {   
            d := e.Next()   //获取e任务的下一个，赋值给d
            //在ticker中的任务列表tasks删除e
            t.tasks.Remove(e)       
            e = d   //删除后，d 赋值给e
            continue
        }
        //如果任务时间是在现在之前的，传送到通道，重试次数会触发重试并且放在后面
        if !task.At.After(now) {
            d := e.Next()
            t.tasks.Remove(e)
            e = d 
            task.Retries--      //任务的重试次数减少一次
            
            // 发送任务到 fire chan *Task
            fire <- task
            
            // 如果任务重试次数大于0，At时间也延长,并且放到后面
            if task.Retries > 0 {
                task.At = task.At.Add(time.Duration(task.Interval) * time.Second)
                t.tasks.PushBack(task)
            }
        // 在这之后的就下一个
        } else {
            e = e.Next()
        }
    }
    //直到获取到任务为空为止
    end := t.tasls.Len()
    // 输出任务减少，从开始到结束
    if begin != end {
        t.Logger.Printf("T-I tasks: %d -> %d\n", begin, end)
    }
    t.Unlock()
}
```

## 使用`socketconfig`配置监听`socket`

```golang
socket := ctx.SocketConfig

server, err := net.Listen(socket.Network, socket.URL)

//最后都会close
defer func() {
    server.Close()
    ctx.Status.Close()  //redis.client
    ctx.Queue.Close()   //redis队列
    ctx.Logger.Println("S-I", "server closed")
}()
```

# 在开启`StartServer`中继续开启 `contextRoutine`

```golang
go contextRoutine(ctx, queueMessage, tickerTask)
```

## `contextRoutine`,传入一个是接受`msg`的通道，一个是接收`Ticker.Tasks`中出来的`Task`的通道

```golang
type messageTask struct {
	message []byte
	session *session.Session
	id      byte
	detail  []byte
	action  int
}

func contextRoutine(ctx *context.Context, q <-chan *queue.Message, t <-chan *ticker.Task)
    //循环,接收
    for {
        select {
        // ticker.Tasks传递出来的task经过t发出来的任务
        case    task, ok := <-t:
            if !ok {
                return
            }
            //把task的内容Payload实例化变为messageTask
            m := task.Payload.(*messageTask)
            
            // 如果这个任务重试次数小于1，关闭这个消息任务对应session
            if task.Retries < 1 {
                closeConnection(ctx, m.session)
            } else {
                //如果还有重试次数, 把消息发给session.Data通道中
                m.session.Data <- m.message
            }
        //收到订阅来的message
        case m, ok := <-q
            if !ok {
                return
            }
            //处理消息
            handleQueueMessage(ctx, m)
        }
    }

```

## `closeConnection`关闭连接

```golang
// 断开这个session的连接, Stop sync.Once 只会做一次
func closeConnetcion(ctx *context.Context, s *session.Session) {
    s.Stop.Do(func() {
        s.Close()
        ctx.Remove(s)
    })
}

// 关闭session中两个通道
func (s *Session) Close() {
	close(s.Data)   
	close(s.Ack)
	for _, task := range s.Tasks {
		task.Deleted = true
	}
}
//把 context中的客户端连接，clients中对应*session.Session移除
func (ctx *Context) Remove(s *session.Session) {
    delete(ctx.clients, s)
}
```

## `handleQueueMessage(ctx, m)`处理订阅消息

```golang
func handleQueueMessage(ctx *context.Context, m *queue.Message) {
    //处理消息的detail
    detail, _: json.Marshal(m.Detail)
    //处理动作,action, 使用notify
    switch m.Action {
    case actionDelivery:
		notify(ctx, m.Sn, m.Action, detail, validateDelivery)
	........
	default:
		ctx.Logger.Println("Q-E", m.Sn, "action", m.Action)
    }
}

func notify(ctx *context.Context, sn string, action int, detail []byte, check func(*context.Context, string, []byte) bool) {
    //大写字母
	upper := strings.ToUpper(sn)
	//检测detail内容是否符合其action对应的内容格式
	if check != nil && !check(ctx, upper, detail) {
		return
	}
	// 判断ctx.clients是否有此机台， 有就返回其*session.Session对象
	s := ctx.Session(upper)
	if s == nil {
		ctx.Logger.Println("Q-W", upper, "discard without session")
		return
	}
	// 发送消息，内容包含 *session.Session, sn, action, detail
	sendNormal(ctx, s, upper, action, detail)
}

//判断这个机台号是否在ctx.clients中,有就返回其 *session.Session对象
func (ctx *Context) Session(sn string) *session.Session {
    for k := range ctx.clients {
        if k.Sn == sn {
            return k    //这里返回的是clients中的键,键为 *session.Session类型
        }
    }
    return nil
}
```

## `sendNormal` 把msg根据协议，打包发送到session的Data通道里，如果有retry，把还需要发送到\*ticker.Ticker.tasks列表中

```golang
func sendNormal(ctx *context.Context, s *session.Session, sn string, action int, detail []byte) {
	send(ctx, s, sn, action, detail, retry: true, verbose:true)
}

func send(ctx *context.Context, s *session.Session. sn string, action int, detail []byte, retry bool, v bool) {
    id := s.Out //session的编号
    s.Out++
    
    b := encode(s, detail)  //根据Session的charset改变，返回detail
    
    //获取协议,并进行消息打包
    p := ctx.SocketConfig.Protocol
    m := protocol.Pack(p, id, sn, action, b, s.Key)
    
    // 继续细节
    if v {
		ctx.Logger.Println("S>I", sn, id, action, string(detail), 0)
	}
	
	// 判断ctx.clients是否有这个*session.Session
	if ctx.Has(s) {
	    s.Data <- m     // 有的话，就把消息往这个session的Data chan里发送
	}
	// 如果是重新发送
	// 实例一个ticker.Task,放入 ctx.Ticker.
	if retry {
	    task := &ticker.Task{
			At:       time.Now().Add(time.Duration(p.Timeout) * time.Second),
			Interval: p.Timeout,
			Retries:  p.Fails,
			Payload: &messageTask{
				message: m,
				session: s,
				id:      id,
				detail:  detail,
				action:  action,
			},
		}
		//在session的Tasks中添加键为id,值为*ticker.Task的实例
		s.Add(id, task)     
		// 把*ticker.Task实例添加到ctx.Ticker.tasks的列表中
		ctx.Ticker.Add(task)
	}
}
```

# `server.Accept()`,接受socket连接

```golang
// new 一个session 对象
func New(conn net.Conn) *Session {
	return &Session{
		Conn:    conn,
		Data:    make(chan []byte),
		Ack:     make(chan byte),
		Tasks:   make(map[byte]*ticker.Task),
		Charset: "utf-8",
	}
}

for {
    conn, err := server.Accept()
    
    s := session.New(conn)  //用连接去实例化session
    
    go receiveRoutine(ctx, s)
    go sendRoutine(ctx, s)
}
```

## `go receiveRoutine(ctx, s)` 接受socket信息

```golang
func receiveRoutine(ctx *context.Context, s *session.Session) {
    //最终关闭这个session
    defer closeConnection(ctx, s)
    
    // 获取协议
    p := ctx.SocketConfig.Protocol
    
    //分配内存空间
    read := make([]byte, p.Maximum*2)
    for {
        // 连接超时设置，过时conn.Read返回空
		deadline := time.Now().Add(time.Second * time.Duration(p.Timeout*p.Fails))
        s.Conn.SetReadDeadline(deadline)
        // 读取Conn内容到read
        n, err := s.Conn.Read(read)
        if err != nil {
            return
        }
        //如果 n 有内容
        for end :=0; end < n; {
            // 使用协议解包，读取    
   			message, offset, err := protocol.Unpack(p, read, end, n, func(sn string, id byte) (string, error) {
   			        // 初始化，key为空
                    if len(s.Key) == 0 {
                        s.Sn = sn
                        if err := addSession(ctx, s); err != nil{
                            return
                        } else if id == s.In {
                            return
                        }
                    }
                    // socket传过来的id 设置为 s.In
                    s.In = id
                    return s.Key, nil
        })
        // end < n 条件不合符就会推出
        end = offset
		if err != nil {
			ctx.Logger.Println("S-W", s.Sn, "discard", err)
		} else {
			if ctx.Verbose || message.Action != actionStatus {
				ctx.Logger.Println("S<I", message.Sn, message.ID, message.Action, string(message.Details))
			}
			handleClientMessage(ctx, s, message)
		}
    }
}
```

## 解包协议`Unpack`

```golang
func Unpack(p *Protocol, buffer []byte, begin int, end int, keyHandler func(string, byte) (string, error)) (*Message, int, error) {
	if p.Minimum > end-begin {
		// discard less than minimum length
		return nil, end, fmt.Errorf("too short: %s", buffer[begin:end])
	}

	offset := match(buffer, begin, end, []byte(p.Header))
	if offset < 0 {
		// discard without header
		return nil, end, fmt.Errorf("no header: %s", buffer[begin:end])
	}
	if p.Minimum > end-offset {
		// discard less than minimum length
		return nil, end, fmt.Errorf("too short: %s", buffer[begin:end])
	}
	length, err := fromHex(buffer, offset, p.Length)
	if err != nil || offset+length+p.Plain.Begin > end {
		// discard length error
		return nil, end, fmt.Errorf("length: %s", buffer[begin:end])
	}

	messageEnd := offset + length + p.Plain.Begin
	id, err := fromHex(buffer, offset, p.ID)
	if err != nil {
		// discard with error id
		return nil, messageEnd, fmt.Errorf("id: %s", buffer[begin:end])
	}
	sn := strings.ToUpper(string(slice(buffer, offset, p.Sn)))
	key, err := keyHandler(sn, byte(id))
	if err != nil {
		return nil, messageEnd, fmt.Errorf("%s: %s", err.Error(), buffer[begin:end])
	}

	plain := buffer[offset+p.Plain.Begin : messageEnd]
	sign := slice(buffer, offset, p.Sign)
	if !bytes.Equal(sign, getSign(plain, key)) {
		// discard with signature error
		return nil, messageEnd, fmt.Errorf("sign: %s", buffer[begin:end])
	}
	timestamp, err := fromHex(buffer, offset, p.Timestamp)
	if err != nil {
		// discard with timestamp error
		return nil, messageEnd, fmt.Errorf("timestamp: %s", buffer[begin:end])
	}
	action, err := fromHex(buffer, offset, p.Action)
	if err != nil {
		// discard with action error
		return nil, messageEnd, fmt.Errorf("action: %s", buffer[begin:end])
	}
	detail := buffer[offset+p.Detail.Begin : messageEnd]

	// message success
	return &Message{
		ID:        byte(id),
		Sn:        sn,
		Action:    action,
		Timestamp: timestamp,
		Details:   detail,
	}, messageEnd, nil
}
```

## `receive`连接时，执行`addSession`,`s.key`为空才会执行, 请求系统获取key，并添加到ctx.clients

```golang
func addSession(ctx *context.Context, s *session.Session) error {
	if old := ctx.Session(s.Sn); old != nil {
		closeConnection(ctx, old)
	}

	if len(ctx.SocketConfig.Secret) > 0 {
		s.Key = ctx.SocketConfig.Secret
	} else {
		key, err := device.QueryKey(ctx.DeviceConfig, s.Sn)
		if err != nil {
			ctx.Logger.Println("R>E", "query key", s.Sn, err)
			return err
		}
		s.Key = key
	}

	if device.IsProduct(s.Sn) {
		operator, err := device.QueryOperator(ctx.DeviceConfig, s.Sn)
		if err != nil {
			ctx.Logger.Println("R>E", "query operator", s.Sn, err)
			return err
		}
		s.Operator = operator
	}

	ctx.Add(s)
	return nil
}
```

## 处理好信息，添加到clients中后, 使用`handleClientMessage`

```golang
//传入 上下文，session对象，和解包的message
func handleClientMessage(ctx *context.Context, s *session.Session, m *protocol.Message)
    switch m.Action {
    case case actionStatus:
		handleOrderMessage(ctx, s, m)
	default .....
    }

func handleOrderMessage(ctx *context.Context, s *session.Session, m *protocol.Message)
{
    items, err := validator.Valid(m.Details, map[string]validator.Item{
		"position": {Required: true, Type: "array[string]"},
		"card":     {Required: true, Type: "string"},
		"posted":   {Required: false, Type: "number", Between: "0,1"},
	})
	if err != nil {
		sendAck(ctx, s, m, "order: "+err.Error())
		return
	}
	// ctx.Status是一个 redis.Client 连接redis
	if err = ctx.Status.UpdateAt(m.Sn); err != nil {
		sendAck(ctx, s, m, "order: "+err.Error())
		return
	}
	addMessageQueueWhileAck(ctx, s, m, items)
}

func addMessageQueueWhileAck(ctx *context.Context, s *session.Session, m *protocol.Message, detail interface{}) {
	addMessageQueueWhileAckWithVerbose(ctx, s, m, detail, verbose: true)
}


func addMessageQueueWhileAckWithVerbose(ctx *context.Context, s *session.Session, m *protocol.Message, detail interface{}, verbose bool) {
	if err := addMessageQueue(ctx, m.Action, m.Sn, detail); err != nil {
		sendAckWithVerbose(ctx, s, m, "server error: queue is not exists", verbose)
	} else {
		sendAckWithVerbose(ctx, s, m, "", verbose)
	}
}


//主要是接受了信息，根据action 使用ctx.Queue.Publish发布
func addMessageQueue(ctx *context.Context, action int, sn string, detail interface{}) error {
	if ctx.Queue == nil {
		return errors.New("queue is nil")
	}
	m := &queue.Message{
		Sn:     strings.ToUpper(sn),
		Action: action,
		Detail: detail,
	}
	if err := ctx.Queue.Publish(m, ctx.Verbose || (m.Action != actionStatus)); err == nil {
		return err
	}
	return nil
}

// 发送ack
func sendAckWithVerbose(ctx *context.Context, s *session.Session, m *protocol.Message, err string, v bool) {
	if len(err) > 0 {
		ctx.Logger.Println("S-E", m.Sn, err, string(m.Details))
	}
	ack := &ack{
		ID:      m.ID,
		Message: err,
	}
	body, _ := json.Marshal(ack)
    // 根据协议包起来， s.Data <-m 把消息发送到s.Data中
	send(ctx, s, m.Sn, actionAck, body, false, v)
}
```

# `sendRoutine` 回复socket连接

```golang
func sendRoutine(ctx *context.Context, s *session.Session) {
    defer closeConnection(ctx, s)
    for {
        // 接受上面收到信息后发送过来的s.Data信息，ack回复
        select {
        case m, ok := <-s.Data:
            if !ok {
                return
            }
            if len(m) > 0 {
                if _,err := s.Conn.Write(m); err != nil {
                    return
                }
            }
        //收到回复正确的Ack后，删除重试的任务，s.Remove(id) id 为任务id
        case id, ok := <-s.Ack:
            if !ok {
                return
            }
            s.Remove(id)
        }
    }
}

// Remove removes the task at session s with id.
func (s *Session) Remove(id byte) {
	if task, ok := s.Tasks[id]; ok {
		task.Deleted = true
	}
}
```

