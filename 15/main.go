package main

import (
	"context"
	"fmt"
	"image"
	"math"
	"runtime"
	"sort"
	"sync"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	m := make(Map)

	input := aocutil.InputString()
	lines := aocutil.SplitLines(input)
	for _, line := range lines {
		const fstr = "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d"
		var sensorAt, beaconAt Pt
		aocutil.Sscanf(line, fstr, &sensorAt.X, &sensorAt.Y, &beaconAt.X, &beaconAt.Y)

		m[beaconAt] = Beacon{}
		m[sensorAt] = Sensor{
			NearestBeacon:  beaconAt,
			BeaconDistance: mdist(sensorAt, beaconAt),
		}
	}

	ctx := context.Background()
	part1(ctx, m, 2_000_000)
	part2(ctx, m, 4_000_000)
}

func part1(ctx context.Context, m Map, y int) {
	min := Pt{math.MinInt, y}
	max := Pt{math.MaxInt, y}

	countCh := make(chan int, 1)
	sensorRanges(ctx, m, min, max, func(rs []Range) {
		r := rs[0]
		countCh <- r.To.X - r.From.X
	})

	fmt.Println("part 1:", <-countCh)
}

func part2(ctx context.Context, m Map, max int) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	distressPtCh := make(chan Pt, 1)

	sensorRanges(ctx, m, Pt{0, 0}, Pt{max, max}, func(rs []Range) {
		// Find any gaps within our merged ranges.
		for i, curr := range rs[1:] {
			// If the previous range ends before the current range
			// starts, then we have a gap.
			if prev := rs[i]; prev.To.X+1 < curr.From.X {
				distressPtCh <- Pt{prev.To.X + 1, prev.To.Y}
				cancel()
				return
			}
		}
	})

	select {
	case distressPt := <-distressPtCh:
		fmt.Println("part 2:", tuningFreq(distressPt), distressPt)
	default:
		fmt.Println("part 2: no distress signal found")
	}
}

const distressC = 4_000_000

func tuningFreq(pt Pt) int {
	return (pt.X * distressC) + pt.Y
}

func sensorRanges(ctx context.Context, m Map, min, max Pt, fn func([]Range)) {
	var wg sync.WaitGroup
	defer wg.Wait()

	in := make(chan [2]int)
	defer close(in)

	sensors := m.Sensors()
	beacons := m.Beacons()

	ncpu := aocutil.Min2(runtime.NumCPU(), max.Y-min.Y+1)
	for i := 0; i < ncpu; i++ {
		var ranges []Range
		var merged []Range

		do := func(y int) {
			// Zero out our variables.
			ranges = ranges[:0]
			merged = merged[:0]

			// Precalculate the Range for this row.
			for _, sensor := range sensors {
				// We can use a trick here: we can calculate how much Range a
				// single sensor has for a single row by subtracting its maximum
				// distance to our Y, both of which are already calculated.
				yDist := abs(y - sensor.At.Y)
				if sensor.Data.BeaconDistance < yDist {
					continue
				}

				xDist := sensor.Data.BeaconDistance - yDist
				ranges = append(ranges, Range{
					Pt{aocutil.Max2(sensor.At.X-xDist, min.X), y},
					Pt{aocutil.Min2(sensor.At.X+xDist, max.X), y},
				})
			}

			// Account for beacons too.
			for _, beacon := range beacons {
				if beacon.At.Y == y && min.X <= beacon.At.X && beacon.At.X <= max.X {
					ranges = append(ranges, Range{beacon.At, beacon.At})
				}
			}

			// Sort the ranges in the order that they start so that we can
			// merge them.
			sort.Slice(ranges, func(i, j int) bool {
				return ranges[i].From.X < ranges[j].From.X
			})

			// Merge the ranges.
			merged = append(merged, ranges[0])
			for _, r := range ranges[1:] {
				last := &merged[len(merged)-1]
				// If the current range starts before the previous range
				// ends, then we need to merge them.
				if r.From.X <= last.To.X {
					// Extend the previous range to the end of the
					// current range.
					last.To.X = aocutil.Max2(last.To.X, r.To.X)
				} else {
					// Otherwise, we can just append the current range.
					merged = append(merged, r)
				}
			}

			fn(merged)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case yrange, ok := <-in:
					if !ok {
						return
					}
					for y := yrange[0]; y <= yrange[1]; y++ {
						select {
						case <-ctx.Done():
							return
						default:
							do(y)
						}
					}
				}
			}
		}()
	}

	const batch = 10_000
	for y1 := min.Y; y1 <= max.Y; y1 += batch {
		y2 := aocutil.Min2(y1+batch, max.Y)

		select {
		case <-ctx.Done():
			return
		case in <- [2]int{y1, y2}:
			// ok
		}
	}
}

// Pt is a point in 2D space.
type Pt = image.Point

// Range is a range of points.
type Range struct {
	From, To Pt
}

// Map is a map of the area.
type Map map[Pt]Data

// MapPt is a point in 2D space with data.
type MapPt[T Data] struct {
	At   Pt
	Data T
}

// Beacons returns the list of sensors.
func (m Map) Beacons() []MapPt[Beacon] {
	var pts []MapPt[Beacon]
	for pt, data := range m {
		if beaecon, ok := data.(Beacon); ok {
			pts = append(pts, MapPt[Beacon]{pt, beaecon})
		}
	}
	return pts
}

// Sensors returns the list of sensors.
func (m Map) Sensors() []MapPt[Sensor] {
	var pts []MapPt[Sensor]
	for pt, data := range m {
		if sensor, ok := data.(Sensor); ok {
			pts = append(pts, MapPt[Sensor]{pt, sensor})
		}
	}
	return pts
}

type Data interface {
	Data() byte
}

func (Beacon) Data() byte { return 'B' }
func (Sensor) Data() byte { return 'S' }

type Beacon struct{}

type Sensor struct {
	NearestBeacon  Pt
	BeaconDistance int
}

func mdist(a, b Pt) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
