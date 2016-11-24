/*
 * 			==============================
 *   			  PLEASE READ THIS
 * 			==============================
 *
 * This IS NOT an example of how to use this library
 * or implement a plugin system for your application.
 *
 * From a purely functional point of view, this test file
 * is 100% accurate, however there are some aspects that
 * are against the goplugin convention, or does not have
 * any type safety or compile checking.
 *
 */

package main

import (
	"github.com/jsimnz/goplugin"
)

type myTestPlugin interface {
	Test(string)
}

const example_plugin = goplugin.Type(iota)

var p = testPlugin{"asd"}

type testPlugin struct {
	a string
}

func (tp testPlugin) Test() {
	//return 1 + 1
}

func init() {
	//goplugin.RegisterPlugin((p).(myTestPlugin))
	goplugin.RegisterPlugin(p)
}
