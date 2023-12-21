package main

import (
	"libdb.so/aoc-2023/aocutil"
)

// ModuleSystemMetrics is a collection of metrics for a module system.
type ModuleSystemMetrics struct {
	LoPulses int
	HiPulses int
}

// Accumulate adds two module system metrics together.
func (m ModuleSystemMetrics) Accumulate(other ModuleSystemMetrics) ModuleSystemMetrics {
	return ModuleSystemMetrics{
		LoPulses: m.LoPulses + other.LoPulses,
		HiPulses: m.HiPulses + other.HiPulses,
	}
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
