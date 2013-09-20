package underscore

import (
	"fmt"
	"math"
	"math/rand"
)

type Underscore struct {}
type T interface{}
type eachlistiterator func(T,T,T) bool
type mapiterator func(T,T,T) T

const EachContinue bool = false
const EachBreak    bool = true


// Is a given value an array?
func IsString (obj T) bool {
	v,_ := obj.(string)
	return v != ""
}

// Is a given value an array?
func IsArray (obj T) bool {
	_,ok := obj.([]T)
	return ok // v != nil
}
func IsArrayEach (obj T, idx T, list T) bool {
	return IsArray(obj)
}

// Is a given value an array?
func IsArrayOfMaps (obj T) bool {
	v,_ := obj.([]map[T]T)
	return v != nil
}


// Is a given variable a map
func IsMap (obj T) bool {
	v,_ := obj.(map[T]T) 
	return v != nil
}

// Is a given value an array?
func IsFunction(obj T) bool {
	v,_ := obj.(func(T,T,T)T)
	return v != nil
}


// Is a given array, string, or object empty?
// An "empty" object has no enumerable own-properties.
func IsEmpty (obj T) bool {
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




func Each(elemslist_or_map T, iterator eachlistiterator ) {
	if elemslist_or_map == nil || IsEmpty(elemslist_or_map) {
		return
	} 

	if IsArray( elemslist_or_map ) {
		for i,elem := range elemslist_or_map.([]T) {
			if iterator(elem, i, elemslist_or_map.([]T)) == EachBreak {
				return
			}
		}

	} else if IsArrayOfMaps( elemslist_or_map ) {
		for i,elem := range elemslist_or_map.([]map[T]T) {
			if iterator(elem, i, elemslist_or_map.([]map[T]T)) == EachBreak {
				return
			}
		}

	} else if IsMap( elemslist_or_map ) {
		for k,v := range elemslist_or_map.(map[T]T) {
			if iterator(v,k,elemslist_or_map.(map[T]T)) == EachBreak {
				return
			}
		}
	} else {
		fmt.Printf("Each isnt doing anything useful with first arg, %v\n",elemslist_or_map)
	}
	
}


// Return the results of applying the iterator to each element.
func Map(obj T, iterator func(T,T,T) T) []T {
	results := make([]T,0)
	if obj == nil {
		return results
	}
	Each(obj, func (value T, index T, list T) bool {
		if v := iterator(value,index,list); v != nil {
			results = append( results, v )
		}
		return EachContinue
	})
	return results
}
var Collect func (obj T, iterator func(T,T,T) T) []T = Map


const ReduceError = "Reduce of empty array with no initial value"

// **Reduce** builds up a single result from a list of values, aka `inject`,
// or `foldl`. 
func Reduce (obj []T, iterator func(T,T,T,T) T, memo ...T) (T,string) {
	initial := len(memo) > 0
	if obj == nil {
		obj = make([]T,0)
	}
	Each(obj, func (value T, index T, list T) bool {
		if !initial {
			memo[0] = value
			initial = true
		} else {
			memo[0] = iterator(memo[0], value, index, list)
		}
		return EachContinue
	})
	if !initial {
		return nil,ReduceError
	}
	return memo[0],""
}

var Inject func (obj []T, iterator func(T,T,T,T) T, memo ...T) (T,string) = Reduce
var FoldL  func (obj []T, iterator func(T,T,T,T) T, memo ...T) (T,string) = Reduce


// The right-associative version of reduce, also known as `foldr`.
func ReduceRight (obj []T, iterator func(T,T,T,T) T, memo ...T) (T,string) {
	initial := len(memo) > 0
	if obj == nil {
		obj = make([]T,0)
	}
	length := len(obj)
	Each(obj, func (value T, index T, list T) bool {
		length = length - 1
		index = length
		if !initial {
			memo[0] = list.([]T)[index.(int)]
			initial = true
		} else {
			memo[0] = iterator(memo[0], list.([]T)[index.(int)], index, list)
		}
		return EachContinue
	})
	if !initial {
		return nil,ReduceError
	}
	return memo[0],""
}

var FoldR  func (obj []T, iterator func(T,T,T,T) T, memo ...T) (T,string) = ReduceRight

func IdentityEach ( val T, index T, list T ) bool {
	return val == val
}

func IdentityIsTruthy( val T, index T, list T ) bool {
	v,ok := val.(bool)
	if ok {
		return  v
	}
	sv,sok := val.(string)
	if sok {
		return  sv != ""
	}
	iv,iok := val.(int)
	if iok {
		return  iv != 0
	}
	return val != nil
}

func Identity ( val T, index T, list T ) T {
	return val
}

// Determine if at least one element in the object matches a truth test.
// Aliased as `some`.
func Any (obj []T, opt_predicate ...func(val T,index T, list T)bool ) bool {
	var predicate func(T,T,T)bool
	if len(opt_predicate) == 0 {
		predicate = IdentityEach
	} else {
		predicate = opt_predicate[0]
	}	
	anyresult := false
	if obj == nil {
		return anyresult
	}

	eachFunc := func (value T, index T, list T) bool {
		if anyresult {
			return EachBreak
		}
		anyresult = predicate(value, index, list)
		if anyresult {
			return EachBreak
		}
		return EachContinue
	}
	Each(obj, eachFunc)
	return anyresult
}
var Some func(obj []T, opt_predicate ...func(val T,index T, list T)bool ) bool = Any


// Return the first value which passes a truth test. 
// Aliased as `detect`.
func Find (obj []T, predicate func(T,T,T) bool ) T {
	var result T
	Any(obj, func (value T, index T, list T) bool {
      if predicate(value, index, list) {
        result = value
        return EachBreak
      }
		return EachContinue
	})
	return result
}

var Detect func(obj []T, iterator func(T,T,T) bool ) T  = Find




// Return all the elements that pass a truth test.
// Aliased as `select`.
func Filter (obj []T, iterator eachlistiterator ) []T {
	results := make([]T,0)
	if obj == nil || len(obj) == 0 {
		return results
	}
	Each(obj, func (value T, index T, list T) bool {
      if iterator(value, index, list) {
			results = append(results , value)
		}
		return EachContinue
	})
	return results;
}

var Select func(obj []T, iterator eachlistiterator ) []T = Filter


// Return all the elements for which a truth test fails.
func Reject (obj []T, iterator eachlistiterator ) []T {
	return Filter(obj, func(value T, index T, list T) bool {
		return !iterator(value, index, list)
	})
}


func Every (obj []T, opt_iterator ...eachlistiterator ) bool {
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
	Each(obj, func (value T, index T, list T) bool {
		result = result && iterator(value, index, list)
		if ! result {
			return EachBreak
		}
		return EachContinue
	})
	return result
}

var All func(obj []T, opt_iterator ...eachlistiterator ) bool = Every


// Determine if the array or object contains a given value (using `==`).
// Aliased as `include`.
func Contains (obj []T, target T, opt_comparator ...func(T,T)bool) bool {
	if obj == nil {
		return false
	}
	var comparator func(T,T)bool
	if len(opt_comparator) > 0 {
		comparator = opt_comparator[0]
	}
	return Any(obj, func (value T, index T, list T) bool {
		if comparator != nil {
			return comparator(value,target)
		} else {
			return value == target
		}
	})
}

var Include func(obj []T, target T, opt_comparator ...func(T,T)bool) bool = Contains



// Invoke a method (with arguments) on every item in a collection.
//func Invoke(obj []T, method func(...T)) {
//	var args = slice.call(arguments, 2);
//	var isFunc = _.isFunction(method);
//	return Map(obj, func(value T) {
//		return (isFunc ? method : value[method]).apply(value, args);
//	});
//}


// Convenience version of a common use case of `map`: fetching a property.
func Pluck(obj []T, targetvalue T) []T {
	return Map(obj, func(testvalue T, index T , origlist T) T { 
		if IsMap(testvalue) {
			return testvalue.(map[T]T)[targetvalue]
		}
		if targetvalue == testvalue {
			return testvalue
		}
		return nil
	} )
}



// Convenience version of a common use case of `filter`: selecting only objects
// containing specific `key:value` pairs.
func Where(obj []T, attrs map[T]T, optReturnFirstFound ...bool) T {
	var returnFirstFound bool
	if len(optReturnFirstFound) == 0 {
		returnFirstFound = false
	} else {
		returnFirstFound = optReturnFirstFound[0]
	}
	if IsEmpty(attrs) {
		return make([]T,0)
	}
	if returnFirstFound {
		return Find(obj, func(value T, key T, list T) bool {
			for k,v := range attrs {
				if v != value.(map[T]T)[k] {
					return false
				}
			}
			return true
		})

	} else {
		return Filter(obj, func(value T, key T, list T) bool {
			for k,v := range attrs {
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
	return Where(obj,attrs,true)
}

/* 
TODO punting on these until I found a generic or better Go way of impl max/min
SEE https://groups.google.com/forum/#!topic/golang-nuts/f32UN1TYiAI
// Return the maximum element or (element-based computation).
_.max = function(obj, iterator, context) {
// Return the minimum element (or element-based computation).
_.min = function(obj, iterator, context) {
*/

// Shuffle an array, using the modern version of the
// [Fisher-Yates shuffle](http://en.wikipedia.org/wiki/Fisher–Yates_shuffle).
func Shuffle(obj []T) []T {
	shuffled := make([]T,len(obj))
	indices := rand.Perm(len(obj))
	for i,idx := range indices {
		shuffled[i],shuffled[idx] = obj[idx],obj[i]
	}
	return shuffled
}

// An internal function to generate lookup iterators
func lookupIterator (value T) func(obj T, idx T, list T) T {
	if IsFunction(value) { 
		//fmt.Printf("lookupIterator got a func\n")
		return value.(func(obj T, idx T, list T)T)
	}
	//fmt.Printf("lookupIterator didnt get a func\n")
	return func(obj T, idx T, list T) T {
		//fmt.Printf("inlookup iterator, got obj %v, idx %v, list %v\n",obj,idx,list)
		if IsMap(obj) {
			return obj.(map[T]T)[value]
		}
		return obj
	}
}


// An internal function used for aggregate "group by" operations.
func group (behavior func( result map[T]T, k T, v T)  ) func(o T,v T) map[T]T {
	return func(obj T, value T) map[T]T {
		result := make(map[T]T,0)
		var iterator func(T,T,T) T
		if value == nil {
			iterator = Identity
		} else {
			iterator = lookupIterator(value)
		}

		Each(obj, func(value T, index T, list T) bool {
			key := iterator(value, index, obj)
			//_,ok := result[key]
			behavior(result, key, value)
			return EachContinue
		})
		return result
	}
}


// Shortcut function for checking if an object has a given property directly
// on itself (in other words, not on a prototype).
func Has (obj T, key T) bool {
	_,ok := obj.(map[T]T)[key]
	return ok
}

// Groups the object's values by a criterion. Pass either a string attribute
// to group by, or a function that returns the criterion.
var GroupBy = group(func(result map[T]T, key T, value T) {
	if key == nil {
		return
	}
	if Has(result,key) {
		//fmt.Printf("in group, got res %v, key %v, val %v\n\n",result,key,value)
		slice := result[key].([]T)
		slice = append(slice , value )
		result[key] = slice
	} else {
		slice := make([]T,1)
		slice[0] = value
		result[key] = slice
	}
})

var IndexBy = group( func(result map[T]T, key T, value T) {
	if key == nil {
		return
	}
	result[key] = value
})

var CountBy = group( func(result map[T]T, key T, value T) {
	if key == nil {
		return
	}
	if Has(result,key) {
      //fmt.Printf("in group, got res %v, key %v, val %v\n\n",result,key,value)
      result[key] = (result[key]).(int) + 1
   } else {
      result[key] = 1
   }
})

func SortedIndex (array T, obj T, lessThan func(T,T)bool, opt_iterator ...func(T,T,T)T) int {
	var value T
	var iterator func(T,T,T) T
	_,isArrayOfMaps := array.([]map[T]T)
	_,isArray := array.([]T)
	if !(isArrayOfMaps || isArray) {
		fmt.Printf("Error: can't find sorted index of a non-list object, got %v\n",array)
		return math.MinInt64
	}
	if len(opt_iterator) > 0 {
		iterator = opt_iterator[0]
		value = iterator(obj,nil,nil)
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
		mid := uint(low + high) >> 1
		if iterator != nil {
			if isArrayOfMaps {
				otherValue = iterator(array.([]map[T]T)[mid],nil,nil)
			} else if isArray {
				otherValue = iterator(array.([]T)[mid],nil,nil)
			}
		} else {
			if isArrayOfMaps {
				otherValue = array.([]map[T]T)[mid]
			} else if isArray {
				otherValue = array.([]T)[mid]
			}
		}		
		if lessThan(otherValue,value) {
			low = int(mid) + 1 
		} else {
			high = int(mid)
		}
	}
	return low
}

// Safely create a real, live array from anything iterable.
func ToArray( obj T ) []T {
	if obj == nil {
		return make([]T,0)
	}
	if IsArray(obj) || IsArrayOfMaps(obj) || IsString(obj) {
		return Map(obj,Identity)
	}
	if IsMap(obj) {
		return Values(obj.(map[T]T))
	}
	fmt.Printf("Error: ToArray, got something I dont know what to do with %v\n",obj)
	return nil
}

//Return the number of elements in an object.
func Size(obj T) int {
	if IsEmpty( obj ) {
		return 0
	}	
	if IsArrayOfMaps( obj ) {
		return len(obj.([]map[T]T))
	}
	if IsArray( obj ) {
		return len(obj.([]T))
	}
	if IsMap( obj ) {
		return len(obj.(map[T]T))
	}
	if IsString( obj ) {
		return len(obj.(string))
	} else {
		fmt.Printf("TypeError (Size): what is this? %v\n",obj)
		return math.MinInt64
	}
	
}


// Array Functions

// Get the first element of an array. Passing **n** will return the first N
// values in the array. Aliased as `head` and `take`. The **guard** check
// allows it to work with `_.map`.
func FirstN (array []T, n int, opt_guard ...bool) []T {
	if array == nil {
		return nil
	}
	if n == 0 {
		retval := make([]T,0)
		retval = append(retval,array[0])
		return retval
	} else if len(opt_guard) > 0 && opt_guard[0] {
		retval := make([]T,0)
		retval = append(retval,array[0])
		return retval
	}
	if n > len(array) {
		return array[:]
	} else {
		return array[0:n]
	}
}

var HeadN func(array []T, n int, opt_guard ...bool) []T = FirstN
var TakeN func(array []T, n int, opt_guard ...bool) []T = FirstN

func First(array []T) T {
	if array == nil {
		return nil
	}
	return array[0]
}

var Head func(array []T) T = First
var Take func(array []T) T = First


func Initial(array []T , opt_n ...int) []T {
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


func Last(array []T , opt_n ...int) []T {
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
	return array[(arraylen-n):]
}


// Returns everything but the first entry of the array. Aliased as `tail` and `drop`.
// Especially useful on the arguments object. Passing an **n** will return
// the rest N values in the array. The **guard**
// check allows it to work with `_.map`.
func Rest (array []T) []T {
	if array == nil {
		return nil
	}
	dst := make([]T,len(array)-1)
	copy(dst, array[1:])
	return dst
}

var Tail func(array []T) []T = Rest
var Drop func(array []T) []T = Rest

func Compact(array []T) []T {
	return Filter(array,IdentityIsTruthy)
}

// Internal implementation of a recursive `flatten` function.
func flatten(input []T, shallow bool, output []T) []T {
	//if shallow && Every(input, IsArrayEach) {
	//	fmt.Printf("everything is array and shallow: %v\n",input)
	//	output = append(output,input...)
	//	return output
	//}
	Each(input, func(value T, idx T, list T) bool  {
		if IsArray(value) {
			if shallow {
				//fmt.Printf("output before: %v\n",output)
				output = append(output,value.([]T)...)
				//fmt.Printf("output after : %v\n",output)
			} else {
				output = flatten(value.([]T), shallow, output)
			}
		} else {
			//fmt.Printf("is not array %v output before: %v\n",value, output)
			output = append(output,value)
		}
		return EachContinue
	})
	return output
}

// Flatten out an array, either recursively (by default), or just one level.
func Flatten (array []T, opt_shallow ...bool) []T {
	var shallow bool
	if len(opt_shallow) > 0 {
		shallow = opt_shallow[0]
	}
   return flatten(array, shallow, make([]T,0))
}


// Take the difference between one array and a number of other arrays.
// Only the elements present in just the first array will remain.
func Difference (toRemove []T, opt_from ...[]T) []T {
	if len(opt_from) == 0 {
		return make([]T,0)
	}
	var rest []T = make([]T,0)
	for _,from := range opt_from {
		rest = flatten(from,true,rest)
	}
	return Filter(toRemove, func(value T, idx T, list T) bool { 
		return ! Contains(rest, value) 
	})
}

func Without (toRemove []T, opt_from ...T) []T {
	if len(opt_from) == 0 {
		return make([]T,0)
	}
	var rest []T = make([]T,0)
	for _,from := range opt_from {
		rest = append(rest,from)
	}
	return Difference( toRemove, rest )
}

// Produce a duplicate-free version of the array. If the array has already
// been sorted, you have the option of using a faster algorithm.
// Aliased as `unique`.
func Uniq(list T, isSorted T /*bool or func*/, opt_iterator ...T) []T {
	var array []T
	var arrayofmaps []map[T]T
	isAM := IsArrayOfMaps(list)
	isA  := IsArray(list)
	if isAM {
		arrayofmaps = list.([]map[T]T)
	} else if isA {
		array = list.([]T)
	}

	var iterator mapiterator
	var comparator func(T,T) bool
	if IsFunction(isSorted) {
		iterator = isSorted.(mapiterator)
		isSorted = false
		if len(opt_iterator) > 0 {
			comparator = opt_iterator[0].(func(T,T)bool)
		}
	} else if len(opt_iterator) > 0 {
		iterator =   opt_iterator[0].(func(T,T,T)T)
		comparator = opt_iterator[1].(func(T,T)bool)
	}
	var initialA []T
	if iterator != nil {
		if isA {
			initialA = Map(array, iterator)
		} else if isAM {
			initialA = Map(arrayofmaps, iterator)
		}
	} else {
		if isA {
			initialA = array
		}
	}
   results := make([]T,0)
   seen := make([]T,0)
	if isA {
		Each(initialA, func(value T, index T, list T) bool {
			if isSorted.(bool) {
				if index == 0 || seen[ len(seen) - 1] != value {
					seen = append(seen,value)
					results = append( results, array[index.(int)])
				}
			} else if ! Contains(seen, value) {
			  seen = append( seen , value)
			  results = append( results, array[index.(int)])
			}
			return EachContinue
		})
	}
	if isAM {
		Each(arrayofmaps, func(value T, index T, list T) bool {
			if isSorted.(bool) {
				if index == 0 || seen[ len(seen) - 1] != value {
					seen = append(seen,value)
					results = append( results, array[index.(int)])
				}
			} else if ! Contains(seen, value, comparator) {
			  seen = append( seen , value)
			  results = append( results, value )
			}
			return EachContinue
		})
	}
   return results
}

var Unique func(list T, isSorted T /*bool or func*/, opt_iterator ...T) []T  = Uniq

// Produce an array that contains the union: each distinct element from all of
// the passed-in arrays.
func Union (opt_array ...T) []T {
	return Uniq(Flatten(opt_array, true),false)
}

// Produce an array that contains every item shared between all the
// passed-in arrays.
func Intersection(lessThan func(T,T)bool,opt_array ...T) []T {
	rest := Uniq(Flatten(Rest(opt_array),true),false)
	return Filter(Uniq(opt_array[0],false), func(this T, idx T, list T) bool {
		return Every(rest, func(that T, idx2 T,list T) bool {
			return IndexOf(rest, this,lessThan) != -1
		})
	})
}

// Return the position of the first occurrence of an
// item in an array, or -1 if the item is not included in the array.
// If the array is large and already in sort order, pass `true`
// for **isSorted** to use binary search.
func IndexOf (array []T, item T, lessThan func(T,T) bool, isSorted ...bool) int {
	if array == nil{
		return -1
	}
	length := len(array)
	if length == 0 {
		return -1
	}
	i := 0
	// do binary search if isSorted = true
   if len(isSorted) > 0 && isSorted[0] {
		i = SortedIndex(array, item, lessThan )
		if array[i] == item {
			return i
		} else {
			return -1
		}
   }
	for i,v := range array {
		if v == item {
			return i
		}
	}
	return -1
}

func LastIndexOf (array []T, item T, from ...int) int {
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




func Zip (arrays ...[]T ) []T {
	if arrays == nil || len(arrays) == 0 {
		return make([]T,0)
	}
	var length int = 0
	var tmplength int = 0
	var num_arrays int
	for _,array := range arrays {
		num_arrays += 1
		tmplength = len(array)
		if tmplength > length {
			length = tmplength
		}
	}
	//var retval [][]T
	retval := make([]T,length)
	for i := 0 ; i < length; i++ {
		zipped := make([]T,num_arrays)
		for j,array := range arrays {
			if i < len(array) {
				zipped[j] = array[i]
			}
		}
		retval[i] = zipped
	}
	return retval	
}

func Object( pairs_or_two_arrays ...[]T ) map[T]T {
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
			for i := 0 ; i < length ; i ++  {
				retval[ kvpairs[i].([]T)[0] ] = kvpairs[i].([]T)[1]
			}
		} else {
			for i := 0 ; i < length ; i += 2 {
				retval[ kvpairs[i] ] = kvpairs[i+1]
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
	for i := 0 ; i < length ; i ++ {
		retval[ keys[i] ] = values[i]	
	}
	return retval
}

// Generate an integer Array containing an arithmetic progression. A port of
// Underscore's range() which is a port of the native Python `range()` function. See
// [the Python documentation](http://docs.python.org/library/functions.html#range).
func Range (start_stop_and_step ...int) []int {
	if start_stop_and_step == nil {
		return make([]int,0)
	}
	var start int
	var stop int
	var step int
	argslength := len(start_stop_and_step)
	if argslength == 3 {
		start,stop,step = start_stop_and_step[0],start_stop_and_step[1],start_stop_and_step[2]
	} else if argslength == 2 {
		start,stop = start_stop_and_step[0],start_stop_and_step[1]
		step = 1
	} else if argslength == 1 {
		stop = start_stop_and_step[0]
		start = 0
		step = 1
	}
	
	if step == 0 {
		step = 1
	}

   length := int(math.Max(math.Ceil(float64(stop - start) / float64(step)), 0))
   idx := 0
	retval := make([]int,length)

	for idx < length {
		retval[idx] = start
      start += step
		idx += 1
    }

    return retval
}


// Map Functions

// Retrieve the names of a maps keys
func Keys(obj map[T]T) []T {
	retval := make([]T,0)
	if obj == nil {
		return retval
	}
	for key,_ := range obj {
		retval = append(retval, key)
	}
	return retval
}

// Retrieve the values of a maps keys
func Values(obj map[T]T) []T {
	retval := make([]T,0)
	if obj == nil {
		return retval
	}
	keys := Keys(obj)
	for _,key := range keys {
		retval = append(retval, obj[key])
	}
	return retval
}

// Convert an object into a list of `[key, value]` pairs.
func Pairs (obj map[T]T) []T {
	keys := Keys(obj)
	length := len(keys)
	pairs := make([]T,length)
	for i := 0; i < length; i++ {
		pairs[i] = []T{keys[i], obj[keys[i]] }
	}
	return pairs
}

// Function Functions

// Partially apply a function by creating a version that has had some of its
// arguments pre-filled, without changing its dynamic `this` context.
func Partial (fn func(...T) T , savedArgs ...T) func(...T) T {
	return func( laterArgs ...T ) T {
		//args := make([]T, len(savedArgs) + len(laterArgs))
		args := append( savedArgs, laterArgs... )
		//for i,v := range savedArgs {
			//args[i] = v
		//}
		//for j,v := range laterArgs {
			//args[ len(savedArgs) + j ]= v
		//}
      return fn( args...)
	}
}
