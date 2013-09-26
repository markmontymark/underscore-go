package underscore

import (
	"./lib/asserts"
	"fmt"
	"sort"
	"testing"
)

type IntSlice []int

func (this IntSlice) Len() int           { return len(this) }
func (this IntSlice) Less(a, b int) bool { return this[a] < this[b] }
func (this IntSlice) Swap(a, b int)      { this[a], this[b] = this[b], this[a] }

func TestEach(t *testing.T) {

	sliceCollector := make([]string, 0)
	silly := func(elem T, i T, list T) bool {
		sliceCollector = append(sliceCollector, fmt.Sprintf("each: elem %d, index %v, list %v,", elem, i, list))
		return EachContinue
	}

	Each([]T{0, 3, 2, 1}, silly)

	asserts.Equals(t, "Test Each(0,3,2,4)", "[each: elem 0, index 0, list [0 3 2 1], each: elem 3, index 1, list [0 3 2 1], each: elem 2, index 2, list [0 3 2 1], each: elem 1, index 3, list [0 3 2 1],]", fmt.Sprintf("%v", sliceCollector))
}

func TestMap(t *testing.T) {

	add2 := func(elem T, i T, list T) T {
		return elem.(int) + 2
	}

	identityValueMap := func(value T, key T, obj T) T { return value }
	identityKeyMap := func(value T, key T, obj T) T { return key }
	identityNestedKeyMap := func(value T, key T, obj T) T { return value.(map[T]T)["a"] }
	identityNestedKeyShorterMap := func(value T, key T, obj T) T { return value.(map[T]T)["b"] }

	asserts.Equals(t, "testing map with array",
		"[2 5 4 3]", fmt.Sprintf("%v", Map([]T{0, 3, 2, 1}, add2)))
	asserts.Equals(t, "testing collect with array",
		"[2 5 4 3]", fmt.Sprintf("%v", Collect([]T{0, 3, 2, 1}, add2)))

	amap := map[T]T{"a": 3, "b": 2, "c": 1}
	list := make([]T, 0)
	list = append(list, amap)
	list = append(list, map[T]T{"a": 1, "d": 4, "e": 5})

	asserts.Equals(t, "testing collectmap with map[string]int ",
		"[3 2 1]",
		fmt.Sprintf("%v", Collect(amap, identityValueMap)))

	asserts.Equals(t, "testing collectmap with map[string]int ",
		"[map[a:3 b:2 c:1] map[a:1 d:4 e:5]]",
		fmt.Sprintf("%v", Collect(list, identityValueMap)))

	asserts.Equals(t, "testing collectmap with map[string]int ",
		"[a b c]",
		fmt.Sprintf("%v", Collect(amap, identityKeyMap)))

	asserts.Equals(t, "testing collectmap with map[string]int ",
		"[0 1]",
		fmt.Sprintf("%v", Collect(list, identityKeyMap)))

	asserts.Equals(t, "testing collectmap with map[string]int ",
		"[3 1]",
		fmt.Sprintf("%v", Collect(list, identityNestedKeyMap)))

	asserts.Equals(t, "testing collectmap with map[string]int ",
		"[2]",
		fmt.Sprintf("%v", Collect(list, identityNestedKeyShorterMap)))
}

func TestReduce(t *testing.T) {

	v, err := Reduce(
		[]T{1, 2, 3},
		func(sum T, num T, i T, list T) T { return sum.(int) + num.(int) },
		0)
	if err != "" {
		fmt.Printf("FAIL: %s\n", err)
		return
	}
	asserts.IntEquals(t, "can reduce sum up an array", 6, v.(int))

	v, err = Reduce(
		[]T{1, 2, 3},
		func(sum T, num T, i T, list T) T { return sum.(int) * num.(int) },
		3)
	if err != "" {
		fmt.Printf("FAIL: %s\n", err)
		return
	}
	asserts.IntEquals(t, "can reduce multiply up an array", 18, v.(int))

	v, err = Inject(
		[]T{1, 2, 3},
		func(sum T, num T, i T, list T) T { return sum.(int) * num.(int) },
		3)
	if err != "" {
		fmt.Printf("FAIL: %s\n", err)
		return
	}
	asserts.IntEquals(t, "can inject multiply up an array", 18, v.(int))

	v, err = Inject(
		[]T{1, 2, 3},
		func(sum T, num T, i T, list T) T { return sum.(int) * num.(int) },
		3)
	if err != "" {
		fmt.Printf("FAIL: %s\n", err)
		return
	}
	asserts.IntEquals(t, "can foldl multiply up an array", 18, v.(int))

}

func TestReduceRight(t *testing.T) {

	v, err := ReduceRight(
		[]T{"2", "3", "4"},
		func(sum T, num T, i T, list T) T { return sum.(string) + "," + num.(string) },
		"")
	if err != "" {
		fmt.Printf("FAIL: %s\n", err)
		return
	}
	asserts.Equals(t, "can ReduceRight divide up an array", ",4,3,2", v.(string))
}

func TestFind(t *testing.T) {
	array := []T{1, 2, 3, 4}
	v := Find(array, func(n T, i T, list T) bool {
		return n.(int) > 2
	})
	asserts.IntEquals(t, "should return first found `value`", 3, v.(int))

	v = Find(array, func(n T, i T, list T) bool { return false })
	asserts.Nil(t, "should return `nil` if `value` is not found", v)
}
func TestDetect_asFindAlias(t *testing.T) {
	array := []T{1, 2, 3, 4}
	v := Detect(array, func(n T, i T, list T) bool {
		return n.(int) > 2
	})
	asserts.IntEquals(t, "should return first found `value`", 3, v.(int))

	v = Detect(array, func(n T, i T, list T) bool { return false })
	asserts.Nil(t, "should return `nil` if `value` is not found", v)
}

func TestFilter(t *testing.T) {
	array := []T{1, 2, 3, 4}
	v := Filter(array, func(n T, i T, list T) bool {
		return n.(int) > 2
	})
	asserts.Equals(t, "should return last two values: 3 4", fmt.Sprintf("%v", []T{3, 4}), fmt.Sprintf("%v", v))
}
func TestSelect_asFilterAlias(t *testing.T) {
	array := []T{1, 2, 3, 4}
	v := Select(array, func(n T, i T, list T) bool {
		return n.(int) > 2
	})
	asserts.Equals(t, "should return last two values: 3 4", fmt.Sprintf("%v", []T{3, 4}), fmt.Sprintf("%v", v))
}

func TestReject(t *testing.T) {
	array := []T{1, 2, 3, 4}
	v := Reject(array, func(n T, i T, list T) bool {
		return n.(int) > 2
	})
	asserts.Equals(t, "should return first two values: 1 2", fmt.Sprintf("%v", []T{1, 2}), fmt.Sprintf("%v", v))
}

func TestEvery(t *testing.T) {
	array := []T{1, 2, 3, 4}
	v := Every(array, func(n T, i T, list T) bool {
		return n.(int) < 5
	})
	asserts.Equals(t, "should return true as all values: 1 2 3 4 are less than 5", "true", fmt.Sprintf("%v", v))

	v = Every(array, func(n T, i T, list T) bool {
		return n.(int) < 4
	})
	asserts.Equals(t, "should return false as not all values are < 4", "false", fmt.Sprintf("%v", v))
}

func TestAll(t *testing.T) {
	array := []T{1, 2, 3, 4}
	v := All(array, func(n T, i T, list T) bool {
		return n.(int) < 5
	})
	asserts.Equals(t, "should return true as all values: 1 2 3 4 are less than 5", "true", fmt.Sprintf("%v", v))

	v = All(array, func(n T, i T, list T) bool {
		return n.(int) < 4
	})
	asserts.Equals(t, "should return false as not all values are < 4", "false", fmt.Sprintf("%v", v))
}

func TestAny(t *testing.T) {
	array := []T{1, 2, 3, 4}
	v := Any(array, func(n T, i T, list T) bool {
		return n.(int) < 5
	})
	asserts.Equals(t, "should return true as at least one value is less than 5", "true", fmt.Sprintf("%v", v))

	v = Any(array, func(n T, i T, list T) bool {
		return n.(int) < 4
	})
	asserts.Equals(t, "should return true as at least one value is less than 4", "true", fmt.Sprintf("%v", v))

	v = Any(array, func(n T, i T, list T) bool {
		return n.(int) < 3
	})
	asserts.Equals(t, "should return true as at least one value is less than 3", "true", fmt.Sprintf("%v", v))

	v = Any(array, func(n T, i T, list T) bool {
		return n.(int) < 2
	})
	asserts.Equals(t, "should return true as at least one value is less than 2", "true", fmt.Sprintf("%v", v))

	v = Any(array, func(n T, i T, list T) bool {
		return n.(int) < 1
	})
	asserts.Equals(t, "should return false as no value are less than 1", "false", fmt.Sprintf("%v", v))
}

func TestInclude(t *testing.T) {
	array := []T{1, 2, 3, 4}
	v := Include(array, 1)
	asserts.Equals(t, "should return true as array contains a 1", "true", fmt.Sprintf("%v", v))

	v = Include(array, 2)
	asserts.Equals(t, "should return true as array contains a 1", "true", fmt.Sprintf("%v", v))

	v = Include(array, 3)
	asserts.Equals(t, "should return true as array contains a 1", "true", fmt.Sprintf("%v", v))

	v = Include(array, 4)
	asserts.Equals(t, "should return true as array contains a 1", "true", fmt.Sprintf("%v", v))

	v = Include(array, 5)
	asserts.Equals(t, "should return false as array doesnt contain a 5", "false", fmt.Sprintf("%v", v))
}

func TestInvoke(t *testing.T) {
	list := []T{IntSlice{5, 1, 7}, IntSlice{3, 2, 1}}
	result := Invoke(list, func(item T, args ...T) T {
		sort.Sort(item.(IntSlice))
		return item
	})
	asserts.Equals(t, "each array sorted", fmt.Sprintf("%v", result),
		"[[1 5 7] [1 2 3]]")
}

func TestPluck(t *testing.T) {
	people := []T{"name", "moe", "age", 30, "name", "curly", "age", 50}
	v := Pluck(people, "name")
	asserts.Equals(t, "pulls names out of objects",
		"[name name]",
		fmt.Sprintf("%v", v))

	v = Pluck(people, 30)
	asserts.Equals(t, "pulls 30 out of list",
		"[30]",
		fmt.Sprintf("%v", v))

	people = make([]T, 0)
	people = append(people, map[T]T{"name": "moe", "age": 30})
	people = append(people, map[T]T{"name": "curly", "age": 50})
	v = Pluck(people, "name")
	asserts.Equals(t, "pulls names out of objects",
		"[moe curly]",
		fmt.Sprintf("%v", v))
}

func TestWhere(t *testing.T) {

	list := make([]T, 0)
	list = append(list, map[T]T{"a": 1, "b": 2})
	list = append(list, map[T]T{"a": 2, "b": 2})
	list = append(list, map[T]T{"a": 1, "b": 3})
	list = append(list, map[T]T{"a": 1, "b": 4})

	v := Where(list, map[T]T{"a": 1})
	asserts.Equals(t, "Find objects with key a:1", "3", fmt.Sprintf("%v", len(v.([]T))))
	asserts.Equals(t, "Last found has a b:4", "4", fmt.Sprintf("%v", v.([]T)[len(v.([]T))-1].(map[T]T)["b"]))

	v = Where(list, map[T]T{"b": 2})
	asserts.Equals(t, "Find objects with b:2", "2", fmt.Sprintf("%v", len(v.([]T))))

	v = Where(list, map[T]T{"b": 2}, true)
	asserts.Equals(t, "Find objects with b:2", "map[a:1 b:2]", fmt.Sprintf("%v", v))
}

func TestWhereOOP(t *testing.T) {
/*
	list := make([]T, 0)
	list = append(list, map[T]T{"a": 1, "b": 2})
	list = append(list, map[T]T{"a": 2, "b": 2})
	list = append(list, map[T]T{"a": 1, "b": 3})
	list = append(list, map[T]T{"a": 1, "b": 4})
*/
	list := []T{ map[T]T{"a": 1, "b": 2}, map[T]T{"a": 2, "b": 2}, map[T]T{"a": 1, "b": 3}, map[T]T{"a": 1, "b": 4} }

	v := New(list).Chain().Where(map[T]T{"a": 1}).Value()
	asserts.Equals(t, "Find objects with key a:1", "3", fmt.Sprintf("%v", len(v.([]T))))
	asserts.Equals(t, "Last found has a b:4", "4", fmt.Sprintf("%v", v.([]T)[len(v.([]T))-1].(map[T]T)["b"]))

	v2 := New(list).Chain().Where(map[T]T{"b": 2}).Value()
	asserts.Equals(t, "Find objects with b:2", "2", fmt.Sprintf("%v", len(v2.([]T))))

	v3 := New(list).Chain().Where(map[T]T{"b": 2}, true).Value()
	asserts.Equals(t, "Find first object with b:2", "map[a:1 b:2]", fmt.Sprintf("%v", v3))

	v4 := New(list).Chain().Where(map[T]T{"b": 2}, true).Value()
	asserts.Equals(t, "Find first object with b:2 when chained", "map[a:1 b:2]", fmt.Sprintf("%v", v4))
}

func TestFindWhere(t *testing.T) {
	list := make([]T, 0)
	list = append(list, map[T]T{"a": 1, "b": 2})
	list = append(list, map[T]T{"a": 2, "b": 2})
	list = append(list, map[T]T{"a": 1, "b": 3})
	list = append(list, map[T]T{"a": 1, "b": 4})
	v := FindWhere(list, map[T]T{"a": 1})
	asserts.Equals(t, "Find first object with key a:1", "map[a:1 b:2]", fmt.Sprintf("%v", v))
}

func TestMax(t *testing.T) {
	list := []int{2, 3, 4, 9, 5, 6, 7, 8}
	asserts.Equals(t, "Find max element in array", "9", fmt.Sprintf("%v", MaxInt(list...)))
}

func TestMin(t *testing.T) {
	list := []int{2, 3, 4, 9, 5, 6, 7, 8}
	asserts.Equals(t, "Find min element in array", "2", fmt.Sprintf("%v", MinInt(list...)))
}

func TestSortBy(t *testing.T) {
	people := []map[T]T{{"name": "curly", "age": 50}, {"name": "moe", "age": 30}}
	peopleSorted := SortBy(people,
		func(obj, b, c T) T { return obj.(map[T]T)["age"] },
		func(a, b *map[T]T) bool {
			return (*a)["criteria"].(int) < (*b)["criteria"].(int)
		})
	asserts.Equals(t, "stooges sorted by age, plucking just 'name'",
		fmt.Sprintf("%v", Pluck(peopleSorted, "name")), "[moe curly]")

	list := []T{nil, 4, 1, nil, 3, 2}
	asserts.Equals(t, "SortBy with nil values",
		fmt.Sprintf("%v",
			SortBySorter(list, Identity, func(a, b *map[T]T) bool {
				return (*a)["criteria"].(int) < (*b)["criteria"].(int)
			})),
		"[1 2 3 4]")

	words := []T{"one", "two", "three", "four", "five"}
	sorted := SortBySorter(words,
		func(obj, b, c T) T { return len(obj.(string)) },
		func(a, b *map[T]T) bool { return (*a)["criteria"].(int) < (*b)["criteria"].(int) })
	asserts.Equals(t, "sorted by length", fmt.Sprintf("%v", sorted),
		"[one two four five three]")

	type Pair struct {
		x, y int
	}

	collection := []T{
		&Pair{1, 1}, &Pair{1, 2},
		&Pair{1, 3}, &Pair{1, 4},
		&Pair{1, 5}, &Pair{1, 6},
		&Pair{2, 1}, &Pair{2, 2},
		&Pair{2, 3}, &Pair{2, 4},
		&Pair{2, 5}, &Pair{2, 6},
		&Pair{3, 1}, &Pair{3, 2},
		&Pair{3, 3}, &Pair{3, 4},
		&Pair{3, 5}, &Pair{3, 6},
	}

	actual := SortBySorter(collection,
		Identity,
		func(a, b *map[T]T) bool {
			if (*a)["criteria"].(*Pair).x == (*b)["criteria"].(*Pair).x {
				return (*a)["criteria"].(*Pair).y < (*b)["criteria"].(*Pair).y
			}
			return (*a)["criteria"].(*Pair).x < (*b)["criteria"].(*Pair).x
		})

	asserts.Equals(t, "sortby should be stable",
		fmt.Sprintf("%v", actual), fmt.Sprintf("%v", collection))
}

func TestGroupBy(t *testing.T) {

	data := GroupBy([]T{1, 2, 3, 4, 5, 6, 1}, func(obj, key, val T) T {
		//fmt.Printf("group by func got obj %v, key %v, val %v\n", obj,key,val)
		//fmt.Printf("group by func, returning %v\n", (val.(int) % 2) )
		return obj.(int) % 2
	})
	asserts.Equals(t, "group ints ", "parity map[1:[1 3 5 1] 0:[2 4 6]]",
		fmt.Sprintf("parity %v", data))
	asserts.Equals(t, "group evens ", "[2 4 6]", fmt.Sprintf("%v", data[0]))

	data2 := []T{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
	grouped := GroupBy(data2, func(obj T, key T, val T) T { return len(obj.(string)) })

	asserts.Equals(t, "grouping words of length 3",
		fmt.Sprintf("%v", grouped[3]), "[one two six ten]")
	asserts.Equals(t, "grouping words of length 4",
		fmt.Sprintf("%v", grouped[4]), "[four five nine]")
	asserts.Equals(t, "grouping words of length 5",
		fmt.Sprintf("%v", grouped[5]), "[three seven eight]")

	data3 := []map[T]T{{"a": 1, "b": 2}, {"b": 3}, {"a": 4, "c": 5}, {"a": 1, "b": 7, "c": 8}}
	grouped3 := GroupBy(data3, func(obj T, key T, val T) T {
		return obj.(map[T]T)["a"]
	})
	asserts.Equals(t, "group by an object keys value", "map[1:[map[a:1 b:2] map[a:1 b:7 c:8]] 4:[map[a:4 c:5]]]", fmt.Sprintf("%v", grouped3))
}

//TODO add more 'indexBy' tests from collections.js
func TestIndexBy(t *testing.T) {
	data := []map[T]T{{"a": 1, "b": 2}, {"b": 3}, {"a": 4, "c": 5}, {"a": 1, "b": 7, "c": 8}}
	grouped := IndexBy(data, func(obj T, key T, val T) T {
		return obj.(map[T]T)["a"]
	})
	asserts.Equals(t, "index by an object keys value", "map[1:map[a:1 b:7 c:8] 4:map[a:4 c:5]]", fmt.Sprintf("%v", grouped))
}

//TODO add more 'countBy' tests from collections.js
func TestCountBy(t *testing.T) {
	data := []map[T]T{{"a": 1, "b": 2}, {"b": 3}, {"a": 4, "c": 5}, {"a": 1, "b": 7, "c": 8}}
	grouped := CountBy(data, func(obj T, key T, val T) T {
		return obj.(map[T]T)["a"]
	})
	asserts.Equals(t, "count by an object keys value", "map[1:2 4:1]", fmt.Sprintf("%v", grouped))
}

func TestSortedIndex(t *testing.T) {

	intLessThan := func(a T, b T) bool { return a.(int) < b.(int) }
	numbers := []T{10, 20, 30, 40, 50}
	num := 35
	indexForNum := SortedIndex(numbers, num, intLessThan)
	asserts.IntEquals(t, "35 should be inserted at index 3", 3, indexForNum)

	indexFor30 := SortedIndex(numbers, 30, intLessThan)
	asserts.IntEquals(t, "30 should be inserted at index 2", 2, indexFor30)

	objects := []map[T]T{{"x": 10}, {"x": 20}, {"x": 30}, {"x": 40}}
	iterator := func(obj T, idx T, list T) T { return obj.(map[T]T)["x"] }
	asserts.IntEquals(t, "sorted index with object list", 2,
		SortedIndex(objects, map[T]T{"x": 25}, intLessThan, iterator))
	asserts.IntEquals(t, "sorted index with object list, take 2", 3,
		SortedIndex(objects, map[T]T{"x": 35}, intLessThan, iterator))
}

func TestShuffle(t *testing.T) {
	list := []T{2, 3, 4, 9, 5, 6, 7, 8}
	shuffledlist := Shuffle(list)
	asserts.IntEquals(t, "Find max element in array",
		Max(intLessThan, list...).(int),
		Max(intLessThan, shuffledlist...).(int))

	asserts.True(t, "sort orig list and shuffled list",
		Every(list, func(item, idx, list T) bool {
			return Contains(shuffledlist, item)
		}))
}

func TestSample(t *testing.T) {
	numbers := Range(10)
	all_sampled := Sample(numbers, 10)
	asserts.True(t, "contains the same members before and after sample, take 1",
		Every(all_sampled, func(val, idx, list T) bool {
			return Contains(numbers, val)
		}))
	asserts.True(t, "contains the same members before and after sample ,take 2",
		Every(numbers, func(val, idx, list T) bool {
			return Contains(all_sampled.([]T), val)
		}))

	all_sampled2 := Sample(numbers, 20)
	asserts.True(t, "also works when sampling more objects than are present, take 1",
		Every(all_sampled2, func(val, idx, list T) bool {
			return Contains(numbers, val)
		}))
	asserts.True(t, "also works when sampling more objects than are present,take 2",
		Every(numbers, func(val, idx, list T) bool {
			return Contains(all_sampled2.([]T), val)
		}))

	/*
	   ok(_.contains(numbers, _.sample(numbers)), 'sampling a single element returns something from the array');
	   strictEqual(_.sample([]), undefined, 'sampling empty array with no number returns undefined');
	   notStrictEqual(_.sample([], 5), [], 'sampling empty array with a number returns an empty array');
	   notStrictEqual(_.sample([1, 2, 3], 0), [], 'sampling an array with 0 picks returns an empty array');
	   deepEqual(_.sample([1, 2], -1), [], 'sampling a negative number of picks returns an empty array');
	   ok(_.contains([1, 2, 3], _.sample({a: 1, b: 2, c: 3})), 'sample one value from an object');
	*/
}

func TestToArray(t *testing.T) {
	a := []T{1, 2, 3}
	asserts.Equals(t, "Clone an array", fmt.Sprintf("%v", a), fmt.Sprintf("%v", ToArray(a)))
	b := map[T]T{"one": 1, "two": 2, "three": 3}
	numbers := ToArray(b)
	asserts.Equals(t, "object flattened into array", "[1 2 3]", fmt.Sprintf("%v", numbers))
}

func TestSize(t *testing.T) {
	data := Size(map[T]T{"a": 1, "b": 4, "c": 6})
	asserts.IntEquals(t, "size of a map", 3, data)

	data1 := Size([]map[T]T{{"a": 1, "b": 4, "c": 6}})
	asserts.IntEquals(t, "size of a map", 1, data1)

	data10 := Size([]map[T]T{})
	asserts.IntEquals(t, "size of a map", 0, data10)

	data2 := Size([]T{"a", "b", "c", 1})
	asserts.IntEquals(t, "size of a varied list ", 4, data2)

	data20 := Size([]T{})
	asserts.IntEquals(t, "size of an empty list ", 0, data20)
}
