package underscore

import (
	"./lib/asserts"
	"fmt"
	"testing"
)

//type IntSlice []int
//func (this IntSlice) Len () int { return len(this) }
//func (this IntSlice) Less ( a , b int) bool { return this[a] < this[b] }
//func (this IntSlice) Swap( a , b int)  { this[a],this[b] = this[b],this[a] }

func TestFirst(t *testing.T) {
	asserts.IntEquals(t, "can pull out the first element of an array", First([]T{1, 2, 3}).(int), 1)
	asserts.Equals(t, "can pull out the first element of an array",
		fmt.Sprint( FirstN([]T{1, 2, 3}, 2)), "[1 2]")

	asserts.Equals(t, "can pull out the zeroth element of an array",
		fmt.Sprint( FirstN([]T{1, 2, 3}, 0)), "[1]")

	asserts.Equals(t, "can pull out too many items out",
		fmt.Sprint( FirstN([]T{1, 2, 3}, 5)), "[1 2 3]")

	asserts.Equals(t, "works well with map",
		fmt.Sprint( Map([]T{[]T{1, 2, 3}, []T{1, 2, 3}}, func(obj T, idx T, list T) T { return First(obj.([]T)) })), "[1 1]")
	asserts.Equals(t, "works well with nil",
		fmt.Sprint( First(nil)), "<nil>")
}

func TestRest(t *testing.T) {
	asserts.Equals(t, "can pull out rest of items",
		fmt.Sprint( Rest([]T{1, 2, 3})), "[2 3]")

	asserts.Equals(t, "works well with map",
		fmt.Sprint( Map([]T{[]T{1, 2, 3}, []T{1, 2, 3}}, func(obj T, idx T, list T) T { return Rest(obj.([]T)) })), "[[2 3] [2 3]]")

	asserts.Equals(t, "works well with nil",
		fmt.Sprint( Rest(nil)), "[]")
}

func TestInitial(t *testing.T) {
	asserts.Equals(t, "can pull out the initial elements of an array",
		fmt.Sprint( Initial([]T{1, 2, 3})),
		"[1 2]")
	asserts.Equals(t, "can pull out the first element of an array",
		fmt.Sprint( Initial([]T{1, 2, 3}, 2)), "[1 2]")

	asserts.Equals(t, "can pull out the zeroth element of an array",
		fmt.Sprint( Initial([]T{1, 2, 3}, 0)), "[]")

	asserts.Equals(t, "can pull out the first element of an array",
		fmt.Sprint( Initial([]T{1, 2, 3}, 1)), "[1]")

	asserts.Equals(t, "can pull out all elements of an array",
		fmt.Sprint( Initial([]T{1, 2, 3}, 3)), "[1 2 3]")

	asserts.Equals(t, "can pull out too many out",
		fmt.Sprint( Initial([]T{1, 2, 3}, 4)), "[1 2 3]")

	asserts.Equals(t, "can pull out too many items out",
		fmt.Sprint( Initial([]T{1, 2, 3}, 5)), "[1 2 3]")

	asserts.Equals(t, "works well with map",
		fmt.Sprint( Map([]T{[]T{1, 2, 3}, []T{1, 2, 3}}, func(obj T, idx T, list T) T { return Initial(obj.([]T), 2) })), "[[1 2] [1 2]]")
	asserts.Equals(t, "works well with nil",
		fmt.Sprint( Initial(nil)), "[]")
}

func TestLast(t *testing.T) {
	asserts.Equals(t, "can pull out the initial elements of an array",
		fmt.Sprint( Last([]T{1, 2, 3})),
		"[3]")
	asserts.Equals(t, "can pull out the last 2 element of an array",
		fmt.Sprint( Last([]T{1, 2, 3}, 2)), "[2 3]")

	asserts.Equals(t, "can pull out the zeroth element of an array",
		fmt.Sprint( Last([]T{1, 2, 3}, 0)), "[]")

	asserts.Equals(t, "can pull out the first element of an array",
		fmt.Sprint( Last([]T{1, 2, 3}, 1)), "[3]")

	asserts.Equals(t, "can pull out all elements of an array",
		fmt.Sprint( Last([]T{1, 2, 3}, 3)), "[1 2 3]")

	asserts.Equals(t, "can pull out too many out",
		fmt.Sprint( Last([]T{1, 2, 3}, 4)), "[1 2 3]")

	asserts.Equals(t, "can pull out too many items out",
		fmt.Sprint( Last([]T{1, 2, 3}, 5)), "[1 2 3]")

	asserts.Equals(t, "works well with map",
		fmt.Sprint( Map([]T{[]T{1, 2, 3}, []T{1, 2, 3}}, func(obj T, idx T, list T) T { return Last(obj.([]T), 2) })), "[[2 3] [2 3]]")

	asserts.Equals(t, "works well with nil",
		fmt.Sprint( Initial(nil)), "[]")
}

func TestCompact(t *testing.T) {
	asserts.Equals(t, "can trim out all falsy values",
		fmt.Sprint( Compact([]T{0, 1, false, 2, false, 3})), "[1 2 3]")
	asserts.Equals(t, "can trim out all falsy values",
		fmt.Sprint( Compact(nil)), "[]")
}

func TestFlatten(t *testing.T) {
	list := []T{1, []T{2}, []T{3, []T{[]T{[]T{4}}}}}
	asserts.Equals(t, "can flatten nested arrays", fmt.Sprint( Flatten(list)), "[1 2 3 4]")
	asserts.Equals(t, "can shallowly flatten nested arrays", fmt.Sprint( Flatten(list, true)),
		"[1 2 3 [[[4]]]]")
	list2 := []T{[]T{1}, []T{2}, []T{3}, []T{[]T{4}}}
	asserts.Equals(t, "can shallowly flatten arrays containing only other arrays",
		fmt.Sprint( Flatten(list2, true)), "[1 2 3 [4]]")
	asserts.Equals(t, "can flatten arrays containing only other arrays",
		fmt.Sprint( Flatten(list2)), "[1 2 3 4]")

	list3 := []T{[]T{"a", "b", "c"}, []T{"d", "e", "f", "g"}}
	asserts.Equals(t, "can flatten arrays containing only other arrays",
		fmt.Sprint( Flatten(list3)), "[a b c d e f g]")
	asserts.Equals(t, "can shallow flatten arrays containing only other arrays",
		fmt.Sprint( Flatten(list3, true)), "[a b c d e f g]")
}

func TestWithout(t *testing.T) {

	list := []T{1, 2, 1, 0, 3, 1, 4}
	asserts.Equals(t, "can remove arbitrary ints from a list",
		fmt.Sprint( Without(list, 0, 1)), "[2 3 4]")

	asserts.Equals(t, "can remove all instances from an array of args from a list",
		fmt.Sprint( Without(list, []T{0, 1})), "[2 3 4]")

	list2 := []T{1, 2, "1", 0, 3, 1, 4}
	asserts.Equals(t, "can remove all instances of func args from an object",
		fmt.Sprint( Without(list2, 0, 1)), "[2 1 3 4]")

	asserts.Equals(t, "can remove all instances from an array of args from an object",
		fmt.Sprint( Without(list2, []T{0, 1})), "[2 1 3 4]")

	asserts.Equals(t, "can remove all instances from an array of args from an object",
		fmt.Sprint( Without(list2, []T{0, "1"})), "[1 2 3 1 4]")

	asserts.Equals(t, "can remove all instances from an array of args from an object",
		fmt.Sprint( Without(list2, []T{0, "1", 1})), "[2 3 4]")

	/* Kludgy comparator?, go can't compare maps via identity, so stringifying each */
	listm := []T{map[T]T{"one": 1}, map[T]T{"two": 2}}
	retval := Without(listm, map[T]T{"one": 1}, func(a, b T) bool {
		return fmt.Sprint( a) == fmt.Sprint( b)
	})
	asserts.Equals(t, "uses real object identity for comparisons.", fmt.Sprint( retval), "[map[two:2]]")
	asserts.True(t, "ditto", len(retval) == 1)
}

func TestUniq(t *testing.T) {
	list := []T{1, 2, 1, 3, 1, 4}
	asserts.Equals(t, "can find the unique values of an unsorted array",
		fmt.Sprint( Uniq(list, false)), "[1 2 3 4]")

	list2 := []T{1, 1, 1, 2, 2, 3}
	asserts.Equals(t, "can find the unique values of a sorted array faster",
		fmt.Sprint( Uniq(list2, true)), "[1 2 3]")

	list3 := []map[T]T{{"name": "moe"}, {"name": "curly"}, {"name": "larry"}, {"name": "curly"}}
	iterator := func(value T, key T, list T) T { return value.(map[T]T)["name"] }
	comparator := func(a T, b T) bool { return a.(map[T]T)["name"] == b.(map[T]T)["name"] }
	asserts.Equals(t, "can find the unique values of an array using a custom iterator",
		fmt.Sprint( Uniq(list3, false, iterator, comparator)), "[map[name:moe] map[name:curly] map[name:larry]]")
	asserts.Equals(t, "can find the unique values of an array using a custom iterator",
		fmt.Sprint( Map(Uniq(list3, false, iterator, comparator), iterator)), "[moe curly larry]")

}

func TestIntersection(t *testing.T) {
	strLessThan := func(this T, that T) bool {
		return this.(string) < that.(string)
	}
	stooges := []T{"moe", "curly", "larry"}
	leaders := []T{"moe", "groucho"}
	asserts.Equals(t, "can take the set intersection of two arrays",
		fmt.Sprint( Intersection(strLessThan, stooges, leaders)), "[moe]")

	theSixStooges := []T{"moe", "moe", "curly", "curly", "larry", "larry"}
	asserts.Equals(t, "returns a duplicate-free array",
		fmt.Sprint( Intersection(strLessThan, theSixStooges, leaders)), "[moe]")
}

func TestUnion(t *testing.T) {
	result := Union([]T{1, 2, 3}, []T{2, 30, 1}, []T{1, 40})
	asserts.Equals(t, "takes the union of a list of arrays", fmt.Sprint( result), "[1 2 3 30 40]")

	result2 := Union([]T{1, 2, 3}, []T{2, 30, 1}, []T{1, 40, []T{1}})
	asserts.Equals(t, "takes the union of a list of nested arrays", fmt.Sprint( result2), "[1 2 3 30 40 [1]]")

	result3 := Union(nil, []T{1, 2, 3})
	asserts.Equals(t, "takes the union of nil and a list of arrays", fmt.Sprint( result3), "[<nil> 1 2 3]")
}

func TestDifference(t *testing.T) {
	result := Difference([]T{1, 2, 3}, IdentityComparator, []T{2, 30, 40})
	asserts.Equals(t, "takes the difference of two arrays",
		fmt.Sprint( result), "[1 3]")

	result2 := Difference([]T{1, 2, 3, 4}, IdentityComparator, []T{2, 30, 40}, []T{1, 11, 111})
	asserts.Equals(t, "takes the difference of three arrays",
		fmt.Sprint( result2), "[3 4]")
}

func TestZip(t *testing.T) {

	names := []T{"moe", "larry", "curly"}
	ages := []T{30, 40, 50}
	leaders := []T{true}

	stooges := Zip(names, ages, leaders)
	asserts.Equals(t, "zipped together arrays of different lengths",
		fmt.Sprint( stooges), "[[moe 30 true] [larry 40 <nil>] [curly 50 <nil>]]")

	stooges2 := Zip([]T{"moe", 30, "stooge 1"}, []T{"larry", 40, "stooge 2"}, []T{"curly", 50, "stooge 3"})
	asserts.Equals(t, "zipped pairs",
		fmt.Sprint( stooges2), "[[moe larry curly] [30 40 50] [stooge 1 stooge 2 stooge 3]]")

	// In the case of difference lengths of the tuples undefineds
	// should be used as placeholder
	stooges3 := Zip([]T{"moe", 30}, []T{"larry", 40}, []T{"curly", 50, "extra data"})
	asserts.Equals(t, "zipped pairs with empties",
		fmt.Sprint( stooges3), "[[moe larry curly] [30 40 50] [<nil> <nil> extra data]]")
	empty := Zip([]T{})
	asserts.Equals(t, "unzipped empty", fmt.Sprint( empty), "[]")

	empty2 := Zip([]T{})
	asserts.Equals(t, "unzipped empty2", fmt.Sprint( empty2), "[]")
}

func TestObject(t *testing.T) {
	result := Object([]T{"moe", "larry", "curly"}, []T{30, 40, 50})
	shouldBe := map[T]T{"moe": 30, "larry": 40, "curly": 50}
	asserts.Equals(t, "two arrays zipped together into an object",
		fmt.Sprint( result), fmt.Sprint( shouldBe))

	result2 := Object([]T{"one", 1, "two", 2, "three", 3})
	shouldBe2 := map[T]T{"one": 1, "two": 2, "three": 3}
	asserts.Equals(t, "an array of pairs zipped together into an object",
		fmt.Sprint( result2), fmt.Sprint( shouldBe2))

	result3 := Object([]T{[]T{"one", 1}, []T{"two", 2}, []T{"three", 3}})
	shouldBe3 := map[T]T{"one": 1, "two": 2, "three": 3}
	asserts.Equals(t, "an array of pairs zipped together into an object",
		fmt.Sprint( result3), fmt.Sprint( shouldBe3))

	asserts.Nil(t, "handles nils", Object(nil))

	stooges := map[T]T{"moe": 30, "larry": 40, "curly": 50}
	asserts.Equals(t, "an object converted to pairs and back to an object",
		fmt.Sprint( stooges), fmt.Sprint( Object(Pairs(stooges))))
}

func TestIndexOf(t *testing.T) {
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

func TestLastIndexOf(t *testing.T) {
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

func TestRange(t *testing.T) {
	asserts.Equals(t, "range with no arguments generates an empty array",
		fmt.Sprint( Range()), "[]")
	asserts.Equals(t, "range with 0 as a first argument generates an empty array",
		fmt.Sprint( Range(0)), "[]")
	asserts.Equals(t, "range with a single positive argument generates an array of elements 0,1,2,...,n-1",
		fmt.Sprint( Range(4)), "[0 1 2 3]")
	asserts.Equals(t, "range with two arguments a &amp; b, a&lt;b generates an array of elements a,a+1,a+2,...,b-2,b-1",
		fmt.Sprint( Range(5, 8)), "[5 6 7]")
	asserts.Equals(t, "range with two arguments a &amp; b, b&lt;a generates an empty array",
		fmt.Sprint( Range(8, 5)), "[]")
	asserts.Equals(t, "range with three arguments a &amp; b &amp; c, c &lt; b-a, a &lt; b generates an array of elements a,a+c,a+2c,...,b - (multiplier of a) &lt; c",
		fmt.Sprint( Range(3, 10, 3)), "[3 6 9]")
	asserts.Equals(t, "range with three arguments a &amp; b &amp; c, c &gt; b-a, a &lt; b generates an array with a single element, equal to a",
		fmt.Sprint( Range(3, 10, 15)), "[3]")
	asserts.Equals(t, "range with three arguments a &amp; b &amp; c, a &gt; b, c &lt; 0 generates an array of elements a,a-c,a-2c and ends with the number not less than b",
		fmt.Sprint( Range(12, 7, -2)), "[12 10 8]")
	asserts.Equals(t, "final example in the Python docs",
		fmt.Sprint( Range(0, -10, -1)), "[0 -1 -2 -3 -4 -5 -6 -7 -8 -9]")

}
