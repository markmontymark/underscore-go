package underscore

import (
	"fmt"
	"testing"
)



func TestEach( t *testing.T ) {
	fmt.Printf(" did we get this far? ")

	silly := func( elem T, i int, list []T ) bool {
		fmt.Printf("in silly with elem %d, index %d, list %v\n", elem, i, list)
		return EachContinue
	}

	sillymap := func( val T, key T, elem map[T]T ) bool {
		fmt.Printf("in silly with val %d, key %v, elem %v\n", val, key, elem)
		return EachContinue
	}

	EachArray([]T{0,3,2,1}, silly )
	EachStruct(map[T]T{"a":1,"b":2,"c":3}, sillymap )
}
