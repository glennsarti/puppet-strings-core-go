package puppet_strings_go

import (
	"fmt"
	"regexp"
	"strings"
)

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
	// fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")
	// fmt.Println(content)
	// fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")

	ds := newDocstring()
	ds.parse(content)

	return ds
}

func newDocstring() Docstring {
	return Docstring{
		Tags: make([]DocstringTag, 0),
	}
}


func (ds *Docstring) isTagDirective(tagName string) bool {
	// TODO:
	// list = %w(attribute endgroup group macro method scope visibility)
	// list.include?(tag_name)
	return false
}

func (ds *Docstring) parse(content string) (err error) {
	//sr := NewStringReader(content)
	if content == "" { return }
	lines := strings.Split(content, "\n")
	numLines := len(lines)
	lines = append(lines, "")

	indentRegex := regexp.MustCompile(`\A\s*`)
	emptyRegex := regexp.MustCompile(`\A[^\s]*\z`)
	metaTagRegex := regexp.MustCompile(`\A@(!)?((?:\w\.?)+)(?:\s+(.*))?\z`)

	lastIndent := -1
	origIndent := 0
	directive := false
	lastLine := ""
	tagName := ""
	tagLineBuf := make([]string, 0)

	for index, line := range(lines) {
		fmt.Printf("|%s|\n",line)
		indent := indentRegex.FindStringIndex(line)[1]
		empty := emptyRegex.MatchString(line)
		done := index == numLines

		// if tag_name && (((indent < orig_indent && !empty) || done ||
		//     (indent == 0 && !empty)) || (indent <= last_indent && line =~ META_MATCH))
		if tagName != "" && (((indent < origIndent && !empty) || done || (indent == 0 && !empty)) || (indent <= lastIndent && metaTagRegex.MatchString(line))) {
			if directive || ds.isTagDirective(tagName) {
				// TODO:
				//     directive = create_directive(tag_name, buf)
				//     if directive
				//       docstring << parse_content(directive.expanded_text).chomp
				//     end
			} else {
				t, err := ds.createTag(tagName, tagLineBuf)
				if err != nil {
					fmt.Println("Error creating tag:")
					fmt.Println(err)
				} else {
					ds.Tags = append(ds.Tags, *t)
				}
			}
			tagName = ""
			tagLineBuf = make([]string, 0)
			directive = false
			origIndent = 0
		}

		// # Found a meta tag
		//if line =~ META_MATCH
		if m := metaTagRegex.FindStringSubmatch(line); m != nil {
			directive = (m[1] != "")
			tagName = m[2]
			tagLineBuf = append(tagLineBuf, m[3])
		// elsif tag_name && indent >= orig_indent && !empty
		} else if tagName != "" && indent >= origIndent && !empty {
			if origIndent == 0 { origIndent = indent }
			// Extra data added to the tag on the next line
			if emptyRegex.MatchString(lastLine) { tagLineBuf = append(tagLineBuf, "") }
			tagLineBuf = append(tagLineBuf, regexp.MustCompile(fmt.Sprintf("\\A[ \\t]{%d}", origIndent)).ReplaceAllString(line, ""))

		// elsif !tag_name
		} else if tagName == "" {
			ds.Text += line + "\n"
		}

		lastIndent = indent
		lastLine = line
	}

	//fmt.Printf("%d %s %b", lastIndent, lastLine, string(directive))
	return nil
}

func (ds *Docstring) createTag(tagName string, lines []string) (tag *DocstringTag, err error) {

	switch tagName {
	case "example":
		return ds.parseTagWithTitleAndText(tagName, lines)

	case "return", "raise":
		return ds.parseTagWithTypes(tagName, lines)
	}

	return &DocstringTag{},nil
}
















// # Creates a new directive using the registered {#library}
//     # @return [Tags::Directive] the directive object that is created
//     def create_directive(tag_name, tag_buf)
//       if library.has_directive?(tag_name)
//         dir = library.directive_create(tag_name, tag_buf, self)
//         if dir.is_a?(Tags::Directive)
//           @directives << dir
//           dir
//         end
//       else
//         log.warn "Unknown directive @!#{tag_name}" +
//                  (object ? " in file `#{object.file}` near line #{object.line}" : "")
//         nil
//       end
//     rescue Tags::TagFormatError
//       log.warn "Invalid directive format for @!#{tag_name}" +
//                (object ? " in file `#{object.file}` near line #{object.line}" : "")
//       nil
//     end




//-----


// func ParseDocstring(content string) Docstring {
// 	fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")
// 	fmt.Println(content)
// 	fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")

// 	ds := newDocstring()
// 	ds.parse(content)

// 	return ds
// }

// func newDocstring() Docstring {
// 	return Docstring{
// 		Tags: make([]DocstringTag, 0),
// 	}
// }

// func (ds *Docstring) parse(content string) {
// 	sr := NewStringReader(content)
// 	for {
// 		if sr.IsEOF() { break }
// 		line, _, _ := sr.NextLine()
// 		line = ds.trimDocstringLine(line)

// 		if ds.isATagLine(line) {
// 			tagText := ds.consumeTagText(sr)
// 			tagName, tagInfo := ds.extractTagInfo(line)

// 			// https://github.com/lsegal/yard/blob/b589fa0dc0a21f3304da949dc418e0dc0032182b/lib/yard/tags/default_factory.rb#L101
// 			// https://github.com/lsegal/yard/blob/main/lib/yard/tags/library.rb#L312
// 			switch tagName {
// 			case "raise":
// 				ds.Tags = append(ds.Tags, ds.parseRaiseTag(tagInfo, tagText))
// 			case "example":
// 				ds.Tags = append(ds.Tags, ds.parseExampleTag(tagInfo, tagText))
// 			case "return":
// 				ds.Tags = append(ds.Tags, ds.parseExampleTag(tagInfo, tagText))
// 			default:
// 				fmt.Printf("ERR Unkown Tag name '%s'\n", tagName)
// 			}
// 		} else {
// 			ds.Text += line
// 		}
// 	}
// }

// func (ds *Docstring) isATagLine(line string) bool {
// 	// TODO: A bit naive but should work
// 	return strings.HasPrefix(line, "@")
// }

// func (ds *Docstring) extractTagInfo(line string) (name string, info string) {
// 	// Could be slow using regex?
// 	regex := regexp.MustCompile(`\A@(?P<Name>[^\s]*)[\s]*(?P<Info>.*)\z`)
// 	sm := regex.FindStringSubmatch(line)

// 	return sm[1], sm[2]
// }

// func (ds *Docstring) trimDocstringLine(line string) string {
// 	// Trim any leading hash (comment chars) and whitespace
// 	if strings.HasPrefix(line,"#") { line = line[1:]}
// 	return strings.TrimSpace(line)
// }

// func (ds *Docstring) consumeTagText(sr StringReader) string {
// 	text := ""
// 	for {
// 		if sr.IsEOF() { return text }
// 		line, start, end := sr.PeekNextLine()
// 		line = ds.trimDocstringLine(line)
// 		if ds.isATagLine(line) {
// 			sr.SetPos(start)
// 			return text
// 		} else {
// 			text += line
// 			sr.SetPos(end + 1)
// 		}
// 	}
// }

// // define_tag "Raises",             :raise,       :with_types
// // TODO What about types?!?!?
// func (ds *Docstring) parseRaiseTag(info string, text string) DocstringTag {
// 	return DocstringTag{
// 		TagName: "raise",
// 		Text: info,
// 	}
// }

// //define_tag "Example",            :example, :with_title_and_text
// func (ds *Docstring) parseExampleTag(info string, text string) DocstringTag {
// 	// TODO: what if info and/or text aren't specified
// 	return DocstringTag{
// 		TagName: "example",
// 		Text: text,
// 		Name: info,
// 	}
// }



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


