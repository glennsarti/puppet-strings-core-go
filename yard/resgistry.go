package yard

type (
	// Registry blah
	Registry struct {
		All []interface{}
	}
)

// NewRegistry Creates a new empty registry
func NewRegistry() *Registry {
	return &Registry{}
}
