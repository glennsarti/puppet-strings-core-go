package yard

type (
	// DataTypeAlias Puppet Strings YARD data_type_alias
	DataTypeAlias struct {
		Name    string
		File    string
		Line    int
		AliasOf string `json:"alias_of"`
	}
)
