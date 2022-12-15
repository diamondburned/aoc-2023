package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"runtime"
	"sort"
	"time"

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

	part1(m, 2_000_000)
	start := time.Now()
	part2(m, 4_000_000)
	log.Println("part 2 took", time.Since(start))
}

func part1(m Map, y int) {
	sensors := m.NearbySensorsForRow(y)
	beacons := m.BeaconsAtRow(y)

	beaconMap := aocutil.NewSet[int](len(beacons))
	for _, pt := range beacons {
		beaconMap[pt.At.X] = struct{}{}
	}

	minX := int(math.MaxInt64)
	maxX := int(math.MinInt64)

	// We need to find the minumum and maximum X that our sensors cover from the
	// wanted y.
	for _, sensor := range sensors {
		minX = aocutil.Min2(minX, sensor.At.X-sensor.Data.BeaconDistance)
		maxX = aocutil.Max2(maxX, sensor.At.X+sensor.Data.BeaconDistance)
	}

	coverage := aocutil.NewSet[int](int(maxX - minX))
	for x := minX; x <= maxX; x++ {
		if beaconMap.Has(x) {
			continue
		}

		for _, sensor := range sensors {
			if sensorCovers(sensor, Pt{x, y}) {
				coverage[x] = struct{}{}
				break
			}
		}
	}

	var count int
	for range coverage {
		count++
	}
	fmt.Println("part 1:", count)
}

// sensorCovers returns true if the given point is within the region that covers
// from the sensor to the nearest beacon.
func sensorCovers(sensor MapPt[Sensor], pt Pt) bool {
	d1 := mdist(sensor.At, pt)
	d2 := sensor.Data.BeaconDistance
	return d1 <= d2
}

func part2(m Map, max int) {
	type coverage struct{ From, To int }

	sensors := make([]MapPt[Sensor], 0, len(m))
	for pt, data := range m {
		if sensor, ok := data.(Sensor); ok {
			sensors = append(sensors, MapPt[Sensor]{pt, sensor})
		}
	}

	// Sort so that the sensors with the largest beacon distance are first.
	sort.Slice(sensors, func(i, j int) bool {
		return sensors[i].Data.BeaconDistance < sensors[j].Data.BeaconDistance
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inputCh := make(chan coverage) // y
	outputCh := make(chan Pt)

	go func() {
		defer cancel()

		var ranges []coverage
		var merged []coverage

		do := func(y int) {
			// Zero out our variables.
			ranges = ranges[:0]
			merged = merged[:0]

			// Precalculate the coverage for this row.
			for _, sensor := range sensors {
				// We can use a trick here: we can calculate how much
				// coverage a single sensor has for a single row by
				// subtracting its maximum distance to our Y, both of which
				// are already calculated.
				yDist := abs(y - sensor.At.Y)
				if sensor.Data.BeaconDistance < yDist {
					continue
				}

				xDist := sensor.Data.BeaconDistance - yDist
				x1 := aocutil.Max2(sensor.At.X-xDist, 0)
				x2 := aocutil.Min2(sensor.At.X+xDist, max)
				ranges = append(ranges, coverage{x1, x2})
			}

			// Account for beacons too.
			for _, pt := range m.BeaconsAtRow(y) {
				if 0 <= pt.At.X && pt.At.X <= max {
					ranges = append(ranges, coverage{pt.At.X, pt.At.X})
				}
			}

			// Sort the ranges in the order that they start so that we can
			// merge them.
			sort.Slice(ranges, func(i, j int) bool {
				return ranges[i].From < ranges[j].From
			})

			// Merge the ranges.
			merged = append(merged, ranges[0])
			for _, r := range ranges[1:] {
				last := &merged[len(merged)-1]
				// If the current range starts before the previous range
				// ends, then we need to merge them.
				if r.From <= last.To {
					// Extend the previous range to the end of the
					// current range.
					last.To = aocutil.Max2(last.To, r.To)
				} else {
					// Otherwise, we can just append the current range.
					merged = append(merged, r)
				}
			}

			// Find any gaps within our merged ranges.
			prev := merged[0]
			for _, curr := range merged[1:] {
				// If the previous range ends before the current range
				// starts, then we have a gap.
				if prev.To+1 < curr.From {
					// We'll just assume that this gap is what we want
					// and is the size of 1.
					select {
					case <-ctx.Done():
					case outputCh <- Pt{prev.To + 1, y}:
					}
					return
				}
			}
		}

		for {
			select {
			case <-ctx.Done():
				return
			case yr := <-inputCh:
				for y := yr.From; y <= yr.To; y++ {
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

	var distressPos Pt

	batch := max / runtime.NumCPU()
search:
	for y := int(0); y < max; y += batch {
		select {
		case <-ctx.Done():
			break search
		case distressPos = <-outputCh:
			break search
		case inputCh <- coverage{y, aocutil.Min2(y+batch, max)}:
			// ok
		}
	}

	fmt.Println("part 2:", tuningFreq(distressPos), distressPos)
}

const distressC = 4_000_000

func tuningFreq(pt Pt) int {
	return (pt.X * distressC) + pt.Y
}

type Pt struct{ X, Y int }

type Map map[Pt]Data

type MapPt[T Data] struct {
	At   Pt
	Data T
}

func (m Map) DrawMap(w io.Writer, start, end Pt) {
	for y := start.Y; y <= end.Y; y++ {
		for x := start.X; x <= end.X; x++ {
			p := Pt{x, y}
			if data, ok := m[p]; ok {
				w.Write([]byte{data.Data()})
			} else {
				w.Write([]byte("."))
			}
		}
		w.Write([]byte{'\n'})
	}
}

// NearbySensorsForRow returns the list of coordinates of the sensors that are
// near the given row. A sensor is near if its signal reaches the given row.
func (m Map) NearbySensorsForRow(y int) []MapPt[Sensor] {
	var pts []MapPt[Sensor]
	for pt, data := range m {
		if sensor, ok := data.(Sensor); ok {
			if pt.Y-sensor.BeaconDistance <= y && y <= pt.Y+sensor.BeaconDistance {
				pts = append(pts, MapPt[Sensor]{pt, sensor})
			}
		}
	}
	return pts
}

// BeaconsAtRow returns the list of coordinates of the beacons that are at the
// given row.
func (m Map) BeaconsAtRow(y int) []MapPt[Beacon] {
	var pts []MapPt[Beacon]
	for pt, data := range m {
		if beacon, ok := data.(Beacon); ok && pt.Y == y {
			pts = append(pts, MapPt[Beacon]{pt, beacon})
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
