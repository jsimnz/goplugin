/*
 * This file was AUTO GENERATED by GoPlugin
 * to create an externally available interface
 * for the defined plugins
 *
 * The interface relies on C and CGO to 'export'
 * user defined plugins via a simple C-Style
 * interface
 *
 * To read more visit:
 * www.goplugin.org/how-it-works#external-plugin
 */
package main

import "C"

// AUTO GENERATED BY GOPLUGIN
// External interface for PluginAdd Name method
// exported via CGO

//export _Type
func _Type() uint16 {
	return uint16(p.Type())
}

// AUTO GENERATED BY GOPLUGIN
// External interface for PluginAdd Name method
// exported via CGO

//export _PluginAdd_Name_GoPlugin
func _PluginAdd_Name_GoPlugin() string {
	return p.Name()
}

// AUTO GENERATED BY GOPLUGIN
// External interface for PluginAdd Compute method
// exported via CGO

//export _PluginAdd_Compute_GoPlugin
func _PluginAdd_Compute_GoPlugin(arg1 int, arg2 int) int {
	return p.Compute(arg1, arg2)
}