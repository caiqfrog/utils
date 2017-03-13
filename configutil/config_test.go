package configutil

import (
	"fmt"
	"testing"
)

type Test struct {
	B         bool       `cfg:"b,true"`
	Int16     int16      `cfg:"i16,notnull,16"`
	Uint32    uint32     `cfg:"u32,32"`
	Float64   []float64  `cfg:"f64s,,0.1,0.5"`
	String    string     `cfg:"hello|world,cbd"`
	Inline    TestInline `cfg:",inline,"`
	StringPtr *string    `cfg:"world"`
}

type TestInline struct {
	V1 string `cfg:"inline1,notnull,"`
	V2 int64  `cfg:"inline2, notnull,"`
}

type TestNull struct {
	Uint64 uint64 `cfg:"u64,notnull,"`
}

type TestInlineWrapper struct {
	InlinePtr *TestInline `cfg:",inline,"`
}

func Test_Load(t *testing.T) {
	val := Test{}
	if err := Load(`./test.cfg`, &val); nil != err {
		panic(err)
	}
	fmt.Println(val)
}

func Test_LoadNull(t *testing.T) {
	val := TestNull{}
	if err := Load(`./test.cfg`, &val); nil == err {
		panic(`notnull did not work`)
	}
}

func Test_Inline(t *testing.T) {
	val := TestInlineWrapper{}
	if err := Load(`./test.cfg`, &val); nil != err {
		panic(err)
	}
	fmt.Println(val.InlinePtr)
}
