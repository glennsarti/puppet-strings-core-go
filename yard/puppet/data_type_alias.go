package puppet

type (
	// DataTypeAlias Puppet Strings YARD data_type_alias
	DataTypeAlias struct {
		Name    string `json:"name"`
		File    string `json:"file"`
		Line    int    `json:"line"`
		AliasOf string `json:"alias_of"`
	}
)
