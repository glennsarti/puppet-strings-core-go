package puppet_strings_go

// https://github.com/lsegal/yard/blob/main/lib/yard/tags/default_tag.rb
type DefaultDocstringTag struct {
	DocstringTag
	Defaults []string `json:"defaults,omitempty"`
}
