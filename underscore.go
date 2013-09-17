package underscore

import (
	"fmt"
	"math"
	"math/rand"
)

type Underscore struct {}

type T interface{}

const EachContinue bool = false
const EachBreak    bool = true


// Is a given value an array?
func IsString (obj T) bool {
	v,_ := obj.(string)
	return v != ""
}

// Is a given value an array?
func IsArray (obj T) bool {
	v,_ := obj.([]T)
	return v != nil
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


type eachlistiterator func(T,T,T) bool


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
		fmt.Printf("Each isnt doing anything useful\n")
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
func Contains (obj []T, target T) bool {
	if obj == nil {
		return false
	}
	return Any(obj, func (value T, index T, list T) bool {
		return value == target
	})
}

var Include func(obj []T, target T) bool = Contains



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


