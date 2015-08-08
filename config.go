package goplugin

// Configuration for PluginManagers
type PluginManagerConfig struct {
	// Directory to load plugins from
	Dir string
}

// Configuration for a PluginInterface
// definition
type PluginInterfaceConfig struct {
	// An easily used indentifier to match
	// plugins to defined interfaces
	Identifier Type
	// The defined interface which plugins
	// must implement
	Interface interface{}
	// A struct to act as a factory for creating
	// new instances of a loaded plugin of a certain
	// type
	Factory interface{}

	// Lits of methods required by the interface
	methods []string
}
