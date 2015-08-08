package main

import (
	"github.com/jsimnz/example/plugin"
	"github.com/jsimnz/goplugin"
)

var (
	cfg = goplugin.PluginManagerConfig{
		Interface: (*ApplicationPlugin)(nil),

		Dir: "some/dir",
	}
	pluginMgr = goplugin.NewPluginManager(cfg)
)

func init() {
	pluginMgr.RegisterInterface(goplugin.InterfaceConfig{
		Identifier: APPLICATION_PLUGIN,
		Interface:  (*ApplicationPlugin)(nil),
		Factory:    ApplicationPluginFactory{},
	})

	pluginMgr.Init()
}

func main() {
	a := 4
	b := 2
	for _, p := range pluginMgr.Plugins() {
		plugin := (ApplicationPlugin)(p)
		fmt.Println("Running plugin:", plugin.Name())
		fmt.Println("Plugin return:", plugin.Compute(a, b))
	}
}
