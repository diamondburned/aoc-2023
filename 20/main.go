package main

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

// Signal is a signal sent between modules.
type Signal bool

const (
	Lo Signal = false
	Hi Signal = true
)

func (s Signal) String() string {
	if s == Lo {
		return "░"
	}
	return "▓"
}

// ModuleID is a unique identifier for a module.
type ModuleID string

type ModuleSystemMetrics struct {
	LoPulses int
	HiPulses int
}

func (m ModuleSystemMetrics) Accumulate(other ModuleSystemMetrics) ModuleSystemMetrics {
	return ModuleSystemMetrics{
		LoPulses: m.LoPulses + other.LoPulses,
		HiPulses: m.HiPulses + other.HiPulses,
	}
}

// Modules contains all modules in a module system.
type Modules map[ModuleID]Module

// AsGraphviz returns a Graphviz representation of the module system.
func (m Modules) AsGraphviz() string {
	var b strings.Builder
	b.WriteString("digraph {\n")
	b.WriteString("button\n")

	moduleName := func(m Module) string {
		if m == nil {
			return ""
		}
		str := m.String()
		str, _, _ = strings.Cut(str, " ")
		str = strings.Trim(str, ":")
		return str
	}

	for _, module := range m {
		for _, sink := range module.Sinks() {
			fmt.Fprintf(&b,
				"%q -> %q;\n",
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

// Subsystem returns a subsystem of modules that are connected to a given sink.
func (m Modules) Subsystem(sink ModuleID) Modules {
	bfs := aocutil.AcyclicBFS(sink, func(id ModuleID) []ModuleID {
		return m.FindModulesWithSink(id)
	})
	subsystem := Modules{}
	for id := range bfs {
		subsystem[id] = m[id]
	}
	return subsystem
}

// ModuleSystem is a collection of modules represented as a map from module
// identifiers to modules.
type ModuleSystem struct {
	Modules  Modules
	deferred []deferredSignal
}

// NewModuleSystem creates a new module system.
func NewModuleSystem(modules Modules) *ModuleSystem {
	return &ModuleSystem{
		Modules: modules,
	}
}

type deferredSignal struct {
	source ModuleID
	sink   ModuleID
	signal Signal
}

// Send sends a signal to a module.
func (s *ModuleSystem) Send(source, sink ModuleID, signal Signal) {
	s.deferred = append(s.deferred, deferredSignal{source, sink, signal})
}

// Tick causes all modules to process all signals sent to them.
func (s *ModuleSystem) Tick() aocutil.Iter[ModuleSystemMetrics] {
	return func(yield func(ModuleSystemMetrics) bool) {
		for len(s.deferred) > 0 {
			sending := s.deferred
			s.deferred = nil

			var metrics ModuleSystemMetrics
			for _, signal := range sending {
				// log.Printf("%s -%s-> %s", signal.source, signal.signal, signal.sink)
				module := s.Modules[signal.sink]
				module.Send(signal.source, signal.signal)

				switch signal.signal {
				case Lo:
					metrics.LoPulses++
				case Hi:
					metrics.HiPulses++
				}
			}

			if !yield(metrics) {
				return
			}
		}
	}
}

// Settled returns true if no more signals need to be sent.
func (s *ModuleSystem) Settled() bool {
	return len(s.deferred) == 0
}

// // Propagate sends and ticks until no more signals need to be sent.
// // It returns the number of signals sent.
// func (s *ModuleSystem) Propagate(source, sink ModuleID, signal Signal) ModuleSystemMetrics {
// 	s.Send(source, sink, signal)
// 	var metrics ModuleSystemMetrics
// 	for m := range s.Tick() {
// 		metrics = metrics.Accumulate(m)
// 	}
// 	return metrics
// }

// Module represents a communication module.
type Module interface {
	fmt.Stringer
	// ID returns the module's identifier.
	ID() ModuleID
	// Send sends a signal to the module.
	Send(source ModuleID, signal Signal)
	// Sinks returns the module's sinks.
	Sinks() []ModuleID
}

var (
	_ Module = (*FlipFlop)(nil)
	_ Module = (*Conjunction)(nil)
	_ Module = (*Broadcaster)(nil)
	_ Module = (*DummyModule)(nil)
)

// Broadcaster is the simplest module. It sends all signals it receives to all
// of its sinks.
type Broadcaster struct {
	id     ModuleID
	system *ModuleSystem
	sinks  []ModuleID
}

func (b *Broadcaster) ID() ModuleID {
	return b.id
}

func (b *Broadcaster) Send(_ ModuleID, s Signal) {
	b.send(s)
}

func (b *Broadcaster) send(s Signal) {
	for _, sink := range b.sinks {
		b.system.Send(b.id, sink, s)
	}
}

func (b *Broadcaster) Sinks() []ModuleID {
	return b.sinks
}

func (b *Broadcaster) String() string {
	sinks := aocutil.Map(b.sinks, func(sink ModuleID) string { return string(sink) })
	return fmt.Sprintf("broadcaster(%s) -> %s", b.id, strings.Join(sinks, ", "))
}

// FlipFlop is a flip-flop module,
type FlipFlop struct {
	Broadcaster
	state Signal
}

func (f *FlipFlop) ID() ModuleID {
	return f.id
}

func (f *FlipFlop) Send(_ ModuleID, s Signal) {
	if s == Lo {
		f.state = !f.state
		f.send(f.state)
	}
}

func (f *FlipFlop) String() string {
	sinks := aocutil.Map(f.sinks, func(sink ModuleID) string { return string(sink) })
	return fmt.Sprintf("flipflop(%s) -> %s: %s", f.id, strings.Join(sinks, ", "), f.state)
}

// Conjunction is a conjunction module.
type Conjunction struct {
	Broadcaster
	sources []ModuleID
	states  []Signal
}

func (c *Conjunction) ID() ModuleID {
	return c.id
}

func (c *Conjunction) Send(source ModuleID, s Signal) {
	c.init()

	i := slices.Index(c.sources, source)
	if i == -1 {
		log.Panicf("unknown source %q", source)
	}

	c.states[i] = s
	if c.AllHi() {
		c.send(Lo)
	} else {
		c.send(Hi)
	}
}

// AllHi returns true if all of the conjunction's sources are high.
func (c *Conjunction) AllHi() bool {
	c.init()
	for _, s := range c.states {
		if s == Lo {
			return false
		}
	}
	return true
}

func (c *Conjunction) init() {
	if c.sources == nil {
		c.sources = c.system.Modules.FindModulesWithSink(c.id)
		c.states = make([]Signal, len(c.sources))
	}
}

func (c *Conjunction) String() string {
	c.init()
	states := make([]string, len(c.states))
	for i, state := range c.states {
		states[i] = fmt.Sprintf("%s: %s", c.sources[i], state)
	}
	return fmt.Sprintf("conjunction(%s): {%s}", c.id, strings.Join(states, ""))
}

// DummyModule is a module that does nothing but store the signal it receives.
type DummyModule struct {
	id ModuleID
	s  Signal
}

func (d *DummyModule) ID() ModuleID {
	return d.id
}

func (d *DummyModule) Send(_ ModuleID, s Signal) {
	d.s = s
}

func (d *DummyModule) Sinks() []ModuleID {
	return nil
}

func (d *DummyModule) String() string {
	return fmt.Sprintf("dummy(%s): %s", d.id, d.s)
}

// Button is a module that sends a signal when pressed.
type Button struct {
	system *ModuleSystem
	sink   ModuleID
}

func (b *Button) ID() ModuleID {
	return "button"
}

func (b *Button) Send(source ModuleID, s Signal) {
	if s == Hi {
		return
	}
	b.system.Send("button", b.sink, s)
}

func (b *Button) Sinks() []ModuleID {
	return []ModuleID{b.sink}
}

func (b *Button) Push() {
	b.system.Send("button", b.sink, Lo)
}

// func (b *Button) Propagate(s Signal) ModuleSystemMetrics {
// 	return b.system.Propagate("button", b.sink, s)
// }

func (b *Button) String() string {
	return fmt.Sprintf("button -> %s", b.sink)
}

type InputModules struct {
	*ModuleSystem
	Button *Button
}

func parseInput(input string) InputModules {
	lines := aocutil.SplitLines(input)
	system := &ModuleSystem{Modules: make(map[ModuleID]Module, len(lines))}
	var first ModuleID

	for _, line := range lines {
		source, sinksPart, _ := strings.Cut(line, " -> ")
		sinks := strings.Split(sinksPart, ", ")

		base := Broadcaster{
			id:     ModuleID(strings.Trim(source, "%&")),
			system: system,
			sinks:  make([]ModuleID, len(sinks)),
		}
		for i, sink := range sinks {
			base.sinks[i] = ModuleID(sink)
		}

		var module Module
		switch {
		case strings.HasPrefix(source, "%"):
			module = &FlipFlop{Broadcaster: base}
		case strings.HasPrefix(source, "&"):
			module = &Conjunction{Broadcaster: base}
		default:
			module = &base
		}

		system.Modules[module.ID()] = module

		if first == "" {
			first = module.ID()
		}
	}

	button := &Button{system: system, sink: "broadcaster"}
	system.Modules[button.ID()] = button

	for _, module := range system.Modules {
		// Fill in any missing sinks.
		for _, sink := range module.Sinks() {
			if _, ok := system.Modules[sink]; !ok {
				system.Modules[sink] = &DummyModule{id: sink}
			}
		}
		if initer, ok := module.(interface{ init() }); ok {
			initer.init()
		}
	}

	return InputModules{
		ModuleSystem: system,
		Button:       button,
	}
}

func part1(input string) int {
	system := parseInput(input)
	var metrics ModuleSystemMetrics
	for cycle := 1; cycle <= 1000; cycle++ {
		system.Button.Push()
		for m := range system.Tick() {
			metrics = metrics.Accumulate(m)
		}
	}
	return metrics.HiPulses * metrics.LoPulses
}

func part2(input string) int {
	system := parseInput(input)

	// rx will have 1 source with 4 sources.
	rxsrcIDs := system.Modules.FindModulesWithSink("rx")
	if len(rxsrcIDs) != 1 {
		log.Println("you are probably not running the real input")
		return -1
	}

	rxsrc := system.Modules[rxsrcIDs[0]].(*Conjunction)
	aocutil.Assertf(len(rxsrc.sources) == 4, "expected 4 sources for %s", rxsrc.ID())

	// Run until each of the rx conjunctions have all hi states.
	rxsrcCycles := make([]int, 4)
	for cycle := 1; slices.Contains(rxsrcCycles, 0); cycle++ {
		system.Button.Send("", Lo)
		for _ = range system.Tick() {
			for i, state := range rxsrc.states {
				if state != Hi || rxsrcCycles[i] != 0 {
					continue
				}

				log.Printf(
					"rx conjunction %q took %d cycles to have all hi states",
					rxsrc.sources[i], cycle)

				rxsrcCycles[i] = cycle
			}
		}
	}

	// The number of cycles that it'll take for all of the rx conjunctions to
	// all have hi states is the LCM of the number of cycles it takes each of
	// them to have all hi states.
	lcm := aocutil.LCM(rxsrcCycles...)
	return lcm
}
