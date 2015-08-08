package goplugin

import (
	"fmt"
	"reflect"
	"strings"
)

// format the filename of a plugin to its
// idomatic version
// ex. test_plugin -> TestPlugin
func formatPluginName(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = strings.Title(str)
	return strings.Replace(str, " ", "", -1)
}

func getExportedFnName(pluginName, fnName string) string {
	return fmt.Sprintf("_%s_%s_GoPlugin", pluginName, fnName)
}

func getFactoryFnFieldName(fnName string) string {
	return fnName + "Fn"
}

// check if a struct val implements some interface
func implements(val reflect.Type, iface reflect.Type) bool {
	ifaceType := iface.Elem()
	return val.Implements(ifaceType)
}
