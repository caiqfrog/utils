package echoutil

import (
	"fmt"
	"testing"

	_ "net/http/pprof"

	"reflect"

	"github.com/labstack/echo"
)

type TestController struct{}

func (self TestController) HandleTest1() (method, path string, handle echo.HandlerFunc) {
	method = "POST"
	path = "/\\/:abcd"
	handle = func(c echo.Context) error {
		return nil
	}
	return
}

func (self *TestController) HandleTest2() (method, path string, handle func(context echo.Context) error) {
	method = "POST"
	path = "/\\/:whell"
	handle = func(c echo.Context) error {
		return nil
	}
	return
}

func (self TestController) Empty() string {
	return "hello world"
}

func TestRegister(t *testing.T) {
	e := echo.New()

	RegisterHandler(e, "/hhh\\/", new(TestController))
	if err := RegisterHandler(e, "hhh", new(TestController)); nil != err {
		fmt.Println(err)
	}

	fmt.Println(e.Routes())
}

func BenchmarkRegister(b *testing.B) {
	b.ReportAllocs()

	val := reflect.ValueOf(new(TestController))

	for i := 0; i < b.N; i++ {
		getControllerContext(`hhh`, val)
	}
}
