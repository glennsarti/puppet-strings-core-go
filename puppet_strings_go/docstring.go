package puppet_strings_go

import (
	"fmt"
	"regexp"
	"strings"
)

type Docstring struct {
	Text string
	Tags []DocstringTag
	OptionTags []OptionsDocstringTag
}

func newDocstring() Docstring {
	return Docstring{
		Tags: make([]DocstringTag, 0),
		OptionTags: make([]OptionsDocstringTag, 0),
	}
}

func ParseDocstring(content string) Docstring {
	ds := newDocstring()
	ds.parse(content)
	return ds
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
				if err := ds.createTag(tagName, strings.Join(tagLineBuf, "\n")); err != nil {
					fmt.Println("Error creating tag:")
					fmt.Println(err)
				}
			}
			tagName = ""
			tagLineBuf = make([]string, 0)
			directive = false
			origIndent = 0
		}

		// # Found a meta tag
		if m := metaTagRegex.FindStringSubmatch(line); m != nil {
			directive = (m[1] != "")
			tagName = m[2]
			tagLineBuf = append(tagLineBuf, m[3])
		} else if tagName != "" && indent >= origIndent && !empty {
			if origIndent == 0 { origIndent = indent }
			// Extra data added to the tag on the next line
			if emptyRegex.MatchString(lastLine) { tagLineBuf = append(tagLineBuf, "") }
			tagLineBuf = append(tagLineBuf, regexp.MustCompile(fmt.Sprintf("\\A[ \\t]{%d}", origIndent)).ReplaceAllString(line, ""))
		} else if tagName == "" {
			if ds.Text != "" { ds.Text += "\n"}
			ds.Text += line
		}

		lastIndent = indent
		lastLine = line
	}

	//fmt.Printf("%d %s %b", lastIndent, lastLine, string(directive))
	return nil
}

func (ds *Docstring) createTag(tagName string, text string) (err error) {
	switch tagName {
	case "example":
		if t, err := ds.parseTagWithTitleAndText(tagName, text); err == nil {
			ds.Tags = append(ds.Tags, *t)
		} else { return err }

	case "return", "raise":
		if t, err := ds.parseTagWithTypes(tagName, text); err == nil {
			ds.Tags = append(ds.Tags, *t)
		} else { return err }

	case "param":
		if t, err := ds.parseTagWithTypesAndName(tagName, text); err == nil {
			ds.Tags = append(ds.Tags, *t)
		} else { return err }

	case "option":
		if t, err := ds.parseTagWithOptions(tagName, text); err == nil {
			ds.OptionTags = append(ds.OptionTags, *t)
		} else { return err }

	// Custom Puppet Stings tags
	case "enum":
		// Enum tag is just an Option tag with a different tagname
		// Ref - https://github.com/puppetlabs/puppet-strings/blob/main/lib/puppet-strings/yard/tags/enum_tag.rb
		if t, err := ds.parseTagWithOptions(tagName, text); err == nil {
			ds.OptionTags = append(ds.OptionTags, *t)
		} else { return err }

	default:
		//return errors.New(fmt.Sprintf("Unknown tag '@%s'", tagName))
		fmt.Printf("!!! Unknown tag '@%s'\n", tagName)
		return nil
	}

// - abstract
// - api
// - attr
// - attr_reader
// - attr_writer
// - author
// - deprecated
// x example
// - note
// x option
// - overload
// x param
// - private
// x raise
// x return
// - see
// - since
// - todo
// - version
// - yield
// - yieldparam
// - yieldreturn
	return nil
}
