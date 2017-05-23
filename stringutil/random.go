package stringutil

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"time"

	MathRand "math/rand"
)

// 随机字符集类型
type RandomType uint

const (
	RandomDigit   RandomType = 1 << iota // 0-9
	RandomLower                          // a-z
	RandomUpper                          // A-Z
	RandomExtend                         // !@#$%^&*_+-=()[]{}?|
	CharsetDigit  = `0123456789`
	CharsetLower  = `abcdefghijklmnopqrstuvwxyz`
	CharsetUpper  = `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
	CharsetExtend = `!@#$%^&*_+-=()[]{}?|`
)

var _RandomArray = []RandomType{RandomDigit, RandomLower, RandomUpper, RandomExtend}
var _CharsetArray = []string{CharsetDigit, CharsetLower, CharsetUpper, CharsetExtend}

/**
 * 生成随机字符串
 */
func Random(count int, types ...RandomType) string {
	//
	typ := uint(0)
	for _, val := range types {
		typ |= uint(val)
	}
	if 0 == typ {
		// 未设置类型，默认为0-9a-zA-Z
		typ = uint(RandomDigit) | uint(RandomLower) | uint(RandomUpper)
	}
	//
	set := ""
	for idx, val := range _RandomArray {
		if 0 != (uint(val) & typ) {
			set += _CharsetArray[idx]
		}
	}
	//
	buf := make([]byte, count*2)
	// 真随机
	offset, _ := rand.Read(buf)

	bbuf := bytes.NewBuffer(buf)
	// 根据时间种子获得的伪随机
	if offset != count {
		seeder := MathRand.New(MathRand.NewSource(time.Now().UnixNano()))
		for i := offset; i < count; i += 8 {
			i64 := seeder.Int63()
			// FIXME 错误处理
			binary.Write(bbuf, binary.BigEndian, i64)
		}
	}
	//
	lset := len(set)
	result := make([]byte, count)

	for i := 0; i < count; i++ {
		x := uint16(0)
		binary.Read(bbuf, binary.BigEndian, &x)
		result[i] = set[int(x)%lset]
	}
	return string(result)
}

/**
 * 生成所有支持类型的随机字符串
 */
func RandomFull(count int) string {
	return Random(count, _RandomArray...)
}
