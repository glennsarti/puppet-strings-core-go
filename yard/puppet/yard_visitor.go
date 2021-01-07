package puppet

import (
	"regexp"

	"github.com/glennsarti/puppet-strings-core-go/yard"
	"github.com/lyraproj/puppet-parser/parser"
)

type (
	// YardVisitor A Puppet AST visitor to extract YARD compatabile expressions
	YardVisitor struct {
		registry    *yard.Registry
		lineIndexes []int
		content     string
	}
)

// Visit Traverse a Puppet AST Program
func (v YardVisitor) Visit(programExpr *parser.Program, registry *yard.Registry) {
	v.registry = registry

	// TODO: I don't like that I have to copy the code from puppet-parser to extract
	// out the line information, but here we are.
	v.lineIndexes = append(make([]int, 0, 32), 0)
	v.content = programExpr.String()
	rdr := parser.NewStringReader(v.content)
	for c, _ := rdr.Next(); c != 0; c, _ = rdr.Next() {
		if c == '\n' {
			v.lineIndexes = append(v.lineIndexes, rdr.Pos())
		}
	}

	var emptyPath []parser.Expression
	programExpr.AllContents(emptyPath, v.visitor)
}

func (v YardVisitor) visitor(path []parser.Expression, e parser.Expression) {
	switch e.(type) {
	case *parser.TypeAlias:
		{
			v.visitTypeAlias(&path, e.(*parser.TypeAlias))
		}
	case *parser.PlanDefinition:
		{
			v.visitPlanDefinition(&path, e.(*parser.PlanDefinition))
		}
	default:
		{
			// xType := fmt.Sprintf("%T", e)
			// log.Println("UNKNOWN ", xType, " ", e.String())
		}
	}
}

func (v YardVisitor) getLine(lineNum int) string {
	if lineNum < 0 || lineNum > len(v.lineIndexes) {
		return ""
	}

	startIndex := v.lineIndexes[lineNum]
	if lineNum < len(v.lineIndexes) {
		return v.content[startIndex : v.lineIndexes[lineNum+1]-1] // -1 is to strip the trailing '/n'
	}
	return v.content[startIndex:len(v.content)]
}

func (v YardVisitor) extractComments(fromLine int, loc *(parser.Locator)) string {
	comment := regexp.MustCompile(`^\s*#+\s?`)

	docString := ""
	for lineNum := fromLine - 1; lineNum >= 0; lineNum-- {
		lineText := v.getLine(lineNum)
		newLineText := comment.ReplaceAllString(lineText, "")
		if len(lineText) == len(newLineText) {
			break
		}

		if len(docString) == 0 {
			docString = newLineText
		} else {
			docString = newLineText + "\n" + docString
		}
	}

	return docString
}

func (v YardVisitor) visitTypeAlias(path *[]parser.Expression, e *parser.TypeAlias) {
	o := DataTypeAlias{
		Name:    e.Name(),
		File:    e.File(),
		Line:    e.Line(),
		AliasOf: e.Type().String(),
		// TODO: Docstring
	}

	v.registry.All = append(v.registry.All, o)
}

func (v YardVisitor) visitPlanDefinition(path *[]parser.Expression, e *parser.PlanDefinition) {
	o := PuppetPlan{
		Name:   e.Name(),
		File:   e.File(),
		Line:   e.Line(),
		Source: e.String(),
		// TODO: defaults
		DocString: yard.CreateDocString(v.extractComments(e.Line()-1, e.Locator())), // Puppet Line Numbers are 1-based
	}

	// TODO: defaults
	v.registry.All = append(v.registry.All, o)
}
