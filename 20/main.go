package main

import (
	"log"
	"slices"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	aocutil.Run(part1, part2)
}

type InputModules struct {
	*ModuleSystem
	Button *Button
}

func parseInput(input string) InputModules {
	lines := aocutil.SplitLines(input)
	system := &ModuleSystem{Modules: make(map[ModuleID]Module, len(lines))}

	for _, line := range lines {
		source, sinksPart, _ := strings.Cut(line, " -> ")
		sinks := strings.Split(sinksPart, ", ")

		base := baseModule{
			id:     ModuleID(strings.Trim(source, "%&")),
			sinks:  aocutil.Map(sinks, func(sink string) ModuleID { return ModuleID(sink) }),
			system: system,
		}

		var module Module
		switch {
		case strings.HasPrefix(source, "%"):
			module = &FlipFlop{baseModule: base}
		case strings.HasPrefix(source, "&"):
			module = &Conjunction{baseModule: base}
		default:
			module = &Broadcaster{baseModule: base}
		}

		system.Modules[module.ID()] = module
	}

	button := &Button{system: system, sink: "broadcaster"}
	system.Modules[button.ID()] = button

	// Fill in any missing sinks.
	for _, module := range system.Modules {
		for _, sink := range module.Sinks() {
			if _, ok := system.Modules[sink]; !ok {
				system.Modules[sink] = &DummyModule{
					baseModule: baseModule{id: sink, system: system},
				}
			}
		}
	}

	// Initialize any modules that need it.
	for _, module := range system.Modules {
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
	rxsrcCycles := make([]int, len(rxsrc.sources))

	// Run until each of the rx conjunctions have all hi states.
	for cycle := 1; slices.Contains(rxsrcCycles, 0); cycle++ {
		system.Button.Push()
		for _ = range system.Tick() {
			for i, state := range rxsrc.states {
				if state == Hi && rxsrcCycles[i] == 0 {
					log.Printf(
						"rx conjunction %q took %d cycles to have all hi states",
						rxsrc.sources[i], cycle,
					)
					rxsrcCycles[i] = cycle
				}
			}
		}
	}

	// The number of cycles that it'll take for all of the rx conjunctions to
	// all have hi states is the LCM of the number of cycles it takes each of
	// them to have all hi states.
	return aocutil.LCM(rxsrcCycles...)
}
