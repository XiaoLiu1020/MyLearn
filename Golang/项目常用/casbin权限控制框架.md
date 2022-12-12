- [Golang访问权限控制框架`casbin`](#golang访问权限控制框架casbin)
- [特点:](#特点)
- [基本工作原理](#基本工作原理)
- [开始使用](#开始使用)
	- [创建`model.conf`](#创建modelconf)
	- [创建`policy`文件-这里使用本地`policy.csv`](#创建policy文件-这里使用本地policycsv)
	- [基本使用`main.go`](#基本使用maingo)
	- [其他模型的教程文档\&例子模板](#其他模型的教程文档例子模板)
- [高阶使用](#高阶使用)
	- [超级管理员实现](#超级管理员实现)
	- [`g = _, _` 的用法](#g--_-_-的用法)
	- [多个 `g _, _`](#多个-g-_-_)
	- [多层角色](#多层角色)
	- [`domain`领域](#domain领域)
	- [动态控制可读可写 -- 比如根据时间控制](#动态控制可读可写----比如根据时间控制)
	- [动态初始化 `model.conf`(代码内实现model)](#动态初始化-modelconf代码内实现model)
- [使用`gorm-adapter`--策略储存`policy`](#使用gorm-adapter--策略储存policy)
	- [基本使用](#基本使用)
- [自己项目中的运用`rbac.go`](#自己项目中的运用rbacgo)

# Golang访问权限控制框架`casbin`
官方文档: https://www.bookstack.cn/read/Casbin-zh/1.md

以下部分内容来自官方文档

# 特点:
- 支持自定义请求的格式，默认的请求格式为`{subject, object, action}`。 
- 具有访问控制模型`model`和`策略policy`两个核心概念。 
- 支持`RBAC中的多层角色继承`，不止主体可以有角色，资源也可以具有角色。 
- 支持超级用户，如 `root 或 Administrator`，超级用户可以不受授权策略的约束访问任意资源。 
- 支持`多种内置的操作符`，如 keyMatch，方便对路径式的资源进行管理，如 `/foo/bar 可以映射到 /foo*`
- Casbin 不做的事情:
  - 身份认证 authentication（即验证用户的用户名、密码），casbin只负责访问控制。应该有其他专门的组件负责身份认证，然后由casbin进行访问控制，二者是相互配合的关系。 
  - 管理用户列表或角色列表。 Casbin 认为由项目自身来管理用户、角色列表更为合适， 用户通常有他们的密码，但是 Casbin 的设计思想并不是把它作为一个存储密码的容器。 而是存储RBAC方案中用户和角色之间的映射关系。

# 基本工作原理
Casbin 中, 访问控制模型被抽象为基于 `PERM (Policy, Effect, Request, Matcher)` 的一个文件。

这些都能在`model.conf`定义

因此，切换或升级项目的授权机制与修改配置一样简单。 您可以通过组合可用的模型来定制您自己的访问控制模型。

# 开始使用
## 创建`model.conf`
`model` 语法: https://www.bookstack.cn/read/Casbin-zh/7.md#Request%E5%AE%9A%E4%B9%89
```conf
# 请求定义
[request_definition]
r = sub,obj,act
# sub ——> 想要访问资源的用户角色(Subject)——请求实体
# obj ——> 访问的资源(Object)
# act ——> 访问的方法(Action: get、post...)


# 策略定义
# 策略(.csv文件p的格式，定义的每一行为policy rule;p,p2为policy rule的名字。)
[policy_definition]
p = sub,obj,act
# p2 = sub,act 表示sub对所有资源都能执行act


# 组定义
[role_definition]
g = _, _
# g = _,_定义了用户——角色，角色——角色的映射关系，前者是后者的成员，拥有后者的权限。
# _,_表示用户，角色/用户组


# 策略效果
[policy_effect]
e = some(where (p.eft == allow))
# 上面表示有任意一条 policy rule 满足, 则最终结果为 allow；p.eft它可以是allow或deny，它是可选的，默认是allow

# 匹配器
[matchers]
#m = r.sub == p.sub && r.obj == p.obj && r.act == p.act

# 分组的权限校验
m = g(r.sub, p.sub) && g(r.obj, p.obj) && r.act == p.act

# 上面模型文件规定了权限由sub,obj,act三要素组成，只有在策略列表中有和它完全相同的策略时，该请求才能通过。
```

## 创建`policy`文件-这里使用本地`policy.csv`

```csv
p,liu,data1,read
p,yan,data2,write
p,leader,data3,read

g,liu,leader    
```
## 基本使用`main.go`
```golang
package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
)

func check(e *casbin.Enforcer, subject, object, action string) {
	// 权限校验
	e.EnableLog(false)
	ok, err := e.Enforce(subject, object, action)
	if err!= nil {
        fmt.Println("err: ", err)
		return
    }
	if ok {
		fmt.Printf("%s CAN %s %s \n", subject, action, object)
	} else {
		fmt.Printf("%s CANNOT %s %s\n", subject, action, object)
	}
}

// 参考文档: https://www.bookstack.cn/read/Casbin-zh/3.md
// 访问例子example: https://www.bookstack.cn/read/Casbin-zh/6.md
// model 语法: https://www.bookstack.cn/read/Casbin-zh/7.md
// policy可以存储在mysql, gorm框架可以使用以下适配器: https://github.com/casbin/gorm-adapter

func main() {
	e, err := casbin.NewEnforcer("./model.conf", "./policy.csv")
	if err!= nil {
        panic(err)
    }
	check(e, "liu", "data1", "read")
	check(e, "yan", "data2", "write")
	check(e, "liu", "data2", "read")	// false
	check(e, "liu", "data3", "read")
	check(e, "liu", "data3", "write")	// false
	check(e, "yan", "data3", "read")	// false
	
	// 查看Roles
	fmt.Println("=================================")
	roles := e.GetAllRoles()
	fmt.Println("roles.length :", len(roles))
	for _, role := range roles {
		fmt.Println("role: ", role)     // leader, 因为policy文件有 g, liu, leader; liu 属于leader组
	}
}
```

## 其他模型的教程文档&例子模板
可以用其他教程: https://www.bookstack.cn/read/Casbin-zh/4.md#PERM%E5%85%83%E6%A8%A1%E5%9E%8B%20(Policy,%20Effect,%20Request,%20Matcher)

例子: 访问例子example: https://www.bookstack.cn/read/Casbin-zh/6.md

# 高阶使用
部分参考: https://blog.csdn.net/qq_39280718/article/details/126743310
## 超级管理员实现
```bash
[matchers]
# 允许sub == root 可以操作任何
e = r.sub == p.sub && r.obj == p.obj && r.act == p.act || r.sub == "root"
```
## `g = _, _` 的用法
`g = _, _` 定义了用户——角色或角色——角色的映射关系，前者是后者的成员，拥有后者的权限。
```casbin
# model.conf
[matchers]
# 组用户 g(r.sub, p.sub)
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```
这里csv中定义分了两组
```casbincsv
p, admin, data, read
p, admin, data, write
p, developer, data, read

g, zxp, admin
g, zhang, developer
```
## 多个 `g _, _`
`model更改`
```casbin
[role_definition]
g=_,_
g2=_,_

[matchers]
m = g(r.sub, p.sub) && g2(r.obj, p.obj) && r.act == p.act
# 只改这两个，其他不变
```
police举例
```casbincsv
p, admin, data1, read
p, admin, data1, write
p, admin, data2, read
p, admin, data2, write
p, developer, data2, read
p, developer, data2, write
p, developer, data1, read
g, zxp, admin
g, zhang, developer
g2, data1.data, data1
g2, data2.data, data2
```

## 多层角色
只更改`policy`, zxp属于两个角色
```casbincsv
p, senior, data, write
p, developer, data, read

g, zxp, senior
g, senior, developer
g, zhang, developer
```

## `domain`领域
就是多了个`dom`去使用,其他一样
```casbin
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _,_,_
# g2 = _,_,_ 表示用户, 角色/用户组, 域(也就是租户)


[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
```

## 动态控制可读可写 -- 比如根据时间控制
该模式下,可以不需要`policy`文件,因为只涉及带request_definition的时间
```casbin
# model.conf, 用到了Hour
[matchers]
m = r.sub.Hour >= 5 && r.sub.Hour < 20 || r.sub.Name == r.obj.Owner
```
```golang
type Object struct {
  Name  string
  Owner string
}

type Subject struct {
  Name string
  Hour int
}

func check(e *casbin.Enforcer, sub Subject, obj Object, act string) {
  ok, _ := e.Enforce(sub, obj, act)
  if ok {
    fmt.Printf("%s CAN %s %s at %d:00\n", sub.Name, act, obj.Name, sub.Hour)
  } else {
    fmt.Printf("%s CANNOT %s %s at %d:00\n", sub.Name, act, obj.Name, sub.Hour)
  }
}

func main() {
  e, err := casbin.NewEnforcer("./model.conf", "./policy.csv")
  if err != nil {
    log.Fatalf("NewEnforecer failed:%v\n", err)
  }

  //r.sub.Hour < 18 || r.sub.Name == r.obj.Owner 这两个满足一个就可读
  o := Object{"data", "zxp"}
    
    
  s1 := Subject{"zxp", 10}
  check(e, s1, o, "read")//可读

  s2 := Subject{"zhang", 10}
  check(e, s2, o, "read")//可读

  s4 := Subject{"zhang", 20}
  check(e, s4, o, "read")//不可读
    
}
```

## 动态初始化 `model.conf`(代码内实现model)
```golang
package main

import (
	"fmt"
	"log"

	fileadapter "github.com/casbin/casbin/persist/file-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

func check(e *casbin.Enforcer, sub, obj, act string) {
	ok, _ := e.Enforce(sub, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s CANNOT %s %s\n", sub, act, obj)
	}
}

//第一种写法
func main() {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act") // AddDef  增加定义, r=request, p=policy, e=effect, m=matchers
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "r.sub == g.sub && r.obj == p.obj && r.act == p.act")

	a := fileadapter.NewAdapter("./policy.csv")
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}

	check(e, "zxp", "data1", "read")
	check(e, "zhang", "data2", "write")
	check(e, "zxp", "data1", "write")
	check(e, "zxp", "data2", "read")
}

//第二种写法
// func main() {
// 	text := `
//   [request_definition]
//   r = sub, obj, act

//   [policy_definition]
//   p = sub, obj, act

//   [policy_effect]
//   e = some(where (p.eft == allow))

//   [matchers]
//   m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
//   `

// 	m, _ := model.NewModelFromString(text)
// 	a := fileadapter.NewAdapter("./policy.csv")
// 	e, _ := casbin.NewEnforcer(m, a)
// }
```

# 使用`gorm-adapter`--策略储存`policy`
参考文档: https://github.com/casbin/gorm-adapter

更多adapter: https://www.bookstack.cn/read/Casbin-zh/10.md#MySQL%20adapter
## 基本使用 
```golang
package main

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Initialize a Gorm adapter and use it in a Casbin enforcer:
	// The adapter will use the MySQL database named "casbin".
	// If it doesn't exist, the adapter will create it automatically.
	// You can also use an already existing gorm instance with gormadapter.NewAdapterByDB(gormInstance)
	a, _ := gormadapter.NewAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/") // Your driver and data source.
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf", a)
	
	// Or you can use an existing DB "abc" like this:
	// The adapter will use the table named "casbin_rule".
	// If it doesn't exist, the adapter will create it automatically.
	// a := gormadapter.NewAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/abc", true)

	// Load the policy from DB.
	e.LoadPolicy()
	
	// Check the permission.
	e.Enforce("alice", "data1", "read")
	
	// Modify the policy.
	// e.AddPolicy(...)
	// e.RemovePolicy(...)
	
	// Save the policy back to DB.
	e.SavePolicy()
}
```

# 自己项目中的运用`rbac.go`
里面结合了jwt使用
```golang
package middlewares

import (
	"github.com/casbin/casbin"
	"github.com/labstack/echo/v4"
	"saiyaoyun.com/piece/cache"
	"saiyaoyun.com/piece/components/e"
	"saiyaoyun.com/piece/model"
)

type RbacConf struct {
	ModelFile  string `json:"model_file"`
	PolicyFile string `json:"policy_file"`
	Debug      bool   `json:"debug"`
}

type Rbac struct {
	*RbacConf
	*casbin.Enforcer
}

func NewRbac(config *RbacConf) *Rbac {
	return &Rbac{
		RbacConf: config,
		Enforcer: casbin.NewEnforcer(config.ModelFile, config.PolicyFile, config.Debug),
	}
}

func Authorize(cache *cache.RedisConn, dao *model.Dao, r *Rbac) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			obj := c.Request().URL.Path
			act := c.Request().Method
			claims, err := GetUserClaims(c, cache, dao) //  判断是否授权, 成功后获取Role角色-- subject
			if err != nil {
				return err
			}
			sub := claims.Role
			if ok := r.Enforce(sub, obj, act); !ok {
				return e.NewServerError(e.ErrorNoPermissions)
			}
			return next(c)
		}
	}
}
// 下面是 获取上下文的user并且 判断是否授权, 授权成功则继续判断角色访问
func GetUserClaims(c echo.Context, cache *cache.RedisConn, dao *model.Dao) (*JwtCustomClaims, error) {
  user := c.Get("user")
  if user == nil {
    return nil, e.NewServerError(e.ErrorLoginRequired)
  }
  switch user.(type) {
  case *jwt.Token:
    token := user.(*jwt.Token)
    claims := token.Claims.(*JwtCustomClaims)
    if !verifyUserState(cache, dao, claims.ID) {
      return nil, e.NewServerError(e.ErrorUserInvalid)
    }
    return claims, nil
  }
  return nil, e.NewServerError(e.ErrorTokenInvalid)
}
```

