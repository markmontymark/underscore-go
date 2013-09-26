package underscore

import (
	"./lib/asserts"
	"fmt"
	"strings"
	"testing"
)

func TestChainingMapFlattenReduce(t *testing.T) {
	lyrics := []T{"I'm a lumberjack and I'm okay", "I sleep all night and I work all day", "He's a lumberjack and he's okay", "He sleeps all night and he works all day"}
	counts := New(lyrics).Chain().
		Map(func(line, idx, list T) T { return strings.Split(line.(string), "") }).
		Flatten().
		Reduce(func(hash, k, ignorel, ignore2 T) T {
		if _, ok := hash.(map[T]int)[k]; ok {
			hash.(map[T]int)[k] += 1
		} else {
			hash.(map[T]int)[k] = 1
		}
		return hash
	},
		map[T]int{}).
		Value()
	aCount, _ := counts.(map[T]int)["a"]
	eCount, _ := counts.(map[T]int)["e"]
	asserts.True(t, "Counted all the letters in the song", aCount == 16 && eCount == 10)
}

func TestSelectRejectSortBy(t *testing.T) {
	numbers1 := []T{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	numbers2 := New(numbers1).Chain().
		Select(func(n, idx, list T) bool { return n.(int)%2 == 0 }).
		Reject(func(n, idx, list T) bool { return n.(int)%4 == 0 }).
		SortBySorter(
		func(n, idx, list T) T { return -1 * n.(int) },
		func(a, b *map[T]T) bool {
			return (*a)["criteria"].(int) < (*b)["criteria"].(int)
		}).
		Value()
	asserts.Equals(t, "filtered and reversed the numbers",
		fmt.Sprintf("%v", numbers2), "[10 6 2]")
}

func TestChainingWorksInSmallStages(t *testing.T) {
	o := New([]T{1, 2, 3, 4}).Chain()
	asserts.Equals(t, "first two elems",
		fmt.Sprintf("%v", o.Filter(func(i, d, l T) bool { return i.(int) < 3 }).Value()), "[1 2]")

	asserts.Equals(t, "last two elems",
		fmt.Sprintf("%v", o.Filter(func(i, d, l T) bool { return i.(int) > 2 }).Value()), "[3 4]")

	asserts.Equals(t, "which of the last two elems is odd",
		fmt.Sprintf("%v",
			o.Filter(func(i, d, l T) bool { return i.(int) > 2 }).
				Filter(func(i, d, l T) bool { return i.(int)%2 == 0 }).
				Value()), "[4]")

	asserts.Equals(t, "which of the last two elems is odd",
		fmt.Sprintf("%v",
			o.Filter(func(i, d, l T) bool { return i.(int) < 5 }).
				Filter(func(i, d, l T) bool { return i.(int)%2 == 1 }).
				Value()), "[1 3]")
}

func TestReverseConcatUnshiftPopMap(t *testing.T) {
	numbers1 := []T{1, 2, 3, 4, 5}
	numbers2 := New(numbers1).
		Chain().
		Reverse().
		Concat([]T{5, 5, 5}).
		Unshift(17).
		Pop().
		Map(func(n, i, list T) T { return n.(int) * 2 }).
		Value()
	asserts.Equals(t, "can chain together array functions",
		fmt.Sprintf("%v", numbers2), "[34 10 8 6 4 2 10 10]")
}
