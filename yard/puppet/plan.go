package puppet

import (
	"github.com/glennsarti/puppet-strings-core-go/yard"
)

type (
	// PuppetPlan Puppet Strings YARD data_type_alias
	PuppetPlan struct {
		Name      string          `json:"name"`
		File      string          `json:"file"`
		Line      int             `json:"line"`
		Source    string          `json:"source"`
		DocString *yard.DocString `json:"docstring,omitempty"`
	}
)
