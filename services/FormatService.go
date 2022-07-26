/*
 Author: Kernel.Huang
 Mail: kernelman79@gmail.com
 Date: 3/18/21 1:01 PM
*/
package services

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

type FormatServices struct {
	from interface{}
	to   interface{}
}

var FormatService = new(FormatServices)

func (format *FormatServices) Base64ToString() string {
	return format.FromBase64(format.to).ToString()
}

func (format *FormatServices) Base64ToByte() []byte {
	return format.FromBase64(format.to).ToByte()
}

func (format *FormatServices) ToByte() []byte {
	switch value := format.to.(type) {

	case []byte:
		return value
	case string:
		return []byte(value)
	case []string:
		return []byte(strings.Join(value, ","))
	case int:
		return []byte(strconv.Itoa(value))
	case map[string]interface{}:
		toByte, err := json.Marshal(value)
		if err != nil {
			log.Println("Map interface to byte error: ", err)
		}
		return toByte
	case interface{}:
		toByte, err := json.Marshal(value)
		if err != nil {
			log.Println("Interface to byte error: ", err)
		}
		return toByte
	}

	return nil
}

func (format *FormatServices) ToRune() []rune {
	switch value := format.to.(type) {

	case string:
		return []rune(value)
	case []string:
		return []rune(strings.Join(value, ","))
	}

	return nil
}

func (format *FormatServices) ToBool() bool {
	switch value := format.to.(type) {

	case string:
		toBool, err := strconv.ParseBool(value)
		if err != nil {
			log.Println("String to bool error: ", err)
		}
		return toBool
	case int:
		toBool, err := strconv.ParseBool(strconv.Itoa(value))
		if err != nil {
			log.Println("Int to bool error: ", err)
		}
		return toBool
	}

	return false
}

func (format *FormatServices) ToString() string {
	switch value := format.to.(type) {

	case string:
		return value
	case []string:
		return strings.Join(value, ",")
	case []byte:
		return string(value)
	case int:
		return strconv.Itoa(value)
	case int64:
		return strconv.FormatInt(value, 10)
	case []rune:
		return string(value)
	case map[string]interface{}:
		toByte, err := json.Marshal(value)
		if err != nil {
			log.Println("Map interface to string error: ", err)
		}
		return string(toByte)
	}

	return ""
}

func (format *FormatServices) ToInt() int {
	switch value := format.to.(type) {

	case int:
		return value
	case string:
		str, err := strconv.Atoi(value)
		if err != nil {
			log.Println("String to int error: ", err)
		}
		return str
	case []byte:
		bytes, err := strconv.Atoi(string(value))
		if err != nil {
			log.Println("[]byte to int error: ", err)
		}
		return bytes
	}

	return 0
}

func (format *FormatServices) ToUint() uint {
	switch value := format.to.(type) {

	case uint:
		return value
	case []byte:
		bytes, err := strconv.Atoi(string(value))
		if err != nil {
			log.Println("[]byte to uint error: ", err)
		}
		return uint(bytes)
	case int:
		return uint(value)
	case string:
		str, err := strconv.Atoi(value)
		if err != nil {
			log.Println("String to uint error: ", err)
		}
		return uint(str)
	}

	return 0
}

func (format *FormatServices) ToInt64() int64 {
	switch value := format.to.(type) {

	case int64:
		return value
	case string:
		str, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Println("String to int64 error: ", err)
		}
		return str
	case []byte:
		bytes, err := strconv.ParseInt(string(value), 10, 64)
		if err != nil {
			log.Println("[]byte to int64 error: ", err)
		}
		return bytes
	}

	return 0
}

func (format *FormatServices) ToBase64() string {
	switch value := format.to.(type) {

	case string:
		return base64.StdEncoding.EncodeToString([]byte(value))
	case []byte:
		return base64.StdEncoding.EncodeToString(value)
	}

	return ""
}

func (format *FormatServices) ToArrayArrayInt() [][]int {
	switch final := format.to.(type) {

	case map[string]interface{}:
		var convert [][]int
		for _, v := range final {
			convert = append(convert, format.FromInterface(v).ToArrayInt())
		}
		return convert
	default:
		log.Println("%s%T", "Unknown type: ", final)
	}
	return nil
}

func (format *FormatServices) ToArrayInt() []int {
	switch final := format.to.(type) {

	case []string:
		convert := make([]int, len(final))
		for k, v := range final {
			convert[k] = format.FromInterface(v).ToInt()
		}
		return convert
	default:
		log.Println("%s%T", "Unknown type: ", final)
	}
	return nil
}

func (format *FormatServices) ToArrayString() []string {
	switch value := format.to.(type) {

	case []int:
		convert := make([]string, len(value))
		for _, v := range value {
			convert = append(convert, string(v))
		}
		return convert
	case string:
		return strings.Split(value, ",")
	case int:
		return strings.Split(strconv.Itoa(value), ",")
	case []string:
		return value
	case []byte:
		return strings.Split(string(value), ",")
	case []rune:
		return strings.Split(string(value), ",")
	}

	return nil
}

func (format *FormatServices) ToMapStringByte() map[string][]byte {
	switch value := format.to.(type) {

	case map[string][]byte:
		return value
	case map[string]string:
		convert := make(map[string][]byte)

		for k, v := range value {
			convert[k] = []byte(v)
		}
		return convert
	case map[string]int:
		convert := make(map[string][]byte)

		for k, v := range value {
			convert[k] = []byte(strconv.Itoa(v))
		}
		return convert
	}

	return nil
}

func (format *FormatServices) ToMapIntInterface() map[int]interface{} {
	switch value := format.to.(type) {

	case []interface{}:
		convert := make(map[int]interface{})
		for k, v := range value {
			convert[k] = v
		}
		return convert
	}

	return nil
}

func (format *FormatServices) ToMapStringInterface() map[string]interface{} {
	switch value := format.to.(type) {

	case []interface{}:
		convert := make(map[string]interface{})

		for k, v := range value {
			convert[string(k)] = v
		}
		return convert

	case []map[string]interface{}:
		newMap := make(map[string]interface{})

		for k, v := range value {
			newMap[string(k)] = v
		}
		return newMap

	case []map[int]interface{}:
		newMap := make(map[string]interface{})

		for k, v := range value {
			newMap[string(k)] = v
		}
		return newMap

	case map[string]interface{}:
		return value
	}

	return nil
}

func (format *FormatServices) ToArrayInterface() []interface{} {
	switch value := format.to.(type) {

	case []interface{}:
		return value
	}

	return nil
}

func (format *FormatServices) InterfaceToArray() []string {
	switch value := format.to.(type) {

	case []interface{}:
		var convert []string

		for _, param := range value {
			convert = append(convert, param.(string))
		}

		return convert
	}

	return nil
}

/**
 * 待测试
 */
func (format *FormatServices) ToArrayMapStringInterface() []map[string]interface{} {
	switch value := format.to.(type) {

	case []interface{}:
		newMap := make([]map[string]interface{}, len(value))

		for k, v := range value {
			newMap[k] = format.FromInterface(v).ToMapStringInterface()
		}
		return newMap

	case []map[string]interface{}:
		return value
	}

	return nil
}

func (format *FormatServices) ToMapStringString() map[string]string {
	switch value := format.to.(type) {

	case map[string]interface{}:
		convert := make(map[string]string)

		for k, v := range value {
			switch v := v.(type) {

			case int:
				convert[k] = string(v)
			case string:
				convert[k] = v
			}
		}
		return convert

	case map[string]string:
		return value

	case map[string][]byte:
		convert := make(map[string]string)

		for k, v := range value {
			convert[k] = string(v)
		}
		return convert

	case map[string]int:
		convert := make(map[string]string)

		for k, v := range value {
			convert[k] = strconv.Itoa(v)
		}
		return convert

	default:
		log.Println("%s%T", "Unknown type: ", value)
	}

	return nil
}

func (format *FormatServices) FromString(from string) *FormatServices {
	format.to = from
	return format
}

func (format *FormatServices) FromArrayString(from []string) *FormatServices {
	format.to = from
	return format
}

func (format *FormatServices) FromByte(from []byte) *FormatServices {
	format.to = from
	return format
}

func (format *FormatServices) FromInt(from int) *FormatServices {
	format.to = from
	return format
}

func (format *FormatServices) FromInt64(from int64) *FormatServices {
	format.to = from
	return format
}

func (format *FormatServices) FromRune(from []rune) *FormatServices {
	format.to = from
	return format
}

func (format *FormatServices) FromBase64(from interface{}) *FormatServices {
	switch value := from.(type) {

	case string:
		decodeBase64, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			log.Println("Decode the base64 string failed: ", err)
		}
		format.to = decodeBase64
	case []byte:
		decodeBase64, err := base64.StdEncoding.DecodeString(string(value))
		if err != nil {
			log.Println("Decode the base64 string failed: ", err)
		}
		format.to = decodeBase64
	}

	return format
}

func (format *FormatServices) FromMapStringByte(from map[string][]byte) *FormatServices {
	convert := make(map[string]string)

	for k, v := range from {
		convert[k] = string(v)
	}

	format.to = convert
	return format
}

func (format *FormatServices) FromMapStringInterface(from map[string]interface{}) *FormatServices {
	format.to = from
	return format
}

func (format *FormatServices) FromInterface(from interface{}) *FormatServices {
	format.to = from
	return format
}

func (format *FormatServices) FromBsonM(from interface{}) *FormatServices {
	format.to = from
	return format
}
