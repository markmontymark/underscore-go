package underscore

import (
	//"fmt"
)

type Underscore struct {}

type T interface{}

const EachContinue bool = false
const EachBreak    bool = true


// TODO, so far, I've morphed _.each into EachArray and EachStruct, looking for a way to switch  on elems/elem and then can merge the two

func Each(elems []T, iterator func(T,int,[]T) bool ) {
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



