package goplugin

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/jsimnz/dl"
)

/*
 int MY_TEST_VAR;
*/
import "C"

func setTestVar(x int) {
	C.MY_TEST_VAR = C.int(x)
}

type Type uint16

// Interface required to be considered a plugin
// The only methods defined are those that make
// it possible for GoPlugin to register and understand
// how to treat plugins.
//
// So far this just includes the Type() method which
// informs a PluginManager what kind of plugin has been
// loaded
type Plugin interface {
	// Type() Type
}

// A helper struct to be embedded to more easily
// create plugin Factories
type PluginFactory struct {
	TypeFn func() Type
}

type plugin struct {
	lib      *dl.DL                 // reference to the loaded shared lib
	name     string                 // name of the plugin
	iFace    *PluginInterfaceConfig // reference to the Plugin interface definition
	instance reflect.Value          // loaded instance of the plugin
}

// Manages the definitions, loadings, initializations,
// and access of plugins
type PluginManager struct {
	// Defined plugin interfaces
	interfaces map[Type]PluginInterfaceConfig
	// Loaded plugins that have been mapped to types
	pluginsByType map[Type][]*plugin
	// Loaded plugins that have been mapped by their name
	pluginsByName map[string]*plugin
	// Configuration settings
	cfg PluginManagerConfig
}

// Creata a new instance of the plugin manager with a given
// PluginMangerConfig
func NewPluginManager(cfg PluginManagerConfig) (*PluginManager, error) {
	if cfg.Dir == "" {
		if dir := os.Getenv("PLUGIN_DIR"); dir != "" {
			cfg.Dir = dir
		} else {
			return nil, errors.New("Config directory cannot be empty")
		}
	}
	pm := new(PluginManager)
	pm.interfaces = make(map[Type]PluginInterfaceConfig)
	pm.pluginsByName = make(map[string]*plugin)
	pm.pluginsByType = make(map[Type][]*plugin)
	pm.cfg = cfg

	return pm, nil
}

// Register a plugin to be used by your application.
// Requires a type to identify plugins by, and interface
// and factory to generate plugins of that interface with
func (pm *PluginManager) RegisterInterface(cfg PluginInterfaceConfig) error {
	// ensure the identifier is unique and hasn't been registered
	// before
	if _, exists := pm.interfaces[cfg.Identifier]; exists {
		return errors.New("Plugin interface already defined")
	}

	if reflect.TypeOf(cfg.Factory).Kind() != reflect.Struct {
		return errors.New("Factory is not a struct")
	}

	// ensure the Factory implements the interface
	if !implements(reflect.TypeOf(cfg.Factory), reflect.TypeOf(cfg.Interface)) {
		return errors.New("Given factory does not implement the registered plugin interface")
	}

	// ensure the Factory implements the goplugin.Plugin interface,
	// which means by extension means the Registered Interface
	// implements the goplugin.Plugin interface
	if _, ok := (cfg.Factory).(Plugin); !ok {
		return errors.New("Defined plugin interface does not implement the goplugin.Plugin interface")
	}

	// load all the method signatures of the interface
	// config
	factory := reflect.TypeOf(cfg.Factory)
	for i := 0; i < factory.NumMethod(); i++ {
		method := factory.Method(i)
		cfg.methods = append(cfg.methods, method.Name)
	}

	// save the plugin interface definition
	pm.interfaces[cfg.Identifier] = cfg
	return nil
}

func RegisterPlugin(plugin interface{}) {
	setTestVar(1)
}

// Initialze the PluginManager once all the
// PluginInterfaces have been defined, and
// load the plugins
func (pm *PluginManager) Init() {
	err := pm.LoadPlugins()
	if err != nil {
		panic(err)
	}
}

// Return an array of plugins of the given Type t
func (pm PluginManager) Plugins(t Type) []interface{} {
	if plugins, exists := pm.pluginsByType[t]; exists {
		plugin_ifaces := make([]interface{}, 0, len(plugins))
		for _, p := range plugins {
			plugin_ifaces = append(plugin_ifaces, p.instance.Interface())
		}
		return plugin_ifaces
	} else {
		return make([]interface{}, 0) // empty array
	}
}

// Return an individual plugin instance by name
func (pm PluginManager) Plugin(name string) interface{} {
	if p, exists := pm.pluginsByName[name]; exists {
		return p.instance.Interface()
	} else {
		return nil
	}
}

// Load all the plugins that can be found in
// configured plugin Dir
func (pm *PluginManager) LoadPlugins() error {
	return pm.loadPlugins()
}

// Load all the plugins in the Config.Dir location
func (pm *PluginManager) loadPlugins() error {
	files, err := ioutil.ReadDir(pm.cfg.Dir)
	if err != nil {
		return err
	}

	// Make sure all the files were going to load
	// are indeed compiled shared libs of plugins
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".so") {
			err = pm.loadPlugin(pm.cfg.Dir, file.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// load a plugins shared library and integrate it into
// the registered plugin types/interfaces
func (pm *PluginManager) loadPlugin(path, name string) error {
	libPath := path + "/" + name
	lib, err := dl.Open(libPath, dl.RTLD_LAZY)
	if err != nil {
		return err
	}
	// remove the .so from the file name
	name = name[:len(name)-3]

	p := &plugin{
		name: formatPluginName(name),
		lib:  lib,
	}

	// get the type of the plugin
	// and its associated interface definition
	t := p.getType()
	if _, exists := pm.interfaces[t]; !exists {
		return errors.New("Plugin loaded with an undefined type")
	}
	iface := pm.interfaces[t]
	p.iFace = &iface

	// ensure it implements the interface as defined
	// by the type and create a new instance of this plugin
	// via the defined factory, and map
	err = p.bootstrap()
	if err != nil {
		return err
	}

	// save our plugin
	err = pm.savePlugin(p)
	if err != nil {
		return err
	}

	// we're good!
	return nil
}

// save our loaded and bootstrapped plugins
// to the plugin manager
func (pm *PluginManager) savePlugin(p *plugin) error {
	// save the plugin by name
	if _, exists := pm.pluginsByName[p.name]; exists {
		return errors.New("Plugin already exists with that name")
	}
	pm.pluginsByName[p.name] = p

	// save the plugin by type
	// either add it to the existing slice
	// or create a new one
	if plugins, ok := pm.pluginsByType[p.iFace.Identifier]; ok {
		plugins = append(plugins, p)
		pm.pluginsByType[p.iFace.Identifier] = plugins
	} else {
		plugins = []*plugin{p}
		pm.pluginsByType[p.iFace.Identifier] = plugins
	}

	return nil
}

// Grab the type of plugin from the export _PluginType function
func (p plugin) getType() Type {
	var typeFn func() uint16
	fnName := "_Type"
	p.lib.Sym(fnName, &typeFn)
	if typeFn == nil {
		panic("Failed to load plugin: Missing _Type function")
	}
	return Type(typeFn())
}

// Start loading all the required function symbols
// for the given interface
// Return an error if there are methods missing that
// are required for the interface
func (p *plugin) bootstrap() error {
	// create a new instance of the plugin via the
	// factory
	p.instance = reflect.New(reflect.TypeOf(p.iFace.Factory))
	instanceVal := reflect.ValueOf(p.instance.Interface()).Elem()

	// get list of methods
	for _, method := range p.iFace.methods {
		methodExported := getExportedFnName(p.name, method)
		factoryMethodName := getFactoryFnFieldName(method)
		factoryMethod := instanceVal.FieldByName(factoryMethodName)
		if !factoryMethod.IsValid() {
			return errors.New("Plugin factory missing method")
		}
		factoryMethodPtr := factoryMethod.Addr()

		// attach factory method to exported method
		p.lib.Sym(methodExported, factoryMethodPtr.Interface())
		if factoryMethodPtr.IsNil() {
			return errors.New("Could not get exported plugin method symbol")
		}
	}

	return nil
}
