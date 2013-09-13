package underscore

import (
	//"fmt"
)

type Underscore struct {}

type T interface{}

const EachContinue bool = false
const EachBreak    bool = true


// TODO, so far, I've morphed _.each into EachArray and EachStruct, looking for a way to switch  on elems/elem and then can merge the two

type eachlistiterator func(T,int,[]T) bool

func Each(elems []T, iterator eachlistiterator ) {
	if elems == nil {
		return
	} 
	for i,elem := range elems {
		if iterator(elem, i, elems) == EachBreak {
			return;
		}
	}
}

func EachMap(elem map[T]T, iterator func(T,T,map[T]T) bool ) {
	if elem == nil {
		return
	} 
	for k,v := range elem {
		if iterator(v,k,elem) == EachBreak {
			return;
		}
	}
}


// Return the results of applying the iterator to each element.
func Map(obj []T, iterator func(T,int,[]T) T) []T {
	results := make([]T,0)
	if obj == nil {
		return results
	}
	Each(obj, func (value T, index int, list []T) bool {
		if v := iterator(value,index,list); v != nil {
			results = append( results, v )
		}
		return EachContinue
	})
	return results
}


func MapMap(objlist []map[T]T, iterator func(T,T,map[T]T) T) []T {
	results := make([]T,0)
	if objlist == nil {
		return results
	}

	for _,obj := range objlist {
		EachMap( obj, func (value T, key T, origobj map[T]T ) bool {
			if v := iterator(value,key,origobj) ; v != nil {
				results = append( results, v)
			}
			return EachContinue
		})
	}
	return results
}

var Collect func (obj []T, iterator func(T,int,[]T) T) []T = Map
var CollectMap func (obj []map[T]T, iterator func(T,T,map[T]T) T ) []T = MapMap

const ReduceError = "Reduce of empty array with no initial value"

// **Reduce** builds up a single result from a list of values, aka `inject`,
// or `foldl`. 
func Reduce (obj []T, iterator func(T,T,int,[]T) T, memo ...T) (T,string) {
	initial := len(memo) > 0
	if obj == nil {
		obj = make([]T,0)
	}
	Each(obj, func (value T, index int, list []T) bool {
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

var Inject func (obj []T, iterator func(T,T,int,[]T) T, memo ...T) (T,string) = Reduce
var FoldL  func (obj []T, iterator func(T,T,int,[]T) T, memo ...T) (T,string) = Reduce


// The right-associative version of reduce, also known as `foldr`.
func ReduceRight (obj []T, iterator func(T,T,int,[]T) T, memo ...T) (T,string) {
	initial := len(memo) > 0
	if obj == nil {
		obj = make([]T,0)
	}
	length := len(obj)
	Each(obj, func (value T, index int, list []T) bool {
		length = length - 1
		index = length
		if !initial {
			memo[0] = list[index]
			initial = true
		} else {
			memo[0] = iterator(memo[0], list[index], index, list)
		}
		return EachContinue
	})
	if !initial {
		return nil,ReduceError
	}
	return memo[0],""
}

var FoldR  func (obj []T, iterator func(T,T,int,[]T) T, memo ...T) (T,string) = ReduceRight

func IdentityEach ( val T, index int, list[]T ) bool {
	return val == val
}

// Determine if at least one element in the object matches a truth test.
// Aliased as `some`.
func Any (obj []T, opt_iterator ...func(val T,index int, list[]T)bool ) bool {
	var iterator func(T,int, []T)bool
	if len(opt_iterator) == 0 {
		iterator = IdentityEach
	} else {
		iterator = opt_iterator[0]
	}	
	anyresult := false
	if obj == nil {
		return anyresult
	}

	eachFunc := func (value T, index int, list []T) bool {
		if anyresult {
			return EachBreak
		}
		anyresult = iterator(value, index, list)
		if anyresult {
			return EachBreak
		}
		return EachContinue
	}
	Each(obj, eachFunc)
	return anyresult
}
var Some func(obj []T, opt_iterator ...func(val T,index int, list[]T)bool ) bool = Any


// Return the first value which passes a truth test. 
// Aliased as `detect`.
func Find (obj []T, iterator func(T,int,[]T) bool ) T {
	var result T
	Any(obj, func (value T, index int, list []T) bool {
		//fmt.Printf("in Any(%v,%v,...)\n",value,index)
      if iterator(value, index, list) {
			//fmt.Printf("in Any iterator, it returned true for value = %v\n",value)
        result = value
        return EachBreak
      }
		//fmt.Printf("in Any iterator, else returned value so continue n",value)
		return EachContinue
	})
	return result
}

var Detect func(obj []T, iterator func(T,int,[]T) bool ) T  = Find




// Return all the elements that pass a truth test.
// Aliased as `select`.
func Filter (obj []T, iterator eachlistiterator ) []T {
	results := make([]T,0)
	if obj == nil || len(obj) == 0 {
		return results
	}
	Each(obj, func (value T, index int, list []T) bool {
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
	return Filter(obj, func(value T, index int, list []T) bool {
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
	Each(obj, func (value T, index int, list []T) bool {
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
	return Any(obj, func (value T, index int, list []T) bool {
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
	return Map(obj, func(testvalue T, index int , origlist[]T) T { 
		if targetvalue  == testvalue {
			return testvalue
		}
		return nil
	} )
}

// Convenience version of a common use case of `map`: fetching a property.
func PluckMap(obj []map[T]T, targetkey T) []T {
  return MapMap(obj, func(value T, testkey T, origobj map[T]T) T { 
    if targetkey == testkey {
      return value
    }
    return nil
  } )
} 
