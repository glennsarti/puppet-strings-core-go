package yard

type (
	// ParamTag blah
	ParamTag struct {
		TagName string   `json:"tag_name"`
		Name    string   `json:"name,omitempty"`
		Text    string   `json:"text,omitempty"`
		Types   []string `json:"types,omitempty"`
	}
)

// NewParamTag blah
func NewParamTag(content string) *ParamTag {
	tag := ParamTag{
		TagName: "param",
	}
	tag.fromString(content)
	return &tag
}

func (p *ParamTag) fromString(content string) {
	// TODO: Do something
	// https://github.com/lsegal/yard/blob/main/lib/yard/tags/default_factory.rb#L39
}
