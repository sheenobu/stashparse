package stashparse

import (
	"bytes"
	"testing"
)

var OneInput = `
input {
	tcp {
		port => 8080
	}
}
`

var TwoInput = `
input {
	tcp {
		port => 8080
	}
	stdin {

	}
}
`

var TwoInputWithBranch = `
input {
	tcp {
		port => 8080
		t => [ "X" ]
	}
	if true {
		stdin {

		}
	}
}
`

var InputOutput = `
input {
	tcp {
		port => 8080
	}
}

output {
	file {
		path => /tmp/out.txt
	}
}
`
var IfElse = `
filter {
	if expr {
		tcp {
			port => 8080
		}
	} else {
		elasticsearch {

		}
	}
}

`

func assertNotNill(t *testing.T, val interface{}, name string) bool {
	if val == nil {
		t.Errorf("Expected %s to not be nil", name)
		return false
	}
	return true
}

func assertNill(t *testing.T, val interface{}, name string) bool {
	if val != nil {
		t.Errorf("Expected %s to be nil", name)
		return false
	}
	return true
}

func assertError(t *testing.T, err error, name string) bool {
	if err != nil {
		t.Errorf("Expected error for %s to be empty, is %s", err, name)
		return false
	}
	return true
}

func TestIfElse(t *testing.T) {
	var config Config
	r := bytes.NewReader([]byte(IfElse))
	decoder := NewDecoder(r)
	err := decoder.Decode(&config)

	assertError(t, err, "decoder.Decode")
	assertNill(t, config.Input, "config.Input")
	assertNill(t, config.Output, "config.Output")

	if assertNotNill(t, config.Filter, "config.Filter") {
		if typ := config.Filter.Type(); typ != OP_SECTION {
			t.Errorf("Expected config.Filter.Type() to be OP_SECTION, is %d", typ)
		}

		if len(config.Filter.Children()) != 2 {
			t.Errorf("Expected length of config.Filter.Children to be not 2, is not")
		} else {
			ifb := config.Filter.Children()[0].(*Branch)

			if typ := ifb.Type(); typ != OP_BRANCH {
				t.Errorf("Expected config.Filter.Children[0].Type() to be OP_BRANCH, is %d", typ)
			}

			tcp := ifb.Children()[0].(*Plugin)

			if typ := tcp.Type(); typ != OP_PLUGIN {
				t.Errorf("Expected config.Filter.Children[0].Children[0].Type() to be OP_PLUGIN, is %d", typ)
			}

			if name := tcp.Name; name != "tcp" {
				t.Errorf("Expected config.Filter.Children[0].Children[0].Name to be tcp, is %d", name)
			}

			if port, ok := tcp.Config["port"]; !ok || port != "8080" {
				t.Errorf("Expected port to be 8080, is %s", port)
			}

			elseb := config.Filter.Children()[1].(*Branch)

			if typ := elseb.Type(); typ != OP_BRANCH {
				t.Errorf("Expected config.Filter.Children[1].Type() to be OP_BRANCH, is %d", typ)
			}

			ela := elseb.Children()[0].(*Plugin)

			if typ := ela.Type(); typ != OP_PLUGIN {
				t.Errorf("Expected config.Filter.Children[1].Children[0].Type() to be OP_PLUGIN, is %d", typ)
			}

			if name := ela.Name; name != "elasticsearch" {
				t.Errorf("Expected config.Filter.Children[1].Children[0].Name to be elasticsearch, is %d", name)
			}

		}
	}
}

func TestParseOneInputOnly(t *testing.T) {
	var config Config
	r := bytes.NewReader([]byte(OneInput))
	decoder := NewDecoder(r)
	err := decoder.Decode(&config)

	assertError(t, err, "decoder.Decode")

	if assertNotNill(t, config.Input, "config.Input") {

		if typ := config.Input.Type(); typ != OP_SECTION {
			t.Errorf("Expected config.Input.Type() to be OP_SECTION, is %d", typ)
		}

		if len(config.Input.Children()) != 1 {
			t.Errorf("Expected length of config.Input.Children to be not 1, is not")
		} else {

			tcp := config.Input.Children()[0].(*Plugin)

			if typ := tcp.Type(); typ != OP_PLUGIN {
				t.Errorf("Expected config.Input.Type() to be OP_PLUGIN, is %d", typ)
			}

			if port, ok := tcp.Config["port"]; !ok || port != "8080" {
				t.Errorf("Expected port to be 8080, is %s", port)
			}
		}
	}
}

func TestParseTwoInputs(t *testing.T) {
	var config Config
	r := bytes.NewReader([]byte(TwoInput))
	decoder := NewDecoder(r)
	err := decoder.Decode(&config)

	if err != nil {
		t.Errorf("Expected decoding to succeed, got %s", err)
	}

	if config.Input == nil {
		t.Errorf("Expected config.Input to be not nil, is nil")
	} else {
		if typ := config.Input.Type(); typ != OP_SECTION {
			t.Errorf("Expected config.Input.Type() to be OP_SECTION, is %d", typ)
		}

		if i := len(config.Input.Children()); i != 2 {
			t.Errorf("Expected length of config.Input.Children to be 2, is %d", i)
		} else {

			tcp := config.Input.Children()[0].(*Plugin)

			if tcp.Name != "tcp" {
				t.Errorf("Expected tcp as first element")
			}

			if typ := tcp.Type(); typ != OP_PLUGIN {
				t.Errorf("Expected config.Input.Type() to be OP_PLUGIN, is %d", typ)
			}

			if port, ok := tcp.Config["port"]; !ok || port != "8080" {
				t.Errorf("Expected port to be 8080, is %s", port)
			}

			stdin := config.Input.Children()[1].(*Plugin)

			if stdin.Name != "stdin" {
				t.Errorf("Expected stdin as second element")
			}

			if typ := stdin.Type(); typ != OP_PLUGIN {
				t.Errorf("Expected config.Input.Type() to be OP_PLUGIN, is %d", typ)
			}

			if port, ok := tcp.Config["port"]; !ok || port != "8080" {
				t.Errorf("Expected port to be 8080, is %s", port)
			}
		}
	}
}

func TestParseTwoInputsWithBranch(t *testing.T) {
	var config Config
	r := bytes.NewReader([]byte(TwoInputWithBranch))
	decoder := NewDecoder(r)
	err := decoder.Decode(&config)

	if err != nil {
		t.Errorf("Expected decoding to succeed, got %s", err)
	}

	if config.Input == nil {
		t.Errorf("Expected config.Input to be not nil, is nil")
	} else {
		if typ := config.Input.Type(); typ != OP_SECTION {
			t.Errorf("Expected config.Input.Type() to be OP_SECTION, is %d", typ)
		}

		if i := len(config.Input.Children()); i != 2 {
			t.Errorf("Expected length of config.Input.Children to be 2, is %d", i)
		} else {

			tcp := config.Input.Children()[0].(*Plugin)

			if tcp.Name != "tcp" {
				t.Errorf("Expected tcp as first element")
			}

			if typ := tcp.Type(); typ != OP_PLUGIN {
				t.Errorf("Expected config.Input.Type() to be OP_PLUGIN, is %d", typ)
			}

			if port, ok := tcp.Config["port"]; !ok || port != "8080" {
				t.Errorf("Expected port to be 8080, is %s", port)
			}

			ifbranch := config.Input.Children()[1].(*Branch)

			if ifbranch.Name != "if" {
				t.Errorf("Expected if as second element, got %s", ifbranch.Name)
			}

			if typ := ifbranch.Type(); typ != OP_BRANCH {
				t.Errorf("Expected config.Input.Type() to be OP_BRANCH, is %d", typ)
			}

			if expr := ifbranch.Expression(); expr != "true" {
				t.Errorf("Expected expression to be true, is %s", expr)
			}
		}
	}
}

func TestInputOutput(t *testing.T) {
	var config Config
	r := bytes.NewReader([]byte(InputOutput))
	decoder := NewDecoder(r)
	err := decoder.Decode(&config)

	if err != nil {
		t.Errorf("Expected decoding to succeed, got %s", err)
	}

	if config.Input == nil {
		t.Errorf("Expected config.Input to be not nil, is nil")
	} else {
		if typ := config.Input.Type(); typ != OP_SECTION {
			t.Errorf("Expected config.Input.Type() to be OP_SECTION, is %d", typ)
		}

		if i := len(config.Input.Children()); i != 1 {
			t.Errorf("Expected length of config.Input.Children to be 1, is %d", i)
		} else {

			tcp := config.Input.Children()[0].(*Plugin)

			if tcp.Name != "tcp" {
				t.Errorf("Expected tcp as first element")
			}

			if typ := tcp.Type(); typ != OP_PLUGIN {
				t.Errorf("Expected config.Input.Type() to be OP_PLUGIN, is %d", typ)
			}

			if port, ok := tcp.Config["port"]; !ok || port != "8080" {
				t.Errorf("Expected port to be 8080, is %s", port)
			}
		}
	}

	if config.Output == nil {
		t.Errorf("Expected config.Output to be not nil, is nil")
	} else {
		if typ := config.Output.Type(); typ != OP_SECTION {
			t.Errorf("Expected config.Output.Type() to be OP_SECTION, is %d", typ)
		}

		if i := len(config.Output.Children()); i != 1 {
			t.Errorf("Expected length of config.Output.Children to be 1, is %d", i)
		} else {

			file := config.Output.Children()[0].(*Plugin)

			if file.Name != "file" {
				t.Errorf("Expected tcp as first element")
			}

			if typ := file.Type(); typ != OP_PLUGIN {
				t.Errorf("Expected config.Output.Type() to be OP_PLUGIN, is %d", typ)
			}

			if path, ok := file.Config["path"]; !ok || path != "/tmp/out.txt" {
				t.Errorf("Expected path to be /tmp/out.txt, is %s", path)
			}
		}
	}
}
