package httputil

import "io"
import "encoding/json"

/**
 * 读取reader，并转换为json格式
 */
func JSONReader(reader io.Reader, ptr interface{}) error {
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(ptr); nil != err && io.EOF != err {
		return err
	}
	return nil
}
