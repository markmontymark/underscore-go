package underscore

import (
	"./lib/asserts"
	"fmt"
	"math"
	"testing"
)

// TODO: missing Underscore ctor tests
// TODO: missing Identity tests

func TestRandom(t *testing.T) {
	array := Range(1000)
	min := math.Pow(2, 31)
	max := math.Pow(2, 62)

	asserts.True(t, "should produce a random number greater than or equal to the minimum number",
		Every(array, func(a, b, c T) bool {
			r := RandomFloat64(min, max)
			return r >= min && r <= max
		}))

	asserts.True(t, "should produce a random number when passed `Number.MAX_VALUE`",
		Some(array, func(a, b, c T) bool {
			r := RandomFloat64(math.MaxFloat64)
			return r > 0.0
		}))
}

// TODO: missing UniqueId tests

func TestTimes(t *testing.T) {
	vals := []T{}
	Times(3, func(i ...T) T {
		vals = append(vals, i...)
		return nil
	})
	asserts.Equals(t, "is 0 indexed", fmt.Sprint( vals), "[0 1 2]")

	vals2 := []T{}
	New(3).Times(func(i ...T) T {
		vals2 = append(vals2, i...)
		return i[0]
	})
	asserts.Equals(t, "works as a wrapper", fmt.Sprint( vals2), "[0 1 2]")

	// collects return values
	asserts.Equals(t, "collects return values",
		fmt.Sprint( New(3).Times(func(i ...T) T { return i[0] })), "[0 1 2]")

	asserts.Equals(t, "zero times retval is empty array",
		fmt.Sprint( Times(0, New(nil).Identity)), "[]")
	asserts.Equals(t, " -1 times retval is empty array",
		fmt.Sprint( Times(-1, New(nil).Identity)), "[]")
	asserts.Equals(t, " -Infinity times retval is empty array",
		fmt.Sprint( Times(int(math.Inf(-1)), New(nil).Identity)),
		"[]")
}

// XXX: missing mixin tests -- might not do...no prototype to add to in Go
// XXX: missing escape,unescape tests -- might not do...Go and encoding/ package
// XXX: missing template tests -- might not do...Go has its own template package

// TODO: missing result tests
