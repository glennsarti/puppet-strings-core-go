package puppet

import (
	"log"
	"strings"

	"github.com/glennsarti/puppet-strings-core-go/yard"
	"github.com/lyraproj/puppet-parser/parser"
)

// Parse Parses a puppet language file
func Parse(fileName string, content string, tasksMode bool, registry *yard.Registry) {
	parseOpts := make([]parser.Option, 0)
	if strings.HasSuffix(fileName, `.epp`) {
		parseOpts = append(parseOpts, parser.EppMode)
	}
	if tasksMode {
		parseOpts = append(parseOpts, parser.TasksEnabled)
	}

	actualParser := parser.CreateParser(parseOpts...)
	programExpr, err := actualParser.Parse(fileName, string(content), false)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// xType := fmt.Sprintf("%T", programExpr)
	// log.Println("!!! ", xType)

	yv := YardVisitor{}
	yv.Visit(programExpr.(*parser.Program), registry)

	// var visitor parser.PathVisitor = func(path []parser.Expression, e parser.Expression) {
	// 	xType := fmt.Sprintf("%T", e)
	// 	log.Println("!!! ", xType, " ", e.String())
	// }

	// var emptyPath []parser.Expression
	// programExpr.AllContents(emptyPath, visitor)
}
