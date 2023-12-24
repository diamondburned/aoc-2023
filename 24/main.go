package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	. "libdb.so/aoc-2023/aocutil"
)

func main() {
	Run(part1, part2)
}

type FlyingHailstone struct {
	// For a point:
	// - x(t) = P.X + t * V.X
	// - y(t) = P.Y + t * V.Y
	// - z(t) = P.Z + t * V.Z
	// Therefore:
	// - t = (x - P.X) / V.X
	// - t = (y - P.Y) / V.Y
	// - t = (z - P.Z) / V.Z

	P Point3D[float64]
	V Point3D[float64]
}

func parseInput(input string) []FlyingHailstone {
	const f = `%d, %d, %d @ %d, %d, %d`
	lines := SplitLines(input)
	hailstones := make([]FlyingHailstone, len(lines))
	for i, line := range lines {
		var x, y, z, vx, vy, vz int
		fmt.Sscanf(line, f, &x, &y, &z, &vx, &vy, &vz)
		p := Pt3(float64(x), float64(y), float64(z))
		v := Pt3(float64(vx), float64(vy), float64(vz))
		hailstones[i] = FlyingHailstone{
			P: p,
			V: v,
		}
	}
	return hailstones
}

func (h FlyingHailstone) ToLine() Line3D[float64] {
	return Line3D[float64]{
		Start: h.P,
		End:   h.P.Add(h.V),
	}
}

func part1(input string) int {
	hailstones := parseInput(input)

	boundsMinValue := 7.0
	boundsMaxValue := 27.0
	if len(hailstones) >= 300 {
		// real input so use real bounds
		boundsMinValue = 200_000_000_000_000.0
		boundsMaxValue = 400_000_000_000_000.0
	}

	boundsMin := Pt[float64](boundsMinValue, boundsMinValue)
	boundsMax := Pt[float64](boundsMaxValue, boundsMaxValue)

	var count int

	for a, b := range PairCombinations(hailstones) {
		a2d := a.ToLine().RemoveZ()
		b2d := b.ToLine().RemoveZ()

		i, intersection := a2d.RayIntersection(b2d)
		if intersection == NoIntersection || !i.WithinInclusive(boundsMin, boundsMax) {
			continue
		}

		count++
	}

	return count
}

func part2(input string) int {
	/*
		Let (p0, v0) be the position and velocity of the rock that we'll throw
		so that it'll hit all the hailstones at t[i], where i is the index of
		the hailstone.

		Pick the first 3 hailstones, and let their positions be p[0], p[1], p[2]
		and their velocities be v[0], v[1], v[2].

		p0 + t[0] * v0 = p[0] + t[0] * v[0]
		p0 + t[1] * v0 = p[1] + t[1] * v[1]
		p0 + t[2] * v0 = p[2] + t[2] * v[2]

		implies

		(p0 - p[0]) + t[0] * (v0 - v[0]) = 0
		(p0 - p[1]) + t[1] * (v0 - v[1]) = 0
		(p0 - p[2]) + t[2] * (v0 - v[2]) = 0

		implies

		(p0 - p[0]) = t[0] * (v[0] - v0)
		(p0 - p[1]) = t[1] * (v[1] - v0)
		(p0 - p[2]) = t[2] * (v[2] - v0)

		You can expand this to 3 equations in 3 unknowns, with the unknowns
		being p0.X, p0.Y, p0.Z, v0.X, v0.Y, v0.Z, t[0], t[1], t[2].

		Flattening this system of equations, we get:

		p0.X + t[0] * v0.X = p[0].X + t[0] * v[0].X
		p0.Y + t[0] * v0.Y = p[0].Y + t[0] * v[0].Y
		p0.Z + t[0] * v0.Z = p[0].Z + t[0] * v[0].Z
		p0.X + t[1] * v0.X = p[1].X + t[1] * v[1].X
		p0.Y + t[1] * v0.Y = p[1].Y + t[1] * v[1].Y
		p0.Z + t[1] * v0.Z = p[1].Z + t[1] * v[1].Z
		p0.X + t[2] * v0.X = p[2].X + t[2] * v[2].X
		p0.Y + t[2] * v0.Y = p[2].Y + t[2] * v[2].Y
		p0.Z + t[2] * v0.Z = p[2].Z + t[2] * v[2].Z
	*/

	hailstones := parseInput(input)

	// Create a Qalculate input.
	eqns := make([]string, 0, 9)
	vars := []string{"px", "py", "pz", "vx", "vy", "vz"}
	for i, h := range hailstones[:3] {
		eqns = append(eqns,
			fmt.Sprintf(`(px + (vx * t%d) == %.0f + (%.0f * t%d))`, i, h.P.X, h.V.X, i),
			fmt.Sprintf(`(py + (vy * t%d) == %.0f + (%.0f * t%d))`, i, h.P.Y, h.V.Y, i),
			fmt.Sprintf(`(pz + (vz * t%d) == %.0f + (%.0f * t%d))`, i, h.P.Z, h.V.Z, i))
		vars = append(vars, fmt.Sprintf("t%d", i))
	}

	expr := fmt.Sprintf(
		strings.Join([]string{
			`var('%s')`,
			`print(solve([%s], [%s])[0])`,
		}, "\n"),
		strings.Join(vars, " "),
		strings.Join(eqns, ", "),
		strings.Join(vars, ", "))

	out := sage(expr)
	var p Point3D[int]
	Sscanf(out,
		"[px == %d, py == %d, pz == %d",
		&p.X, &p.Y, &p.Z)

	return p.X + p.Y + p.Z
}

func sage(input string) string {
	var b strings.Builder
	cmd := exec.Command("sage", "--nodotsage", "-c", input)
	cmd.Stdout = &b
	cmd.Stderr = os.Stderr
	E1(cmd.Run())
	return b.String()
}
