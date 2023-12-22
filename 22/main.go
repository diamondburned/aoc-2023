package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"slices"
	"strconv"

	"github.com/sourcegraph/conc/iter"
	"github.com/tidwall/pinhole"
	. "libdb.so/aoc-2023/aocutil"
)

func main() {
	Run(part1, part2)
}

type Point3D struct {
	X, Y, Z int
}

func (p Point3D) String() string {
	return fmt.Sprintf("(%d,%d,%d)", p.X, p.Y, p.Z)
}

func (p Point3D) Add(q Point3D) Point3D {
	return Point3D{p.X + q.X, p.Y + q.Y, p.Z + q.Z}
}

func (p Point3D) Sub(q Point3D) Point3D {
	return Point3D{p.X - q.X, p.Y - q.Y, p.Z - q.Z}
}

type Rect3D struct {
	Min, Max Point3D // both inclusive
}

func (r Rect3D) String() string {
	return fmt.Sprintf("%s~%s", r.Min, r.Max)
}

func (r Rect3D) Canon() Rect3D {
	if r.Min.X > r.Max.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Min.Y > r.Max.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	if r.Min.Z > r.Max.Z {
		r.Min.Z, r.Max.Z = r.Max.Z, r.Min.Z
	}
	return r
}

// ContainsPt returns true if the point is contained in the rectangle.
func (r Rect3D) ContainsPt(p Point3D) bool {
	return true &&
		r.Min.X <= p.X && p.X <= r.Max.X &&
		r.Min.Y <= p.Y && p.Y <= r.Max.Y &&
		r.Min.Z <= p.Z && p.Z <= r.Max.Z
}

// Contains returns true if the other rectangle is contained in the rectangle.
func (r Rect3D) Contains(other Rect3D) bool {
	return true &&
		r.Min.X <= other.Min.X && other.Max.X <= r.Max.X &&
		r.Min.Y <= other.Min.Y && other.Max.Y <= r.Max.Y &&
		r.Min.Z <= other.Min.Z && other.Max.Z <= r.Max.Z
}

// Overlaps returns true if the other rectangle overlaps the rectangle.
func (r Rect3D) Overlaps(other Rect3D) bool {
	return true &&
		other.Min.X <= r.Max.X && r.Min.X <= other.Max.X &&
		other.Min.Y <= r.Max.Y && r.Min.Y <= other.Max.Y &&
		other.Min.Z <= r.Max.Z && r.Min.Z <= other.Max.Z
}

// Intersect returns the intersection of the two rectangles.
func (r Rect3D) Intersect(other Rect3D) Rect3D {
	return Rect3D{
		Min: Point3D{
			X: max(r.Min.X, other.Min.X),
			Y: max(r.Min.Y, other.Min.Y),
			Z: max(r.Min.Z, other.Min.Z),
		},
		Max: Point3D{
			X: min(r.Max.X, other.Max.X),
			Y: min(r.Max.Y, other.Max.Y),
			Z: min(r.Max.Z, other.Max.Z),
		},
	}
}

func (r Rect3D) Translate(p Point3D) Rect3D {
	return Rect3D{
		Min: r.Min.Add(p),
		Max: r.Max.Add(p),
	}
}

func (r Rect3D) MoveToZ(z int) Rect3D {
	delta := z - r.Min.Z
	return r.Translate(Point3D{Z: delta})
}

func (r Rect3D) Size() int {
	return 1 *
		(r.Max.X - r.Min.X + 1) *
		(r.Max.Y - r.Min.Y + 1) *
		(r.Max.Z - r.Min.Z + 1)
}

type Brick struct {
	ID int
	Rect3D
}

// IsOnGround returns true if the brick is on the ground.
func (b Brick) IsOnGround() bool {
	return b.Min.Z == 0
}

// WithinZ returns true if the brick is within the given Z value.
func (b Brick) WithinZ(z int) bool {
	return b.Min.Z <= z && z <= b.Max.Z
}

func (b Brick) String() string {
	return fmt.Sprintf("%4d %s", b.ID, b.Rect3D)
}

type Bricks []Brick

func (b Bricks) Sort() {
	slices.SortFunc(b, func(a, b Brick) int {
		if a.Max.Z != b.Max.Z {
			return a.Max.Z - b.Max.Z
		}
		if a.Min.Z != b.Min.Z {
			return a.Min.Z - b.Min.Z
		}
		if a.Max.Y != b.Max.Y {
			return a.Max.Y - b.Max.Y
		}
		if a.Min.Y != b.Min.Y {
			return a.Min.Y - b.Min.Y
		}
		if a.Max.X != b.Max.X {
			return a.Max.X - b.Max.X
		}
		if a.Min.X != b.Min.X {
			return a.Min.X - b.Min.X
		}
		return 0
	})
}

// AtZ returns an iterator that yields all bricks at the given Z value.
func (bs Bricks) AtZ(z int, excludes ...Brick) Iter[Brick] {
	return func(yield func(Brick) bool) {
		for _, b := range bs {
			if b.WithinZ(z) && !slices.Contains(excludes, b) && !yield(b) {
				return
			}
		}
	}
}

// AtZAndOverlapsWith returns an iterator that yields all bricks at the given Z
// value that overlap the other rectangle.
func (bs Bricks) AtZAndOverlapsWith(z int, other Rect3D, excludes ...Brick) Iter[Brick] {
	return func(yield func(Brick) bool) {
		for _, b := range bs {
			if b.WithinZ(z) && b.Overlaps(other) && !slices.Contains(excludes, b) && !yield(b) {
				return
			}
		}
	}
}

// DropBrick drops the brick on the ground. It lowers the brick so that its Z is
// as low as it can be without colliding with any other brick.
func (bs Bricks) DropBrick(drop Brick) Brick {
	var z int
	for z = drop.Min.Z - 1; z > 0; z-- {
		if bs.AtZAndOverlapsWith(z, drop.MoveToZ(z)).Any() {
			break
		}
	}
	z++ // undo
	drop.Rect3D = drop.MoveToZ(z)
	return drop // TODO
}

// DropAllBricks drops all bricks on the ground.
// It returns the number of bricks that were dropped.
func (bs Bricks) DropAllBricks() int {
	var droppedCount int
	for i, b := range bs {
		dropped := bs.DropBrick(b)
		if dropped != b {
			bs[i] = dropped
			droppedCount++
		}
	}
	return droppedCount
}

func (bs Bricks) Render() image.Image {
	const strokeSize = 0.05

	var maxPt Point3D
	for _, b := range bs {
		maxPt.X = max(maxPt.X, b.Max.X)
		maxPt.Y = max(maxPt.Y, b.Max.Y)
		maxPt.Z = max(maxPt.Z, b.Max.Z)
		maxPt.X = max(maxPt.X, Abs(b.Min.X))
		maxPt.Y = max(maxPt.Y, Abs(b.Min.Y))
		maxPt.Z = max(maxPt.Z, Abs(b.Min.Z))
	}

	worldBounds := float64(max(maxPt.X, maxPt.Y, maxPt.Z)) * 2
	cubeSize := 1 / worldBounds
	log.Printf("worldBounds=%v cubeSize=%v", worldBounds, cubeSize)

	world := pinhole.New()
	colors := []color.Color{
		color.RGBA{0xff, 0x00, 0x00, 0xff},
		color.RGBA{0x00, 0xff, 0x00, 0xff},
		color.RGBA{0x00, 0x00, 0xff, 0xff},
	}
	for i, b := range bs {
		const fontScale = 0.5 * strokeSize
		world.Begin()
		world.DrawString(
			(float64(b.Min.X)+(float64(b.Max.X-b.Min.X)/2))/worldBounds/fontScale,
			(float64(b.Min.Y)+(float64(b.Max.Y-b.Min.Y)/2))/worldBounds/fontScale,
			(float64(b.Min.Z)+(float64(b.Max.Z-b.Min.Z)/2))/worldBounds/fontScale,
			strconv.Itoa(b.ID),
		)
		world.Scale(fontScale, fontScale, fontScale)
		world.Colorize(colors[i%len(colors)])
		world.End()

		log.Println("drawing at",
			float64(b.Min.X)/worldBounds-cubeSize/2,
			float64(b.Min.Y)/worldBounds-cubeSize/2,
			float64(b.Min.Z)/worldBounds-cubeSize/2,
			float64(b.Max.X)/worldBounds+cubeSize/2,
			float64(b.Max.Y)/worldBounds+cubeSize/2,
			float64(b.Max.Z)/worldBounds+cubeSize/2,
		)

		world.Begin()
		world.DrawCube(
			float64(b.Min.X)/worldBounds-cubeSize/2,
			float64(b.Min.Y)/worldBounds-cubeSize/2,
			float64(b.Min.Z)/worldBounds-cubeSize/2,
			float64(b.Max.X)/worldBounds+cubeSize/2,
			float64(b.Max.Y)/worldBounds+cubeSize/2,
			float64(b.Max.Z)/worldBounds+cubeSize/2,
		)
		world.Colorize(colors[i%len(colors)])
		world.End()
	}
	// world.Translate(-0.2, 0, -0.35)
	// world.Translate(-0.1, 0.5, -0.35)
	world.Translate(0, 0.65, -0.35)
	world.Rotate(-math.Pi/2.25, 0, 0)

	return world.Image(5000, 5000, &pinhole.ImageOptions{
		BGColor:   color.White,
		LineWidth: strokeSize,
		Scale:     1,
	})
}

// BricksOnTopOf returns an iterator that yields all bricks that are on top of the
// given brick.
func (bs Bricks) BricksOnTopOf(brick Brick) Iter[Brick] {
	moved := brick.Translate(Point3D{Z: +1})
	return bs.AtZAndOverlapsWith(moved.Max.Z, moved)
}

// UnsupportedBricksOnTopOf returns an iterator that yields all bricks that are
// on top of the given brick and are not supported by any other bricks.
// The bricks given in the excludes list are not considered as supporting
// bricks.
func (bs Bricks) UnsupportedBricksOnTopOf(b Brick, excludes ...Brick) Iter[Brick] {
	// take the slow way
	onTop := bs.BricksOnTopOf(b).All()
	if len(onTop) == 0 {
		return SliceIter(onTop)
	}

	// Search if there are any bricks that are at the same level as this
	// one and is suppporting the bricks on top of this one.
	excludes = append(slices.Clone(excludes), b)
	for neighboringBrick := range bs.AtZ(b.Max.Z, excludes...) {
		neighborMoved := neighboringBrick.Translate(Point3D{Z: +1})
		// We delete the bricks on top that this neighboring brick is
		// supporting.
		onTop = slices.DeleteFunc(onTop, func(top Brick) bool {
			return top.Overlaps(neighborMoved)
		})
		if len(onTop) == 0 {
			// We made sure that all bricks on top of the current brick
			// are supported by other bricks.
			break
		}
	}

	return SliceIter(onTop)
}

// // RemoveSupportedBricks accepts a list of bricks and removes all bricks that
// // are supported by some other bricks.
// func (bs Bricks) RemoveSupportedBricks(bricks []Brick) []Brick {}

// BrickIsDisintegrable returns true if the given brick can be disintegrated.
func (bs Bricks) BrickIsDisintegrable(b Brick) bool {
	return !bs.UnsupportedBricksOnTopOf(b).Any()
}

// AllDisintegrable returns an iterator that yields all bricks that can be
// disintegrated.
func (bs Bricks) AllDisintegrable() Iter[Brick] {
	return func(yield func(Brick) bool) {
		for _, b := range bs {
			if bs.BrickIsDisintegrable(b) && !yield(b) {
				return
			}
		}
	}
}

// AllNonDisintegrable returns an iterator that yields all bricks that cannot be
// disintegrated.
func (bs Bricks) AllNonDisintegrable() Iter[Brick] {
	return func(yield func(Brick) bool) {
		for _, b := range bs {
			if !bs.BrickIsDisintegrable(b) && !yield(b) {
				return
			}
		}
	}
}

func parseInput(input string) Bricks {
	lines := SplitLines(input)
	bricks := make(Bricks, len(lines))
	for i, line := range lines {
		const f = `%d,%d,%d~%d,%d,%d`
		b := Brick{ID: i}
		fmt.Sscanf(line, f,
			&b.Min.X, &b.Min.Y, &b.Min.Z,
			&b.Max.X, &b.Max.Y, &b.Max.Z)
		b.Rect3D = b.Rect3D.Canon()
		bricks[i] = b
	}
	bricks.Sort()
	return bricks
}

func part1(input string) int {
	bricks := parseInput(input)
	bricks.DropAllBricks()
	OpenImage(bricks.Render())
	return bricks.AllDisintegrable().Count()
}

func part2(input string) int {
	bricks := parseInput(input)
	bricks.DropAllBricks()

	allFallen := iter.Map(bricks.AllNonDisintegrable().All(), func(b *Brick) int {
		// Clone this so we can start just absolutely disintegrating bricks.
		simulation := slices.Clone(bricks)
		// Disintegrate this brick.
		simulation = slices.Delete(simulation, b.ID, b.ID+1)
		// Start simulating until no more bricks would fall.
		var fell int
		for {
			n := simulation.DropAllBricks()
			fell += n
			if n == 0 {
				break
			}
		}
		return fell
	})
	return Sum(allFallen)
}
