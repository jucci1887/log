/*
 Author: Kernel.Huang
 Mail: kernelman79@gmail.com
 File: main
 Date: 3/18/21 1:01 PM
*/

package services

import (
	"github.com/bitly/go-simplejson"
	"log"
)

type JsonServices struct {
	json  *simplejson.Json
	types map[string]func()
}

var JsonService = &JsonServices{}

var JsonTypes = map[string]func(){
	"int":       JsonService.Int,
	"bool":      JsonService.Bool,
	"map":       JsonService.Map,
	"array":     JsonService.Array,
	"string":    JsonService.String,
	"interface": JsonService.Interface,
}

func (js *JsonServices) New(data []byte) *JsonServices {
	newJson, err := simplejson.NewJson(data)
	if err != nil {
		log.Println("Gets the JSON format failed: ", err)
	}
	js.json = newJson

	return js
}

func (js *JsonServices) Get(key string, types string) *JsonServices {
	js.json = js.json.Get(key)
	JsonTypes[types]()

	return js
}

func (js *JsonServices) End() *FormatServices {
	return FormatService
}

func (js *JsonServices) Int() {
	to, err := js.json.Int()
	js.HandleTypesError(err)
	FormatService.to = to
}

func (js *JsonServices) Bool() {
	to, err := js.json.Bool()
	js.HandleTypesError(err)
	FormatService.to = to
}

func (js *JsonServices) String() {
	to, err := js.json.String()
	js.HandleTypesError(err)
	FormatService.to = to
}

func (js *JsonServices) Map() {
	to, err := js.json.Map()
	js.HandleTypesError(err)
	FormatService.to = to
}

func (js *JsonServices) Array() {
	to, err := js.json.Array()
	js.HandleTypesError(err)
	FormatService.to = to
}

func (js *JsonServices) Interface() {
	FormatService.to = js.json.Interface()
}

func (js *JsonServices) HandleTypesError(err error) {
	if err != nil {
		log.Println("Get the Json value type failed: ", err)
	}
}
