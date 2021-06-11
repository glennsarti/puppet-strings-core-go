package puppet_strings_go

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	typelistOpeningChars = "[({<"
	typelistClosingChars = ">})]"
)

func (ds *Docstring) parseTagWithTitleAndText(tagName string, lines []string) (tag *DocstringTag, err error) {
	title, desc, err := ds.extractTitleAndDescFromLines(lines)
	if (err != nil) { return nil, err }
	return &DocstringTag{
		TagName: tagName,
		Text: desc,
		Name: title,
	}, nil
}

func (ds *Docstring) parseTagWithTypes(tagName string, lines []string) (tag *DocstringTag, err error) {
	name, types, text, err := ds.extractTypesAndNameFromText(lines, typelistOpeningChars, typelistClosingChars)
	if err != nil { return nil, err }
	if name != "" { return nil, errors.New(fmt.Sprintf("Cannot specify a name before type list for '@%s'", tagName))}

	return &DocstringTag{
		TagName: tagName,
		Text: text,
		Types: types,
	}, nil
}

func (ds *Docstring) extractTitleAndDescFromLines(lines []string) (title string, desc string, err error) {
	if len(lines) == 0 { return "","", errors.New("Missing text for a tag") }
	title = ""
	desc = ""

	if len(lines) == 1 { return strings.TrimSpace(lines[0]), desc, nil }

	if regexp.MustCompile(`\A[ \t]\z`).MatchString(lines[0]) {
		return title, strings.Join(lines[1:],"\n"), nil
	} else {

		title = strings.TrimSpace(lines[0])
		desc = strings.Join(lines[1:],"\n")
		// Strip any double, or more, spaces
		desc = regexp.MustCompile(`[ ]{2,}`).ReplaceAllString(desc, " ")
		desc = strings.TrimSpace(desc)
		return
	}
}

func (ds *Docstring) extractTypesAndNameFromText(lines []string, openingTypes string, closingTypes string) (before string, types []string, text string, err error) {
	before, list, text, err := ds.extractTypesAndNameFromTextUnStripped(lines, openingTypes, closingTypes)
	if err != nil { return before, list, text, err }

	for i, e := range list {
		list[i] = strings.TrimSpace(e)
	}

	return strings.TrimSpace(before), list, strings.TrimSpace(text), err
}


func (ds *Docstring) extractTypesAndNameFromTextUnStripped(lines []string, openingTypes string, closingTypes string) (before string, types []string, text string, err error) {
	return "", make([]string, 0), "", nil
}
