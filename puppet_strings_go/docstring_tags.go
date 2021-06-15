package puppet_strings_go

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func (ds *Docstring) typelistOpeningChars() []rune {
	return []rune{'"', '[', '(', '{', '<'}
}

func (ds *Docstring) typelistClosingChars() []rune {
	return[]rune{'>', '}', ')', ']', '"'}
}

func (ds *Docstring) whiteSpaceRunes() []rune {
	return[]rune{0, ' ', '\t'}
}

func (ds *Docstring) emptyRunes() []rune {
	return[]rune{}
}

const (
	methodNameMatch = `[a-zA-Z_]\w*[!?=]?|[-+~]\@|<<|>>|=~|===?|![=~]?|<=>|[<>]=?|\*\*|[-\/+%^&*~` + "`" + `|]|\[\]=?`
)

func (ds *Docstring) parseTagWithTitleAndText(tagName string, text string) (tag *DocstringTag, err error) {
	title, desc, err := ds.extractTitleAndDescFromText(text)
	if (err != nil) { return nil, err }
	return &DocstringTag{
		TagName: tagName,
		Text: desc,
		Name: title,
	}, nil
}

func (ds *Docstring) parseTagWithTypes(tagName string, text string) (tag *DocstringTag, err error) {
	name, types, text, err := ds.extractTypesAndNameFromText(text, ds.typelistOpeningChars(), ds.typelistClosingChars())
	if err != nil { return nil, err }
	if name != "" { return nil, errors.New(fmt.Sprintf("Cannot specify a name before type list for '@%s'", tagName))}

	return &DocstringTag{
		TagName: tagName,
		Text: text,
		Types: types,
	}, nil
}

func (ds *Docstring) parseTagWithOptions(tagName string, text string) (tag *DocstringTag, err error) {
	// name, remainText, err := ds.extractNameFromText(text)
	// if err != nil { return nil, err }

      // def parse_tag_with_options(tag_name, text)
      //   name, text = *extract_name_from_text(text)
      //   OptionTag.new(tag_name, name, parse_tag_with_types_name_and_default(tag_name, text))
      // end



	// name, types, text, err := ds.extractTypesAndNameFromText(text, ds.typelistOpeningChars(), ds.typelistClosingChars())
	// if err != nil { return nil, err }
	// if name != "" { return nil, errors.New(fmt.Sprintf("Cannot specify a name before type list for '@%s'", tagName))}

	return &DocstringTag{
		TagName: tagName,
		Text: text,
	}, nil
}

func (ds *Docstring) parseTagWithTypesAndName(tagName string, text string) (tag *DocstringTag, err error) {
	name, types, text, err := ds.extractTypesAndNameFromText(text, ds.typelistOpeningChars(), ds.typelistClosingChars())
	if err != nil { return nil, err }
	return &DocstringTag{
		TagName: tagName,
		Text: text,
		Types: types,
		Name: name,
	}, nil
}

func (ds *Docstring) extractNameFromText(text string) (name string, remainText string, err error) {
	sr := NewStringReader(text)

	ds.consumeWhiteSpace(sr, false)
	name = ds.consumeUntilWhiteSpaceOrRune(sr, ds.emptyRunes())

	return name, sr.UntilEnd(), nil
}

func (ds *Docstring) extractTitleAndDescFromText(text string) (title string, desc string, err error) {
	if len(text) == 0 { return "","", errors.New("Missing text for a tag") }
	title = ""
	desc = ""

	// TODO: This could be done better using a string reader
	lines := strings.SplitN(text, "\n", 2)
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

func (ds *Docstring) extractTypesAndNameFromText(text string, openingTypes []rune, closingTypes []rune) (before string, types []string, remainText string, err error) {
	before, list, remainText, err := ds.extractTypesAndNameFromTextUnStripped(text, openingTypes, closingTypes)
	if err != nil { return before, list, text, err }

	for i, e := range list {
		list[i] = strings.TrimSpace(e)
	}

	return strings.TrimSpace(before), list, strings.TrimSpace(remainText), err
}

func (ds *Docstring) consumeWhiteSpace(sr StringReader, consumeLF bool) {
	for {
		c, _ := sr.Peek()

		switch {
		case c == 0:
			return
		case c == '\n':
			if consumeLF { sr.Next() } else { return }
		case includes(ds.whiteSpaceRunes(),c):
			sr.Next()
		default:
			return
		}
	}
}

func (ds *Docstring) consumeUntilWhiteSpaceOrRune(sr StringReader, runes []rune) (string) {
	startPos := sr.Pos()
	for {
		c, endPos := sr.Peek()
		if includes(ds.whiteSpaceRunes(), c) || includes(runes, c) {
			return sr.SubString(startPos, endPos)
		}
		sr.Next()
	}
}

func (ds *Docstring) consumeUntilRune(sr StringReader, r rune) (string) {
	startPos := sr.Pos()
	for {
		c, endPos := sr.Peek()
		switch c {
		case 0,'\n':
			return sr.SubString(startPos, endPos)
		case r:
			_, endPos = sr.Next()
			return sr.SubString(startPos, endPos)
		}
		sr.Next()
	}
}

func (ds *Docstring) consumeTypes(sr StringReader, openingTypes []rune, closingTypes []rune) ([]string) {
	depth := 0
	list := make([]string, 0)
	startPos := sr.Pos()

	// TODO: Can we compile this only once?
	mnmRegex := regexp.MustCompile(methodNameMatch)

	for {
		ds.consumeWhiteSpace(sr, false)
		c, _ := sr.Next()

		switch {
		case c == 0:
			return list
		case c == ',' && depth == 0:
			list = append(list, sr.SubString(startPos, sr.Pos() - 1))
			startPos = sr.Pos()
		case c == '\'' || c == '"':
			// YARD doesn't do any interpolation so it's plain literal strings
			ds.consumeUntilRune(sr, c)
		case includes(openingTypes, c):
			depth += 1
		case includes(closingTypes, c):
			if depth == 0 {
				if startPos != sr.Pos() -1 {
					list = append(list, sr.SubString(startPos, sr.Pos() - 1))
				}
				return list
			} else {
				depth -= 1
			}
		case c == '=':
			n, _ := sr.Peek()
			// Hash rockets trip up the closing '>' tag so skip by
			if n == '>' {
				sr.Advance(2)
			}
		case c == '#':
			if m := mnmRegex.FindStringSubmatch(sr.PeekUntilEnd()); m != nil {
				// TODO: Should really advance in byte length not number of characters
				// but the regex is really only looking for single byte chars its mostly safe
				sr.Advance(len(m[0]))
			}
		}
	}
}

func (ds *Docstring) extractTypesAndNameFromTextUnStripped(text string, openingTypes []rune, closingTypes []rune) (before string, types []string, remainText string, err error) {
	if (len(text) == 0) { return "", make([]string, 0), "", nil }
	sr := NewStringReader(text)

	foundTypes := false
	before = ""
	after := ""

	for {
		ds.consumeWhiteSpace(sr, false)
		c, _ := sr.Next()
		if c == 0 || c == '\n' { break }

		if !foundTypes && includes(openingTypes, c) {
			if t := ds.consumeTypes(sr, openingTypes, closingTypes); len(t) > 0 {
				types = t
				foundTypes = true
				continue
			}
		}
		if before == "" {
			before = string(c) + ds.consumeUntilWhiteSpaceOrRune(sr, openingTypes)
			continue
		}
		after = string(c)
		break
	}

	if !foundTypes {
		after = text
		before = ""
	} else {
		after = after + sr.UntilEnd()
	}

	return before, types, after, nil
}

func includes(arr []rune, val rune) bool {
	for _, item := range arr {
		if item == val { return true }
	}
	return false
}
