package underscore

import (
	//"fmt"
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
	//fmt.Printf("Is Array ok = %v\n",ok)
	return v != nil
}

// Is a given variable a map
func IsMap (obj T) bool {
	v,_ := obj.(map[T]T) 
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
	return true
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

	} else if IsMap( elemslist_or_map ) {
		for k,v := range elemslist_or_map.(map[T]T) {
			if iterator(v,k,elemslist_or_map.(map[T]T)) == EachBreak {
				return
			}
		}
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


/*
// Convenience version of a common use case of `filter`: selecting only objects
// containing specific `key:value` pairs.
func Where(obj map[T]T, attrs map[T]T, returnFirstFound bool) []T {
	if IsEmpty(attrs) {
		return make([]T,0)
	}
	if returnFirstFound {
		return Find(obj, func(value map[T]T, index int, list[]T) bool {
			for k,v := range attrs {
				if v != value[k] {
					return false
				}
			}
			return true
		})

	} else {
		return Filter(obj, func(value map[T]T, index int, list[]T) bool {
			for k,v := range attrs {
				if v != value[k] {
					return false
				}
			}
			return true
		})
	}
}
*/
