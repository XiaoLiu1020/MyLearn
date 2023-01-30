- [`Gin` 使用`validator`](#gin-使用validator)
- [`gin.binding` 自定义校验](#ginbinding-自定义校验)
- [用法](#用法)
- [使用](#使用)
	- [自定义校验的函数](#自定义校验的函数)
	- [实现对应tag的错误返回信息](#实现对应tag的错误返回信息)
	- [错误捕获使用](#错误捕获使用)

# `Gin` 使用`validator`
基本使用可以查看推荐博客:  Gin框架使用validator进行数据校验及自定义翻译器

https://blog.csdn.net/wxl095/article/details/124533618

# `gin.binding` 自定义校验

参考包地址: `"github.com/gin-gonic/gin", "github.com/gin-gonic/gin/binding"`

有时候我们需要对参数进行统一规范的校验, 在 `gin`框架基础上设置自定义的校验器, 进行 `binding`的时候只需要给 `Request Struct`的 `tag`加上特定标识

就可以使用 `ctx.ShouldBing(param)`完成统一过滤校验了

还使用了 `go-playground/validator`库 进行校验,加入到 `binding`里面

# 用法

```golang
type Param struct {
	TagName    string `json:"tagName" binding:"required,name"`             // 文档规范里面的字段校验
	OtherField string `json:"otherField" binding:"required,max=10,min=2"` // 其他个别字段的长度校验
}
```

# 使用

```go
    ...
    import (
  "github.com/gin-gonic/gin/binding"
     "github.com/go-playground/validator/v10"
    )
    ...
 if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
  RegisterValidation(v)
 }
 ...
 
  
func RegisterValidation(v v *validator.Validate) { 
 // 注册自定义字段校验
 registerCustomFields(v)
}

// registerCustomFields 注册自定义字段校验
func registerCustomFields(v *validator.Validate) {
    _ = v.RegisterValidation("name", stringCountValidator(2, 30))    // 名字类字段
    _ = v.RegisterValidation("desc", stringCountValidator(2, 100))   // 描述类字段
    _ = v.RegisterValidation("text", stringCountValidator(0, 1000))  // 长文本字段
    _ = v.RegisterValidation("remark", stringCountValidator(0, 500)) // 备注字段
    _ = v.RegisterValidation("role", stringCountValidator(2, 30))    // 角色名称
    _ = v.RegisterValidation("tag", tagName)                         // 标签名字段
    _ = v.RegisterValidation("usergroup", userGroupName)             // 用户组名字段
    _ = v.RegisterValidation("username", userName)                   // 用户名字段
    _ = v.RegisterValidation("account", account)                     // 账号字段
    _ = v.RegisterValidation("nickname", nickName)                   // 姓名字段
    _ = v.RegisterValidation("password", password)                   // 密码字段
}
```

## 自定义校验的函数

```golang
// stringValidator 字符长度校验
func stringCountValidator(min, max int) func(fl validator.FieldLevel) bool {
    return func(fl validator.FieldLevel) bool {
        value, ok := fl.Field().Interface().(string)
        if !ok {
            return false
        }
        return checkStringCount(value, min, max)
    }
}

func checkStringCount(s string, min int, max int) bool {
    count := utf8.RuneCountInString(s)
    if count == 0 {
        // 不校验必填项 由required标签负责
        return true
        }
    if count < min || count > max {
        return false
    }
    return true
}
// tagName 封装多一层判断是否是标签名称
var tagName validator.Func = func(fl validator.FieldLevel) bool {
    value, ok := fl.Field().Interface().(string)
    if !ok {
        return false
    }
    if len(value) == 0 {
        return true
    }
    b := regTag.MatchString(value)
        if !b {
            return false
        }
    return checkStringCount(value, 2, 10)
}

// password 判断是否是密码,强度验证
var password validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if len(value) == 0 {
		return true
	}

	var (
		hasLetter = false
		hasNumber = false
		hasSymbol = false
	)

	for _, s := range value {
		switch {
		case unicode.IsUpper(s) || unicode.IsLower(s):
			hasLetter = true
		case unicode.IsNumber(s):
			hasNumber = true
		case unicode.IsPunct(s) || unicode.IsSymbol(s):
			hasSymbol = true
		default:
			return false
		}
	}
	count := 0
	if hasSymbol {
		count += 1
	}
	if hasLetter {
		count += 1
	}
	if hasNumber {
		count += 1
	}

	if count < 2 {
		return false
	}
	return checkStringCount(value, 8, 72)
}

```

## 实现对应tag的错误返回信息

```go
import (
    zhtranslations "github.com/go-playground/validator/v10/translations/zh"
    ut "github.com/go-playground/universal-translator"
    "github.com/go-playground/locales/zh"
    "github.com/go-playground/validator/v10"
)

var zhTranslations ut.Translator

func RegisterValidation(v *validator.Validate) {
    // 这里注册中文转换
	registerTranslator(v)
}

// registerTranslator 定义中文校验转换信息
func registerTranslator(v *validator.Validate) {
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
        // json信息处理
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
    
    // 新增zh转换
	zhTrans := zh.New()
	uniTranslator := ut.New(zhTrans, zhTrans)   // unicode translator
	zhTranslator, _ = uniTranslator.GetTranslator("zh")
    // 把中文translator 注册进 validator
	_ = zhtranslations.RegisterDefaultTranslations(v, zhTranslator)
	registerCustomFieldTranslator(v)
}
    // 注册字段
func registerCustomFieldTranslator(v *validator.Validate) {
	translations := []struct {
		tag             string
		translation     string
		override        bool
		customRegisFunc validator.RegisterTranslationsFunc
		customTransFunc validator.TranslationFunc
	}{
		{
			tag:         "name",
			translation: "{0}为2-30个字符",
			override:    false,
		},
		{
			tag:         "desc",
			translation: "{0}为2-100个字符",
			override:    false,
		},
		{
			tag:         "text",
			translation: "{0}最多1000个字符",
			override:    false,
		},
		{
			tag:         "remark",
			translation: "{0}最多500个字符",
			override:    false,
		},
		{
			tag:         "role",
			translation: "{0}为2-30个字符",
			override:    false,
		},
		{
			tag:         "tag",
			translation: "{0}不符合规则",
			override:    false,
		},
		{
			tag:         "usergroup",
			translation: "{0}不符合规则",
			override:    false,
		},
		{
			tag:         "username",
			translation: "{0}不符合规则",
			override:    false,
		},
		{
			tag:         "account",
			translation: "{0}不符合规则",
			override:    false,
		},
		{
			tag:         "nickname",
			translation: "{0}不符合规则",
			override:    false,
		},
		{
			tag:         "password",
			translation: "{0}不符合规则",
			override:    false,
		},
	}
	for _, t := range translations {
		err := v.RegisterTranslation(t.tag, zhTranslator, registrationFunc(t.tag, t.translation, t.override), translateFunc)
		if err != nil {
			return
		}
	}
}

// registrationFunc 用于确认添加对应的 tag, translation, override 关系
func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {

	return func(ut ut.Translator) (err error) {

		if err = ut.Add(tag, translation, override); err != nil {
			return
		}

		return

	}

}
// translateFunc 运行的翻译函数
func translateFunc(ut ut.Translator, fe validator.FieldError) string {
    // T translate
	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		return fe.(error).Error()   // 返回validator.FieldError 错误
	}

	return t
}

// translateValidateError 错误返回
func translateValidateError(err error) string {
	if err == nil {
		return ""
	}

	errors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}
	msg := make([]string, 0)
	for _, e := range errors {
		msg = append(msg, e.Translate(zhTranslator))
	}
	return strings.Join(msg, "; ")
}
```
## 错误捕获使用

需要使用自己自定义`errors`, 使用 `import ""github.com/pkg/errors"`

```go

// 验证捕获错误
func ErrResp(c *gin.Context, err error) {
    switch t := errors.Cause(err).(type) {
    case validator.ValidationErrors:
        c.JSON(http.StatusOK, enum.InvalidArgumentWithDetail(translateValidateError(val)))
    }
}

// ErrorCode
type ErrorCode struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

InvalidArgument  = ErrorCode{Code: 3, Msg: "参数错误"}

// InvalidArgumentWithDetail 携带具体错误信息的参数错误码
func InvalidArgumentWithDetail(msg string) ErrorCode {
	return ErrorCode{
		Code: InvalidArgument.Code,
		Msg:  msg,
	}
}
```
