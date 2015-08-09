package stashparse

import (
	"fmt"
	"github.com/looplab/fsm"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Decoder is the entity that can decode a configuration file into the Config structure
type Decoder struct {
	reader io.Reader
}

// NewDecoder constructs a decoder using the given reader
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r,
	}
}

// Decode decodes the input reader into the config pointer
func (d *Decoder) Decode(config *Config) error {

	tokenBuffer := ""
	expressionBuffer := ""
	conditionalNesting := 0
	errChan := make(chan error)
	conditionName := ""
	var lastOp Operation = nil

	parserState := fsm.NewFSM(
		"section",
		fsm.Events{
			{Name: "input", Src: []string{"section"}, Dst: "section_start"},
			{Name: "output", Src: []string{"section"}, Dst: "section_start"},
			{Name: "filter", Src: []string{"section"}, Dst: "section_start"},

			{Name: "{", Src: []string{"section_start"}, Dst: "entity"},
			{Name: "}", Src: []string{"entity"}, Dst: "section"},

			{Name: "if", Src: []string{"entity"}, Dst: "cond_start"},
			{Name: "else", Src: []string{"entity"}, Dst: "cond_start"},
			{Name: "{", Src: []string{"cond_start"}, Dst: "entity"},

			{Name: "driver", Src: []string{"entity"}, Dst: "entity_start"},
			{Name: "{", Src: []string{"entity_start"}, Dst: "driver"},
			{Name: "}", Src: []string{"driver"}, Dst: "entity"},

			{Name: "key", Src: []string{"driver"}, Dst: "key_start"},
			{Name: "=>", Src: []string{"key_start"}, Dst: "value"},
			{Name: "val", Src: []string{"value"}, Dst: "driver"},
		},
		fsm.Callbacks{
			"leave_cond_start": func(e *fsm.Event) {
				conditionalNesting++
			},
			"leave_entity": func(e *fsm.Event) {
				if e.Dst == "section" {
					if conditionalNesting != 0 {
						conditionalNesting--
						e.Cancel()
					}
				}
				if lastOp != nil && e.Event == "}" {
					lastOp = lastOp.Parent()
				}
			},
			"leave_driver": func(e *fsm.Event) {
				if lastOp != nil && e.Event == "}" {
					lastOp = lastOp.Parent()
				}
			},
			"enter_section": func(e *fsm.Event) {
				if conditionalNesting != 0 {
					errChan <- fmt.Errorf("Empty if/else branch")
				}
			},
			"enter_cond_start": func(e *fsm.Event) {
				conditionName = e.Event
			},
		},
	)

	var buffer []byte = make([]byte, 2)

	var err error
	var n int
	var lastKey string = ""

	r := d.reader
	for n, err = r.Read(buffer); err == nil && n > 0; n, err = r.Read(buffer) {
		for i := 0; i < n; i++ {
			runeValue, width := utf8.DecodeRuneInString(string(buffer[i:]))
			if (parserState.Current() != "value" && unicode.IsSpace(runeValue)) || (parserState.Current() == "value" && (uint32(runeValue) == '\n')) {
				if tokenBuffer != "" {
					if strings.HasPrefix(parserState.Current(), "section") {
						if parserState.Can(tokenBuffer) {
							parserState.Event(tokenBuffer)
						} else {
							return fmt.Errorf("Unexpected Token %s", tokenBuffer)
						}
						if tokenBuffer != "}" && tokenBuffer != "{" {
							if tokenBuffer == "input" {
								config.Input = NewSection(lastOp, "input")
								lastOp = config.Input
							} else if tokenBuffer == "output" {
								config.Output = NewSection(lastOp, "output")
								lastOp = config.Output
							} else if tokenBuffer == "filter" {
								config.Filter = NewSection(lastOp, "filter")
								lastOp = config.Filter
							}
						}
					} else if strings.HasPrefix(parserState.Current(), "entity") {
						if parserState.Can(tokenBuffer) {
							parserState.Event(tokenBuffer)
						} else {
							if parserState.Can("driver") {
								l := NewPlugin(lastOp, tokenBuffer)
								lastOp.Add(l)
								lastOp = l
								parserState.Event("driver")
							} else {
								return fmt.Errorf("Unexpected Token %s", tokenBuffer)
							}
						}
					} else if strings.HasPrefix(parserState.Current(), "driver") {
						if parserState.Can(tokenBuffer) {
							parserState.Event(tokenBuffer)
						} else {
							lastKey = tokenBuffer
							parserState.Event("key")
						}
					} else if strings.HasPrefix(parserState.Current(), "key") {
						if parserState.Can(tokenBuffer) {
							parserState.Event(tokenBuffer)
						} else {
							return fmt.Errorf("Unexpected Token %s", tokenBuffer)
						}
					} else if strings.HasPrefix(parserState.Current(), "value") {
						lastOp.(*Plugin).Set(lastKey, strings.Trim(tokenBuffer, " "))
						parserState.Event("val")
					} else if strings.HasPrefix(parserState.Current(), "cond") {
						if parserState.Can(tokenBuffer) {

							l := NewBranch(lastOp, conditionName, expressionBuffer)
							lastOp.Add(l)
							lastOp = l

							parserState.Event(tokenBuffer)
							expressionBuffer = ""
						} else if parserState.Current() == "cond_start" {
							expressionBuffer += tokenBuffer
						}
					}

				}
				tokenBuffer = ""
			} else {
				tokenBuffer += string(buffer[i : i+width])

			}
		}
	}

	select {
	case e := <-errChan:
		return e
	default:
	}

	if err != nil && err != io.EOF {
		return err
	}

	if parserState.Current() != "section" {
		return fmt.Errorf("Error parsing config. Unclosed entity")
	}

	return nil
}
