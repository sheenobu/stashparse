package stashparse

// Plugin defines a plugin used in the configuration file
type Plugin struct {
	Name     string
	children []Operation
	Config   map[string]string
	parent   Operation
}

// NewPlugin creates a new plugin operation
func NewPlugin(parent Operation, name string) *Plugin {
	return &Plugin{
		name,
		make([]Operation, 0),
		make(map[string]string),
		parent,
	}
}

// Operation implementation

// Parent returns the parent operation
func (p *Plugin) Parent() Operation {
	return p.parent
}

// Add adds a child operation to the branch
func (p *Plugin) Add(op Operation) {
	p.children = append(p.children, op)
}

// Children returns the list of operations
func (p *Plugin) Children() []Operation {
	return p.children
}

// Type returns the operation type, OP_PLUGIN
func (p *Plugin) Type() OperationType {
	return OP_PLUGIN
}

// --

// Set adds a new configuration value to the plugin
func (p *Plugin) Set(name string, val string) {
	p.Config[name] = val
}
