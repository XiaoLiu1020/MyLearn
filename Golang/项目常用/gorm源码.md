[TOC]

# `gorm`源码解析

# `gorm`读取`demo`

```go
import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

type Product struct {
  gorm.Model
  Code string
  Price uint
}


func main() {
    db, err := gorm.Open("mysql", "user:password@(ip:port)/dbname?charset=utf8&parseTime=True&loc=Local")
  defer db.Close()

  var product Product
  db.First(&product, 1) // find product with id 1
}
```

代码流程：

1.  连接数据库`open`
2.  使用`db`句柄读取`Product`数据表内容

## 1. `gorm.Open`

会先判断后面参数`...args`, `switch args[0] case string, case SQLCommon`，

就是直接连接或者直接使用sql语句

返回的`db`结构体

```go
// 初始化DB结构体
db = &DB{
		db:        dbSQL,
		logger:    defaultLogger,
		callbacks: DefaultCallback,
		dialect:   newDialect(dialect, dbSQL),
	}
	db.parent = db
	...
	// Send a ping to make sure the database connection is alive.
	// 利用刚建立的single db 尝试连接数据库
	if d, ok := dbSQL.(*sql.DB); ok {
		if err = d.Ping(); err != nil && ownDbSQL {
			d.Close()
		}
	}


// Gorm中使用的DB对象
type DB struct {
    sync.RWMutex                // 锁
    Value        interface{}
    Error        error
    RowsAffected int64

    // single db
    db                SQLCommon  // 原生db.sql对象，包含query相关的原生方法
    blockGlobalUpdate bool
    logMode           logModeValue
    logger            logger
    search            *search      // 保存搜索的条件where, limit, group，比如调用db.clone()时，会指定search
    values            sync.Map

    // global db
    parent        *DB
    callbacks     *Callback        // 当前sql绑定的函数调用链
    dialect       Dialect           // 不同数据库适配注册sql.db
    singularTable bool
}
```

`gorm`数据库连接

![img](https://pic1.zhimg.com/80/v2-619573624faba96337c354059e29996c_720w.jpg)

1.  打开数据库连接（ps：此处叫连接可能有点不合适，源代码注释是：Open initialize a new db connection, need to import driver first）
2.  `设置db.parent为当前db`
3.  `设置SQLCommon`（会委派给go自带db）

## 2. 进行查询`First`

```go
var product Product
	db.First(&product, 1) // find product with id 1

// First find first record that match given conditions, order by primary key
func (s *DB) First(out interface{}, where ...interface{}) *DB {
	newScope := s.NewScope(out)
	newScope.Search.Limit(1)
	return newScope.Set("gorm:order_by_primary_key", "ASC").
		inlineCondition(where...)
		.callCallbacks(s.parent.callbacks.queries).db
}
```

会操作`DB`使用`s.NewScope(out)`，　生成`newScope`

`NewScope`函数如下:

```go
// NewScope create a scope for current operation
func (s *DB) NewScope(value interface{}) *Scope {
	dbClone := s.clone()
	dbClone.Value = value
	scope := &Scope{db: dbClone, Value: value}
	if s.search != nil {
		scope.Search = s.search.clone()
	} else {
		scope.Search = &search{}
	}
	return scope
}
```

其中调用`dbClone := s.clone()`，生成一个新的`db`

然后再根据传入的`out`，也就是查询`model`，生成一个对应的`scope`

`Scope`结构体如下:

```go
// 包含每一个sql操作的相关信息
type Scope struct {
    Search          *search            // 检索条件在1中是同一个对象
    Value           interface{}     // 保存实体类
    SQL             string            // sql语句
    SQLVars         []interface{}
    db              *DB                // DB对象
    instanceID      string
    primaryKeyField *Field
    skipLeft        bool
    fields          *[]*Field        // 字段
    selectAttrs     *[]string
}

```

这个克隆的db实例，包裹在Scope里面。`在刚才First方法里面，也就是First方法内有效`。所以，业务代码持有的总是最原始的db实例，即通过gorm.Open出来的db实例。

假如，业务代码继续其他db操作。gorm的其他方法（如Find/First/Update等）都会再克隆一个db，“包裹”在scope里面，进行操作。

创建之间调用`newScope.Search.Limit(1)`

使用了`scope`里面的`search`对象，结构如下：

```go
// search 对象存放了所有查询的条件 从名字就能看出来 有where or having 各种条件
type search struct {
    db               *DB
    whereConditions  []map[string]interface{}
    orConditions     []map[string]interface{}
    notConditions    []map[string]interface{}
    havingConditions []map[string]interface{}
    joinConditions   []map[string]interface{}
    initAttrs        []interface{}
    assignAttrs      []interface{}
    selects          map[string]interface{}
    omits            []string
    orders           []interface{}
    preload          []searchPreload
    offset           interface{}
    limit            interface{}
    group            string
    tableName        string
    raw              bool
    Unscoped         bool
    ignoreOrderQuery bool
}
```

调用`Limit`方法很简单：

```go
func (s *search) Limit(limit interface{}) *search {
	s.limit = limit
	return s
}
```

最后调用 `return newScope.Set("gorm:order_by_primary_key", "ASC").inlineCondition(where...).callCallbacks(s.parent.callbacks.queries).db`

很好理解:

​	`Set`会把`name`传入到`s.values.Store(name, value)`中，这里`gorm:order_by_primary_key`是`name`，`ASC`就是`value`

`inlineCondition(where...)` 就是把条件语句`where...`存入`scope`的`whereConditions`中

### `callCallbacks`

callCallback是`逐步对多个Callback发起call，也就是按顺序调用callbacks`。每个Callback做一件事情，比如读取数据库值mapping到struct，级联读取其他值。这样好处是：

1.  callback设计比较简单，做一件事（看下面源码，就指定其实是就是具有相同签名的函数）
2.  callbacks拓展性好，即s.parent.callbacks.queries, s.parent.callbacks.queries, s.parent.callbacks.deletes等执行过程可随意扩展

Callback struct源码：

```go
// Callback is a struct that contains all CRUD callbacks
//   Field `creates` contains callbacks will be call when creating object
//   Field `updates` contains callbacks will be call when updating object
//   Field `deletes` contains callbacks will be call when deleting object
//   Field `queries` contains callbacks will be call when querying object with query methods like Find, First, Related, Association...
//   Field `rowQueries` contains callbacks will be call when querying object with Row, Rows...
//   Field `processors` contains all callback processors, will be used to generate above callbacks in order
type Callback struct {
	creates    []*func(scope *Scope) 
	updates    []*func(scope *Scope)
	deletes    []*func(scope *Scope)
	queries    []*func(scope *Scope)
	rowQueries []*func(scope *Scope)
	processors []*CallbackProcessor
}
```

函数:

```go
// 循环调用传入的functions
func (scope *Scope) callCallbacks(funcs []*func(s *Scope)) *Scope {
    defer func() {
        if err := recover(); err != nil {
            if db, ok := scope.db.db.(sqlTx); ok {
                db.Rollback()
            }
            panic(err)
        }
    }()
    // 使用for循环 调用回调函数
    for _, f := range funcs {
        (*f)(scope)
        if scope.skipLeft {
            break
        }
    }
    return scope
}
```

## 最后查询方法是那个`callback`调用的回调函数

1.  `queryCallback` 方法组成sql语句 调用database/sql 中的query方法在上一篇分析中可以看到 循环rows结果获取数据
2.  `prepareQuerySQL`方法主要是组成sql语句的方法 通过反射获取字段名表明等属性

```go
func queryCallback(scope *Scope) {
    if _, skip := scope.InstanceGet("gorm:skip_query_callback"); skip {
        return
    }

    //we are only preloading relations, dont touch base model
    if _, skip := scope.InstanceGet("gorm:only_preload"); skip {
        return
    }

    defer scope.trace(NowFunc())

    var (
        isSlice, isPtr bool
        resultType     reflect.Type
        results        = scope.IndirectValue()
    )
    // 找到排序字段
    if orderBy, ok := scope.Get("gorm:order_by_primary_key"); ok {
        if primaryField := scope.PrimaryField(); primaryField != nil {
            scope.Search.Order(fmt.Sprintf("%v.%v %v", scope.QuotedTableName(), scope.Quote(primaryField.DBName), orderBy))
        }
    }

    if value, ok := scope.Get("gorm:query_destination"); ok {
        results = indirect(reflect.ValueOf(value))
    }

    if kind := results.Kind(); kind == reflect.Slice {
        isSlice = true
        resultType = results.Type().Elem()
        results.Set(reflect.MakeSlice(results.Type(), 0, 0))

        if resultType.Kind() == reflect.Ptr {
            isPtr = true
            resultType = resultType.Elem()
        }
    } else if kind != reflect.Struct {
        scope.Err(errors.New("unsupported destination, should be slice or struct"))
        return
    }
    // 准备查询语句
    scope.prepareQuerySQL()

    if !scope.HasError() {
        scope.db.RowsAffected = 0
        if str, ok := scope.Get("gorm:query_option"); ok {
            scope.SQL += addExtraSpaceIfExist(fmt.Sprint(str))
        }
        // 调用database/sql 包中的query来查询
        if rows, err := scope.SQLDB().Query(scope.SQL, scope.SQLVars...); scope.Err(err) == nil {
            defer rows.Close()

            columns, _ := rows.Columns()
            // 循环rows 组成对象
            for rows.Next() {
                scope.db.RowsAffected++

                elem := results
                if isSlice {
                    elem = reflect.New(resultType).Elem()
                }

                scope.scan(rows, columns, scope.New(elem.Addr().Interface()).Fields())

                if isSlice {
                    if isPtr {
                        results.Set(reflect.Append(results, elem.Addr()))
                    } else {
                        results.Set(reflect.Append(results, elem))
                    }
                }
            }

            if err := rows.Err(); err != nil {
                scope.Err(err)
            } else if scope.db.RowsAffected == 0 && !isSlice {
                scope.Err(ErrRecordNotFound)
            }
        }
    }
}
func (scope *Scope) prepareQuerySQL() {
    // 如果是rwa 则组织sql语句
    if scope.Search.raw {
        scope.Raw(scope.CombinedConditionSql())
    } else {
        // 组织select 语句
        // scope.selectSQL() 组织select 需要查询的字段
        // scope.QuotedTableName() 获取表名
        // scope.CombinedConditionSql()组织条件语句
        scope.Raw(fmt.Sprintf("SELECT %v FROM %v %v", scope.selectSQL(), scope.QuotedTableName(), scope.CombinedConditionSql()))
    }
    return
}
```

# 总结

这篇文章从一个最简单的where条件和first函数入手了解Gorm主体的流程和主要的对象。其实可以看出Gorm的本质：

1.  创建DB对象，注册mysql连接
2.  创建对象 通过tag设置一些主键，外键等
3.  通过where或者其他比如group having 等设置查询的条件
4.  通过first函数最终生成sql语句
5.  调用database/sql 中的方法通过mysql驱动真正的查询数据
6.  通过反射来组成对象或者是数组对象提供使用

