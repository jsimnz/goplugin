package main

import (
	"fmt"

	"github.com/jsimnz/goplugin"
	"github.com/jsimnz/goplugin/example/application/plugin_api"
)

var (
	cfg = goplugin.PluginManagerConfig{
		// Interface: (*ApplicationPlugin)(nil),

		Dir: "../plugins",
	}
	pluginMgr *goplugin.PluginManager
)

func main() {
	pluginMgr, err := goplugin.NewPluginManager(cfg)
	if err != nil {
		panic(err)
	}

	pluginMgr.RegisterInterface(goplugin.PluginInterfaceConfig{
		Identifier: plugin_api.APPLICATION_PLUGIN,
		Interface:  (*plugin_api.ApplicationPlugin)(nil),
		Factory:    plugin_api.ApplicationPluginFactory{},
	})

	pluginMgr.Init()

	a := 4
	b := 2
	for _, p := range pluginMgr.Plugins(plugin_api.APPLICATION_PLUGIN) {
		plugin := (p).(plugin_api.ApplicationPlugin)
		fmt.Println("Running plugin:", plugin.Name())
		fmt.Println("Plugin return:", plugin.Compute(a, b))
	}
}
