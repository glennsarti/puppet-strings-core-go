// https://github.com/goruby/goruby/blob/43c93d4f2611780d923d8616b9aed61bd013384b/parser/interface.go#L61
// https://github.com/goruby/goruby/blob/43c93d4f2611780d923d8616b9aed61bd013384b/object/kernel.go#L186
package ruby

import (
	"container/list"
	"fmt"
	"go/token"
	"io/ioutil"
	"log"

	"github.com/glennsarti/puppet-strings-core-go/yard"
	"github.com/goruby/goruby/ast"
	"github.com/goruby/goruby/parser"
)

type (
	// YardVisitor blah
	YardVisitor struct {
		Registry *yard.Registry
		Root     *ast.Program
		//lineIndexes []int
		//content     string
	}
)

// Parse Parses a puppet language file
func Parse(fileName string, content string, registry *yard.Registry) {
	file, _ := ioutil.ReadFile(fileName)

	prog, _ := parser.ParseFile(token.NewFileSet(), fileName, file, 0)

	yv := YardVisitor{
		Registry: registry,
		Root:     prog,
	}
	ast.Walk(&yv, prog)
}

// Visit blah
func (v *YardVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		// Can't do anything with nil
		return v
	}
	//xType := fmt.Sprintf("%T", node)
	//out := fmt.Sprintf("Evaluating: %s: %s", xType, node.String())
	//out := fmt.Sprintf("Evaluating: %s %s", xType, node.TokenLiteral())
	//out := fmt.Sprintf("%s", xType)

	switch node.(type) {
	// case *ast.ExpressionStatement:
	// 	{
	// 		e.evalExpressionStatement(node.(*ast.ExpressionStatement))
	// 	}
	// case *ast.ScopedIdentifier:
	// 	{
	// 		e.evalScopedIdentifier(node.(*ast.ScopedIdentifier))
	// 	}
	case *ast.ContextCallExpression:
		{
			v.evalContextCallExpression(node.(*ast.ContextCallExpression))
		}
	// case *ast.Identifier:
	// 	{
	// 	}
	// case *ast.SymbolLiteral:
	// 	{
	// 	}
	default:
		xType := fmt.Sprintf("%T", node)
		log.Println("Ignoring ", xType)
	}
	return v
}

func (v *YardVisitor) evalContextCallExpression(node *ast.ContextCallExpression) {

	// if node.Context != nil {
	// 	log.Println("%%%%%", node.Context.String())
	// }
	log.Println("%%%%%", node.Function.Value)

	path, _ := ast.Path(v.Root, node)
	if path == nil {
		log.Println("Could not get path for ", node)
		return
	}

	if node.Function.Value != "newfunction" {
		log.Println("%%%%%", node.Function.Value)

	}

	return
	rubyModule := v.getRubyModule(path)

	log.Println("CONTEXT: ", rubyModule)

	switch node.Function.Value {
	case "create_function":
		{
			// 4x Function call
			// Puppet::Functions.create_function(:'mymodule::upcase') do
			if rubyModule != "Puppet::Functions" {
				return
			}

			// Possible func 4x
		}
	case "newfunction":
		{
			// 3x Function call
			// module Puppet::Parser::Functions
			//   newfunction(:write_line_to_file) do |args|
			//     filename = args[0]
			//     str = args[1]
			//     File.open(filename, 'a') {|fd| fd.puts str }
			//   end
			// end
			//panic("NOT IMPLEMENTED")
		}
	}

	// // switch node.Function.Value {
	// // case "create_function":{

	// // }
	// // }

	// log.Println("!!!!!!!!!!!! ", node.Function.Value)
	// e.Eval(node.Context)
	// e.evalExpressions(node.Arguments)
	// if node.Block != nil {
	// 	e.Eval(node.Block)
	// }
}

// func (v *YardVisitor) createPuppetFunction(node *ast.ContextCallExpression, context string) {
// }

func (v *YardVisitor) getRubyModule(list *list.List) string {
	context := ""
	for item := list.Front(); item != nil; item = item.Next() {

		switch item.Value.(type) {
		case *ast.ExpressionStatement:
			{
				context += item.Value.(*ast.ExpressionStatement).TokenLiteral()
			}
		case *ast.ScopedIdentifier:
			{
				context += item.Value.(*ast.ScopedIdentifier).TokenLiteral()
			}
		case *ast.ContextCallExpression:
			{
				node := item.Value.(*ast.ContextCallExpression)
				if node.Context != nil {
					context += item.Value.(*ast.ContextCallExpression).Context.TokenLiteral()
				}

			}
		default:
			{
				xType := fmt.Sprintf("%T", item.Value)

				log.Println("getContext Ignoring ", xType)
			}
		}

	}

	return context
}
