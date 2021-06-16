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
		content,
		ds.typelistOpeningChars(),
		ds.typelistClosingChars(),
	)
	if err != nil { t.Errorf("%s: Expected no error but got %s", prefix, err)}
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

func xxTestExtractTypesAndNameFromText1(t *testing.T) {

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

func assertStringArray(t *testing.T, prefix string, actual []string, expect []string) {
	// Check for nil-ness
	if expect == nil && actual == nil { return }
	if expect == nil && actual != nil {
		t.Errorf("%s: Expected list to be nil but got a list", prefix)
		return
	}
	if expect != nil && actual == nil {
		t.Errorf("%s: Expected a list but got nil", prefix)
		return
	}
	if len(actual) != len(expect) {
		t.Errorf(
			"%s: Expected list to be [%s], but got [%s]",
			prefix,
			strings.Join(expect, ", "),
			strings.Join(actual, ", "),
		)
	} else {
		for i, item := range expect {
			if actual[i] != item {
				t.Errorf("%s: Expected list item %d to be '%s' but got '%s'", prefix, i, item, actual[i])
			}
		}
	}
}


func AssertParseTagWithTypesNameAndDefault(
	t *testing.T,
	prefix string,
	content string,
	expectName string,
	expectText string,
	expectTypeList []string,
	expectDefault []string,
) {
	ds := newDocstring()

	tag, err := ds.parseTagWithTypesNameAndDefault(
		"testtag",
		content,
	)
	if err != nil { t.Errorf("%s: Expected no error but got %s", prefix, err); return }
	if tag == nil { t.Errorf("%s: Expected a tag but got nil", prefix); return }

	if tag.Name != expectName {
		t.Errorf("%s: Expected name to be '%s' but got '%s'", prefix, expectName, tag.Name)
	}
	if tag.Text != expectText {
		t.Errorf("%s: Expected text to be '%s' but got '%s'", prefix, expectText, tag.Text)
	}

	assertStringArray(t, prefix + " (Types)", tag.Types, expectTypeList)
	assertStringArray(t, prefix + " (Defaults)", tag.Defaults, expectDefault)
}


func TestParseTagWithTypesNameAndDefault1(t *testing.T) {
	AssertParseTagWithTypesNameAndDefault(
		t,
		"parses a standard type list with name before types (no default)",
		"NAME [x, y, z] description",
		"NAME",
		"description",
		[]string{"x", "y", "z"},
		nil,
	)

	AssertParseTagWithTypesNameAndDefault(
		t,
		"parses a standard type list with name after types (no default)",
		"  [x, y, z] NAME description",
		"NAME",
		"description",
		[]string{"x", "y", "z"},
		nil,
	)

	AssertParseTagWithTypesNameAndDefault(
		t,
		"parses a tag definition with name, typelist and default",
		"  [x, y, z] NAME (default, values) description",
		"NAME",
		"description",
		[]string{"x", "y", "z"},
		[]string{"default", "values"},
	)

	AssertParseTagWithTypesNameAndDefault(
		t,
		"parses a tag definition with name, typelist and default when name is before type list",
		"  NAME [x, y, z] (default, values) description",
		"NAME",
		"description",
		[]string{"x", "y", "z"},
		[]string{"default", "values"},
	)

	AssertParseTagWithTypesNameAndDefault(
		t,
		"allows typelist to be omitted",
		"  NAME (default, values) description",
		"NAME",
		"description",
		nil,
		[]string{"default", "values"},
	)
}
