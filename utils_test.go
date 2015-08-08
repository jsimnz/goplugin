package goplugin

import (
	"reflect"
	"testing"
)

type Tester interface {
	Test()
}

type MyTest struct{}

func (m MyTest) Test() {}

func TestFormatPluginName(t *testing.T) {
	pluginName := formatPluginName("test_plugin")
	expected := "TestPlugin"
	if pluginName != expected {
		t.Errorf("Expected: %v, got: %v", expected, pluginName)
	}
}

func TestGetExportedFnName(t *testing.T) {
	pluginName := "MyPlugin"
	fnName := "MyMethod"
	exportedFnName := getExportedFnName(pluginName, fnName)
	expected := "_MyPlugin_MyMethod_GoPlugin"
	if exportedFnName != expected {
		t.Errorf("Expected: %v, got: %v", expected, exportedFnName)
	}
}

func TestGetFactoryFnFieldName(t *testing.T) {
	fnName := "MyMethod"
	factoryFnFieldName := getFactoryFnFieldName(fnName)
	expected := "MyMethodFn"
	if factoryFnFieldName != expected {
		t.Errorf("Expected: %v, got: %v", expected, factoryFnFieldName)
	}
}

func TestImplements(t *testing.T) {
	iface := (*Tester)(nil)
	val := MyTest{}
	doesImplement := implements(reflect.TypeOf(val),
		reflect.TypeOf(iface))
	if !doesImplement {
		t.Errorf("Expected: true, got: %v", doesImplement)
	}
}
