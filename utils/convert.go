package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"golang.org/x/crypto/ripemd160"
)

/**
 *将int类型的数据转换为[]byte类型
 */
func Int2Byte(num int64) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes(), err
}

/**
 *gob编码序列化
 */
func Encode(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(v)
	return buf.Bytes(), err
}

/**
 *gob编码反序列化
 */
func Decode(data []byte, v interface{}) (interface{}, error) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(v)
	return v, err
}

/**
 *将json格式的数组转换为对应得字符串类型得切片
 */
func JSONArray2String(array string) ([]string, error) {
	var stringSlice []string
	err := json.Unmarshal([]byte(array), &stringSlice)
	return stringSlice, err

}

/**
 *将json格式的数组转换float64数组
 */
func JSONArray2Float(array string) ([]float64, error) {
	var floatSlice []float64
	err := json.Unmarshal([]byte(array), &floatSlice)
	return floatSlice, err
}

/**
 *sha256计算
 */
func Hash256(data []byte) []byte {
	sha := sha256.New()
	sha.Write(data)
	return sha.Sum(nil)
}

/**
 *ripemd160计算
 */
func HashRipemd160(data []byte) []byte {
	hashripe := ripemd160.New()
	hashripe.Write(data)
	return hashripe.Sum(nil)
}