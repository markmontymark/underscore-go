package underscore

import (
	//"fmt"
)

type Underscore struct {}

type T interface{}

const EachContinue bool = false
const EachBreak    bool = true


// TODO, so far, I've morphed _.each into EachArray and EachStruct, looking for a way to switch  on elems/elem and then can merge the two

func EachArray(elems []T, iterator func(T,int,[]T) bool ) {
	if elems == nil {
		return
	} 
	for i,elem := range elems {
		if iterator(elem, i, elems) == EachBreak {
			return;
		}
	}
}

func EachStruct(elem map[T]T, iterator func(T,T,map[T]T) bool ) {
	if elem == nil {
		return
	} 
	for k,v := range elem {
		if iterator(v,k,elem) == EachBreak {
			return;
		}
	}
}

