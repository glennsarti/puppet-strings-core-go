package yard

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type (
	// DocString blah
	DocString struct {
		Text string        `json:"text,omitempty"`
		Tags []interface{} `json:"tags,omitempty"`
	}
)

// CreateDocString Creates a DocString struct from a string
func CreateDocString(content string) *DocString {
	output := DocString{}
	output.Parse(content)
	log.Println("1!!!!!!!!!!!!! ", output.Text)

	return &output
}

// Parse stuff
func (d *DocString) Parse(content string) {
	tagRegex := regexp.MustCompile(`(?i)^@(!)?((?:\w\.?)+)(?:\s+(.*))?$`)
	indentRegEx := regexp.MustCompile(`^\s*`)
	emptyRegEx := regexp.MustCompile(`^\s*$`)
	restrictedIndentRegEx := regexp.MustCompile(`^[ \t]*$`)

	// Pre parse

	// Convert the string into an array of lines
	lines := strings.Split(content, "\n")

	// Init state machine
	directive := false // For tags like '!macro'
	tagName := ""
	tagBuf := []string{}
	origIndent := 0
	text := ""
	lastIndent := 0
	lastLine := ""

	if directive {
		// do nothing
	}
	// if len(tagBuf) > 0 {
	// 	// do nothing
	// }
	// if lastIndent > 0 {
	// 	// do nothing
	// }
	// if len(lastLine) > 0 {
	// 	// do nothing
	// }

	for li, line := range lines {

		log.Println("PARSING LINE: ", line)

		metaMatch := tagRegex.FindAllStringSubmatch(line, -1)
		lineIsATag := len(metaMatch) > 0
		indent := len(indentRegEx.FindString(line))
		empty := emptyRegEx.MatchString(line)
		done := li >= (len(lines) - 1)
		inATag := len(tagName) > 0

		// if tag_name && (((indent < orig_indent && !empty) || done ||
		// 		(indent == 0 && !empty)) || (indent <= last_indent && line =~ META_MATCH))
		isOutdented := indent < origIndent && !empty
		isOutdentedTag := indent <= lastIndent && lineIsATag
		leftAlignedText := indent == 0 && !empty
		if inATag && (isOutdented || done || leftAlignedText || isOutdentedTag) {
			// 	buf = tag_buf.join("\n")
			// 	if directive || tag_is_directive?(tag_name)
			if directive {
				// We don't support directives so ignore it
				// 		directive = create_directive(tag_name, buf)
				// 		if directive
				// 			docstring << parse_content(directive.expanded_text).chomp
				// 		end
				// 	else
			} else {
				// create_tag(tag_name, buf)
				d.createTag(tagName, strings.Join(tagBuf, "\n"))
			}
			// 	end

			// 	tag_name = nil
			tagName = ""
			inATag = false
			// 	tag_buf = []
			tagBuf = []string{}
			// 	directive = false
			directive = false
			// 	orig_indent = 0
			origIndent = 0
		}
		// end

		// # Found a meta tag
		if lineIsATag {
			// We only allow a single match
			directive = metaMatch[0][1] == "!"
			tagName = metaMatch[0][2]
			tagBuf = []string{metaMatch[0][3]}
		} else if inATag && indent >= origIndent && !empty {
			// elsif tag_name && indent >= orig_indent && !empty
			// 	orig_indent = indent if orig_indent == 0
			if origIndent == 0 {
				origIndent = indent
			}
			// 	# Extra data added to the tag on the next line
			// 	last_empty = last_line =~ /^[ \t]*$/ ? true : false
			lastEmpty := restrictedIndentRegEx.MatchString(lastLine)
			// 	tag_buf << '' if last_empty
			if lastEmpty {
				tagBuf = append(tagBuf, "")
			}
			// 	tag_buf << line.gsub(/^[ \t]{#{orig_indent}}/, '')
			undentedLine := regexp.MustCompile(fmt.Sprintf("^[ \\t]%d}", origIndent)).ReplaceAllString(line, "")
			tagBuf = append(tagBuf, undentedLine)

		} else if !inATag {
			// elsif !tag_name
			// 	# Regular docstring text
			// 	docstring << line
			// 	docstring << "\n"
			text += line + "\n"
		}
		// end

		// last_indent = indent
		lastIndent = indent
		// last_line = line
		lastLine = line
	}

	// Post

	// Setup the struct
	d.Text = strings.TrimSuffix(text, "\n")

}

func (d *DocString) createTag(tagName string, buf string) {
	// Yard Tag List
	// https://rubydoc.info/gems/yard/file/docs/Tags.md#taglist
	// Ruby code - Tag Library
	// https://github.com/lsegal/yard/blob/main/lib/yard/tags/library.rb

	switch tagName {
	case "param":
		d.Tags = append(d.Tags, NewParamTag(buf))
	}
	log.Println("Create tag ", tagName)
}
