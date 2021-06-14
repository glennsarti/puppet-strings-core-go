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
	return[]rune{0, ' ', '\t', '\n'}
}

const (
	methodNameMatch = `[a-zA-Z_]\w*[!?=]?|[-+~]\@|<<|>>|=~|===?|![=~]?|<=>|[<>]=?|\*\*|[-\/+%^&*~` + "`" + `|]|\[\]=?`
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
	name, types, text, err := ds.extractTypesAndNameFromText(lines, ds.typelistOpeningChars(), ds.typelistClosingChars())
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

func (ds *Docstring) extractTypesAndNameFromText(lines []string, openingTypes []rune, closingTypes []rune) (before string, types []string, text string, err error) {
	before, list, text, err := ds.extractTypesAndNameFromTextUnStripped(lines, openingTypes, closingTypes)
	if err != nil { return before, list, text, err }

	for i, e := range list {
		list[i] = strings.TrimSpace(e)
	}

	return strings.TrimSpace(before), list, strings.TrimSpace(text), err
}


func (ds *Docstring) consumeWhiteSpace(sr StringReader) {
	for {
		c, _ := sr.Peek()
		switch c {
		case 0:
			return
		case ' ','\t':
			break
		default:
			return
		}
		sr.Next()
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
		ds.consumeWhiteSpace(sr)
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

func (ds *Docstring) extractTypesAndNameFromTextUnStripped(lines []string, openingTypes []rune, closingTypes []rune) (before string, types []string, text string, err error) {
	if (len(lines) == 0) { return "", make([]string, 0), "", nil }
	sr := NewStringReader(lines[0])

	foundTypes := false
	before = ""
	after := ""

	for {
		ds.consumeWhiteSpace(sr)
		c, _ := sr.Next()
		if c == 0 { break }

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
		after = lines[0]
		before = ""
	} else {
		after = after + sr.PeekUntilEnd()
	}

	if len(lines) > 1 {
		after = "\n" + strings.Join(lines[1:], "\n")
	}

	return before, types, after, nil
}

func includes(arr []rune, val rune) bool {
	for _, item := range arr {
		if item == val { return true }
	}
	return false
}
