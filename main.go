// +build go1.15

package main

import (
	"encoding/json"
	"fmt"

	"github.com/glennsarti/puppet-strings-core-go/puppet_strings_go"
)

func main() {
	// filename := "/workspaces/puppet-strings-core-go/tests/fixtures/func4x_1.rb"

	// content, err := ioutil.ReadFile(filename)
	// if err != nil {
	// 	panic(err)
	// }


	// // Parser
	// rb := puppet_strings_go.NewRubyStringsFinder()
	// rb.Find(content)
	// b, err := json.MarshalIndent(rb.StringsObjects, "", "  ")
	// if err != nil {
	// 	fmt.Printf("Error: %s", err)
	// 	return;
	// }
	// fmt.Println("-- JSON --")
	// fmt.Println(string(b))

	content := "An overview for the first overload.\n" +
	"@raise SomeError this is some error\n" +
	"@param param1 The first parameter.\n" +
	"@param param2 The second parameter.\n" +
	"@option param2 [String] :option an option\n" +
	"@option param2 [String] :option2 another option\n" +
	"@param param3 The third parameter.\n" +
	"@param param4 The fourth parameter.\n" +
	"@enum param4 :one Option one.\n" +
	"@enum param4 :two Option two.\n" +
	"@return Returns nothing.\n" +
	"@return [Undef]\n" +
	"@example Calling the function foo\n" +
	"  $result = func4x(1, 'foooo')\n" +
	"\n"

	result := puppet_strings_go.ParseDocstring(content)
	b, err := json.MarshalIndent(result, "", "  ")
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
