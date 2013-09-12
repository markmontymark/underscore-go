package underscore

import (
	"fmt"
	"./lib/asserts"
	"testing"
)



func TestEach( t *testing.T ) {

	sliceCollector := make([]string,0)
	silly := func( elem T, i int, list []T ) bool {
		sliceCollector = append(sliceCollector, fmt.Sprintf("each: elem %d, index %d, list %v,", elem, i, list) )
		return EachContinue
	}

	mapCollector := make([]string,0)
	sillymap := func( val T, key T, elem map[T]T ) bool {
		mapCollector = append( mapCollector, fmt.Sprintf("eachmap: val %d, key %v, elem %v,", val, key, elem) )
		return EachContinue
	}

	Each([]T{0,3,2,1}, silly )
	EachMap(map[T]T{"a":1,"b":2,"c":3}, sillymap )

	asserts.Equals( t, "Test Each(0,3,2,4)", "[each: elem 0, index 0, list [0 3 2 1], each: elem 3, index 1, list [0 3 2 1], each: elem 2, index 2, list [0 3 2 1], each: elem 1, index 3, list [0 3 2 1],]", fmt.Sprintf("%v",sliceCollector) )
	asserts.Equals( t, "Test EachMap({\"a\":1,\"b\":2,\"c\":3})", "[eachmap: val 1, key a, elem map[a:1 b:2 c:3], eachmap: val 2, key b, elem map[a:1 b:2 c:3], eachmap: val 3, key c, elem map[a:1 b:2 c:3],]", fmt.Sprintf("%v",mapCollector) )
}

func TestMap( t *testing.T ) {

	add2slice := func( elem T, i int, list []T ) T {
		return elem.(int) + 2
	}

	add2map := func( value T, key T, obj map[T]T ) T {
		return value.(int) + 2
	}

	identityKeyMap := func( value T, key T, obj map[T]T ) T {
		return key
	}

	//sillymap := func( val T, key T, elem map[T]T ) bool {
		//fmt.Printf("in silly with val %d, key %v, elem %v\n", val, key, elem)
		//return EachContinue
	//}

	asserts.Equals( t, "testing map with array", 
		"[2 5 4 3]",
		fmt.Sprintf("%v", Map([]T{0,3,2,1}, add2slice ,nil ) ))

	asserts.Equals( t, "testing map with map[string]int ", 
		"[3 4 5]",
		fmt.Sprintf("%v", MapMap(map[T]T{"a":1,"b":2,"c":3}, add2map , nil ) ))

	asserts.Equals( t, "testing map with map[string]int ", 
		"[a b c]",
		fmt.Sprintf("%v", MapMap(map[T]T{"a":1,"b":2,"c":3}, identityKeyMap , nil ) ))

	asserts.Equals( t, "testing collect with array", 
		"[2 5 4 3]",
		fmt.Sprintf("%v", Collect([]T{0,3,2,1}, add2slice ,nil ) ))

	asserts.Equals( t, "testing collectmap with map[string]int ", 
		"[3 4 5]",
		fmt.Sprintf("%v", CollectMap(map[T]T{"a":1,"b":2,"c":3}, add2map , nil ) ))

	asserts.Equals( t, "testing collectmap with map[string]int ", 
		"[a b c]",
		fmt.Sprintf("%v", CollectMap(map[T]T{"a":1,"b":2,"c":3}, identityKeyMap , nil ) ))


}


func TestReduce( t *testing.T ) {
	
	v,err :=	Reduce( 
		[]T{1,2,3}, 
		func (sum T, num T, i int, list []T) T { return sum.(int) + num.(int) }, 
		0)
	if err != "" {
		fmt.Printf("FAIL: %s\n", err )
		return
	}
	asserts.IntEquals( t, "can reduce sum up an array", 6,	v.(int))

	v,err =	Reduce( 
		[]T{1,2,3}, 
		func (sum T, num T, i int, list []T) T { return sum.(int) * num.(int) }, 
		3)
	if err != "" {
		fmt.Printf("FAIL: %s\n", err )
		return
	}
	asserts.IntEquals( t, "can reduce multiply up an array", 18,	v.(int))

	v,err =	Inject( 
		[]T{1,2,3}, 
		func (sum T, num T, i int, list []T) T { return sum.(int) * num.(int) }, 
		3)
	if err != "" {
		fmt.Printf("FAIL: %s\n", err )
		return
	}
	asserts.IntEquals( t, "can inject multiply up an array", 18,	v.(int))

	v,err =	Inject( 
		[]T{1,2,3}, 
		func (sum T, num T, i int, list []T) T { return sum.(int) * num.(int) }, 
		3)
	if err != "" {
		fmt.Printf("FAIL: %s\n", err )
		return
	}
	asserts.IntEquals( t, "can foldl multiply up an array", 18,	v.(int))

}

func TestReduceRight( t *testing.T ) {
	
	v,err :=	ReduceRight( 
		[]T{"2","3","4"}, 
		func (sum T, num T, i int, list []T) T { return sum.(string) + "," + num.(string) }, 
		"")
	if err != "" {
		fmt.Printf("FAIL: %s\n", err )
		return
	}
	asserts.Equals( t, "can ReduceRight divide up an array", ",4,3,2",	v.(string))
}

func TestFind( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := Find(array, func(n T, i int, list []T) bool { 
		return n.(int) > 2 
	})
	asserts.IntEquals( t, "should return first found `value`", 3, v.(int))

	v = Find(array, func(n T, i int, list []T) bool { return false })
	asserts.Nil( t, "should return `nil` if `value` is not found", v)
}
func TestDetect_asFindAlias( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := Detect(array, func(n T, i int, list []T) bool { 
		return n.(int) > 2 
	})
	asserts.IntEquals( t, "should return first found `value`", 3, v.(int))

	v = Detect(array, func(n T, i int, list []T) bool { return false })
	asserts.Nil( t, "should return `nil` if `value` is not found", v)
}

func TestFilter( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := Filter(array, func(n T, i int, list []T) bool { 
		return n.(int) > 2 
	})
	asserts.Equals( t, "should return last two values: 3 4", fmt.Sprintf("%v",[]T{3,4}), fmt.Sprintf("%v",v))
}
func TestSelect_asFilterAlias( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := Select(array, func(n T, i int, list []T) bool { 
		return n.(int) > 2 
	})
	asserts.Equals( t, "should return last two values: 3 4", fmt.Sprintf("%v",[]T{3,4}), fmt.Sprintf("%v",v))
}

func TestReject( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := Reject(array, func(n T, i int, list []T) bool { 
		return n.(int) > 2 
	})
	asserts.Equals( t, "should return first two values: 1 2", fmt.Sprintf("%v",[]T{1,2}), fmt.Sprintf("%v",v))
}

func TestEvery( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := Every(array, func(n T, i int, list []T) bool { 
		return n.(int) < 5 
	})
	asserts.Equals( t, "should return true as all values: 1 2 3 4 are less than 5", "true", fmt.Sprintf("%v",v))

	v = Every(array, func(n T, i int, list []T) bool { 
		return n.(int) < 4 
	})
	asserts.Equals( t, "should return false as not all values are < 4", "false", fmt.Sprintf("%v",v))
}

func TestAll( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := All(array, func(n T, i int, list []T) bool { 
		return n.(int) < 5 
	})
	asserts.Equals( t, "should return true as all values: 1 2 3 4 are less than 5", "true", fmt.Sprintf("%v",v))

	v = All(array, func(n T, i int, list []T) bool { 
		return n.(int) < 4 
	})
	asserts.Equals( t, "should return false as not all values are < 4", "false", fmt.Sprintf("%v",v))
}

func TestAny( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := Any(array, func(n T, i int, list []T) bool { 
		return n.(int) < 5 
	})
	asserts.Equals( t, "should return true as at least one value is less than 5", "true", fmt.Sprintf("%v",v))

	v = Any(array, func(n T, i int, list []T) bool { 
		return n.(int) < 4 
	})
	asserts.Equals( t, "should return true as at least one value is less than 4", "true", fmt.Sprintf("%v",v))

	v = Any(array, func(n T, i int, list []T) bool { 
		return n.(int) < 3 
	})
	asserts.Equals( t, "should return true as at least one value is less than 3", "true", fmt.Sprintf("%v",v))

	v = Any(array, func(n T, i int, list []T) bool { 
		return n.(int) < 2 
	})
	asserts.Equals( t, "should return true as at least one value is less than 2", "true", fmt.Sprintf("%v",v))

	v = Any(array, func(n T, i int, list []T) bool { 
		return n.(int) < 1 
	})
	asserts.Equals( t, "should return false as no value are less than 1", "false", fmt.Sprintf("%v",v))
}
