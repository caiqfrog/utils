package configutil

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/robfig/config"
)

//
type ParseType int

const (
	Bool ParseType = iota
	Int
	Uint
	Float
	String
	Slice
	Unsupport
)

/**
 *
 */
func parse(cfg *config.Config, ctx *ConfigFieldContext) error {
	str, err := cfg.String(ctx.section, ctx.key)
	if nil != err {
		// 不能不设置
		if OptionNotNull == ctx.option {
			return ErrNull
		}
		str = ctx.defaultValue
	}

	switch parseType(ctx.value.Kind()) {
	case String:
		ctx.value.SetString(str)
	case Bool:
		ctx.value.SetBool(parseBool(str, ctx.defaultValue))
	case Int:
		ctx.value.SetInt(parseInt(str, ctx.defaultValue))
	case Uint:
		ctx.value.SetUint(parseUint(str, ctx.defaultValue))
	case Float:
		ctx.value.SetFloat(parseFloat(str, ctx.defaultValue))
	case Slice:
		str, err := cfg.String(ctx.section, ctx.key)
		if nil != err {
			str = ctx.defaultValue
		}
		values := strings.Split(str, ",")
		eType := ctx.typ.Elem()
		eLen := len(values)
		eValue := reflect.MakeSlice(ctx.typ, eLen, eLen)
		eTyp := parseType(eType.Kind())
		// 数组
		for i := 0; i < eLen; i++ {
			elem := eValue.Index(i)
			switch eTyp {
			case String:
				elem.SetString(values[i])
			case Bool:
				elem.SetBool(parseBool(values[i], ctx.defaultValue))
			case Int:
				elem.SetInt(parseInt(values[i], ctx.defaultValue))
			case Uint:
				elem.SetUint(parseUint(values[i], ctx.defaultValue))
			case Float:
				elem.SetFloat(parseFloat(values[i], ctx.defaultValue))
			default:
				return ErrUnsupportType
			}
		}

		ctx.value.Set(eValue)
	default:
		return ErrUnsupportType
	}

	return nil
}

func parseType(kind reflect.Kind) ParseType {
	if reflect.String == kind {
		return String
	} else if reflect.Bool == kind {
		return Bool
	} else if reflect.Int == kind || reflect.Int8 == kind || reflect.Int16 == kind || reflect.Int32 == kind || reflect.Int64 == kind {
		return Int
	} else if reflect.Uint == kind || reflect.Uint8 == kind || reflect.Uint16 == kind || reflect.Uint32 == kind || reflect.Uint64 == kind {
		return Uint
	} else if reflect.Float32 == kind || reflect.Float64 == kind {
		return Float
	} else if reflect.Slice == kind {
		return Slice
	}
	return Unsupport
}

func parseBool(value, defaultValue string) bool {
	boolean, err := strconv.ParseBool(value)
	if nil != err {
		boolean, _ = strconv.ParseBool(defaultValue)
	}
	return boolean
}

func parseInt(value, defaultValue string) int64 {
	i64, err := strconv.ParseInt(value, 10, 0)
	if nil != err {
		i64, _ = strconv.ParseInt(defaultValue, 10, 0)
	}
	return i64
}

func parseUint(value, defaultValue string) uint64 {
	u64, err := strconv.ParseUint(value, 10, 0)
	if nil != err {
		u64, _ = strconv.ParseUint(defaultValue, 10, 0)
	}
	return u64
}

func parseFloat(value, defaultValue string) float64 {
	f64, err := strconv.ParseFloat(value, 0)
	if nil != err {
		f64, _ = strconv.ParseFloat(defaultValue, 0)
	}
	return f64
}

func parseValue(typ ParseType, value, defaultValue string) interface{} {
	switch typ {
	case String:
		return value
	case Bool:
		return parseBool(value, defaultValue)
	case Int:
		return parseInt(value, defaultValue)
	case Uint:
		return parseUint(value, defaultValue)
	case Float:
		return parseFloat(value, defaultValue)
	default:
		return Unsupport
	}
}
