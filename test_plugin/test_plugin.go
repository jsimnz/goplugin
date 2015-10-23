/*
 * 			==============================
 *   			  PLEASE READ THIS
 * 			==============================
 *
 * This IS NOT an example of how to use this library
 * and implement a plugin system for your application.
 *
 * From a purely functional point of view, this text file
 * is 100% accurate, however there are some aspects that 
 * are against the goplugin convention, or does not have
 * any type safety or compile checking.
 *
 */

package main

import (
	"github.com/jsimnz/goplugin"
)

const example_plugin = goplugin.Type(iota)

var p = testPlugin{}

type testPlugin struct{}

func (tp testPlugin) Test() {}

func init() {
	goplugin.RegisterPlugin(testPlugin{})
}
