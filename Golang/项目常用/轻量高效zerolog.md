- [`zerolog`](#zerolog)
- [性能](#性能)
- [基本使用](#基本使用)
- [高级使用](#高级使用)
  - [`StackError` 跟踪所有报错](#stackerror-跟踪所有报错)
  - [使用其他`outpus`](#使用其他outpus)
  - [内嵌字段`sub dictionary`](#内嵌字段sub-dictionary)
  - [改变默认的`field names`](#改变默认的field-names)
  - [增加显示调用的行号`line number to log`](#增加显示调用的行号line-number-to-log)
  - [抽样日志 --每几条打印多少日志](#抽样日志---每几条打印多少日志)
  - [添加`Hooks`](#添加hooks)
  - [导入到标准包`log`中](#导入到标准包log中)
  - [`Related Projects`](#related-projects)
- [`Gin`使用`zerolog`](#gin使用zerolog)
  - [基本原理](#基本原理)
  - [基本使用](#基本使用-1)

# `zerolog` 
官方文档: https://github.com/rs/zerolog

# 性能
反正看`Benchmark`很猛

https://github.com/rs/zerolog#benchmarks

# 基本使用
```golang

package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Print("hello world") // default: debug
	log.Info().Msg("hello world")
	log.Error().Msg("hello world")

	// Fields
	flog := log.Info().Fields(map[string]interface{}{"name": "liukaitao"})
	flog.Msg("hello world2")

	log.Info().Msg("hello world2")

	// Use Str Float
	log.Info().Str("try_string", "true").Msg("")

	// use Enabled
	if e := log.Info(); e.Enabled() {
		// Compute log output only if enabled.
		value := "bar"
		e.Str("foo", value).Msg("some debug message")
	}
}

// Output: {"time":1516134303,"level":"debug","message":"hello world"}
```
# 高级使用
## `StackError` 跟踪所有报错
```golang
package main

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	err := outer()
	log.Error().Stack().Err(err).Msg("")
}

func inner() error {
	return errors.New("seems we have an error here")
}

func middle() error {
	err := inner()
	if err != nil {
		return err
	}
	return nil
}

func outer() error {
	err := middle()
	if err != nil {
		return err
	}
	return nil
}
```
报错信息
`{"level":"error","stack":[{"func":"inner","line":"20","source":"main.go"},{"func":"middle","line":"24","source":"main.go"},{"func":"outer","line":"32","source":"main.go"},{"func":"main","line":"15","source":"main.go"},{"func":"mai
n","line":"250","source":"proc.go"},{"func":"goexit","line":"1571","source":"asm_amd64.s"}],"error":"seems we have an error here","time":1670847252}
`

## 使用其他`outpus`
```golang
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logger.Info().Str("foo", "bar").Msg("hello world")
```

## 内嵌字段`sub dictionary`
```golang
log.Info().
    Str("foo", "bar").
    Dict("dict", zerolog.Dict().
        Str("bar", "baz").
        Int("n", 1),
    ).Msg("hello world")

// Output: {"level":"info","time":1494567715,"foo":"bar","dict":{"bar":"baz","n":1},"message":"hello world"}
```

## 改变默认的`field names`
```golang
zerolog.TimestampFieldName = "t"
zerolog.LevelFieldName = "l"
zerolog.MessageFieldName = "m"

log.Info().Msg("hello world")

// Output: {"l":"info","t":1494567715,"m":"hello world"}
```
## 增加显示调用的行号`line number to log`
```golang
log.Logger = log.With().Caller().Logger()
log.Info().Msg("hello world")

// Output: {"level": "info", "message": "hello world", "caller": "/go/src/your_project/some_file:21"}
```

## 抽样日志 --每几条打印多少日志
```golang
    sampled := log.Sample(&zerolog.BasicSampler{N: 2})
    sampled.Info().Msg("will be logged every 2 messages")
    
    // Output: {"time":1494567715,"level":"info","message":"will be logged every 10 messages"}
    
    // Will let 5 debug messages per period of 1 second.
    // Over 5 debug message, 1 every 100 debug messages are logged.
    // Other levels are not sampled.
    sampled := log.Sample(zerolog.LevelSampler{
        DebugSampler: &zerolog.BurstSampler{
            Burst: 5,
            Period: 1*time.Second,
            NextSampler: &zerolog.BasicSampler{N: 100},
            },
	    })
    sampled.Debug().Msg("hello world")

// Output: {"time":1494567715,"level":"debug","message":"hello world"}
```

## 添加`Hooks`
```golang
type SeverityHook struct{}

func (h SeverityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
    if level != zerolog.NoLevel {
        e.Str("severity", level.String())
    }
}

hooked := log.Hook(SeverityHook{})
hooked.Warn().Msg("")

// Output: {"level":"warn","severity":"warn"}
```

## 导入到标准包`log`中
```golang
log := zerolog.New(os.Stdout).With().
    Str("foo", "bar").
    Logger()

stdlog.SetFlags(0)
stdlog.SetOutput(log)

stdlog.Print("hello world")

// Output: {"foo":"bar","message":"hello world"}
```
## `Related Projects`
grpc-zerolog: Implementation of grpclog.LoggerV2 interface using zerolog

overlog: Implementation of Mapped Diagnostic Context interface using zerolog

zerologr: Implementation of logr.LogSink interface using zerolog

# `Gin`使用`zerolog`
参考项目: https://github.com/gin-contrib/logger

## 基本原理
实现了 `Gin`的`log interface{}`类型, `SetLogger(opts ...Option) gin.HandlerFunc` 方法

参考: https://github.com/gin-contrib/logger/blob/master/logger.go

## 基本使用
```golang
package main

import (
  "fmt"
  "net/http"
  "regexp"
  "time"

  "github.com/gin-contrib/logger"
  "github.com/gin-contrib/requestid"
  "github.com/gin-gonic/gin"
  "github.com/rs/zerolog"
  "github.com/rs/zerolog/log"
)

var rxURL = regexp.MustCompile(`^/regexp\d*`)

func main() {
  r := gin.New()

  // Add a logger middleware, which:
  //   - Logs all requests, like a combined access and error log.
  //   - Logs to stdout.
  // r.Use(logger.SetLogger())

  // Example pong request.
  r.GET("/pong", logger.SetLogger(), func(c *gin.Context) {
    c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
  })

  // Example ping request.
  r.GET("/ping", logger.SetLogger(
    logger.WithSkipPath([]string{"/skip"}),
    logger.WithUTC(true),
    logger.WithSkipPathRegexp(rxURL),
  ), func(c *gin.Context) {
    c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
  })

  // Example skip path request.
  r.GET("/skip", logger.SetLogger(
    logger.WithSkipPath([]string{"/skip"}),
  ), func(c *gin.Context) {
    c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
  })

  // Example skip path request.
  r.GET("/regexp1", logger.SetLogger(
    logger.WithSkipPathRegexp(rxURL),
  ), func(c *gin.Context) {
    c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
  })

  // Example skip path request.
  r.GET("/regexp2", logger.SetLogger(
    logger.WithSkipPathRegexp(rxURL),
  ), func(c *gin.Context) {
    c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
  })

  // add custom fields.
  r.GET("/id", requestid.New(requestid.WithGenerator(func() string {
    return "foobar"
  })), logger.SetLogger(
    logger.WithLogger(func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
      if trace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
        l = l.With().
          Str("trace_id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()).
          Str("span_id", trace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String()).
          Logger()
      }

      return l.With().
        Str("id", requestid.Get(c)).
        Str("foo", "bar").
        Str("path", c.Request.URL.Path).
        Logger()
    }),
  ), func(c *gin.Context) {
    c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
  })

  // Listen and Server in 0.0.0.0:8080
  if err := r.Run(":8080"); err != nil {
    log.Fatal().Msg("can' start server with 8080 port")
  }
}
```