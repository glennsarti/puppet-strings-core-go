// +build go1.15

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/lyraproj/puppet-parser/json"
	"github.com/lyraproj/puppet-parser/parser"
	"github.com/lyraproj/puppet-parser/pn"
)

// Program to parse and validate a .pp or .epp file
//var validateOnly = flag.Bool("v", false, "validate only")

// var jsonOutput = flag.Bool("j", false, "json output")
var strict = flag.String("s", `off`, "strict (off, warning, or error)")
var tasks = flag.Bool("t", false, "tasks")
var workflow = flag.Bool("w", false, "workflow")

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		pn.Fprintln(os.Stderr, "Usage: parse [options] <pp or epp file to parse>\nValid options are:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fileName := args[0]
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	//strictness := validator.Strict(*strict)

	parseOpts := make([]parser.Option, 0)
	if strings.HasSuffix(fileName, `.epp`) {
		parseOpts = append(parseOpts, parser.EppMode)
	}
	if *tasks {
		parseOpts = append(parseOpts, parser.TasksEnabled)
	}
	if *workflow {
		parseOpts = append(parseOpts, parser.WorkflowEnabled)
	}

	actualParser := parser.CreateParser(parseOpts...)
	expr, err := actualParser.Parse(args[0], string(content), false)
	if err != nil {
		pn.Fprintln(os.Stderr, err.Error())
		// Parse error is always SeverityError
		os.Exit(1)
	}

	// v := validator.ValidatePuppet(expr, strictness)
	// if len(v.Issues()) > 0 {
	// 	severity := issue.Severity(issue.SeverityIgnore)
	// 	for _, i := range v.Issues() {
	// 		pn.Fprintln(os.Stderr, i.String())
	// 		if i.Severity() > severity {
	// 			severity = i.Severity()
	// 		}
	// 	}
	// 	if severity == issue.SeverityError {
	// 		os.Exit(1)
	// 	}
	// }

	var visitor parser.PathVisitor = func(path []parser.Expression, e parser.Expression) {
		xType := fmt.Sprintf("%T", e)
		fmt.Println("!!! ", xType, " ", e.String())
	}

	b := bytes.NewBufferString(``)
	expr.ToPN().Format(b)
	//expr.(*parser.Program).ToPN().ToData()
	var emptyPath []parser.Expression
	expr.AllContents(emptyPath, visitor)
	//expr.lo
	//expr.String()
	fmt.Println(expr.(*parser.Program).Locator().String())
	//pn.Println(b)
	//}
}

func emitJson(value interface{}) {
	b := bytes.NewBufferString(``)
	json.ToJson(value, b)
	pn.Println(b.String())
}
