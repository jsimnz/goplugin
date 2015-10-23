package main

import "C"

//export _TestPlugin_Test_GoPlugin
func _TestPlugin_Test_GoPlugin() {
	p.Test()
}

// //export _TestPlugin_Type_GoPlugin
// func _TestPlugin_Type_GoPlugin() uint16 {
// 	return uint16(p.Type())
// }

//export _Type
func _Type() uint16 {
	return uint16(example_plugin)
}

func main() {}
