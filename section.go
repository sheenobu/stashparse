package stashparse

// Section represents a configuration file section
type Section struct {
	Name     string
	children []Operation
	parent   Operation
}

// NewSection creates a new section option
func NewSection(parent Operation, name string) *Section {
	return &Section{
		name,
		make([]Operation, 0),
		parent,
	}
}

// Operation implementation

// Parent returns the parent operation
func (s *Section) Parent() Operation {
	return s.parent
}

// Children returns the children of the section
func (s *Section) Children() []Operation {
	return s.children
}

// Type returns the type of operation, OP_SECTION
func (s *Section) Type() OperationType {
	return OP_SECTION
}

// Add adds a child operation to the section
func (s *Section) Add(op Operation) {
	s.children = append(s.children, op)
}
