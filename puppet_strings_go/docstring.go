package puppet_strings_go

import "fmt"

type Docstring struct {
	Text string `json:"text,omitempty"`
	Tags []DocstringTag `json:"tags"`
}

type DocstringTag struct {
	TagName string `json:"tag_name,omitempty"`
	Text string `json:"text,omitempty"`
	Types []string `json:"types,omitempty"`
	Name string `json:"name,omitempty"`
}

func ParseDocstring(content string) Docstring {
	fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")
	fmt.Println(content)
	fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")

	return Docstring{}
}

// "puppet_functions": [
// 	{
// 		"name": "func4x_1",
// 		"file": "/workspaces/puppet-strings-core-go/tests/fixtures/func4x_1.rb",
// 		"line": 2,
// 		"type": "ruby4x",

// 		"docstring": {
// 			"text": "An example 4.x function with only one signature.",
// 			"tags": [
// 				{
// 					"tag_name": "param",
// 					"text": "The first parameter.",
// 					"types": [
// 						"Integer"
// 					],
// 					"name": "param1"
// 				},
// 				{
// 					"tag_name": "return",
// 					"text": "Returns nothing.",
// 					"types": [
// 						"Undef"
// 					]
// 				}
// 			]
// 		},
// 		"source": "Puppet::Functions.create_function(:func4x_1) do\n  # @param param1 The first parameter.\n  # @return [Undef] Returns nothing.\n  dispatch :foobarbaz do\n    param          'Integer',       :param1\n  end\nend"


// 		"signatures": [
// 			{
// 				"signature": "func4x_1(Integer $param1)",
// 				"docstring": {
// 					"text": "An example 4.x function with only one signature.",
// 					"tags": [
// 						{
// 							"tag_name": "param",
// 							"text": "The first parameter.",
// 							"types": [
// 								"Integer"
// 							],
// 							"name": "param1"
// 						},
// 						{
// 							"tag_name": "return",
// 							"text": "Returns nothing.",
// 							"types": [
// 								"Undef"
// 							]
// 						}
// 					]
// 				}
// 			}
// 		],
// 	}


