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
func Map(obj []T, iterator func(T,int,[]T) T, context T) []T {
	results := make([]T,0)
	if obj == nil {
		return results
	}
	Each(obj, func (value T, index int, list []T) bool {
		results = append( results, iterator( value, index, list))
		return EachContinue
	})
	return results
}


func MapMap(obj map[T]T, iterator func(T,T,map[T]T) T, context T) []T {
	results := make([]T,0)
	if obj == nil {
		return results
	}
	EachMap(obj, func (value T, key T , obj map[T]T) bool {
		results = append( results, iterator( value, key, obj))
		return EachContinue
	})
	return results
}

var Collect func (obj []T, iterator func(T,int,[]T) T, context T) []T = Map
var CollectMap func (obj map[T]T, iterator func(T,T,map[T]T) T, context T) []T = MapMap

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
// Aliased as `any`.
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

