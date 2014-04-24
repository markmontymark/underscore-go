package underscore

import (
	_ "github.com/markmontymark/asserts"
	_ "fmt"
	"testing"
)

func BenchmarkFirst_a(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		First([]T{1, 2, 3})
	}
}
func BenchmarkFirst_b(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		FirstN([]T{1, 2, 3}, 2)
	}
}
func BenchmarkFirst_c(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		FirstN([]T{1, 2, 3}, 5)
	}
}
func BenchmarkFirst_d(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		Map([]T{[]T{1, 2, 3}, []T{1, 2, 3}}, func(obj T, idx T, list T) T { return First(obj.([]T)) })
	}
}
func BenchmarkFirst_e(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		First(nil)
	}
}

func BenchmarkRest_a(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		Rest([]T{1, 2, 3})
	}
}

func BenchmarkRest_b(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		Map([]T{[]T{1, 2, 3}, []T{1, 2, 3}}, func(obj T, idx T, list T) T { return Rest(obj.([]T)) })
	}
}

func BenchmarkRest_c(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		Rest(nil)
	}
}

func BenchmarkInitial_a(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		Initial([]T{1, 2, 3})
	}
}

func BenchmarkInitial_b(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		Initial([]T{1, 2, 3}, 2)
	}
}

func BenchmarkInitial_c(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		Initial([]T{1, 2, 3}, 0)
	}
}

func BenchmarkInitial_d(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		Initial([]T{1, 2, 3}, 1)
	}
}

func BenchmarkInitial_e(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		Initial([]T{1, 2, 3}, 3)
	}
}

func BenchmarkInitial_f(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		Initial([]T{1, 2, 3}, 4)
	}
}

func BenchmarkInitial_g(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		Initial([]T{1, 2, 3}, 5)
	}
}

func BenchmarkInitial_h(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		Map([]T{[]T{1, 2, 3}, []T{1, 2, 3}}, func(obj T, idx T, list T) T { return Initial(obj.([]T), 2) })
	}
}

func BenchmarkInitial_i(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		Initial(nil)
	}
}


/*
func BenchmarkLast(t *testing.B) {
	asserts.Equals(t, "can pull out the initial elements of an array",
		Last([]T{1, 2, 3})),
		"[3]")
	asserts.Equals(t, "can pull out the last 2 element of an array",
		Last([]T{1, 2, 3}, 2)), "[2 3]")

	asserts.Equals(t, "can pull out the zeroth element of an array",
		Last([]T{1, 2, 3}, 0)), "[]")

	asserts.Equals(t, "can pull out the first element of an array",
		Last([]T{1, 2, 3}, 1)), "[3]")

	asserts.Equals(t, "can pull out all elements of an array",
		Last([]T{1, 2, 3}, 3)), "[1 2 3]")

	asserts.Equals(t, "can pull out too many out",
		Last([]T{1, 2, 3}, 4)), "[1 2 3]")

	asserts.Equals(t, "can pull out too many items out",
		Last([]T{1, 2, 3}, 5)), "[1 2 3]")

	asserts.Equals(t, "works well with map",
		Map([]T{[]T{1, 2, 3}, []T{1, 2, 3}}, func(obj T, idx T, list T) T { return Last(obj.([]T), 2) })), "[[2 3] [2 3]]")

	asserts.Equals(t, "works well with nil",
		Initial(nil)), "[]")
}

func BenchmarkCompact(t *testing.B) {
	asserts.Equals(t, "can trim out all falsy values",
		Compact([]T{0, 1, false, 2, false, 3})), "[1 2 3]")
	asserts.Equals(t, "can trim out all falsy values",
		Compact(nil)), "[]")
}

func BenchmarkFlatten(t *testing.B) {
	list := []T{1, []T{2}, []T{3, []T{[]T{[]T{4}}}}}
	asserts.Equals(t, "can flatten nested arrays", Flatten(list)), "[1 2 3 4]")
	asserts.Equals(t, "can shallowly flatten nested arrays", Flatten(list, true)),
		"[1 2 3 [[[4]]]]")
	list2 := []T{[]T{1}, []T{2}, []T{3}, []T{[]T{4}}}
	asserts.Equals(t, "can shallowly flatten arrays containing only other arrays",
		Flatten(list2, true)), "[1 2 3 [4]]")
	asserts.Equals(t, "can flatten arrays containing only other arrays",
		Flatten(list2)), "[1 2 3 4]")

	list3 := []T{[]T{"a", "b", "c"}, []T{"d", "e", "f", "g"}}
	asserts.Equals(t, "can flatten arrays containing only other arrays",
		Flatten(list3)), "[a b c d e f g]")
	asserts.Equals(t, "can shallow flatten arrays containing only other arrays",
		Flatten(list3, true)), "[a b c d e f g]")
}

func BenchmarkWithout(t *testing.B) {

	list := []T{1, 2, 1, 0, 3, 1, 4}
	asserts.Equals(t, "can remove arbitrary ints from a list",
		Without(list, 0, 1)), "[2 3 4]")

	asserts.Equals(t, "can remove all instances from an array of args from a list",
		Without(list, []T{0, 1})), "[2 3 4]")

	list2 := []T{1, 2, "1", 0, 3, 1, 4}
	asserts.Equals(t, "can remove all instances of func args from an object",
		Without(list2, 0, 1)), "[2 1 3 4]")

	asserts.Equals(t, "can remove all instances from an array of args from an object",
		Without(list2, []T{0, 1})), "[2 1 3 4]")

	asserts.Equals(t, "can remove all instances from an array of args from an object",
		Without(list2, []T{0, "1"})), "[1 2 3 1 4]")

	asserts.Equals(t, "can remove all instances from an array of args from an object",
		Without(list2, []T{0, "1", 1})), "[2 3 4]")

	// Kludgy comparator?, go can't compare maps via identity, so stringifying each 
	listm := []T{map[T]T{"one": 1}, map[T]T{"two": 2}}
	retval := Without(listm, map[T]T{"one": 1}, func(a, b T) bool {
		return a) == b)
	})
	asserts.Equals(t, "uses real object identity for comparisons.", retval), "[map[two:2]]")
	asserts.True(t, "ditto", len(retval) == 1)
}

func BenchmarkPartition(t *testing.B) {
	list := []T{1, 2, 1, 0, 3, 1, 4}
	isOdd := func(elem T) bool {
		return elem.(int) % 2 != 0
	}
	asserts.Equals(t, "Can partition list ",
		Partition(list, isOdd)), "[[1 1 3 1] [2 0 4]]")
}


func BenchmarkUniq(t *testing.B) {
	list := []T{1, 2, 1, 3, 1, 4}
	asserts.Equals(t, "can find the unique values of an unsorted array",
		Uniq(list, false)), "[1 2 3 4]")

	list2 := []T{1, 1, 1, 2, 2, 3}
	asserts.Equals(t, "can find the unique values of a sorted array faster",
		Uniq(list2, true)), "[1 2 3]")

	list3 := []map[T]T{{"name": "moe"}, {"name": "curly"}, {"name": "larry"}, {"name": "curly"}}
	iterator := func(value T, key T, list T) T { return value.(map[T]T)["name"] }
	comparator := func(a T, b T) bool { return a.(map[T]T)["name"] == b.(map[T]T)["name"] }
	asserts.Equals(t, "can find the unique values of an array using a custom iterator",
		Uniq(list3, false, iterator, comparator)), "[map[name:moe] map[name:curly] map[name:larry]]")
	asserts.Equals(t, "can find the unique values of an array using a custom iterator",
		Map(Uniq(list3, false, iterator, comparator), iterator)), "[moe curly larry]")

}

func BenchmarkIntersection(t *testing.B) {
	strLessThan := func(this T, that T) bool {
		return this.(string) < that.(string)
	}
	stooges := []T{"moe", "curly", "larry"}
	leaders := []T{"moe", "groucho"}
	asserts.Equals(t, "can take the set intersection of two arrays",
		Intersection(strLessThan, stooges, leaders)), "[moe]")

	theSixStooges := []T{"moe", "moe", "curly", "curly", "larry", "larry"}
	asserts.Equals(t, "returns a duplicate-free array",
		Intersection(strLessThan, theSixStooges, leaders)), "[moe]")
}

func BenchmarkUnion(t *testing.B) {
	result := Union([]T{1, 2, 3}, []T{2, 30, 1}, []T{1, 40})
	asserts.Equals(t, "takes the union of a list of arrays", result), "[1 2 3 30 40]")

	result2 := Union([]T{1, 2, 3}, []T{2, 30, 1}, []T{1, 40, []T{1}})
	asserts.Equals(t, "takes the union of a list of nested arrays", result2), "[1 2 3 30 40 [1]]")

	result3 := Union(nil, []T{1, 2, 3})
	asserts.Equals(t, "takes the union of nil and a list of arrays", result3), "[<nil> 1 2 3]")
}

func BenchmarkDifference(t *testing.B) {
	result := Difference([]T{1, 2, 3}, IdentityComparator, []T{2, 30, 40})
	asserts.Equals(t, "takes the difference of two arrays",
		result), "[1 3]")

	result2 := Difference([]T{1, 2, 3, 4}, IdentityComparator, []T{2, 30, 40}, []T{1, 11, 111})
	asserts.Equals(t, "takes the difference of three arrays",
		result2), "[3 4]")
}

func BenchmarkZip(t *testing.B) {

	names := []T{"moe", "larry", "curly"}
	ages := []T{30, 40, 50}
	leaders := []T{true}

	stooges := Zip(names, ages, leaders)
	asserts.Equals(t, "zipped together arrays of different lengths",
		stooges), "[[moe 30 true] [larry 40 <nil>] [curly 50 <nil>]]")

	stooges2 := Zip([]T{"moe", 30, "stooge 1"}, []T{"larry", 40, "stooge 2"}, []T{"curly", 50, "stooge 3"})
	asserts.Equals(t, "zipped pairs",
		stooges2), "[[moe larry curly] [30 40 50] [stooge 1 stooge 2 stooge 3]]")

	// In the case of difference lengths of the tuples undefineds
	// should be used as placeholder
	stooges3 := Zip([]T{"moe", 30}, []T{"larry", 40}, []T{"curly", 50, "extra data"})
	asserts.Equals(t, "zipped pairs with empties",
		stooges3), "[[moe larry curly] [30 40 50] [<nil> <nil> extra data]]")
	empty := Zip([]T{})
	asserts.Equals(t, "unzipped empty", empty), "[]")

	empty2 := Zip([]T{})
	asserts.Equals(t, "unzipped empty2", empty2), "[]")
}

func BenchmarkObject(t *testing.B) {
	result := Object([]T{"moe", "larry", "curly"}, []T{30, 40, 50})
	shouldBe := map[T]T{"moe": 30, "larry": 40, "curly": 50}
	asserts.Equals(t, "two arrays zipped together into an object",
		result), shouldBe))

	result2 := Object([]T{"one", 1, "two", 2, "three", 3})
	shouldBe2 := map[T]T{"one": 1, "two": 2, "three": 3}
	asserts.Equals(t, "an array of pairs zipped together into an object",
		result2), shouldBe2))

	result3 := Object([]T{[]T{"one", 1}, []T{"two", 2}, []T{"three", 3}})
	shouldBe3 := map[T]T{"one": 1, "two": 2, "three": 3}
	asserts.Equals(t, "an array of pairs zipped together into an object",
		result3), shouldBe3))

	asserts.Nil(t, "handles nils", Object(nil))

	stooges := map[T]T{"moe": 30, "larry": 40, "curly": 50}
	asserts.Equals(t, "an object converted to pairs and back to an object",
		stooges), Object(Pairs(stooges))))
}

func BenchmarkIndexOf(t *testing.B) {
	intLessThan := func(this T, that T) bool {
		return this.(int) < that.(int)
	}
	numbers := []T{1, 2, 3}
	asserts.IntEquals(t, "can compute indexOf, even without the native function",
		IndexOf(numbers, 2, intLessThan), 1)

	asserts.IntEquals(t, "handles nulls properly",
		IndexOf(nil, 2, intLessThan), -1)

	num := 35
	numbers2 := []T{10, 20, 30, 40, 50}
	asserts.IntEquals(t, "35 is not in the list", IndexOf(numbers2, num, intLessThan, true), -1)
	num = 40
	asserts.IntEquals(t, "40 is in the list", IndexOf(numbers2, num, intLessThan, true), 3)

	numbers3 := []T{1, 40, 40, 40, 40, 40, 40, 40, 50, 60, 70}
	num = 40
	asserts.IntEquals(t, "40 is in the list ", IndexOf(numbers3, num, intLessThan, true), 1)
}

func BenchmarkLastIndexOf(t *testing.B) {
	numbers := []T{1, 0, 1}
	asserts.IntEquals(t, "lastindexof simple test", LastIndexOf(numbers, 1), 2)

	numbers2 := []T{1, 0, 1, 0, 0, 1, 0, 0, 0}
	asserts.IntEquals(t, "can compute lastIndexOf", LastIndexOf(numbers2, 1), 5)
	asserts.IntEquals(t, "lastIndexOf the other element", LastIndexOf(numbers2, 0), 8)

	asserts.IntEquals(t, "handles nulls properly", LastIndexOf(nil, 2), -1)

	numbers3 := []T{1, 2, 3, 1, 2, 3, 1, 2, 3}
	index := LastIndexOf(numbers3, 2, 2)
	asserts.IntEquals(t, "supports the fromIndex argument", index, 1)

}

func BenchmarkRange(t *testing.B) {
	asserts.Equals(t, "range with no arguments generates an empty array",
		Range()), "[]")
	asserts.Equals(t, "range with 0 as a first argument generates an empty array",
		Range(0)), "[]")
	asserts.Equals(t, "range with a single positive argument generates an array of elements 0,1,2,...,n-1",
		Range(4)), "[0 1 2 3]")
	asserts.Equals(t, "range with two arguments a &amp; b, a&lt;b generates an array of elements a,a+1,a+2,...,b-2,b-1",
		Range(5, 8)), "[5 6 7]")
	asserts.Equals(t, "range with two arguments a &amp; b, b&lt;a generates an empty array",
		Range(8, 5)), "[]")
	asserts.Equals(t, "range with three arguments a &amp; b &amp; c, c &lt; b-a, a &lt; b generates an array of elements a,a+c,a+2c,...,b - (multiplier of a) &lt; c",
		Range(3, 10, 3)), "[3 6 9]")
	asserts.Equals(t, "range with three arguments a &amp; b &amp; c, c &gt; b-a, a &lt; b generates an array with a single element, equal to a",
		Range(3, 10, 15)), "[3]")
	asserts.Equals(t, "range with three arguments a &amp; b &amp; c, a &gt; b, c &lt; 0 generates an array of elements a,a-c,a-2c and ends with the number not less than b",
		Range(12, 7, -2)), "[12 10 8]")
	asserts.Equals(t, "final example in the Python docs",
		Range(0, -10, -1)), "[0 -1 -2 -3 -4 -5 -6 -7 -8 -9]")

}
*/
