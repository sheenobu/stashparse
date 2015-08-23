package stashparse

// OperationType represents the type of operation
type OperationType string

const (
	OP_SECTION OperationType = "section"
	OP_PLUGIN  OperationType = "plugin"
	OP_BRANCH  OperationType = "branch"
)

// Operation represents
type Operation interface {

	// Type returns the OperationType of the Operation
	Type() OperationType

	// Children returns the children of the Operation
	Children() []Operation

	// Add adds a new child to the Operation
	Add(Operation)

	// Parent returns the parent of the Operation
	Parent() Operation
}

// Config represents the config file
type Config struct {

	// Input represents the Input section
	Input Operation

	// Filter represeents the Filter section
	Filter Operation

	// Output represents the Output section
	Output Operation
}
