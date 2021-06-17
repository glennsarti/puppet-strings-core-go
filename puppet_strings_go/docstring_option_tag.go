package puppet_strings_go

import (
	"bytes"
)

// https://github.com/lsegal/yard/blob/main/lib/yard/tags/option_tag.rb
type OptionsDocstringTag struct {
	DocstringTag
	Pair *DefaultDocstringTag `json:"pair,omitempty"`
}

func (tag OptionsDocstringTag) MarshalJSON() ([]byte, error) {
	// Custom marshaller due to
	// https://github.com/puppetlabs/puppet-strings/blob/8a4a87e1b44e7581855d2c59ab4402abf438f9d8/lib/puppet-strings/yard/util.rb#L35
	buf := bytes.NewBufferString("{")

	if err := appendMarshalledKeyAndString(buf, "tag_name", tag.TagName, false); err != nil { return nil, err }
	if err := appendMarshalledKeyAndString(buf, "opt_name", tag.Pair.Name, true); err != nil { return nil, err }
	if err := appendMarshalledKeyAndString(buf, "opt_text", tag.Pair.Text, true); err != nil { return nil, err }
	if err := appendOptionalMarshalledKeyAndStringArray(buf, "opt_types", tag.Pair.Types, true); err != nil { return nil, err }
	if err := appendMarshalledKeyAndString(buf, "parent", tag.Name, true); err != nil { return nil, err }

	if err := appendOptionalMarshalledKeyAndString(buf, "text", tag.Text, true); err != nil { return nil, err }
	if err := appendOptionalMarshalledKeyAndStringArray(buf, "types", tag.Types, true); err != nil { return nil, err }
	if err := appendOptionalMarshalledKeyAndString(buf, "name", tag.Name, true); err != nil { return nil, err }

	buf.WriteString("}")
	return buf.Bytes(), nil
}
