package echoutil

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/labstack/echo"
)

// 记录已注册路径，用于判重
var _RegisteredMap = map[string]bool{}

/**
 * 注册Handler
 * @param prefix 前缀
 */
func RegisterHandler(e *echo.Echo, prefix string, controller interface{}) error {
	//
	ctx := []*ControllerContext{}
	//
	var rValue reflect.Value
	if v, is := controller.(reflect.Value); is {
		rValue = v
	} else {
		rValue = reflect.ValueOf(controller)
	}
	//
	ctx = append(ctx, getControllerContext(prefix, rValue)...)
	// 注册
	for _, val := range ctx {
		// 判断是否重复注册
		path := val.Method + ` ` + val.Path
		if _, has := _RegisteredMap[path]; has {
			return errors.New(`Duplicate Handler: ` + path)
		}
		// 注册handler
		e.Match([]string{val.Method}, val.Path, val.Func)
		// 修正echo结尾为/时无法匹配bug
		e.Match([]string{val.Method}, val.Path+`/`, val.Func)
		//
		_RegisteredMap[path] = true
		// 打印
		fmt.Println("register router: ", path)
	}
	//
	return nil
}

//
type _ControllerFunc func() (string, string, func(echo.Context) error)

// 注册上下文
type ControllerContext struct {
	Method string           // 方法，GET/POST等
	Path   string           // restfull路径
	Func   echo.HandlerFunc // Handler方法
}

// 用于规范路径中的斜杠
var _PathReplacer = regexp.MustCompile(`[/\\]+`)

// 获取注册上下文
func getControllerContext(prefix string, value reflect.Value) []*ControllerContext {
	//
	result := []*ControllerContext{}
	//
	mLen := value.NumMethod()
	for i := 0; i < mLen; i++ {
		vMethod := value.Method(i)
		iMethod := vMethod.Interface()
		// 判断是否符合_ControllerFunc格式
		var Method, Path string
		var Func echo.HandlerFunc
		switch fn := iMethod.(type) {
		case func() (string, string, func(context echo.Context) error):
			var _Func func(context echo.Context) error
			Method, Path, _Func = fn()
			Func = echo.HandlerFunc(_Func)
		case func() (string, string, echo.HandlerFunc):
			Method, Path, Func = fn()
		default:
			continue
		}
		// 方法修正为大写
		Method = strings.ToUpper(Method)
		// 格式化路径
		Path = `/` + prefix + `/` + Path
		Path = _PathReplacer.ReplaceAllString(Path, `/`)
		// 去掉结尾的/
		pLen := len(Path)
		if pLen > 1 && Path[pLen-1:] == `/` {
			Path = string(Path[:pLen-1])
		}
		//
		result = append(result, &ControllerContext{
			Method: Method,
			Path:   Path,
			Func:   Func,
		})
	}
	//
	return result
}
