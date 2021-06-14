package puppet_strings_go

import (
	"fmt"
	"strings"
	"testing"
)

// Based on https://github.com/lsegal/yard/blob/main/spec/tags/default_factory_spec.rb

func AssertExtractTypesAndNameFromText(t *testing.T, prefix string, content string, expectBefore string, expectList []string, expectText string) {
	ds := newDocstring()

	before, list, text, err := ds.extractTypesAndNameFromText(
		strings.Split(content, "\n"),
		ds.typelistOpeningChars(),
		ds.typelistClosingChars(),
	)
	if err != nil { t.Errorf("%s: Expected not error but got %s", prefix, err)}
	if before != expectBefore {
		t.Errorf("%s: Expected before to be '%s' but got '%s'", prefix, expectBefore, before)
	}
	if text != expectText {
		t.Errorf("%s: Expected text to be '%s' but got '%s'", prefix, expectText, text)
	}

	if len(list) != len(expectList) {
		t.Errorf(
			"%s: Expected list to be [%s], but got [%s]",
			prefix,
			strings.Join(expectList, ", "),
			strings.Join(list, ", "),
		)
	} else {
		for i, item := range expectList {
			if list[i] != item {
				t.Errorf("%s: Expected list item %d to be '%s' but got '%s'", prefix, i, item, list[i])
			}
		}
	}
}

func TestExtractTypesAndNameFromText1(t *testing.T) {

	AssertExtractTypesAndNameFromText(t,
		"One type",
		"[A]",
		"",
		[]string{"A"},
		"",
	)
	AssertExtractTypesAndNameFromText(t,
		"List of types",
		"[A,B,C]",
		"",
		[]string{"A", "B", "C"},
		"",
	)

	AssertExtractTypesAndNameFromText(t,
		"Ducktypes",
		"[#foo]",
		"",
		[]string{"#foo"},
		"",
	)
	for _, methName := range []string{
		"#foo=",
		"#<<",
		"#<=>",
		"#>>",
		"#==",
		"#===",
		"Array<#<=>>",
		"Array<#==>",
	} {
		AssertExtractTypesAndNameFromText(t,
			fmt.Sprintf("Duck type with special method %s",methName),
			fmt.Sprintf("[%s]",methName),
			"",
			[]string{methName},
			"",
		)
	}
	AssertExtractTypesAndNameFromText(t,
		"Only parses duck types in a type list",
		"#ducktype",
		"",
		[]string{},
		"#ducktype",
	)

	AssertExtractTypesAndNameFromText(t,
		"Text before and after type list",
		" b <String> description",
		"b",
		[]string{"String"},
		"description",
	)
	AssertExtractTypesAndNameFromText(t,
		"Type list in the wrong position",
		"b c <String> description (test)",
		"",
		[]string{},
		"b c <String> description (test)",
	)

	AssertExtractTypesAndNameFromText(t,
		"No types after newline",
		"   \n   [X]",
		"",
		[]string{},
		"[X]",
	)

	AssertExtractTypesAndNameFromText(t,
		"Handles complex list of types",
		" [Test, Array<String, Hash, C>, String]",
		"",
		[]string{
			"Test",
			"Array<String, Hash, C>",
			"String",
		},
		"",
	)

	for _, content := range []string{
		"[a,b,c]",
		"<a,b,c>",
		"(a,b,c)",
		"{a,b,c}",
	} {
		AssertExtractTypesAndNameFromText(t,
			"Handles " + content,
			content,
			"",
			[]string{"a", "b", "c"},
			"",
		)
	}

	AssertExtractTypesAndNameFromText(t,
		"Returns the text before the type list as the last element1",
		"b[x, y, z]",
		"b",
		[]string{"x", "y", "z"},
		"",
	)
	AssertExtractTypesAndNameFromText(t,
		"Returns the text before the type list as the last element2",
		"  ! <x>",
		"!",
		[]string{"x"},
		"",
	)

	AssertExtractTypesAndNameFromText(t,
		"Returns empty result for an empty string",
		"",
		"",
		[]string{},
		"",
	)
	AssertExtractTypesAndNameFromText(t,
		"Returns text unparsed if there is no type list",
		"[]",
		"",
		[]string{},
		"[]",
	)

	AssertExtractTypesAndNameFromText(t,
		"Handles A => B syntax",
		" [Test, Array<String, Hash{A => {B => C}}, C>, String]",
		"",
		[]string{
			"Test",
			"Array<String, Hash{A => {B => C}}, C>",
			"String",
		},
		"",
	)

	AssertExtractTypesAndNameFromText(t,
		"Handles quoted strings",
		" [\"foo, bar\", 'baz, qux', in\"them,iddle\"]",
		"",
		[]string{"\"foo, bar\"", "'baz, qux'", "in\"them,iddle\""},
		"",
	)
}
