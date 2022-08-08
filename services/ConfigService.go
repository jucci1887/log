/*
 Author: Kernel.Huang
 Mail: kernelman79@gmail.com
 File: main
 Date: 3/18/21 1:01 PM
*/

package services

import (
	"github.com/bitly/go-simplejson"
)

type ConfigServices struct {
	from    []byte
	to      string
	types   string
	keyName string
	json    *simplejson.Json
}

var ConfigService = &ConfigServices{}

func (conf *ConfigServices) Get(content []byte) *ConfigServices {
	conf.from = content
	return conf
}

func (conf *ConfigServices) Dir() string {
	conf.keyName = "dir"
	conf.types = "string"
	return conf.result()
}

func (conf *ConfigServices) Name() string {
	conf.keyName = "name"
	conf.types = "string"
	return conf.result()
}

func (conf *ConfigServices) Prefix() string {
	conf.keyName = "prefix"
	conf.types = "string"
	return conf.result()
}

func (conf *ConfigServices) Level() string {
	conf.keyName = "level"
	conf.types = "string"
	return conf.result()
}

func (conf *ConfigServices) Relative() bool {
	conf.keyName = "relative"
	conf.types = "bool"
	source := conf.from

	return JsonService.New(source).Get(conf.keyName, conf.types).End().ToBool()
}

func (conf *ConfigServices) result() string {
	source := conf.from
	conf.to = JsonService.New(source).Get(conf.keyName, conf.types).End().ToString()
	return conf.to
}
