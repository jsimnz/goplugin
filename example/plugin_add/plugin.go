package main

import (
	"github.com/jsimnz/goplugin"
	"github.com/jsimnz/goplugin/example/application/plugin_api"
)

var (
	p = PluginAdd{"add"}
)

//+goplugin
type PluginAdd struct {
	name string
}

func (p PluginAdd) Type() goplugin.Type {
	return plugin_api.APPLICATION_PLUGIN
}

func (p PluginAdd) Name() string {
	return p.name
}

func (p PluginAdd) Compute(a, b int) int {
	return a + b
}

func init() {
	goplugin.RegisterPlugin((plugin_api.ApplicationPlugin)(p))
}

func main() {}
