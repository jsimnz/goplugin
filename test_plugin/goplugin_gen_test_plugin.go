package main

import "C"

//export _TestPlugin_Test_GoPlugin
func _TestPlugin_Test_GoPlugin() {
	p.Test()
}

//export _Type
func _Type() uint16 {
	return uint16(example_plugin)
}

func main() {}
