package main

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

// Broadcaster is the simplest module. It sends all signals it receives to all
// of its sinks.
type Broadcaster struct {
	baseModule
}

func (b *Broadcaster) Send(_ ModuleID, s Signal) { b.send(s) }

func (b *Broadcaster) String() string {
	sinks := aocutil.Map(b.sinks, func(sink ModuleID) string { return string(sink) })
	return fmt.Sprintf("broadcaster(%s) -> %s", b.id, strings.Join(sinks, ", "))
}

// FlipFlop is a flip-flop module,
type FlipFlop struct {
	baseModule
	state Signal
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
	baseModule
	sources []ModuleID
	states  []Signal
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
	baseModule
	s Signal
}

func (d *DummyModule) Send(_ ModuleID, s Signal) {
	d.s = s
}

func (d *DummyModule) String() string {
	return fmt.Sprintf("dummy(%s): %s", d.id, d.s)
}

// Button is a module that sends a signal when pressed.
type Button struct {
	system *ModuleSystem
	sink   ModuleID
}

func (b *Button) ID() ModuleID          { return "button" }
func (b *Button) Sinks() []ModuleID     { return []ModuleID{b.sink} }
func (b *Button) Send(ModuleID, Signal) { b.Push() }

func (b *Button) Push() {
	b.system.Send("button", b.sink, Lo)
}

func (b *Button) String() string {
	return fmt.Sprintf("button -> %s", b.sink)
}
