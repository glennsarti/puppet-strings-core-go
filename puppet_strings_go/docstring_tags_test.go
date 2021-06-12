package puppet_strings_go

import (
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
		t.Errorf("%s: Expected list to have %d items, but got %d", prefix, len(expectList), len(list))
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
		"Handles quoted strings",
		" [\"foo, bar\", 'baz, qux', in\"them,iddle\"]",
		"",
		[]string{"\"foo, bar\"", "'baz, qux'", "in\"them,iddle\""},
		"",
	)
}


// def parse_types(types)
// @f.send(:extract_types_and_name_from_text, types)
// end

// it "handles one type" do
// expect(parse_types('[A]')).to eq [nil, ['A'], ""]
// end

// it "handles a list of types" do
// expect(parse_types('[A, B, C]')).to eq [nil, ['A', 'B', 'C'], ""]
// end

// it "handles ducktypes" do
// expect(parse_types('[#foo]')).to eq [nil, ['#foo'], '']
// end

// %w(#foo= #<< #<=> #>> #== #=== Array<#<=>> Array<#==>).each do |meth|
// it "handles ducktypes with special method name #{meth}" do
// 	expect(parse_types("[#{meth}]")).to eq [nil, [meth], '']
// end
// end

// it "only parses #ducktypes inside brackets" do
// expect(parse_types("#ducktype")).to eq [nil, nil, '#ducktype']
// end

// it "returns the text before and after the type list" do
// expect(parse_types(' b <String> description')).to eq ['b', ['String'], 'description']
// expect(parse_types('b c <String> description (test)')).to eq [nil, nil, 'b c <String> description (test)']
// end

// it "does not allow types to start after a newline" do
// v = parse_types("   \n   [X]")
// expect(v).to eq [nil, nil, "[X]"]
// end

// it "handles a complex list of types" do
// v = parse_types(' [Test, Array<String, Hash, C>, String]')
// expect(v).to include(["Test", "Array<String, Hash, C>", "String"])
// end

// it "handles any of the following start/end delimiting chars: (), <>, {}, []" do
// a = parse_types('[a,b,c]')
// b = parse_types('<a,b,c>')
// c = parse_types('(a,b,c)')
// d = parse_types('{a,b,c}')
// expect(a).to eq b
// expect(b).to eq c
// expect(c).to eq d
// expect(a).to include(['a', 'b', 'c'])
// end

// it "returns the text before the type list as the last element" do
// expect(parse_types('b[x, y, z]')).to eq ['b', ['x', 'y', 'z'], '']
// expect(parse_types('  ! <x>')).to eq ["!", ['x'], '']
// end

// it "returns text unparsed if there is no type list" do
// expect(parse_types('')).to eq [nil, nil, '']
// expect(parse_types('[]')).to eq [nil, nil, '[]']
// end

// it "allows A => B syntax" do
// v = parse_types(' [Test, Array<String, Hash{A => {B => C}}, C>, String]')
// expect(v).to include(["Test", "Array<String, Hash{A => {B => C}}, C>", "String"])
// end

// it "handles quoted values" do
// v = parse_types(' ["foo, bar", \'baz, qux\', in"them,iddle"]')
// expect(v).to include(["\"foo, bar\"", "'baz, qux'", 'in"them,iddle"'])
// end