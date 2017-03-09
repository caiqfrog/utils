package configutil

import (
	"errors"
	"flag"
	"reflect"
	"strings"

	"github.com/robfig/config"
)

const (
	ConfigTag     = `cfg`
	OptionInline  = `inline`  // 需要到下一级去设置
	OptionNotNull = `notnull` // 必须设置值
)

var ErrNotPtr = errors.New(`Not ptr`)
var ErrNotStruct = errors.New(`Not struct`)
var ErrUnsupportType = errors.New(`Unsupport type!`)
var ErrNull = errors.New(`Null`)
var ErrNilInput = errors.New(`input is nil`)

var _Section = ""

func init() {
	Section := flag.String("section", config.DEFAULT_SECTION, "select section")
	_Section = *Section
	flag.Parse()
}

// 配置项上下文
type ConfigFieldContext struct {
	value reflect.Value
	typ   reflect.Type

	section      string
	key          string
	option       string
	defaultValue string
}

/**
 * 加载配置
 * 读取tag中的cfg段数据
 * 格式为 ${key},${option},${defaultValue}/${key},${defaultValue}
 * key的格式为 ${section}|${key} 如果section省略，则取_Section的值
 */
func Load(filepath string, ptr interface{}) error {
	if nil == ptr {
		return ErrNilInput
	}

	val, err := validate(ptr)
	if nil != err {
		return err
	}

	fields, err := getFieldContext(val)
	if nil != err {
		return err
	}

	cfg, err := config.ReadDefault(filepath)
	if nil != err {
		return err
	}
	for _, val := range fields {
		if err := parse(cfg, val); nil != err {
			return err
		}
	}

	return nil
}

/**
 *
 */
func validate(data interface{}) (reflect.Value, error) {
	var result reflect.Value
	switch val := data.(type) {
	case reflect.Value:
		result = val
	case *reflect.Value:
		result = *val
	default:
		result = reflect.ValueOf(data)
		// 入参非指针
		if reflect.Ptr != result.Kind() {
			return result, ErrNotPtr
		}
	}

	for reflect.Ptr == result.Kind() {
		result = result.Elem()
	}
	// 入参非结构体
	if reflect.Struct != result.Kind() {
		return result, ErrNotStruct
	}

	return result, nil
}

/**
 *
 */
func getFieldContext(value reflect.Value) ([]*ConfigFieldContext, error) {

	n := value.NumField()
	typ := value.Type()
	result := make([]*ConfigFieldContext, 0, n)
	//
	for i := 0; i < n; i++ {
		vField := value.Field(i)
		if vField.CanSet() {
			// 解析tag数据
			tField := typ.Field(i)
			tag := tField.Tag.Get(ConfigTag)
			if "" == tag {
				continue
			}

			tags := strings.SplitN(tag, ",", 3)
			var section, key, option, defaultValue string
			switch len(tags) {
			case 1:
				key = tags[0]
			case 2:
				key, defaultValue = tags[0], tags[1]
			default:
				key, option, defaultValue = tags[0], tags[1], tags[2]
			}
			// 如果为inline类型
			option = strings.TrimSpace(option)
			if OptionInline == option {
				for reflect.Ptr == vField.Kind() {
					if 0 == vField.Pointer() {
						vPtr := reflect.New(vField.Type().Elem())
						vField.Set(vPtr)
						vField = vField.Elem()
					}
				}
				f, err := validate(vField)
				if nil != err {
					return nil, err
				}
				v, err := getFieldContext(f)
				if nil != err {
					return nil, err
				}
				result = append(result, v...)
				continue
			}
			key = strings.TrimSpace(key)
			// 未设置key，跳过
			if "" == key || "|" == key {
				continue
			}
			// key格式
			keys := strings.SplitN(key, "|", 2)
			switch len(keys) {
			case 1:
				section = config.DEFAULT_SECTION
			case 2:
				if "" == keys[0] {
					section = _Section
				}
				section = keys[0]
				key = keys[1]
			}

			//for reflect.Ptr == vField.Kind() {
			//	vField = vField.Elem()
			//}
			result = append(result, &ConfigFieldContext{
				value:        vField,
				typ:          vField.Type(),
				section:      section,
				key:          key,
				option:       option,
				defaultValue: defaultValue,
			})
		}
	}

	return result, nil
}
