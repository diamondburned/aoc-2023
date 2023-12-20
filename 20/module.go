package main

import (
	"fmt"
	"slices"
	"strings"
)

// Signal is a signal sent between modules.
type Signal bool

const (
	Lo Signal = false
	Hi Signal = true
)

func (s Signal) String() string {
	if s == Lo {
		return "lo"
	}
	return "hi"
}

// ModuleID is a unique identifier for a module.
type ModuleID string

// Module represents a communication module.
type Module interface {
	fmt.Stringer
	// ID returns the module's identifier.
	ID() ModuleID
	// Sinks returns the module's sinks.
	Sinks() []ModuleID
	// Send sends a signal to the module.
	Send(source ModuleID, signal Signal)
}

var (
	_ Module = (*FlipFlop)(nil)
	_ Module = (*Conjunction)(nil)
	_ Module = (*Broadcaster)(nil)
	_ Module = (*DummyModule)(nil)
)

type baseModule struct {
	id     ModuleID
	sinks  []ModuleID
	system *ModuleSystem
}

func (b *baseModule) ID() ModuleID      { return b.id }
func (b *baseModule) Sinks() []ModuleID { return b.sinks }

func (b *baseModule) send(s Signal) {
	for _, sink := range b.sinks {
		b.system.Send(b.id, sink, s)
	}
}

// Modules contains all modules in a module system.
type Modules map[ModuleID]Module

// AsGraphviz returns a Graphviz representation of the module system.
func (m Modules) AsGraphviz() string {
	var b strings.Builder
	b.WriteString("digraph {\n")
	b.WriteString("  button\n")

	moduleName := func(m Module) string {
		str := m.String()
		str, _, _ = strings.Cut(str, " ")
		str = strings.Trim(str, ":")
		return str
	}

	for _, module := range m {
		for _, sink := range module.Sinks() {
			fmt.Fprintf(&b,
				"  %q -> %q;\n",
				moduleName(module),
				moduleName(m[sink]),
			)
		}
	}

	b.WriteString("}")
	return b.String()
}

// FindModulesWithSink returns all modules that have a given module as a sink.
func (m Modules) FindModulesWithSink(sink ModuleID) []ModuleID {
	var modules []ModuleID
	for _, module := range m {
		for _, s := range module.Sinks() {
			if s == sink {
				modules = append(modules, module.ID())
			}
		}
	}
	slices.Sort(modules)
	return modules
}
