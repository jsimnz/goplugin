package subcommands

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	//"github.com/beyang/go-astquery"
	"github.com/GeertJohan/go.rice"
	"github.com/codegangsta/cli"
	"github.com/yookoala/realpath"
)

type pluginDef struct {
	node     *ast.CallExpr
	name     string // Plugin name
	instance string // Instance variable for plugin
}

type pluginMethod struct {
	name       string        // Method Name
	node       *ast.FuncDecl // ast node representing method definition
	params     []*ast.Field  // List of inputs of method
	results    []*ast.Field  // List of retuned values
	pluginName string
	instance   string
}

func (pm *pluginMethod) addInstance(inst string) {
	pm.instance = inst
}

func (pm pluginMethod) ExportFnName() string {
	return fmt.Sprintf("_%v_%v_GoPlugin",
		strings.Title(pm.pluginName),
		strings.Title(pm.name))
}

func (pm pluginMethod) FormattedParams() string {
	if len(pm.params) > 0 {
		var output string
		for i, p := range pm.params {
			if i == 0 {
				output += fmt.Sprintf("arg%v %v",
					i, typeToString(p.Type))
			} else {
				output += fmt.Sprintf(", arg%v %v",
					i, typeToString(p.Type))
			}
		}

		return output
	}

	return ""
}

var (
	hasMainFunc       bool
	pluginMethods     []*pluginMethod
	registeredPlugins []pluginDef
)

func visit(path string, f os.FileInfo, err error) error {
	if strings.HasSuffix(path, ".go") {
		fmt.Printf("Visited: %s\n", path)
		src, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		parseFile(f.Name(), string(src))
	}
	return nil
}

func parseFile(filename, src string) {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		panic(err)
	}

	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {

		case *ast.CallExpr:
			s := getNodeString(x, fset)
			if strings.HasPrefix(s, "goplugin.RegisterPlugin") {
				fmt.Printf("FOUND RegisterPlugin call... %v\n", s)
				// plugin := pluginDef{
				// 	node: x,
				// 	name: x.
				// }

				//var pluginType string
				if len(x.Args) > 0 {
					switch arg := (x.Args[0]).(type) {
					case *ast.Ident:
						fmt.Printf("%#v\n", arg)
						plugin := pluginDef{
							name:     getIdentType(arg),
							instance: (arg).Name,
							node:     x,
						}
						registeredPlugins = append(registeredPlugins, plugin)

						//case *ast.TypeAssertExpr:
					}
				} else {
					panic(errors.New("Invalid call to RegisterPlugin"))
				}
			}

		case *ast.FuncDecl:
			// check if this function def is a method def
			if x.Recv != nil {
				method := pluginMethod{
					pluginName: getNodeString(x.Recv.List[0].Type, fset),
					name:       x.Name.String(),
					node:       x,
					params:     getMethodParams(x),
					results:    getMethodResults(x),
				}
				pluginMethods = append(pluginMethods, &method)
			} else if strings.HasPrefix(getNodeString(x, fset), "func main()") {
				fmt.Println("FOUND MAIN FUNC!")
				hasMainFunc = true
			}
		}

		return true
	})
}

func getMethodParams(x *ast.FuncDecl) []*ast.Field {
	if x.Type.Params != nil {
		return x.Type.Params.List
	}
	return nil
}

func getMethodResults(x *ast.FuncDecl) []*ast.Field {
	if x.Type.Results != nil {
		return x.Type.Results.List
	}
	return nil
}

func GenerateCmd(c *cli.Context) {
	//root := "/home/jsimnz/Workspace/Go/src/github.com/jsimnz/goplugin/test_plugin/"
	root, err := realpath.Realpath(c.String("path"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Searching path (%v)\n", root)
	filepath.Walk(root, visit)

	if len(registeredPlugins) > 1 {
		fmt.Println("Generate failed!")
		fmt.Println(" - Detected more then one 'RegisterPlugin' call")
		return
	} else if len(registeredPlugins) == 1 {
		fmt.Println("Found registered plugin call!")
	} else {
		fmt.Println("No plugins detected")
		return
	}

	for _, m := range pluginMethods {
		m.addInstance(registeredPlugins[0].instance)
		fmt.Printf("Plugin Name: %v (instance: %v), Method name: %v, Args: %v, Returns: %v\n",
			m.pluginName, m.instance, m.name, m.params[0].Type, m.results)
	}

	templateBox, err := rice.FindBox("templates")
	if err != nil {
		panic(err)
	}
	exportTemplateSrc, err := templateBox.String("cgo_export.tmpl")
	if err != nil {
		panic(err)
	}

	t, err := template.New("export").Parse(exportTemplateSrc)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n====\n")
	t.Execute(os.Stdout, pluginMethods)
}
