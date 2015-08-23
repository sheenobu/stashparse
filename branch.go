package stashparse

// Branch represents an if, else branch
type Branch struct {
	Name       string
	children   []Operation
	expression string
	parent     Operation
}

// NewBranch constructs a branch
func NewBranch(parent Operation, name string, expression string) *Branch {
	return &Branch{
		name,
		make([]Operation, 0),
		expression,
		parent,
	}
}

// Operation implementation

// Parent returns the parent operation
func (b *Branch) Parent() Operation {
	return b.parent
}

// Add adds a child operation to the branch
func (b *Branch) Add(op Operation) {
	b.children = append(b.children, op)
}

// Children returns the list of operations
func (b *Branch) Children() []Operation {
	return b.children
}

// Type returns the operation type, OP_BRANCH
func (b *Branch) Type() OperationType {
	return OP_BRANCH
}

// --

// Expression returns the if/else expression to be evaluated
func (b *Branch) Expression() string {
	return b.expression
}
