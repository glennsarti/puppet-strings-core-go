package ruby

type (
	// PuppetFunction blah
	PuppetFunction struct {
		Name string `json:"name"`
		File string `json:"file"`
		Line int    `json:"line"`
	}
)
