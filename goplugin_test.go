package goplugin

import (
	"os"
	"testing"
)

const (
	example_plugin = Type(iota)
	random_plugin
)

// Example Plugin Interface
type examplePlugin interface {
	Plugin
	Test()
}

type examplePluginFactory struct {
	PluginFactory
	TestFn func()
}

func (plugin examplePluginFactory) Test() {}

func (plugin examplePluginFactory) Type() Type {
	return example_plugin
}

func NewTestPluginManager() *PluginManager {
	pm, _ := NewPluginManager(PluginManagerConfig{
		Dir: "./test_plugin/",
	})

	return pm
}

func NewRegisteredInterface() *PluginManager {
	pm := NewTestPluginManager()
	pm.RegisterInterface(PluginInterfaceConfig{
		Identifier: example_plugin,
		Interface:  (*examplePlugin)(nil),
		Factory:    examplePluginFactory{},
	})

	return pm
}

func TestNewPluginManagerEmptyConfigDir(t *testing.T) {
	p, err := NewPluginManager(PluginManagerConfig{})
	if err == nil {
		t.Errorf("Expected an error when creating plugin manager with empty config dir")
	}
	if p != nil {
		t.Error("Expected returned pluginManager instance to be nil")
	}
}

func TestNewPluginManagerConfigDirEnvVar(t *testing.T) {
	os.Setenv("PLUGIN_DIR", "my/plugins/dir")
	p, err := NewPluginManager(PluginManagerConfig{})
	if err != nil {
		t.Errorf("Expected new plugin manager, got: %v", err)
	}
	if p.cfg.Dir != "my/plugins/dir" {
		t.Errorf("Expected config dir: 'my/plugins/dir', got: %v", p.cfg.Dir)
	}
}

func TestNewPluginManagerWithConfigDir(t *testing.T) {
	p, err := NewPluginManager(PluginManagerConfig{
		Dir: "my/plugins",
	})
	if err != nil {
		t.Error("Expected a newly created pluginManager without errors")
	}
	if p == nil {
		t.Error("Expected a newly created pluginManager, not nil")
	}
}

func TestRegisterInterface(t *testing.T) {
	pm := NewTestPluginManager()

	err := pm.RegisterInterface(PluginInterfaceConfig{
		Identifier: example_plugin,
		Interface:  (*examplePlugin)(nil),
		Factory:    examplePluginFactory{},
	})

	if err != nil {
		t.Errorf("Expected to succesfully register examplePlugin interface, got: %v", err)
	}
}

func TestSavedPluginInterfacesCount(t *testing.T) {
	pm := NewRegisteredInterface()
	numInterfaces := len(pm.interfaces)
	if numInterfaces != 1 {
		t.Errorf("Expected numInterfaces to be 1, got: %v", numInterfaces)
	}
}

func TestSavedPluginInterfaceFactoryStruct(t *testing.T) {
	pm := NewTestPluginManager()

	err := pm.RegisterInterface(PluginInterfaceConfig{
		Identifier: example_plugin,
		Interface:  (*examplePlugin)(nil),
		Factory:    examplePluginFactory{},
	})

	if err != nil {
		t.Errorf("Expected Factory to fail struct type test, got: %v", err)
	}
}

func TestSavedPluginInterfaceByType(t *testing.T) {
	pm := NewRegisteredInterface()
	if _, ok := pm.interfaces[example_plugin]; !ok {
		t.Error("Expected to find saved interfave via Type map lookup")
	}
}

func TestSavedPluginInterfaceNumMethods(t *testing.T) {
	pm := NewRegisteredInterface()
	iface := pm.interfaces[example_plugin]
	if len(iface.methods) != 2 {
		t.Errorf("Expected num methods as 2, got: %v", len(iface.methods))
	}
}

func TestLoadPlugins(t *testing.T) {
	pm := NewRegisteredInterface()
	err := pm.LoadPlugins()
	if err != nil {
		t.Errorf("Expecting no errors, got: %v", err)
	}
}

func TestLoadPluginsInvalidDir(t *testing.T) {
	pm, err := NewPluginManager(PluginManagerConfig{
		Dir: "random/dir",
	})
	pm.RegisterInterface(PluginInterfaceConfig{
		Identifier: example_plugin,
		Interface:  (*examplePlugin)(nil),
		Factory:    examplePluginFactory{},
	})

	if err != nil {
		t.Errorf("Expected creation on plugin manager, got: %v", err)
	}
	err = pm.LoadPlugins()
	if err == nil {
		t.Errorf("Expected error when loading plugins, got: %v", err)
	}
}

func TestLoadPluginsCount(t *testing.T) {
	pm := NewRegisteredInterface()
	pm.RegisterInterface(PluginInterfaceConfig{
		Identifier: example_plugin,
		Interface:  (*examplePlugin)(nil),
		Factory:    examplePluginFactory{},
	})

	err := pm.LoadPlugins()
	if err != nil {
		t.Errorf("Expected to load plugins, got: %v", err)
	}

	if num := len(pm.pluginsByName); num != 1 {
		t.Errorf("Expected 1, got: %v", num)
	}
}

func TestGetPluginsByType(t *testing.T) {
	pm := NewRegisteredInterface()
	pm.Init()
	plugins := pm.Plugins(example_plugin)
	l := len(plugins)
	if l == 0 {
		t.Errorf("Expected returned plugins to be 1, got: %v", l)
	}
}

func TestGetPluginsByTypeCast(t *testing.T) {
	pm := NewRegisteredInterface()
	pm.Init()
	p := pm.Plugins(example_plugin)[0]
	if _, ok := p.(examplePlugin); !ok {
		t.Errorf("Failed to get plugin interface from loaded plugin")
	}
}

func TestGetPluginsByTypeNil(t *testing.T) {
	pm := NewRegisteredInterface()
	pm.Init()
	plugins := pm.Plugins(random_plugin)
	if l := len(plugins); l != 0 {
		t.Errorf("Expected to get an empty plugin array, got: %v", l)
	}
}

func TestGetPluginByName(t *testing.T) {
	pm := NewRegisteredInterface()
	pm.Init()
	p := pm.Plugin("TestPlugin")
	if p == nil {
		t.Error("Expected to get TestPlugin, got nil instead")
	}
}

func TestGetPluginByNameCast(t *testing.T) {
	pm := NewRegisteredInterface()
	pm.Init()
	p := pm.Plugin("TestPlugin")
	if _, ok := p.(examplePlugin); !ok {
		t.Error("Expected TestPlugin to be casted to plugin interface")
	}
}

func TestGetPluginByNameNil(t *testing.T) {
	pm := NewRegisteredInterface()
	pm.Init()
	p := pm.Plugin("NonExistantPlugin")
	if p != nil {
		t.Errorf("Expected to get nil for non existent plugin, got: %v", p)
	}
}
