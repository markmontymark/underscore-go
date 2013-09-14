package underscore

import (
	"./lib/asserts"
	"fmt"
	//"strings"
	"testing"
)



func TestEach( t *testing.T ) {

	sliceCollector := make([]string,0)
	silly := func( elem T, i T, list T ) bool {
		sliceCollector = append(sliceCollector, fmt.Sprintf("each: elem %d, index %v, list %v,", elem, i, list) )
		return EachContinue
	}

	Each([]T{0,3,2,1}, silly )

	asserts.Equals( t, "Test Each(0,3,2,4)", "[each: elem 0, index 0, list [0 3 2 1], each: elem 3, index 1, list [0 3 2 1], each: elem 2, index 2, list [0 3 2 1], each: elem 1, index 3, list [0 3 2 1],]", fmt.Sprintf("%v",sliceCollector) )
}

func TestMap( t *testing.T ) {

	add2 := func( elem T, i T, list T ) T {
		return elem.(int) + 2
	}

	identityValueMap := func( value T, key T, obj T ) T { return value }
	identityKeyMap := func( value T, key T, obj T ) T { return key }
	identityNestedKeyMap := func( value T, key T, obj T ) T { return value.(map[T]T)["a"] }
	identityNestedKeyShorterMap := func( value T, key T, obj T ) T { return value.(map[T]T)["b"] }

	asserts.Equals( t, "testing map with array",     
		"[2 5 4 3]", fmt.Sprintf("%v", Map([]T{0,3,2,1}, add2 ) ))
	asserts.Equals( t, "testing collect with array", 
		"[2 5 4 3]", fmt.Sprintf("%v", Collect([]T{0,3,2,1}, add2) ))

	amap :=  map[T]T{ "a":3,"b":2,"c":1 }
	list := make([]T,0)
	list = append( list, amap )
	list = append( list, map[T]T{ "a":1,"d":4,"e":5 } )

	asserts.Equals( t, "testing collectmap with map[string]int ", 
		"[3 2 1]",
		fmt.Sprintf("%v", Collect( amap , identityValueMap ) ))

	asserts.Equals( t, "testing collectmap with map[string]int ", 
		"[map[a:3 b:2 c:1] map[a:1 d:4 e:5]]",
		fmt.Sprintf("%v", Collect( list, identityValueMap ) ))

	asserts.Equals( t, "testing collectmap with map[string]int ", 
		"[a b c]",
		fmt.Sprintf("%v", Collect( amap , identityKeyMap ) ))

	asserts.Equals( t, "testing collectmap with map[string]int ", 
		"[0 1]",
		fmt.Sprintf("%v", Collect( list, identityKeyMap ) ))

	asserts.Equals( t, "testing collectmap with map[string]int ", 
		"[3 1]",
		fmt.Sprintf("%v", Collect( list, identityNestedKeyMap ) ))

	asserts.Equals( t, "testing collectmap with map[string]int ", 
		"[2]",
		fmt.Sprintf("%v", Collect( list, identityNestedKeyShorterMap ) ))
}


func TestReduce( t *testing.T ) {
	
	v,err :=	Reduce( 
		[]T{1,2,3}, 
		func (sum T, num T, i T, list T) T { return sum.(int) + num.(int) }, 
		0)
	if err != "" {
		fmt.Printf("FAIL: %s\n", err )
		return
	}
	asserts.IntEquals( t, "can reduce sum up an array", 6,	v.(int))

	v,err =	Reduce( 
		[]T{1,2,3}, 
		func (sum T, num T, i T, list T) T { return sum.(int) * num.(int) }, 
		3)
	if err != "" {
		fmt.Printf("FAIL: %s\n", err )
		return
	}
	asserts.IntEquals( t, "can reduce multiply up an array", 18,	v.(int))

	v,err =	Inject( 
		[]T{1,2,3}, 
		func (sum T, num T, i T, list T) T { return sum.(int) * num.(int) }, 
		3)
	if err != "" {
		fmt.Printf("FAIL: %s\n", err )
		return
	}
	asserts.IntEquals( t, "can inject multiply up an array", 18,	v.(int))

	v,err =	Inject( 
		[]T{1,2,3}, 
		func (sum T, num T, i T, list T) T { return sum.(int) * num.(int) }, 
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
		func (sum T, num T, i T, list T) T { return sum.(string) + "," + num.(string) }, 
		"")
	if err != "" {
		fmt.Printf("FAIL: %s\n", err )
		return
	}
	asserts.Equals( t, "can ReduceRight divide up an array", ",4,3,2",	v.(string))
}

func TestFind( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := Find(array, func(n T, i T, list T) bool { 
		return n.(int) > 2 
	})
	asserts.IntEquals( t, "should return first found `value`", 3, v.(int))

	v = Find(array, func(n T, i T, list T) bool { return false })
	asserts.Nil( t, "should return `nil` if `value` is not found", v)
}
func TestDetect_asFindAlias( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := Detect(array, func(n T, i T, list T) bool { 
		return n.(int) > 2 
	})
	asserts.IntEquals( t, "should return first found `value`", 3, v.(int))

	v = Detect(array, func(n T, i T, list T) bool { return false })
	asserts.Nil( t, "should return `nil` if `value` is not found", v)
}

func TestFilter( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := Filter(array, func(n T, i T, list T) bool { 
		return n.(int) > 2 
	})
	asserts.Equals( t, "should return last two values: 3 4", fmt.Sprintf("%v",[]T{3,4}), fmt.Sprintf("%v",v))
}
func TestSelect_asFilterAlias( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := Select(array, func(n T, i T, list T) bool { 
		return n.(int) > 2 
	})
	asserts.Equals( t, "should return last two values: 3 4", fmt.Sprintf("%v",[]T{3,4}), fmt.Sprintf("%v",v))
}

func TestReject( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := Reject(array, func(n T, i T, list T) bool { 
		return n.(int) > 2 
	})
	asserts.Equals( t, "should return first two values: 1 2", fmt.Sprintf("%v",[]T{1,2}), fmt.Sprintf("%v",v))
}

func TestEvery( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := Every(array, func(n T, i T, list T) bool { 
		return n.(int) < 5 
	})
	asserts.Equals( t, "should return true as all values: 1 2 3 4 are less than 5", "true", fmt.Sprintf("%v",v))

	v = Every(array, func(n T, i T, list T) bool { 
		return n.(int) < 4 
	})
	asserts.Equals( t, "should return false as not all values are < 4", "false", fmt.Sprintf("%v",v))
}

func TestAll( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := All(array, func(n T, i T, list T) bool { 
		return n.(int) < 5 
	})
	asserts.Equals( t, "should return true as all values: 1 2 3 4 are less than 5", "true", fmt.Sprintf("%v",v))

	v = All(array, func(n T, i T, list T) bool { 
		return n.(int) < 4 
	})
	asserts.Equals( t, "should return false as not all values are < 4", "false", fmt.Sprintf("%v",v))
}

func TestAny( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := Any(array, func(n T, i T, list T) bool { 
		return n.(int) < 5 
	})
	asserts.Equals( t, "should return true as at least one value is less than 5", "true", fmt.Sprintf("%v",v))

	v = Any(array, func(n T, i T, list T) bool { 
		return n.(int) < 4 
	})
	asserts.Equals( t, "should return true as at least one value is less than 4", "true", fmt.Sprintf("%v",v))

	v = Any(array, func(n T, i T, list T) bool { 
		return n.(int) < 3 
	})
	asserts.Equals( t, "should return true as at least one value is less than 3", "true", fmt.Sprintf("%v",v))

	v = Any(array, func(n T, i T, list T) bool { 
		return n.(int) < 2 
	})
	asserts.Equals( t, "should return true as at least one value is less than 2", "true", fmt.Sprintf("%v",v))

	v = Any(array, func(n T, i T, list T) bool { 
		return n.(int) < 1 
	})
	asserts.Equals( t, "should return false as no value are less than 1", "false", fmt.Sprintf("%v",v))
}

func TestContains( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := Contains(array, 1 )
	asserts.Equals( t, "should return true as array contains a 1", "true", fmt.Sprintf("%v",v))

	v  = Contains(array, 2 )
	asserts.Equals( t, "should return true as array contains a 1", "true", fmt.Sprintf("%v",v))

	v  = Contains(array, 3 )
	asserts.Equals( t, "should return true as array contains a 1", "true", fmt.Sprintf("%v",v))

	v  = Contains(array, 4 )
	asserts.Equals( t, "should return true as array contains a 1", "true", fmt.Sprintf("%v",v))

	v  = Contains(array, 5 )
	asserts.Equals( t, "should return false as array doesnt contain a 5", "false", fmt.Sprintf("%v",v))
}

func TestInclude( t *testing.T ) {
	array := []T{1, 2, 3, 4}
	v := Include(array, 1 )
	asserts.Equals( t, "should return true as array contains a 1", "true", fmt.Sprintf("%v",v))

	v  = Include(array, 2 )
	asserts.Equals( t, "should return true as array contains a 1", "true", fmt.Sprintf("%v",v))

	v  = Include(array, 3 )
	asserts.Equals( t, "should return true as array contains a 1", "true", fmt.Sprintf("%v",v))

	v  = Include(array, 4 )
	asserts.Equals( t, "should return true as array contains a 1", "true", fmt.Sprintf("%v",v))

	v  = Include(array, 5 )
	asserts.Equals( t, "should return false as array doesnt contain a 5", "false", fmt.Sprintf("%v",v))
}


func TestPluck( t *testing.T ) {
	people := []T{ "name" , "moe", "age" , 30, "name" , "curly", "age" , 50}
	v := Pluck(people,"name")
	asserts.Equals( t, "pulls names out of objects",
		"[name name]",
		fmt.Sprintf("%v",v))

	v  = Pluck(people,30)
	asserts.Equals( t, "pulls 30 out of list",
		"[30]",
		fmt.Sprintf("%v",v))

	people = make([]T,0)
	people = append( people, map[T]T{ "name" : "moe", "age" : 30} )
	people = append( people, map[T]T{ "name" : "curly", "age" : 50} )
	v  = Pluck(people,"name")
	asserts.Equals( t, "pulls names out of objects",
		"[moe curly]",
		fmt.Sprintf("%v",v))
}

func TestIsArrayWithArray( t *testing.T ) {
	list := []T{ "name" , "moe", "age" , 30}
	asserts.True( t, "Testing IsArray", IsArray( list) )
}

func TestIsArrayWithString( t *testing.T ) {
	scalar := "name"
	asserts.False( t, "Testing IsArray", IsArray( scalar ) )
}

func TestIsArrayWithMap( t *testing.T ) {
	mapp := make(map[string]int,0)
	asserts.False( t, "Testing IsArray", IsArray( mapp ) )
}

func TestIsStringWithArray( t *testing.T ) {
	list := []T{ "name" , "moe", "age" , 30}
	asserts.False( t, "Testing IsString", IsString( list) )
}

func TestIsStringWithString( t *testing.T ) {
	scalar := "name"
	asserts.True( t, "Testing IsString", IsString( scalar ) )
}

func TestIsStringWithMap( t *testing.T ) {
	mapp := make(map[string]int,0)
	asserts.False( t, "Testing IsString", IsString( mapp ) )
}

func TestIsMapWithArray( t *testing.T ) {
	list := []T{ "name" , "moe", "age" , 30}
	asserts.False( t, "Testing IsMap with array ", IsMap( list) )
}

func TestIsMapWithString( t *testing.T ) {
	scalar := "name"
	asserts.False( t, "Testing IsMap with string", IsMap( scalar ) )
}

func TestIsMapWithMap( t *testing.T ) {
	mapp := make(map[T]T,0)
	asserts.True( t, "Testing IsMap", IsMap( mapp ) )
}


func TestWhere( t *testing.T ) {

	list := make([]T,0)
	list = append(list,map[T]T{"a": 1, "b": 2})
	list = append(list,map[T]T{"a": 2, "b": 2})
	list = append(list,map[T]T{"a": 1, "b": 3})
	list = append(list,map[T]T{"a": 1, "b": 4})

	v := Where(list, map[T]T{"a": 1})
	asserts.Equals( t, "Find objects with key 'a':1", "3", fmt.Sprintf("%v",len(v.([]T))))
	asserts.Equals( t, "Last found has a 'b':4", "4", fmt.Sprintf("%v", v.([]T)[ len(v.([]T)) - 1].(map[T]T)["b"]))

	v  = Where(list,map[T]T{"b":2})
	asserts.Equals( t, "Find objects with 'b':2", "2", fmt.Sprintf("%v",len(v.([]T))))

	v  = Where(list,map[T]T{"b":2},true)
	asserts.Equals( t, "Find objects with 'b':2", "map[a:1 b:2]", fmt.Sprintf("%v",v))
}

func TestFindWhere( t *testing.T ) {

	list := make([]T,0)
	list = append(list,map[T]T{"a": 1, "b": 2})
	list = append(list,map[T]T{"a": 2, "b": 2})
	list = append(list,map[T]T{"a": 1, "b": 3})
	list = append(list,map[T]T{"a": 1, "b": 4})
	v := FindWhere(list, map[T]T{"a": 1})
	asserts.Equals( t, "Find first object with key 'a':1", "map[a:1 b:2]", fmt.Sprintf("%v",v))
}

