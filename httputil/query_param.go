package httputil

import "strconv"

// 获取url中query参数方法
type QueryParamHolder interface {
	QueryParam(string) string
}

/**
 * 从query中获取int64
 */
func Int64QueryParam(holder QueryParamHolder, key string, _default int64) int64 {
	val := holder.QueryParam(key)
	if "" == val {
		return _default
	}
	if i64, err := strconv.ParseInt(val, 10, 0); nil == err {
		return i64
	}
	return _default
}

/**
 * 从query中获取float64
 */
func Float64QueryParam(holder QueryParamHolder, key string, _default float64) float64 {
	val := holder.QueryParam(key)
	if "" == val {
		return _default
	}
	if f64, err := strconv.ParseFloat(val, 10); nil == err {
		return f64
	}
	return _default
}
