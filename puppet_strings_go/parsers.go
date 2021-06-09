package puppet_strings_go

type StringsParser interface {
	Parse(content []byte)
	Result() *AllStringsObjects
}
