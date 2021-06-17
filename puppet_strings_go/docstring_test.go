package puppet_strings_go

import (
	"testing"
)

func assertDocstringTag(
	t *testing.T,
	prefix string,
	expect DocstringTag,
	actual DocstringTag,
) (pass bool) {
	pass = true
	pass = pass && assertString(t, "%s: Expected TagName '%s' but got '%s'", prefix, expect.TagName, actual.TagName)
	pass = pass && assertString(t, "%s: Expected Text '%s' but got '%s'", prefix, expect.Text, actual.Text)
	pass = pass && assertStringArray(t, prefix + "(Types)", actual.Types, expect.Types)
	pass = pass && assertString(t, "%s: Expected Name '%s' but got '%s'", prefix, expect.Name, actual.Name)
	return
}

func assertDefaultDocstringTag(
	t *testing.T,
	prefix string,
	expect *DefaultDocstringTag,
	actual *DefaultDocstringTag,
) (pass bool) {
	if expect == nil && actual == nil { return true }
	if expect == nil && actual != nil {
		t.Errorf("%s: Expected tag to be nil but got a tag", prefix)
		return false
	}
	if expect != nil && actual == nil {
		t.Errorf("%s: Expected a tag but got nil", prefix)
		return false
	}

	pass = assertDocstringTag(t, prefix, expect.DocstringTag, actual.DocstringTag)
	pass = pass && assertStringArray(t, prefix + "(Default Types)", actual.Types, expect.Types)
	return
}

func assertOptionsDocstringTag(
	t *testing.T,
	prefix string,
	expect OptionsDocstringTag,
	actual OptionsDocstringTag,
) (pass bool) {
	pass = assertDocstringTag(t, prefix, expect.DocstringTag, actual.DocstringTag)
	pass = pass && assertDefaultDocstringTag(t, prefix, expect.Pair, actual.Pair)
	return
}


func AssertDocString(
	t *testing.T,
	prefix string,
	content string,
	expectTagCount int,
	expectOptionTagCount int,
	expectText string,
) (ds Docstring) {
	ds = ParseDocstring(content)

	pass := true
	pass = pass && assertInteger(t, "%s: Expected %d tag/s but got %d tag/s", prefix, expectTagCount, len(ds.Tags))
	pass = pass && assertInteger(t, "%s: Expected %d option tag/s but got %d tag/s", prefix, expectOptionTagCount, len(ds.OptionTags))
	pass = pass && assertString(t, "%s: Expected text '%s' but got '%s'", prefix, expectText, ds.Text)

	if !pass { t.FailNow() }
	return
}

func TestParseOptionTag(t *testing.T) {
	pfx := "Option tag with a type"
	ds := AssertDocString(
		t,
		pfx,
		"@option param2 [String] :option an option",
		0,1,"",
	)
	assertOptionsDocstringTag(t, pfx, OptionsDocstringTag{
		DocstringTag: DocstringTag{
			TagName: "option",
			Name: "param2",
		},
		Pair: &DefaultDocstringTag{
			DocstringTag: DocstringTag{
				TagName: "option",
				Name: ":option",
				Text: "an option",
				Types: []string{"String"},
			},
		},
	} , ds.OptionTags[0])

	// From https://github.com/lsegal/yard/blob/main/spec/docstring_spec.rb#L308
	pfx = "handles full @option tags"
	ds = AssertDocString(
		t,
		pfx,
		"@option foo [String] bar (nil) baz",
		0,1,"",
	)
	assertOptionsDocstringTag(t, pfx, OptionsDocstringTag{
		DocstringTag: DocstringTag{
			TagName: "option",
			Name: "foo",
		},
		Pair: &DefaultDocstringTag{
			DocstringTag: DocstringTag{
				TagName: "option",
				Name: "bar",
				Text: "baz",
				Types: []string{"String"},
			},
			Defaults: []string{"nil"},
		},
	} , ds.OptionTags[0])

	pfx = "handles simple @option tags"
	ds = AssertDocString(
		t,
		pfx,
		"@option foo :key bar",
		0,1,"",
	)
	assertOptionsDocstringTag(t, pfx, OptionsDocstringTag{
		DocstringTag: DocstringTag{
			TagName: "option",
			Name: "foo",
		},
		Pair: &DefaultDocstringTag{
			DocstringTag: DocstringTag{
				TagName: "option",
				Name: ":key",
				Text: "bar",
			},
		},
	} , ds.OptionTags[0])
}



// func xxTestParseReturnTag(t *testing.T) {

// 	xx := Docstring{
// 		Text: "docstring text",
// 		Tags: []DocstringTag{
// 			{
// 				TagName: "asdadasd",
// 			},
// 		},
// 	}

// 	// xx := OptionsDocstringTag{
// 	// 	DocstringTag: DocstringTag{
// 	// 		TagName: "tag",
// 	// 	},
// 	// 	Pair: &DefaultDocstringTag{
// 	// 		Defaults: nil,
// 	// 	},
// 	// }

// 	b, err := json.MarshalIndent(xx, "", "  ")
// 	//b, err := json.Marshal(xx)
// 	if err != nil {
// 		fmt.Printf("Error: %s", err)
// 		return;
// 	}
// 	fmt.Println("-- JSON2 --")
// 	fmt.Println(string(b))
// 	fmt.Println("-- JSON2 --")



// 	AssertSingleDocString(
// 		t,
// 		"Return tag with no type",
// 		"@return Returns nothing.\n",
// 		"return",
// 		"Returns nothing.",
// 	)
// }
