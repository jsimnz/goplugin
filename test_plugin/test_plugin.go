package main

import (
	"github.com/jsimnz/goplugin"
)

import "C"

const example_plugin = goplugin.Type(iota)

var (
	p = TestPlugin{}
)

type TestPlugin struct{}

func (tp TestPlugin) Test() {}

func (tp TestPlugin) Type() goplugin.Type {
	return example_plugin
}

func init() {
	goplugin.RegisterPlugin(p)
}

func main() {}
