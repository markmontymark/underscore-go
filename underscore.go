// ported from Underscore.js version 1.5.2, here's its copyright notice for attribution's sake
//
//     Underscore.js 1.5.2
//     http://underscorejs.org
//     (c) 2009-2013 Jeremy Ashkenas, DocumentCloud and Investigative Reporters & Editors
//     Underscore may be freely distributed under the MIT license.
//
package underscore

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
_	"os"
)

// Custom type T is a shorthand for interface{} to save myself from RMI or carpal tunnel syndrome
type T interface{}

// An interface that provides a Value() for use by the Underscore custom type for stringifying its this.wrapped value
type Valuer interface {
	Value() T
}

// An Underscore custom type for OOP-style usage, as opposed to functional use.
// OOP-style example: list:= []T{"i",10,5.3}; fmt.Sprint( New(list).Filter(       func(v T,b T,c T) bool { _,ok:=v.(int);return ok}) ) -> [5]
// vs Functional    : list:= []T{"i",10,5.3}; fmt.Sprint(           Filter( list, func(v T,b T,c T) bool { _,ok:=v.(int);return ok} ) ) -> [5]
type Underscore struct {
	ischained bool
	wrapped   T
	Valuer
	fmt.Stringer
}

// Stringifying when using OOP-style, ie fmt.Sprintf("%v",New([]{"i",1,2})) -> [i 1 2]
// will take whatever is wrapped and pass that to the %v fmt printf specifier
func (this *Underscore) String() string {
	return fmt.Sprint(this.wrapped)
}

// Not just a Yeah Yeah Yeahs song, but also a shorthand for a list of maps
type maps []map[T]T

// Implement Len to allow maps custom type to be sortable with sort.Sort()
func (this maps) Len() int { return len(this) }

// Implement Len to allow maps custom type to be sortable with sort.Sort()
//func (this maps) Less ( a , b int) bool { return this[a]["criteria"] < this[b]["criteria"] }
func (this maps) Swap(a, b int) { this[a], this[b] = this[b], this[a] }

// Functions passed to Each need this signature
// - if what's passed to Each is map[T]T, the signature means: (val T, key T, obj map[T]T) for each key of obj
// - if what's passed to Each is []T, the signature means: (val T, index.(int) T, obj.([]T) T) for each key of obj
type eachlistiterator func(T, T, T) bool

// Functions passed to Map need this signature
type mapiterator func(T, T, T) T

// private constant for use by Each function
// returnning false from a function passed to Each means keep iterating
const eachContinue bool = false

// private constant for use by Each function
// returnning true from a function passed to Each means stop iterating
const eachBreak bool = true

// Create a eference to the Underscore object for use below.
// whatever you passed to New() gets saved in a member variable, wrapped
func New(obj ...T) *Underscore {
	if obj != nil {
		if _, ok := obj[0].(*Underscore); ok {
			return obj[0].(*Underscore)
		}
		un := new(Underscore)
		un.wrapped = obj[0]
		return un
	}
	return new(Underscore)
}

// Mirroring the version we started porting from, 1.5.2
const VERSION string = "1.5.2"

// Collection Functions
// --------------------

// The cornerstone, an `each` implementation, aka `forEach`.
// Handles objects and arrays
func Each(elemslist_or_map T, iterator eachlistiterator) {
	if elemslist_or_map == nil || IsEmpty(elemslist_or_map) {
		return
	}

	if IsArray(elemslist_or_map) {
		for i, elem := range elemslist_or_map.([]T) {
			if iterator(elem, i, elemslist_or_map.([]T)) == eachBreak {
				return
			}
		}

	} else if IsStringArray(elemslist_or_map) {
		for i, elem := range elemslist_or_map.([]string) {
			if iterator(elem, i, elemslist_or_map.([]string)) == eachBreak {
				return
			}
		}

	} else if IsArrayOfMaps(elemslist_or_map) {
		for i, elem := range elemslist_or_map.([]map[T]T) {
			if iterator(elem, i, elemslist_or_map.([]map[T]T)) == eachBreak {
				return
			}
		}

	} else if IsMap(elemslist_or_map) {
		for k, v := range elemslist_or_map.(map[T]T) {
			if iterator(v, k, elemslist_or_map.(map[T]T)) == eachBreak {
				return
			}
		}
	} else {
		fmt.Printf("Each isnt doing anything useful with first arg, %v\n", elemslist_or_map)
	}

}

// Return the results of applying an iterator to each element.
// Aliased as Collect
func Map(obj T, iterator func(T, T, T) T) []T {
	results := make([]T, 0)
	if obj == nil {
		return results
	}
	Each(obj, func(value T, index T, list T) bool {
		if v := iterator(value, index, list); v != nil {
			results = append(results, v)
		}
		return eachContinue
	})
	return results
}

// Return the results of applying an iterator to each element.
// Aka Map
var Collect func(obj T, iterator func(T, T, T) T) []T = Map

func mapMap(obj []map[T]T, iterator func(T, T, T) map[T]T) []map[T]T {
	results := make([]map[T]T, 0)
	if obj == nil {
		return results
	}
	Each(obj, func(value T, index T, list T) bool {
		if v := iterator(value, index, list); v != nil {
			results = append(results, v)
		}
		return eachContinue
	})
	return results
}
func mapForSortBy(obj T, iterator func(T, T, T) map[T]T) []map[T]T {
	results := make([]map[T]T, 0)
	if obj == nil {
		return results
	}
	Each(obj, func(value, index, list T) bool {
		if v := iterator(value, index, list); v != nil && v["criteria"] != nil {
			results = append(results, v)
		}
		return eachContinue
	})
	return results
}

// Error message for Reduce function, which is returned if you end up calling Reduce like this:
// Reduce( []T{},nil), which mirrors what Underscore.js did. I don't know -- I could see returning
// an empty list or nil and not considering either an error at the library level, user-level could
// deem []T{} or nil as "error" and do its own thing
const ReduceError = "Reduce of empty array with no initial value"

// **Reduce** builds up a single result from a list of values
// Aliased as `Inject`
// Aliased as `FoldL`
func Reduce(obj []T, iterator func(T, T, T, T) T, memo ...T) (T, string) {
	initial := len(memo) > 0
	if obj == nil {
		obj = make([]T, 0)
	}
	Each(obj, func(value T, index T, list T) bool {
		if !initial {
			memo[0] = value
			initial = true
		} else {
			memo[0] = iterator(memo[0], value, index, list)
		}
		return eachContinue
	})
	if !initial {
		return nil, ReduceError
	}
	return memo[0], ""
}

// **Inject** builds up a single result from a list of values
// Aliased as `Reduce`
// Aliased as `FoldL`
var Inject func(obj []T, iterator func(T, T, T, T) T, memo ...T) (T, string) = Reduce

// **FoldL** builds up a single result from a list of values
// Aliased as `Reduce`
// Aliased as `Inject`
var FoldL func(obj []T, iterator func(T, T, T, T) T, memo ...T) (T, string) = Reduce

// The right-associative version of reduce
// Aliased `FoldR`
func ReduceRight(obj []T, iterator func(T, T, T, T) T, memo ...T) (T, string) {
	initial := len(memo) > 0
	if obj == nil {
		obj = make([]T, 0)
	}
	length := len(obj)
	Each(obj, func(value T, index T, list T) bool {
		length = length - 1
		index = length
		if !initial {
			memo[0] = list.([]T)[index.(int)]
			initial = true
		} else {
			memo[0] = iterator(memo[0], list.([]T)[index.(int)], index, list)
		}
		return eachContinue
	})
	if !initial {
		return nil, ReduceError
	}
	return memo[0], ""
}

// The right-associative version of reduce
// Aliased as `ReduceRight`
var FoldR func(obj []T, iterator func(T, T, T, T) T, memo ...T) (T, string) = ReduceRight

// Return the first value which passes a truth test.
// Aliased as `Detect`.
func Find(obj []T, predicate func(T, T, T) bool) T {
	var result T
	Any(obj, func(value T, index T, list T) bool {
		if predicate(value, index, list) {
			result = value
			return eachBreak
		}
		return eachContinue
	})
	return result
}

// Return the first value which passes a truth test.
// Aliased as `Find`
var Detect func(obj []T, iterator func(T, T, T) bool) T = Find

// Return all the elements that pass a truth test.
// Aliased as `Select`.
func Filter(obj []T, iterator eachlistiterator) []T {
	results := make([]T, 0)
	if obj == nil || len(obj) == 0 {
		return results
	}
	Each(obj, func(value T, index T, list T) bool {
		if iterator(value, index, list) {
			results = append(results, value)
		}
		return eachContinue
	})
	return results
}

// Return all the elements that pass a truth test.
// Aliased as `Filter`.
var Select func(obj []T, iterator eachlistiterator) []T = Filter

// Return all the elements for which a truth test fails.
func Reject(obj []T, iterator eachlistiterator) []T {
	return Filter(obj, func(value T, index T, list T) bool {
		return !iterator(value, index, list)
	})
}

// Determine whether all of the elements match a truth test.
// Aliased as `All`
func Every(obj T, opt_iterator ...eachlistiterator) bool {
	var iterator eachlistiterator //func(T,int, []T)bool
	if len(opt_iterator) == 0 {
		iterator = IdentityEach
	} else {
		iterator = opt_iterator[0]
	}
	result := true
	if obj == nil {
		return result
	}
	Each(obj.([]T), func(value T, index T, list T) bool {
		result = result && iterator(value, index, list)
		if !result {
			return eachBreak
		}
		return eachContinue
	})
	return result
}

// Determine whether all of the elements match a truth test.
// Aliased as `Every`
var All func(obj T, opt_iterator ...eachlistiterator) bool = Every

// Determine if at least one element in the object matches a truth test.
// Aliased as `Some`.
func Any(obj T, opt_predicate ...func(val, index, list T) bool) bool {
	var predicate func(T, T, T) bool
	if len(opt_predicate) == 0 {
		predicate = IdentityEach
	} else {
		predicate = opt_predicate[0]
	}
	anyresult := false
	if obj == nil {
		return anyresult
	}

	Each(obj.([]T), func(val, index, list T) bool {
		if anyresult {
			return eachBreak
		}
		anyresult = predicate(val, index, list)
		if anyresult {
			return eachBreak
		}
		return eachContinue
	})
	return anyresult
}

// Determine if at least one element in the object matches a truth test.
// Aliased as `Any`.
var Some func(obj T, opt_predicate ...func(val, index, list T) bool) bool = Any

// Determine if the array or object contains a given value (using `==`).
// Aliased as `Include`.
func Contains(obj T, target T, opt_comparator ...func(T, T) bool) bool {
	if obj == nil {
		return false
	}
	var comparator func(T, T) bool
	if len(opt_comparator) > 0 {
		comparator = opt_comparator[0]
	}
	return Any(obj, func(value T, index T, list T) bool {
		if comparator != nil {
			return comparator(value, target)
		} else {
			return value == target
		}
	})
}

// Determine if the array or object contains a given value (using `==`).
// Aliased as `Contains`.
var Include func(obj T, target T, opt_comparator ...func(T, T) bool) bool = Contains

// Invoke a method (with arguments) on every item in a collection.
func Invoke(obj T, method func(this T, thisArgs ...T) T, args ...T) []T {
	//var isFunc = IsFunction(method);
	return Map(obj, func(value T, key T, origObj T) T {
		return method(value, args)
	})
}

// Convenience version of a common use case of `map`: fetching a property.
func Pluck(obj T, targetvalue T) []T {
	return Map(obj, func(testvalue T, index T, origlist T) T {
		if IsMap(testvalue) {
			return testvalue.(map[T]T)[targetvalue]
		}
		if targetvalue == testvalue {
			return testvalue
		}
		return nil
	})
}

// Convenience version of a common use case of `filter`: selecting only objects
// containing specific `key:value` pairs.
func Where(obj []T, attrs map[T]T, optReturnFirstFound ...bool) T {
	var returnFirstFound bool
	if len(optReturnFirstFound) > 0 {
		returnFirstFound = optReturnFirstFound[0]
	}
	if IsEmpty(attrs) {
		return make([]T, 0)
	}
	if returnFirstFound {
		return Find(obj, func(value T, key T, list T) bool {
			for k, v := range attrs {
				if v != value.(map[T]T)[k] {
					return false
				}
			}
			return true
		})

	} else {
		return Filter(obj, func(value T, key T, list T) bool {
			for k, v := range attrs {
				if v != value.(map[T]T)[k] {
					return false
				}
			}
			return true
		})
	}
}

// Convenience version of a common use case of `find`: getting the first object
// containing specific `key:value` pairs.
func FindWhere(obj []T, attrs map[T]T) T {
	return Where(obj, attrs, true)
}

// Internal function for comparing ints
func intLessThan(a T, b T) bool {
	return a.(int) < b.(int)
}

// Return the maximum element or (element-based computation).
func Max(lessThan func(T, T) bool, args ...T) T {
	val := args[0]
	for _, v := range args {
		if !lessThan(v, val) {
			val = v
		}
	}
	return val
}

// Return the maximum element or (element-based computation), int-specific version
func MaxInt(args ...int) int {
	val := args[0]
	for _, v := range args {
		if !intLessThan(v, val) {
			val = v
		}
	}
	return val
}

// Return the minimum element or (element-based computation).
func Min(lessThan func(T, T) bool, args ...T) T {
	val := args[0]
	for _, v := range args {
		if lessThan(v, val) {
			val = v
		}
	}
	return val
}

// Return the minimum element or (element-based computation), int-specific version
func MinInt(args ...int) int {
	val := args[0]
	for _, v := range args {
		if intLessThan(v, val) {
			val = v
		}
	}
	return val
}

// Shuffle an array, using the modern version of the
// [Fisher-Yates shuffle](http://en.wikipedia.org/wiki/Fisher–Yates_shuffle).
func Shuffle(obj []T) []T {
	shuffled := make([]T, len(obj))
	index := 0
	var rand int
	Each(obj, func(val, idx, list T) bool {
		rand = Random(index)
		index += 1
		shuffled[index-1] = shuffled[rand]
		shuffled[rand] = val
		return eachContinue
	})
	return shuffled
}

// Sample **n** random values from a collection.
// If **n** is not specified, returns a single random element.
func Sample(obj T, opt_n ...int) T {
	if IsMap(obj) {
		vals := Values(obj.(map[T]T))
		return vals[Random(len(vals))]
	}
	if opt_n == nil || len(opt_n) == 0 {
		return obj.([]T)[Random(len(obj.([]T)))]
	}
	return Shuffle(obj.([]T))[0:MinInt(opt_n[0], len(obj.([]T)))]
}

// An internal function to generate lookup iterators
func lookupIterator(value T) func(obj, idx, list T) T {
	if IsFunction(value) {
		//fmt.Printf("lookupIterator got a func\n")
		return value.(func(obj, idx, list T) T)
	}
	//fmt.Printf("lookupIterator didnt get a func\n")
	return func(obj, idx, list T) T {
		//fmt.Printf("inlookup iterator, got obj %v, idx %v, list %v\n",obj,idx,list)
		if IsMap(obj) {
			return obj.(map[T]T)[value]
		}
		return obj
	}
}

type sorter struct {
	list    []map[T]T
	orderby func(a, b *map[T]T) bool
}

// Len is part of sort.Interface.
func (s *sorter) Len() int {
	return len(s.list)
}

// Swap is part of sort.Interface.
func (s *sorter) Swap(i, j int) {
	s.list[i], s.list[j] = s.list[j], s.list[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *sorter) Less(i, j int) bool {
	return s.orderby(&s.list[i], &s.list[j])
}

func newSorter(data, orderby T) *sorter {
	this := new(sorter)
	this.list = data.([]map[T]T)
	this.orderby = orderby.(func(a, b *map[T]T) bool)
	return this
}

// A custom type, a list of maps and a (sort)by function, for use in SortBy
type mapSorter struct {
	maps []map[T]T
	by   func(a, b *map[T]T) bool
}

// Implementing Len for the sort.Interface
func (s *mapSorter) Len() int {
	return len(s.maps)
}

// Implementing Swap for the sort.Interface
func (s *mapSorter) Swap(i, j int) {
	s.maps[i], s.maps[j] = s.maps[j], s.maps[i]
}

// Implementing Less for the sort.Interface
// It is implemented by calling the "by" closure in the sorter.
func (s *mapSorter) Less(i, j int) bool {
	return s.by(&s.maps[i], &s.maps[j])
}

// Sort the object's values by a criterion produced by an iterator.
func SortBy(obj, value T, lessThan func(a, b *map[T]T) bool) []T {
	iterator := lookupIterator(value)
	mapped := &mapSorter{
		mapMap(obj.([]map[T]T), func(value, index, list T) map[T]T {
			if value == nil {
				return nil
			}
			return map[T]T{
				"value":    value,
				"index":    index,
				"criteria": iterator(value, index, list),
			}
		}),
		lessThan,
	}
	sort.Sort(mapped)
	return Pluck(mapped.maps, "value")
}

// Sort the object's values by a criterion produced by an iterator.
func SortBySorter(obj, value T, orderby func(a, b *map[T]T) bool) []T {
	iterator := lookupIterator(value)
	mapped := mapForSortBy(obj, func(value, index, list T) map[T]T {
		return map[T]T{
			"value":    value,
			"index":    index,
			"criteria": iterator(value, index, list),
		}
	})
	ss := newSorter(mapped, orderby)
	sort.Sort(ss)
	return Pluck(ss.list, "value")
}

// An internal function used for aggregate "group by" operations.
func group(behavior func(result map[T]T, k T, v T)) func(o T, v T) map[T]T {
	return func(obj T, value T) map[T]T {
		result := make(map[T]T, 0)
		var iterator func(T, T, T) T
		if value == nil {
			iterator = Identity
		} else {
			iterator = lookupIterator(value)
		}

		Each(obj, func(value T, index T, list T) bool {
			key := iterator(value, index, obj)
			//_,ok := result[key]
			behavior(result, key, value)
			return eachContinue
		})
		return result
	}
}

// Groups the object's values by a criterion. Pass either a string attribute
// to group by, or a function that returns the criterion.
var GroupBy = group(func(result map[T]T, key T, value T) {
	if key == nil {
		return
	}
	if Has(result, key) {
		//fmt.Printf("in group, got res %v, key %v, val %v\n\n",result,key,value)
		slice := result[key].([]T)
		slice = append(slice, value)
		result[key] = slice
	} else {
		slice := make([]T, 1)
		slice[0] = value
		result[key] = slice
	}
})

// Indexes the object's values by a criterion, similar to `groupBy`, but for
// when you know that your index values will be unique.
var IndexBy = group(func(result map[T]T, key T, value T) {
	if key == nil {
		return
	}
	result[key] = value
})

// Counts instances of an object that group by a certain criterion. Pass
// either a string attribute to count by, or a function that returns the
// criterion.
var CountBy = group(func(result map[T]T, key T, value T) {
	if key == nil {
		return
	}
	if Has(result, key) {
		//fmt.Printf("in group, got res %v, key %v, val %v\n\n",result,key,value)
		result[key] = (result[key]).(int) + 1
	} else {
		result[key] = 1
	}
})

// Use a comparator function to figure out the smallest index at which
// an object should be inserted so as to maintain order. Uses binary search.
func SortedIndex(array T, obj T, lessThan func(T, T) bool, opt_iterator ...func(T, T, T) T) int {
	var value T
	var iterator func(T, T, T) T
	_, isArrayOfMaps := array.([]map[T]T)
	_, isArray := array.([]T)
	if !(isArrayOfMaps || isArray) {
		fmt.Printf("Error: can't find sorted index of a non-list object, got %v\n", array)
		return math.MinInt64
	}
	if len(opt_iterator) > 0 {
		iterator = opt_iterator[0]
		value = iterator(obj, nil, nil)
	} else {
		value = obj
	}
	low := 0
	high := 0
	if isArrayOfMaps {
		high = len(array.([]map[T]T))
	} else if isArray {
		high = len(array.([]T))
	}
	var otherValue T
	for low < high {
		mid := uint(low+high) >> 1
		if iterator != nil {
			if isArrayOfMaps {
				otherValue = iterator(array.([]map[T]T)[mid], nil, nil)
			} else if isArray {
				otherValue = iterator(array.([]T)[mid], nil, nil)
			}
		} else {
			if isArrayOfMaps {
				otherValue = array.([]map[T]T)[mid]
			} else if isArray {
				otherValue = array.([]T)[mid]
			}
		}
		if lessThan(otherValue, value) {
			low = int(mid) + 1
		} else {
			high = int(mid)
		}
	}
	return low
}

// Safely create a real, live array from anything iterable.
func ToArray(obj T) []T {
	if obj == nil {
		return make([]T, 0)
	}
	if IsArray(obj) || IsArrayOfMaps(obj) || IsString(obj) {
		return Map(obj, Identity)
	}
	if IsMap(obj) {
		return Values(obj.(map[T]T))
	}
	fmt.Printf("Error: ToArray, got something I dont know what to do with %v\n", obj)
	return nil
}

//Return the number of elements in an object.
func Size(obj T) int {
	if IsEmpty(obj) {
		return 0
	}
	if IsArrayOfMaps(obj) {
		return len(obj.([]map[T]T))
	}
	if IsArray(obj) {
		return len(obj.([]T))
	}
	if IsMap(obj) {
		return len(obj.(map[T]T))
	}
	if IsString(obj) {
		return len(obj.(string))
	} else {
		fmt.Printf("TypeError (Size): what is this? %v\n", obj)
		return math.MinInt64
	}
}

// Array Functions

// Get the first n element(s) of an array. Passing **n** will return the first N
// values in the array.
// Aliased as HeadN
// Aliased as TakeN
func FirstN(array []T, n int, opt_guard ...bool) []T {
	if array == nil {
		return nil
	}
	if n == 0 {
		retval := make([]T, 0)
		retval = append(retval, array[0])
		return retval
	} else if len(opt_guard) > 0 && opt_guard[0] {
		retval := make([]T, 0)
		retval = append(retval, array[0])
		return retval
	}
	if n > len(array) {
		return array[:]
	} else {
		return array[0:n]
	}
}

// Get the first n element(s) of an array. Passing **n** will return the first N
// values in the array.
// Aliased as FirstN
// Aliased as TakeN
var HeadN func(array []T, n int, opt_guard ...bool) []T = FirstN

// Get the first n element(s) of an array. Passing **n** will return the first N
// values in the array.
// Aliased as FirstN
// Aliased as HeadN
var TakeN func(array []T, n int, opt_guard ...bool) []T = FirstN

// Get the first element of an array. Passing **n** will return the first N
// values in the array.
// Aliased as Head
// Aliased as Take
func First(array []T) T {
	if array == nil {
		return nil
	}
	return array[0]
}

// Get the first element of an array. Passing **n** will return the first N
// values in the array.
// Aliased as First
// Aliased as Take
var Head func(array []T) T = First

// Get the first element of an array. Passing **n** will return the first N
// values in the array.
// Aliased as First
// Aliased as Head
var Take func(array []T) T = First

// Returns everything but the last entry of the array.
// Passing **n** will return all the values in
// the array, excluding the last N.
func Initial(array []T, opt_n ...int) []T {
	if array == nil {
		return nil
	}
	var n int
	if len(opt_n) > 0 {
		n = opt_n[0]
	} else {
		n = len(array) - 1
	}
	if n > len(array) {
		return array[:]
	}
	return array[0:n]
}

// Get the last element of an array. Passing **n** will return the last N
// values in the array.
func Last(array []T, opt_n ...int) []T {
	if array == nil {
		return nil
	}
	var n int
	arraylen := len(array)
	if len(opt_n) > 0 {
		n = opt_n[0]
	} else {
		n = 1
	}
	if n > arraylen {
		return array[:]
	}
	return array[(arraylen - n):]
}

// Returns everything but the first entry of the array. Aliased as `tail` and `drop`.
// Especially useful on the arguments object. Passing an **n** will return
// the rest N values in the array.
// Aliased as Tail
// Aliased as Drop
func Rest(array []T) []T {
	if array == nil {
		return nil
	}
	dst := make([]T, len(array)-1)
	copy(dst, array[1:])
	return dst
}

// Returns everything but the first entry of the array. Aliased as `tail` and `drop`.
// Especially useful on the arguments object. Passing an **n** will return
// the rest N values in the array.
// Aliased as Rest
// Aliased as Drop
var Tail func(array []T) []T = Rest

// Returns everything but the first entry of the array. Aliased as `tail` and `drop`.
// Especially useful on the arguments object. Passing an **n** will return
// the rest N values in the array.
// Aliased as Rest
// Aliased as Tail
var Drop func(array []T) []T = Rest

// Trim out all falsy values from an array.
func Compact(array []T) []T {
	return Filter(array, IdentityIsTruthy)
}

// Internal implementation of a recursive `flatten` function.
//func flatten(input []T, shallow bool, output []T) []T {
func flatten(input T, shallow bool, output []T) []T {
	//if shallow && Every(input, IsArrayEach) {
	//	fmt.Printf("everything is array and shallow: %v\n",input)
	//	output = append(output,input...)
	//	return output
	//}
	Each(input, func(value T, idx T, list T) bool {
		if IsArray(value) {
			if shallow {
				//fmt.Printf("shallow output before: %v\n",output)
				output = append(output, value.([]T)...)
				//fmt.Printf("shallow output after : %v\n",output)
			} else {
				//fmt.Printf("output before: %v\n",output)
				output = flatten(value.([]T), shallow, output)
				//fmt.Printf("output after : %v\n",output)
			}
		} else if IsStringArray(value) {
			output = flatten(value, shallow, output)
		} else {
			//fmt.Printf("is not array %v output before: %v\n",value, output)
			output = append(output, value)
		}
		return eachContinue
	})
	return output
}

// Flatten out an array, either recursively (by default), or just one level.
//func Flatten (array []T, opt_shallow ...bool) []T {
func Flatten(array T, opt_shallow ...bool) []T {
	var shallow bool
	if len(opt_shallow) > 0 {
		shallow = opt_shallow[0]
	}
	return flatten(array, shallow, make([]T, 0))
}

// Return a version of the array that does not contain the specified value(s).
func Without(toRemove []T, opt_from ...T) []T {
	var comparator func(T, T) bool
	if len(opt_from) == 0 {
		return make([]T, 0)
	} else {
		if v, ok := opt_from[len(opt_from)-1].(func(T, T) bool); ok {
			comparator = v
		}
	}
	var rest []T = make([]T, 0)
	for _, from := range opt_from {
		rest = append(rest, from)
	}
	if comparator == nil {
		return Difference(toRemove, IdentityComparator, rest)
	} else {
		return Difference(toRemove, comparator, rest)
	}

}

// Produce a duplicate-free version of the array. If the array has already
// been sorted, you have the option of using a faster algorithm.
// Aliased as `Unique`.
func Uniq(list T, isSorted T /*bool or func*/, opt_iterator ...T) []T {
	var array []T
	var arrayofmaps []map[T]T
	isAM := IsArrayOfMaps(list)
	isA := IsArray(list)
	if isAM {
		arrayofmaps = list.([]map[T]T)
	} else if isA {
		array = list.([]T)
	}

	var iterator mapiterator
	var comparator func(T, T) bool
	if IsFunction(isSorted) {
		iterator = isSorted.(mapiterator)
		isSorted = false
		if len(opt_iterator) > 0 {
			comparator = opt_iterator[0].(func(T, T) bool)
		}
	} else if len(opt_iterator) > 0 {
		iterator = opt_iterator[0].(func(T, T, T) T)
		comparator = opt_iterator[1].(func(T, T) bool)
	}
	var initialA []T
	if iterator != nil {
		if isA {
			initialA = Map(array, iterator)
		} else if isAM {
			initialA = Map(arrayofmaps, iterator)
		}
	} else if isA {
		initialA = array
	}
	results := make([]T, 0)
	seen := make([]T, 0)
	if isA {
		Each(initialA, func(value T, index T, list T) bool {
			if isSorted.(bool) {
				if index == 0 || seen[len(seen)-1] != value {
					seen = append(seen, value)
					results = append(results, array[index.(int)])
				}
			} else if !Contains(seen, value) {
				seen = append(seen, value)
				results = append(results, array[index.(int)])
			}
			return eachContinue
		})
	}
	if isAM {
		Each(arrayofmaps, func(value T, index T, list T) bool {
			if isSorted.(bool) {
				if index == 0 || seen[len(seen)-1] != value {
					seen = append(seen, value)
					results = append(results, array[index.(int)])
				}
			} else if !Contains(seen, value, comparator) {
				seen = append(seen, value)
				results = append(results, value)
			}
			return eachContinue
		})
	}
	return results
}

// Produce a duplicate-free version of the array. If the array has already
// been sorted, you have the option of using a faster algorithm.
// Aliased as `Uniq`.
var Unique func(list T, isSorted T /*bool or func*/, opt_iterator ...T) []T = Uniq

// Produce an array that contains the union: each distinct element from all of
// the passed-in arrays.
func Union(opt_array ...T) []T {
	return Uniq(Flatten(opt_array, true), false)
}

// Produce an array that contains every item shared between all the
// passed-in arrays.
func Intersection(lessThan func(T, T) bool, opt_array ...T) []T {
	rest := Uniq(Flatten(Rest(opt_array), true), false)
	return Filter(Uniq(opt_array[0], false), func(this T, idx T, list T) bool {
		return Every(rest, func(that T, idx2 T, list T) bool {
			return IndexOf(rest, this, lessThan) != -1
		})
	})
}

// Take the difference between one array and a number of other arrays.
// Only the elements present in just the first array will remain.
func Difference(toRemove []T, comparator func(T, T) bool, opt_from ...[]T) []T {
	if len(opt_from) == 0 {
		return make([]T, 0)
	}
	var rest []T = make([]T, 0)
	for _, from := range opt_from {
		rest = flatten(from, true, rest)
	}
	return Filter(toRemove, func(val, idx, list T) bool {
		return !Contains(rest, val, comparator)
	})
}

// Zip together multiple lists into a single array -- elements that share
// an index go together.
func Zip(arrays ...[]T) []T {
	if arrays == nil || len(arrays) == 0 {
		return make([]T, 0)
	}
	var length int = 0
	var tmplength int = 0
	var num_arrays int
	for _, array := range arrays {
		num_arrays += 1
		tmplength = len(array)
		if tmplength > length {
			length = tmplength
		}
	}
	//var retval [][]T
	retval := make([]T, length)
	for i := 0; i < length; i++ {
		zipped := make([]T, num_arrays)
		for j, array := range arrays {
			if i < len(array) {
				zipped[j] = array[i]
			}
		}
		retval[i] = zipped
	}
	return retval
}

// Converts lists into objects. Pass either a single array of `[key, value]`
// pairs, or two parallel arrays of the same length -- one of keys, and one of
// the corresponding values.
func Object(pairs_or_two_arrays ...[]T) map[T]T {
	if pairs_or_two_arrays == nil {
		return nil
	}
	retval := make(map[T]T)
	if len(pairs_or_two_arrays) == 1 { // got single array of ['k1','v1',k2,v2,...] pairs
		kvpairs := pairs_or_two_arrays[0]
		length := len(kvpairs)
		if length == 0 {
			return retval
		}
		if IsArray(kvpairs[0]) {
			for i := 0; i < length; i++ {
				retval[kvpairs[i].([]T)[0]] = kvpairs[i].([]T)[1]
			}
		} else {
			for i := 0; i < length; i += 2 {
				retval[kvpairs[i]] = kvpairs[i+1]
			}
		}
		return retval
	}
	keys := pairs_or_two_arrays[0]
	values := pairs_or_two_arrays[1]
	length := len(keys)
	if length != len(values) {
		fmt.Printf("Object() Error: got arrays of unequal length\n")
	}
	for i := 0; i < length; i++ {
		retval[keys[i]] = values[i]
	}
	return retval
}

// Return the position of the first occurrence of an
// item in an array, or -1 if the item is not included in the array.
// If the array is large and already in sort order, pass `true`
// for **isSorted** to use binary search.
func IndexOf(array []T, item T, lessThan func(T, T) bool, isSorted ...bool) int {
	if array == nil {
		return -1
	}
	length := len(array)
	if length == 0 {
		return -1
	}
	i := 0
	// do binary search if isSorted = true
	if len(isSorted) > 0 && isSorted[0] {
		i = SortedIndex(array, item, lessThan)
		if array[i] == item {
			return i
		} else {
			return -1
		}
	}
	for i, v := range array {
		if v == item {
			return i
		}
	}
	return -1
}

// Like IndexOf but do the search starting from the end moving backwards
func LastIndexOf(array []T, item T, from ...int) int {
	if array == nil {
		return -1
	}
	var i int
	if from != nil {
		i = from[0]
	} else {
		i = len(array)
	}
	for i > 0 {
		i -= 1
		if array[i] == item {
			return i
		}
	}
	return -1
}

// Generate an integer Array containing an arithmetic progression. A port of
// Underscore's range() which is a port of the native Python `range()` function. See
// [the Python documentation](http://docs.python.org/library/functions.html#range).
func Range(start_stop_and_step ...int) []T {
	if start_stop_and_step == nil {
		return make([]T, 0)
	}
	var start int
	var stop int
	var step int
	argslength := len(start_stop_and_step)
	if argslength == 3 {
		start, stop, step = start_stop_and_step[0], start_stop_and_step[1], start_stop_and_step[2]
	} else if argslength == 2 {
		start, stop = start_stop_and_step[0], start_stop_and_step[1]
		step = 1
	} else if argslength == 1 {
		stop = start_stop_and_step[0]
		start = 0
		step = 1
	}

	if step == 0 {
		step = 1
	}

	length := int(math.Max(math.Ceil(float64(stop-start)/float64(step)), 0))
	idx := 0
	retval := make([]T, length)

	for idx < length {
		retval[idx] = start
		start += step
		idx += 1
	}

	return retval
}

// Function Functions

// Partially apply a function by creating a version that has had some of its
// arguments pre-filled, without changing its dynamic `this` context.
func Partial(fn func(...T) T, savedArgs ...T) func(...T) T {
	return func(laterArgs ...T) T {
		//args := make([]T, len(savedArgs) + len(laterArgs))
		args := append(savedArgs, laterArgs...)
		//for i,v := range savedArgs {
		//args[i] = v
		//}
		//for j,v := range laterArgs {
		//args[ len(savedArgs) + j ]= v
		//}
		return fn(args...)
	}
}

// Memoize an expensive function by storing its results.
func Memoize(fn func(...T) T, opt_hasher ...func(...T) T) func(...T) T {
	memo := map[T]T{}
	var hasher func(...T) T
	if opt_hasher != nil {
		hasher = opt_hasher[0]
	} else {
		hasher = identityHasher
	}
	return func(args ...T) T {
		key := hasher(args...)
		if Has(memo, key) {
			return memo[key]
		}
		memo[key] = fn(args...)
		return memo[key]
	}
}

//delay_.delay(function, wait, *arguments) 
//Much like setTimeout, invokes function after wait milliseconds. If you pass the optional arguments, they will be forwarded on to the function when it is invoked.
//
//var log = _.bind(console.log, console);
//_.delay(log, 1000, 'logged later');
//=> 'logged later' // Appears after one second.
func Delay(fn func(), waitMilliseconds int64, savedArgs ...T) {
	go (func() {
		timer := time.NewTimer( time.Duration( waitMilliseconds * 1000000 ))
		for {
			select {
				case <-timer.C:
				timer.Stop()
				fn()
				break
			}
		}
	})()

/*
		func() {
			//for {
			//<-timer.C
				//timer.Stop()
				//fmt.Fprintf( os.Stdout, "about to call fn\n")
				fn()
			//}
		})
*/
}

//throttle_.throttle(function, wait, [options]) 
//Creates and returns a new, throttled version of the passed function, that, when invoked repeatedly, will only actually call the original function at most once per every wait milliseconds. Useful for rate-limiting events that occur faster than you can keep up with.
//
//By default, throttle will execute the function as soon as you call it for the first time, and, if you call it again any number of times during the wait period, as soon as that period is over. If you'd like to disable the leading-edge call, pass {leading: false}, and if you'd like to disable the execution on the trailing-edge, pass 
//{trailing: false}.
//
//throttled := Throttle(updatePosition, 100)
//$(window).scroll(throttled);
func ThrottleNano(fn func(...T), waitNanoseconds int64, options ...map[string]bool) func(...T) {
	var last int64 //:= time.Now().UnixNano()
	callLeading := true
	callTrailing := true
	var calltimer *time.Timer
	if options != nil {
		if v,ok := options[0]["leading"] ; ok {
			callLeading = v
		}
		if v,ok := options[0]["trailing"] ; ok {
			callTrailing = v
		}
	}
	if callLeading {
		fn()
	}
	return func(args ...T) {
		now := time.Now().UnixNano()
		diff := now - last
		if diff >= waitNanoseconds {
			if calltimer != nil {
				calltimer.Stop()
				calltimer = nil
			}
			last = now
			fn(args...)
		} else {
			if calltimer == nil {
				calltimer = time.NewTimer(time.Duration(waitNanoseconds))
				<-calltimer.C
				last = now
				if callTrailing {
					fn(args...)
				}
				calltimer.Stop()
				calltimer = nil
			}
		}
	}
}

func Throttle(fn func(...T), waitMilliseconds int64, options ...map[string]bool) func(...T) {
	return ThrottleNano(fn, waitMilliseconds * 1000000, options...)
}

// Returns a function that will be executed at most one time, no matter how
// often you call it. Useful for lazy initialization.
func Once(fn func(...T) T) func(...T) T {
	ran := false
	var memo T
	return func(args ...T) T {
		if ran {
			return memo
		}
		ran = true
		memo = fn(args...)
		fn = nil
		return memo
	}
}

// Returns a function that will only be executed after being called N times.
func After(times int, fn func(...T) T) func(...T) T {
	return func(args ...T) T {
		if times < 0 {
			times = 0
		} else {
			times -= 1
		}
		if times < 1 {
			return fn(args...)
		}
		return nil
	}
}

//Returns an int64 timestamp for the current time, using the fastest method available in the runtime. Useful for implementing timing/animation functions.
func Now() int64 {
	return time.Now().Unix()
}

// Returns the first function passed as an argument to the second,
// allowing you to adjust arguments, run code before and after, and
// conditionally execute the original function.
func Wrap(fn func(...T) T, wrapper func(...T) T) func(...T) T {
	return func(args ...T) T {
		wrappedargs := make([]T, 0)
		wrappedargs = append(wrappedargs, fn)
		wrappedargs = append(wrappedargs, args...)
		return wrapper(wrappedargs...)
	}
}

// Returns a function that is the composition of a list of functions, each
// consuming the return value of the function that follows.
func Compose(funcs ...T) func(...T) T {
	return func(args ...T) T {
		for i := len(funcs) - 1; i >= 0; i -= 1 {
			fn := funcs[i].(func(...T) T)
			retval := fn(args...)
			args = nil
			if retval == nil {
				continue
			} else if _, ok := retval.([]T); ok {
				args = retval.([]T)
			} else {
				args = []T{retval}
			}
		}
		return args[0]
	}
}


// Map Functions

// Retrieve the names of a maps keys
func Keys(obj map[T]T) []T {
	retval := make([]T, 0)
	if obj == nil {
		return retval
	}
	for key := range obj {
		retval = append(retval, key)
	}
	return retval
}

// Retrieve the values of a maps keys
func Values(obj map[T]T) []T {
	retval := make([]T, 0)
	if obj == nil {
		return retval
	}
	keys := Keys(obj)
	for _, key := range keys {
		retval = append(retval, obj[key])
	}
	return retval
}

// Convert an object into a list of `[key, value]` pairs.
func Pairs(obj map[T]T) []T {
	keys := Keys(obj)
	length := len(keys)
	pairs := make([]T, length)
	for i := 0; i < length; i++ {
		pairs[i] = []T{keys[i], obj[keys[i]]}
	}
	return pairs
}

func Invert(obj map[T]T) map[T]T {
	result := map[T]T{}
	for k, v := range obj {
		result[v] = k
	}
	return result
}

func Extend(objToExtend map[T]T, args ...T) map[T]T {
	Each(args, func(objToCopy, key, list T) bool {
		if !IsMap(objToCopy) {
			return eachContinue
		}
		for k, v := range objToCopy.(map[T]T) {
			objToExtend[k] = v
		}
		return eachContinue
	})
	return objToExtend
}

// Return a copy of the object only containing the whitelisted properties.
func Pick(obj map[T]T, keysToKeep ...T) map[T]T {
	copy := map[T]T{}
	Each(keysToKeep, func(keyToKeep, index, list T) bool {
		if _, isList := keyToKeep.([]T); isList {
			for _, keyToKeep := range keyToKeep.([]T) {
				if v, ok := obj[keyToKeep]; ok {
					copy[keyToKeep] = v
				}
			}
		} else if v, ok := obj[keyToKeep]; ok {
			copy[keyToKeep] = v
		}
		return eachContinue
	})
	return copy
}

// Return a copy of the object without the blacklisted properties.
func Omit(obj map[T]T, keysToRemove ...T) map[T]T {
	copy := map[T]T{}
	keysToRemove = Flatten(keysToRemove, true)
	for k, v := range obj {
		if !Contains(keysToRemove, k) {
			copy[k] = v
		}
	}
	return copy
}

// Fill in a given object with default properties.
func Defaults(obj map[T]T, args ...T) map[T]T {
	Each(args, func(val, idx, list T) bool {
		if !IsMap(val) {
			return eachContinue
		}
		for k, v := range val.(map[T]T) {
			if _, ok := obj[k]; !ok {
				obj[k] = v
			}
		}
		return eachContinue
	})
	return obj
}

// Create a (not-shallow-cloned if Array or Map) duplicate of an object.
func Clone(obj T) T {
	if IsMap(obj) {
		return Extend(map[T]T{}, obj.(map[T]T))
	}
	if IsArray(obj) {
		return obj.([]T)[:]
	}
	return obj
}

// Invokes interceptor with the obj, and then returns obj.
// The primary purpose of this method is to "tap into" a method chain, in
// order to perform operations on intermediate results within the chain.
func Tap(obj T, fn func(...T) T) T {
	fn(obj)
	return obj
}

func Result(obj, propertyName T) T {
	if obj == nil {
		return nil
	}
	val := obj.(map[T]T)[propertyName]
	if IsFunctionVariadic(val) { // func(...T) T
		return val.(func(...T) T)(obj)
	}
	return val
}

// Add a "chain" function, which will delegate to the wrapper.
func (this *Underscore) Chain() *Underscore {
	this.ischained = true
	return this
}
func (this *Underscore) Max(lessThan func(T, T) bool) *Underscore {
	return this.result(Max(lessThan, this.wrapped.([]T)...))
}
func (this *Underscore) Tap(fn func(...T) T) *Underscore {
	Tap(this.wrapped, fn)
	return this
}
func (this *Underscore) Value() T {
	return this.wrapped
}

func (this *Underscore) IsFinite() *Underscore {
	v, _ := this.wrapped.(float64)
	return this.result(!math.IsNaN(v) && !math.IsInf(v, 1) && !math.IsInf(v, -1))
}

func (this *Underscore) IsNaN() *Underscore {
	v, _ := this.wrapped.(float64)
	return this.result(math.IsNaN(v))
}

func (this *Underscore) Has(key T) *Underscore {
	return this.result(Has(this.wrapped, key))
}

// Utility Functions

func IdentityEach(val T, index T, list T) bool {
	return val == val
}
func identityHasher(val ...T) T {
	return val[0]
}

func IdentityIsTruthy(val T, index T, list T) bool {
	v, ok := val.(bool)
	if ok {
		return v
	}
	sv, sok := val.(string)
	if sok {
		return sv != ""
	}
	iv, iok := val.(int)
	if iok {
		return iv != 0
	}
	return val != nil
}

func Identity(val, index, list T) T {
	return val
}

func IdentityComparator(a, b T) bool {
	return a == b
}

func IdentityMap(val T, index T, list T) map[T]T {
	return val.(map[T]T)
}

// Is a given value an array?
func IsString(obj T) bool {
	v, _ := obj.(string)
	return v != ""
}

// Is a given value an array?
func IsArray(obj T) bool {
	_, ok := obj.([]T)
	return ok // v != nil
}

// Is a given value an array?
func IsStringArray(obj T) bool {
	_, ok := obj.([]string)
	return ok // v != nil
}
func IsArrayEach(obj T, idx T, list T) bool {
	return IsArray(obj)
}

// Is a given value an array?
func IsArrayOfMaps(obj T) bool {
	v, _ := obj.([]map[T]T)
	return v != nil
}

// Is a given variable a map
func IsMap(obj T) bool {
	v, _ := obj.(map[T]T)
	return v != nil
}

// Is a given value an array?
func IsFunction(obj T) bool {
	v, _ := obj.(func(T, T, T) T)
	return v != nil
}

// Is a given value an array?
func IsFunctionVariadic(obj T) bool {
	v, _ := obj.(func(...T) T)
	return v != nil
}

// Is a given array, string, or object empty?
// An "empty" object has no enumerable own-properties.
func IsEmpty(obj T) bool {
	if obj == nil {
		return true
	}
	if IsArray(obj) {
		return len(obj.([]T)) == 0
	}
	if IsMap(obj) {
		return len(obj.(map[T]T)) == 0
	}
	if IsString(obj) {
		return len(obj.(string)) == 0
	}
	return false
}
func (this *Underscore) IsEmpty(obj T) bool {
	return IsEmpty(this.wrapped)
}

// Keep the identity function around for default iterators.
func (this *Underscore) Identity(value ...T) T {
	return value[0]
}

// Shortcut function for checking if an object has a given property directly
// on itself (in other words, not on a prototype).
func Has(obj T, key T) bool {
	_, ok := obj.(map[T]T)[key]
	return ok
}

// Run a function **n** times.
func Times(n int, iterator func(...T) T) []T {
	if n < 0 {
		return []T{}
	}
	collected := make([]T, n)
	for i := 0; i < n; i += 1 {
		collected[i] = iterator(i)
	}
	return collected
}

// Run a function **n** times.
func (this *Underscore) Times(iterator func(...T) T) []T {
	return Times(this.wrapped.(int), iterator)
}

// Return a random integer between min and max (inclusive).
func Random(min int, optmax ...int) int {
	var max int
	if optmax == nil {
		max = min
		min = 0
	} else {
		max = optmax[0]
	}
	return min + int(rand.Float64()*float64(max-min+1))
}

func (this *Underscore) Random(min int, optmax ...int) int {
	return Random(min, optmax...)
}

// Return a random float64 between min and max (inclusive).
func RandomFloat64(min float64, optmax ...float64) float64 {
	var max float64
	if optmax == nil {
		max = min
		min = 0
	} else {
		max = optmax[0]
	}
	return min + rand.Float64()*(max-min+1.0)
}

func (this *Underscore) RandomFloat64(min float64, optmax ...float64) float64 {
	return RandomFloat64(min, optmax...)
}

// OOP-style funcs for Underscore

// OOP-style support, add method to *Underscore, see func Every
// Aliased as Every
func (this *Underscore) All(obj T, opt_iterator ...eachlistiterator) *Underscore {
	return this.result(Every(this.wrapped, opt_iterator...))
}

// OOP-style support, add method to *Underscore, see func Any
// Aliased as Some
func (this *Underscore) Any(opt_predicate ...func(val, index, list T) bool) *Underscore {
	return this.result(Any(this.wrapped, opt_predicate...))
}

// OOP-style support, add method to *Underscore, see func Contains
// Aliased as Include
func (this *Underscore) Contains(target T, opt_comparator ...func(T, T) bool) *Underscore {
	return this.result(Contains(this.wrapped, target, opt_comparator...))
}

// OOP-style support, add method to *Underscore, see func Map
func (this *Underscore) Map(iterator func(T, T, T) T) *Underscore {
	return this.result(Map(this.wrapped, iterator))
}

// OOP-style support, add method to *Underscore, see func Collect
func (this *Underscore) Collect(iterator func(T, T, T) T) *Underscore {
	return this.result(Collect(this.wrapped, iterator))
}

// OOP-style support, add method to *Underscore, see func Flatten
func (this *Underscore) Flatten(opt_shallow ...bool) *Underscore {
	return this.result(Flatten(this.wrapped.([]T), opt_shallow...))
}

// OOP-style support, add method to *Underscore, see func Reduce
func (this *Underscore) Reduce(iterator func(T, T, T, T) T, memo ...T) *Underscore {
	v, _ := Reduce(this.wrapped.([]T), iterator, memo...)
	return this.result(v)
}

// OOP-style support, add method to *Underscore, see func Inject
func (this *Underscore) Inject(iterator func(T, T, T, T) T, memo ...T) *Underscore {
	v, _ := Inject(this.wrapped.([]T), iterator, memo...)
	return this.result(v)
}

// OOP-style support, add method to *Underscore, see func FoldL
func (this *Underscore) FoldL(iterator func(T, T, T, T) T, memo ...T) *Underscore {
	v, _ := FoldL(this.wrapped.([]T), iterator, memo...)
	return this.result(v)
}

// OOP-style support, add method to *Underscore, see func ReduceRight
func (this *Underscore) ReduceRight(iterator func(T, T, T, T) T, memo ...T) *Underscore {
	v, _ := ReduceRight(this.wrapped.([]T), iterator, memo...)
	return this.result(v)
}

// OOP-style support, add method to *Underscore, see func FoldR
func (this *Underscore) FoldR(iterator func(T, T, T, T) T, memo ...T) *Underscore {
	v, _ := FoldR(this.wrapped.([]T), iterator, memo...)
	return this.result(v)
}

// OOP-style support, add method to *Underscore, see func Filter
// Aliased as Select
func (this *Underscore) Filter(iterator eachlistiterator) *Underscore {
	return this.result(Filter(this.wrapped.([]T), iterator))
}

// OOP-style support, add method to *Underscore, see func Select
// Aliased as Filter
func (this *Underscore) Select(iterator eachlistiterator) *Underscore {
	return this.result(Filter(this.wrapped.([]T), iterator))
}

// OOP-style support, add method to *Underscore, see func Reject
func (this *Underscore) Reject(iterator eachlistiterator) *Underscore {
	return this.result(Reject(this.wrapped.([]T), iterator))
}

// OOP-style support, add method to *Underscore, see func SortBy
func (this *Underscore) SortBy(value T, orderby func(a, b *map[T]T) bool) *Underscore {
	return this.result(SortBy(this.wrapped, value, orderby))
}

// OOP-style support, add method to *Underscore, see func SortBySorter
func (this *Underscore) SortBySorter(value T, orderby func(a, b *map[T]T) bool) *Underscore {
	return this.result(SortBySorter(this.wrapped, value, orderby))
}

// OOP-style support, add method to *Underscore, see func Reverse
func (this *Underscore) Reverse() *Underscore {
	a := make([]T, len(this.wrapped.([]T)))
	copy(a, this.wrapped.([]T))
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return this.result(a)
}

// OOP-style support, add method to *Underscore, see func Concat
func (this *Underscore) Concat(array []T) *Underscore {
	a := append(this.wrapped.([]T), array...)
	return this.result(a)
}

// OOP-style support, add method to *Underscore, see func Difference
func (this *Underscore) Difference(comparator func(T, T) bool, opt_from ...[]T) *Underscore {
	return this.result(Difference(this.wrapped.([]T), comparator, opt_from...))
}

// OOP-style support, add method to *Underscore, see func Drop
// Aliased as Rest
// Aliased as Tail
func (this *Underscore) Drop() *Underscore {
	return this.result(Rest(this.wrapped.([]T)))
}

// OOP-style support, add method to *Underscore, see func Unshift
func (this *Underscore) Unshift(elems ...T) *Underscore {
	for _, v := range this.wrapped.([]T) {
		elems = append(elems, v)
	}
	return this.result(elems)
}

// OOP-style support, add method to *Underscore, see func Shift
func (this *Underscore) Shift() *Underscore {
	a, _ := this.wrapped.([]T)
	return this.result(a[0])
}

// OOP-style support, add method to *Underscore, see func Pop
func (this *Underscore) Pop() *Underscore {
	a, _ := this.wrapped.([]T)
	retval := a[:len(a)-1]
	return this.result(retval)
}

// OOP-style support, add method to *Underscore, see func Clone
func (this *Underscore) Clone() *Underscore {
	return this.result(Clone(this.wrapped))
}

// OOP-style support, add method to *Underscore, see func Compacts
func (this *Underscore) Compact() *Underscore {
	v, _ := this.wrapped.([]T)
	return this.result(Compact(v))
}

// OOP-style support, add method to *Underscore, see func Defaults
func (this *Underscore) Defaults(args ...T) *Underscore {
	v, _ := this.wrapped.(map[T]T)
	return this.result(Defaults(v, args...))
}

// OOP-style support, add method to *Underscore, see func Each
func (this *Underscore) Each(iterator eachlistiterator) *Underscore {
	Each(this.wrapped, iterator)
	return this
}

// OOP-style support, add method to *Underscore, see func GroupBy
func (this *Underscore) GroupBy(fn func(v T) map[T]T) *Underscore {
	GroupBy(this.wrapped, fn)
	return this
}

// OOP-style support, add method to *Underscore, see func IndexBy
func (this *Underscore) IndexBy(fn func(v T) map[T]T) *Underscore {
	IndexBy(this.wrapped, fn)
	return this
}

// OOP-style support, add method to *Underscore, see func CountBy
func (this *Underscore) CountBy(fn func(v T) map[T]T) *Underscore {
	CountBy(this.wrapped, fn)
	return this
}

// OOP-style support, add method to *Underscore, see func Every
// Aliased as All
func (this *Underscore) Every(obj T, opt_iterator ...eachlistiterator) *Underscore {
	return this.result(Every(this.wrapped, opt_iterator...))
}

// OOP-style support, add method to *Underscore, see func Extend
func (this *Underscore) Extend(args ...T) *Underscore {
	v, _ := this.wrapped.(map[T]T)
	return this.result(Extend(v, args...))
}

// OOP-style support, add method to *Underscore, see func Detect
// Aliased as Find
func (this *Underscore) Detect(predicate func(T, T, T) bool) *Underscore {
	v, _ := this.wrapped.([]T)
	return this.result(Detect(v, predicate))
}

// OOP-style support, add method to *Underscore, see func Find
// Aliased as Detect
func (this *Underscore) Find(predicate func(T, T, T) bool) *Underscore {
	v, _ := this.wrapped.([]T)
	return this.result(Find(v, predicate))
}

// OOP-style support, add method to *Underscore, see func FindWhere
func (this *Underscore) FindWhere(attrs map[T]T) *Underscore {
	v, _ := this.wrapped.([]T)
	return this.result(FindWhere(v, attrs))
}

// OOP-style support, add method to *Underscore, see func First
// Aliased as Head
// Aliased as Take
func (this *Underscore) First() *Underscore {
	v, _ := this.wrapped.([]T)
	return this.result(First(v))
}

// OOP-style support, add method to *Underscore, see func Head
// Aliased as First
// Aliased as Take
func (this *Underscore) Head() *Underscore {
	v, _ := this.wrapped.([]T)
	return this.result(First(v))
}

// OOP-style support, add method to *Underscore, see func Take
// Aliased as First
// Aliased as Head
func (this *Underscore) Take() *Underscore {
	v, _ := this.wrapped.([]T)
	return this.result(First(v))
}

// OOP-style support, add method to *Underscore, see func FirstN
// Aliased as HeadN
// Aliased as TakeN
func (this *Underscore) FirstN(n int, opt_guard ...bool) *Underscore {
	v, _ := this.wrapped.([]T)
	return this.result(FirstN(v, n, opt_guard...))
}

// OOP-style support, add method to *Underscore, see func HeadN
// Aliased as FirstN
// Aliased as TakeN
func (this *Underscore) HeadN(n int, opt_guard ...bool) *Underscore {
	v, _ := this.wrapped.([]T)
	return this.result(FirstN(v, n, opt_guard...))
}

// OOP-style support, add method to *Underscore, see func Contains
// Aliased as Contains
func (this *Underscore) Include(target T, opt_comparator ...func(T, T) bool) *Underscore {
	return this.result(Contains(this.wrapped, target, opt_comparator...))
}

// OOP-style support, add method to *Underscore, see func IndexOf
func (this *Underscore) IndexOf(item T, lessThan func(T, T) bool, isSorted ...bool) *Underscore {
	return this.result(IndexOf(this.wrapped.([]T), item, lessThan, isSorted...))
}

// OOP-style support, add method to *Underscore, see func Initial
func (this *Underscore) Initial(opt_n ...int) *Underscore {
	v, _ := this.wrapped.([]T)
	return this.result(Initial(v, opt_n...))
}

// OOP-style support, add method to *Underscore, see func Intersection
func (this *Underscore) Intersection(lessThan func(T, T) bool, opt_array ...T) *Underscore {
	v, _ := this.wrapped.([]T)
	return this.result(Intersection(lessThan, append(v, opt_array...)))
}

// OOP-style support, add method to *Underscore, see func Invert
func (this *Underscore) Invert() *Underscore {
	v, _ := this.wrapped.(map[T]T)
	return this.result(Invert(v))
}

// OOP-style support, add method to *Underscore, see func Invoke
func (this *Underscore) Invoke(method func(this T, thisArgs ...T) T, args ...T) *Underscore {
	return this.result(Invoke(this.wrapped, method, args...))
}

// OOP-style support, add method to *Underscore, see func IsFunction
func (this *Underscore) IsFunction() *Underscore {
	return this.result(IsFunction(this.wrapped))
}

// OOP-style support, add method to *Underscore, see func IsFunctionVariadic
func (this *Underscore) IsFunctionVariadic() *Underscore {
	return this.result(IsFunctionVariadic(this.wrapped))
}

// OOP-style support, add method to *Underscore, see func IsArray
func (this *Underscore) IsArray() *Underscore {
	return this.result(IsArray(this.wrapped))
}

// OOP-style support, add method to *Underscore, see func IsArrayOfMaps
func (this *Underscore) IsArrayOfMaps() *Underscore {
	return this.result(IsArrayOfMaps(this.wrapped))
}

// OOP-style support, add method to *Underscore, see func IsMap
func (this *Underscore) IsMap() *Underscore {
	return this.result(IsMap(this.wrapped))
}

// OOP-style support, add method to *Underscore, see func IsString
func (this *Underscore) IsString() *Underscore {
	return this.result(IsString(this.wrapped))
}

// OOP-style support, add method to *Underscore, see func IsStringArray
func (this *Underscore) IsStringArray() *Underscore {
	return this.result(IsStringArray(this.wrapped))
}

// OOP-style support, add method to *Underscore, see func Keys
func (this *Underscore) Keys() *Underscore {
	v, _ := this.wrapped.(map[T]T)
	return this.result(Keys(v))
}

// OOP-style support, add method to *Underscore, see func Last
func (this *Underscore) Last(opt_n ...int) *Underscore {
	v, _ := this.wrapped.([]T)
	return this.result(Last(v, opt_n...))
}

// OOP-style support, add method to *Underscore, see func IndexOf
func (this *Underscore) LastIndexOf(item T, from ...int) *Underscore {
	return this.result(LastIndexOf(this.wrapped.([]T), item, from...))
}

// OOP-style support, add method to *Underscore, see func MaxInt
func (this *Underscore) MaxInt() *Underscore {
	return this.result(MaxInt(this.wrapped.([]int)...))
}

// OOP-style support, add method to *Underscore, see func Min
func (this *Underscore) Min(lessThan func(T, T) bool) *Underscore {
	return this.result(Min(lessThan, this.wrapped.([]T)...))
}

// OOP-style support, add method to *Underscore, see func MinInt
func (this *Underscore) MinInt() *Underscore {
	return this.result(MinInt(this.wrapped.([]int)...))
}

// OOP-style support, add method to *Underscore, see func Object
func (this *Underscore) Object() *Underscore {
	return this.result(Object(this.wrapped.([]T)))
}

// OOP-style support, add method to *Underscore, see func Omit
func (this *Underscore) Omit(keysToRemove ...T) *Underscore {
	return this.result(Omit(this.wrapped.(map[T]T), keysToRemove...))
}

// OOP-style support, add method to *Underscore, see func Pairs
func (this *Underscore) Pairs() *Underscore {
	return this.result(Pairs(this.wrapped.(map[T]T)))
}

// OOP-style support, add method to *Underscore, see func Pick
func (this *Underscore) Pick(keysToKeep ...T) *Underscore {
	return this.result(Pick(this.wrapped.(map[T]T), keysToKeep...))
}

// OOP-style support, add method to *Underscore, see func Pluck
func (this *Underscore) Pluck(targetvalue T) *Underscore {
	return this.result(Pluck(this.wrapped, targetvalue))
}

// OOP-style support, add method to *Underscore, see func Rest
// Aliased as Tail
// Aliased as Drop
func (this *Underscore) Rest() *Underscore {
	return this.result(Rest(this.wrapped.([]T)))
}

// OOP-style support, add method to *Underscore, see func Sample
func (this *Underscore) Sample(opt_n ...int) *Underscore {
	return this.result(Sample(this.wrapped, opt_n...))
}

// OOP-style support, add method to *Underscore, see func Shuffle
func (this *Underscore) Shuffle() *Underscore {
	return this.result(Shuffle(this.wrapped.([]T)))
}

// OOP-style support, add method to *Underscore, see func Size
func (this *Underscore) Size() *Underscore {
	return this.result(Size(this.wrapped))
}

// OOP-style support, add method to *Underscore, see func Some
// Aliased as Any
func (this *Underscore) Some(opt_predicate ...func(val, index, list T) bool) *Underscore {
	return this.result(Any(this.wrapped, opt_predicate...))
}

// OOP-style support, add method to *Underscore, see func Rest
// Aliased as Tail
// Aliased as Drop
func (this *Underscore) Tail() *Underscore {
	return this.result(Rest(this.wrapped.([]T)))
}

// OOP-style support, add method to *Underscore, see func ToArray
func (this *Underscore) ToArray() *Underscore {
	return this.result(ToArray(this.wrapped))
}

// OOP-style support, add method to *Underscore, see func ToArray
func (this *Underscore) Now() *Underscore {
	return this.result(Now())
}

// OOP-style support, add method to *Underscore, see func Union
func (this *Underscore) Union(opt_array ...T) *Underscore {
	args := this.wrapped.([]T)
	return this.result(Union(append(args, opt_array...)))
}

// OOP-style support, add method to *Underscore, see func Uniq
// Aliased as Unique
func (this *Underscore) Uniq(isSorted T /*bool or func*/, opt_iterator ...T) *Underscore {
	return this.result(Uniq(this.wrapped, isSorted, opt_iterator))
}

// OOP-style support, add method to *Underscore, see func Unique
// Aliased as Uniq
func (this *Underscore) Unique(isSorted T /*bool or func*/, opt_iterator ...T) *Underscore {
	return this.result(Unique(this.wrapped, isSorted, opt_iterator))
}

// OOP-style support, add method to *Underscore, see func Values
func (this *Underscore) Values() *Underscore {
	return this.result(Values(this.wrapped.(map[T]T)))
}

// OOP-style support, add method to *Underscore, see func Without
func (this *Underscore) Without(opt_from ...T) *Underscore {
	return this.result(Without(this.wrapped.([]T), opt_from...))
}

// OOP-style support, add method to *Underscore, see func Where
func (this *Underscore) Where(attrs map[T]T, optReturnFirstFound ...bool) *Underscore {
	return this.result(Where(this.wrapped.([]T), attrs, optReturnFirstFound...))
}

// OOP-style support, add method to *Underscore, see func Zip
func (this *Underscore) Zip() *Underscore {
	return this.result(Zip(this.wrapped.([]T)))
}

// Helper function to continue chaining intermediate results.
func (this *Underscore) result(obj T) *Underscore {
	if this.ischained {
		return New(obj).Chain()
	}
	return this
}
