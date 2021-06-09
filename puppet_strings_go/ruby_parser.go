// +build go1.15

package puppet_strings_go

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
)

type RubyStringsFinder struct {
	content []byte
	StringsObjects *AllStringsObjects
}

func NewRubyStringsFinder() *RubyStringsFinder {
	return &RubyStringsFinder{}
}

func (rsf *RubyStringsFinder) Find(content []byte) {
	rp := NewRubyParser()

	rp.RegisterHandler(token_sym_method_call, rsf.visitMethodCall)
	rp.RegisterDefaultHandler(rp.VisitChildrenOfHandler)
	rsf.content = content
	rsf.StringsObjects = &AllStringsObjects{
		PuppetFunctions: make([]PuppetFunction, 0),
	}
	rp.ParseAndVisit(content, nil)
}

func (rsf *RubyStringsFinder) visitMethodCall(p *RubyParser, node *sitter.Node, ctx interface{}) bool {
	// First we need to find the name space and actual method name to determine if we're interested in this
	scope, name := rsf.scopeNameOfMethodCall(p, node)

	// Version 4.x functions
	if scope == "Puppet::Functions" && name == "create_function" {
		rsf.visitPuppet4Function(p, node, ctx)
		return true
	}
	// Version 3.x functions
	// if scope == "Puppet::Parser::Functions" && name == "create_function" {
	// 	return false
	// }

	fmt.Printf("%s %s\n", scope, name)
	return false
}

func (rsf *RubyStringsFinder) scopeNameOfMethodCall(p *RubyParser, node *sitter.Node) (string, string) {
	// It could be a method_call with a call node e.g.
	//   (208) method_call
	//     (206) call
	//       (205) scope_resolution
	//         (92) constant
	//         (49) ::
	//         (92) constant
	//       (11) .
	//       (1) identifier
	if callNode := p.FirstChildWithSymbol(node, token_sym_call); callNode != nil {
		scope := ""
		name := ""
		afterDot := false
		for i := 0; i < int(callNode.ChildCount()); i++ {
			cn := callNode.Child(i)
			fmt.Println(cn.Symbol())
			switch int(cn.Symbol()) {
			case token_sym_scope_resolution:
				scope = cn.Content(rsf.content)
			case token_anon_sym_DOT:
				afterDot = true
			case token_sym_identifier:
				if afterDot { name = cn.Content(rsf.content) }
			}
		}
		return scope, name
	}
	// Could be a plain identifier e.g.
	//       (208) method_call
	//         (1) identifier
	//         (210) argument_list
	//           (251) symbol
	//         (216) do_block
	if ident := p.FirstChildWithSymbol(node, token_sym_identifier); ident != nil {
		return "",ident.Content(rsf.content)
	}

	return "", ""
}

// (144) program
// --(104) comment
// --(208) method_call
// ----(206) call
// ------(205) scope_resolution
// --------(92) constant
// --------(49) ::
// --------(92) constant
// ------(11) .
// ------(1) identifier
// ----(210) argument_list
// ------(51) (
// ------(251) symbol
// ------(10) )
// ----(216) do_block
// ------(38) do
// ------(104) comment
// ------(104) comment
// ------(208) method_call
// --------(1) identifier
// --------(210) argument_list
// ----------(251) symbol
// --------(216) do_block
// ----------(38) do
// ----------(208) method_call
// ------------(1) identifier
// ------------(210) argument_list
// --------------(247) string
// ----------------(125) "
// ----------------(277) "
// --------------(13) ,
// --------------(251) symbol
// ----------(24) end
// ------(24) end

func (rsf *RubyStringsFinder) visitPuppet4Function(p *RubyParser, node *sitter.Node, ctx interface{}) {
	pf := PuppetFunction{
		FuncType: "ruby4x",
		Line: int(node.StartPoint().Row), // TODO: This is zero based. Is this right?
		Docstring: ParseDocstring(rsf.previousComments(node)),
		Source: node.Content(rsf.content),
	}

	rsf.StringsObjects.PuppetFunctions = append(rsf.StringsObjects.PuppetFunctions, pf)
}

func (rsf *RubyStringsFinder) previousComments(node *sitter.Node) string {
	result := ""
	addLf := false

	currentLine := node.StartPoint().Row
  for {
		// If there's no more nodes to check, we're done
		if node.PrevSibling() == nil { return result }
		node = node.PrevSibling()
		// If it's not a comment we're done
		if int(node.Symbol()) != token_sym_comment { return result }
		// If it's not on the previous line we're done
		if node.StartPoint().Row != currentLine - 1 { return result }
		currentLine = node.StartPoint().Row
		if addLf {
			result = node.Content(rsf.content) + "\n" + result
		} else {
			result = node.Content(rsf.content)
			addLf = true
		}
	}
}

