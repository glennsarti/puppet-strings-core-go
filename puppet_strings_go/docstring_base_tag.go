package puppet_strings_go

import (
	"bytes"
	"encoding/json"
)

func appendMarshalledObject(buf *bytes.Buffer, t interface{}, commaPrefix bool) (err error) {
	if commaPrefix { buf.WriteString(",") }
	if s, err := json.Marshal(t); err != nil {
		return err
	} else {
		buf.Write(s)
	}
	return
}

func appendMarshalledKeyAndString(buf *bytes.Buffer, key string, value string, commaPrefix bool) (err error) {
	if commaPrefix { buf.WriteString(",") }
	buf.WriteString("\"" + key + "\":")
	err = appendMarshalledObject(buf, value, false)
	return
}

func appendOptionalMarshalledKeyAndString(buf *bytes.Buffer, key string, value string, commaPrefix bool) (err error) {
	if value == "" { return }
	return appendMarshalledKeyAndString(buf, key, value, commaPrefix)
}

func appendMarshalledKeyAndStringArray(buf *bytes.Buffer, key string, values []string, commaPrefix bool) (err error) {
	if commaPrefix { buf.WriteString(",") }
	buf.WriteString("\"" + key + "\": [")
	for i, v := range values {
		if e := appendMarshalledObject(buf, v, i > 0); e != nil { return err }
	}
	buf.WriteString("]")
	return
}

func appendOptionalMarshalledKeyAndStringArray(buf *bytes.Buffer, key string, values []string, commaPrefix bool) (err error) {
	if values == nil || len(values) == 0 { return }
	return appendMarshalledKeyAndStringArray(buf, key, values, commaPrefix)
}

// https://github.com/lsegal/yard/blob/main/lib/yard/tags/tag.rb
type DocstringTag struct {
	TagName string `json:"tag_name,omitempty"`
	Text string `json:"text,omitempty"`
	Types []string `json:"types,omitempty"`
	Name string `json:"name,omitempty"`
}

func (ds Docstring) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("{")
	if err := appendMarshalledKeyAndString(buf, "text", ds.Text, false); err != nil { return nil, err }

	commaPrefix := false
	buf.WriteString(",\"tags\":[")
	// Append the base tags
	for _, t := range ds.Tags {
		if err := appendMarshalledObject(buf, t, commaPrefix); err != nil { return nil, err}
		commaPrefix = true
	}
	// Append the optional tags
	for _, t := range ds.OptionTags {
		if err := appendMarshalledObject(buf, t, commaPrefix); err != nil { return nil, err}
		commaPrefix = true
	}
	buf.WriteString("]")

	buf.WriteString("}")
	return buf.Bytes(), nil
}
