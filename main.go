// +build go1.15

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/glennsarti/puppet-strings-core-go/puppet_strings_go"
)

func main() {
	filename := "/workspaces/puppet-strings-core-go/tests/fixtures/func4x_1.rb"

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// parser := sitter.NewParser()
	// parser.SetLanguage(ruby.GetLanguage())

	// tree := parser.Parse(nil, content)

	// fmt.Println(tree.RootNode())
	// fmt.Println("-----------------")
	// walkNode(tree.RootNode(), 0)

	rb := puppet_strings_go.NewRubyStringsFinder()
	rb.Find(content)




  b, err := json.MarshalIndent(rb.StringsObjects, "", "  ")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return;
	}
	fmt.Println("-- JSON --")
	fmt.Println(string(b))
}


// -----------


// func walkNode(node *sitter.Node, indent int) {
// 	if (node == nil) { return }

// 	it := ""
// 	for i := 0; i < indent; i++ {
// 		it = it + " "
// 	}
// 	fmt.Printf("%s(%d) %s\n", it, node.Symbol(), node.Type())

// 	for i := 0; i < int(node.ChildCount()); i++ {
// 		walkNode(node.Child(i), indent + 1)
// 	}
// }
