package underscore

import (
	"./lib/asserts"
	"fmt"
	"testing"
)

var fib func ( ...T) T
func init(){
	fib = func(n ...T) T {
		if n[0].(int) < 2 {
			return n[0]
		}
		return fib(n[0].(int) - 1).(int) + fib(n[0].(int) - 2).(int)
	}
}


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

//func TestMax( t *testing.T ) {
	//list := []int{2,3,4,9,5,6,7,8}
	//asserts.Equals( t, "Find max element in array", "9", fmt.Sprintf("%v",MaxInt(list))) 
//}

/*
func TestShuffle( t *testing.T ) {
	list := []T{2,3,4,9,5,6,7,8}
	shuffledlist := Shuffle(list)
	asserts.IntEquals( t, "Find max element in array", len(list), len(shuffledlist))
	asserts.Equals( t, "sort orig list and shuffled list", 
		fmt.Sprintf("%v",sort.Sort(list)),
		fmt.Sprintf("%v",sort.Sort(shuffledlist)))
}
*/

func TestGroupBy( t *testing.T ) {

	data := GroupBy([]T{1, 2, 3, 4, 5, 6,1}, func(obj T,key T,val T) T{ 
		//fmt.Printf("group by func got obj %v, key %v, val %v\n", obj,key,val)
		//fmt.Printf("group by func, returning %v\n", (val.(int) % 2) )
		return obj.(int) % 2 
	})
	asserts.Equals( t, "group ints ", "parity map[1:[1 3 5 1] 0:[2 4 6]]",
		fmt.Sprintf("parity %v",data))
	asserts.Equals( t, "group evens ", "[2 4 6]", fmt.Sprintf("%v",data[0]))


	data2 := []T{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
	grouped := GroupBy( data2, func (obj T,key T,val T) T { return len(obj.(string)) })
 
	asserts.Equals( t, "grouping words of length 3", 
		fmt.Sprintf("%v", grouped[3]), "[one two six ten]")
	asserts.Equals( t, "grouping words of length 4", 
		fmt.Sprintf("%v", grouped[4]), "[four five nine]")
	asserts.Equals( t, "grouping words of length 5", 
		fmt.Sprintf("%v", grouped[5]), "[three seven eight]")

	data3 := []map[T]T{ {"a": 1, "b":2}, {"b":3}, {"a":4,"c":5},{"a":1,"b":7,"c":8} }
	grouped3 := GroupBy( data3, func (obj T,key T,val T) T { 
		return obj.(map[T]T)["a"] 
	})
	asserts.Equals( t, "group by an object key\\'s value", "map[1:[map[a:1 b:2] map[a:1 b:7 c:8]] 4:[map[a:4 c:5]]]",fmt.Sprintf("%v",grouped3))

	grouped4 := IndexBy( data3, func (obj T,key T,val T) T { 
		return obj.(map[T]T)["a"] 
	})
	asserts.Equals( t, "index by an object key\\'s value", "map[1:map[a:1 b:7 c:8] 4:map[a:4 c:5]]",fmt.Sprintf("%v",grouped4))

	grouped5 := CountBy( data3, func (obj T,key T,val T) T { 
		return obj.(map[T]T)["a"] 
	})
	asserts.Equals( t, "count by an object key\\'s value", "map[1:2 4:1]",fmt.Sprintf("%v",grouped5))
}

func TestKeys( t *testing.T ) {
	data := Keys( map[T]T{ "a":1,"b":4,"c":6 } )
	asserts.Equals( t, "keys of a map", "[a b c]",fmt.Sprintf("%v",data))

	nodata := Keys( map[T]T{ } )
	asserts.Equals( t, "keys of an empty map", "[]",fmt.Sprintf("%v",nodata))

	var nildata map[T]T = nil
	nilretval := Keys(nildata)
	asserts.Equals( t, "keys of an empty map", "[]",fmt.Sprintf("%v",nilretval))
}

func TestValues( t *testing.T ) {
	data := Values( map[T]T{ "a":1,"b":4,"c":6 } )
	asserts.Equals( t, "values of a map", "[1 4 6]",fmt.Sprintf("%v",data))

	data2 := Values( map[T]T{ "a":1,"b":1,"c":6 } )
	asserts.Equals( t, "values of a map", "[1 1 6]",fmt.Sprintf("%v",data2))

	nodata := Values( map[T]T{ } )
	asserts.Equals( t, "values of an empty map", "[]",fmt.Sprintf("%v",nodata))

	var nildata map[T]T = nil
	nilretval := Values(nildata)
	asserts.Equals( t, "values of an empty map", "[]",fmt.Sprintf("%v",nilretval))
}

func TestSize( t *testing.T ) {
	data := Size( map[T]T{ "a":1,"b":4,"c":6 } )
	asserts.IntEquals( t, "size of a map", 3,data)

	data1 := Size( []map[T]T{{ "a":1,"b":4,"c":6 }} )
	asserts.IntEquals( t, "size of a map", 1,data1)

	data10 := Size( []map[T]T{} )
	asserts.IntEquals( t, "size of a map", 0,data10)

	data2 := Size( []T{ "a","b","c",1 } )
	asserts.IntEquals( t, "size of a varied list ", 4,data2)

	data20 := Size( []T{ } )
	asserts.IntEquals( t, "size of an empty list ", 0,data20)
}

func TestSortedIndex(t *testing.T) {

	intLessThan := func( a T, b T) bool { return a.(int) < b.(int) }
	numbers := []T{10, 20, 30, 40, 50}
	num := 35
	indexForNum := SortedIndex(numbers, num, intLessThan )
	asserts.IntEquals( t, "35 should be inserted at index 3", 3, indexForNum)

	indexFor30 := SortedIndex(numbers, 30, intLessThan )
	asserts.IntEquals( t, "30 should be inserted at index 2", 2, indexFor30)

	objects := []map[T]T{{"x": 10}, {"x": 20}, {"x": 30}, {"x": 40}}
	iterator := func(obj T,idx T,list T) T { return obj.(map[T]T)["x"] }
	asserts.IntEquals(t,"sorted index with object list", 2, 
		SortedIndex(objects, map[T]T{"x": 25}, intLessThan,iterator))
	asserts.IntEquals(t,"sorted index with object list, take 2", 3, 
		SortedIndex(objects, map[T]T{"x": 35}, intLessThan,iterator))
}


func TestToArray( t *testing.T ){
    a := []T{1,2,3}
    asserts.Equals( t, "Clone an array",fmt.Sprintf("%v",a), fmt.Sprintf("%v",ToArray(a)))
    b := map[T]T{"one" : 1, "two" : 2, "three" : 3}
    numbers := ToArray(b)
    asserts.Equals( t, "object flattened into array", "[1 2 3]", fmt.Sprintf("%v", numbers) )
}


func TestFirst( t *testing.T ) {
		asserts.IntEquals( t, "can pull out the first element of an array", First([]T{1,2,3}).(int), 1 )
		asserts.Equals( t, "can pull out the first element of an array", 
			fmt.Sprintf("%v",FirstN([]T{1,2,3},2)), "[1 2]")

		asserts.Equals( t, "can pull out the zeroth element of an array", 
			fmt.Sprintf("%v",FirstN([]T{1,2,3},0)), "[1]")

		asserts.Equals( t, "can pull out too many items out", 
			fmt.Sprintf("%v",FirstN([]T{1,2,3},5)), "[1 2 3]")

		asserts.Equals( t, "works well with map", 
			fmt.Sprintf("%v",Map( []T{ []T{1,2,3}, []T{1,2,3}}, func(obj T,idx T,list T)T{return First(obj.([]T))})), "[1 1]")
		asserts.Equals( t, "works well with nil", 
			fmt.Sprintf("%v",First(nil)), "<nil>")
}


func TestInitial( t *testing.T ) {
		asserts.Equals( t, "can pull out the initial elements of an array", 
			fmt.Sprintf("%v",Initial([]T{1,2,3})), 
			"[1 2]" )
		asserts.Equals( t, "can pull out the first element of an array", 
			fmt.Sprintf("%v",Initial([]T{1,2,3},2)), "[1 2]")

		asserts.Equals( t, "can pull out the zeroth element of an array", 
			fmt.Sprintf("%v",Initial([]T{1,2,3},0)), "[]")

		asserts.Equals( t, "can pull out the first element of an array", 
			fmt.Sprintf("%v",Initial([]T{1,2,3},1)), "[1]")

		asserts.Equals( t, "can pull out all elements of an array", 
			fmt.Sprintf("%v",Initial([]T{1,2,3},3)), "[1 2 3]")

		asserts.Equals( t, "can pull out too many out", 
			fmt.Sprintf("%v",Initial([]T{1,2,3},4)), "[1 2 3]")

		asserts.Equals( t, "can pull out too many items out", 
			fmt.Sprintf("%v",Initial([]T{1,2,3},5)), "[1 2 3]")

		asserts.Equals( t, "works well with map", 
			fmt.Sprintf("%v",Map( []T{ []T{1,2,3}, []T{1,2,3}}, func(obj T,idx T,list T)T{return Initial(obj.([]T),2)})), "[[1 2] [1 2]]")
		asserts.Equals( t, "works well with nil", 
			fmt.Sprintf("%v",Initial(nil)), "[]")
}


func TestLast( t *testing.T ) {
		asserts.Equals( t, "can pull out the initial elements of an array", 
			fmt.Sprintf("%v",Last([]T{1,2,3})), 
			"[3]" )
		asserts.Equals( t, "can pull out the last 2 element of an array", 
			fmt.Sprintf("%v",Last([]T{1,2,3},2)), "[2 3]")

		asserts.Equals( t, "can pull out the zeroth element of an array", 
			fmt.Sprintf("%v",Last([]T{1,2,3},0)), "[]")

		asserts.Equals( t, "can pull out the first element of an array", 
			fmt.Sprintf("%v",Last([]T{1,2,3},1)), "[3]")

		asserts.Equals( t, "can pull out all elements of an array", 
			fmt.Sprintf("%v",Last([]T{1,2,3},3)), "[1 2 3]")

		asserts.Equals( t, "can pull out too many out", 
			fmt.Sprintf("%v",Last([]T{1,2,3},4)), "[1 2 3]")

		asserts.Equals( t, "can pull out too many items out", 
			fmt.Sprintf("%v",Last([]T{1,2,3},5)), "[1 2 3]")

		asserts.Equals( t, "works well with map", 
			fmt.Sprintf("%v",Map( []T{ []T{1,2,3}, []T{1,2,3}}, func(obj T,idx T,list T)T{return Last(obj.([]T),2)})), "[[2 3] [2 3]]")

		asserts.Equals( t, "works well with nil", 
			fmt.Sprintf("%v",Initial(nil)), "[]")
}

func TestRest( t *testing.T ) {
	asserts.Equals( t, "can pull out rest of items", 
		fmt.Sprintf("%v",Rest([]T{1,2,3})), "[2 3]")

	asserts.Equals( t, "works well with map", 
		fmt.Sprintf("%v",Map( []T{ []T{1,2,3}, []T{1,2,3}}, func(obj T,idx T,list T)T{return Rest(obj.([]T))})), "[[2 3] [2 3]]")

	asserts.Equals( t, "works well with nil", 
		fmt.Sprintf("%v",Rest(nil)), "[]")
}


func TestCompact( t *testing.T ) {
	asserts.Equals( t, "can trim out all falsy values", 
		fmt.Sprintf("%v",Compact([]T{0,1,false,2,false,3})), "[1 2 3]")
	asserts.Equals( t, "can trim out all falsy values", 
		fmt.Sprintf("%v",Compact(nil)), "[]")
}


func TestFlatten( t *testing.T ) {
   list := []T{ 1, []T{2}, []T{3, []T{[]T{[]T{4}}}}}
	asserts.Equals( t, "can flatten nested arrays", fmt.Sprintf("%v",Flatten(list)), "[1 2 3 4]")
   asserts.Equals( t, "can shallowly flatten nested arrays", fmt.Sprintf("%v",Flatten(list, true)),
		"[1 2 3 [[[4]]]]")
   list2 := []T{ []T{1}, []T{2}, []T{3}, []T{[]T{4}}}
   asserts.Equals( t, "can shallowly flatten arrays containing only other arrays",
		fmt.Sprintf("%v",Flatten(list2, true)), "[1 2 3 [4]]")
   asserts.Equals( t, "can flatten arrays containing only other arrays",
		fmt.Sprintf("%v",Flatten(list2)), "[1 2 3 4]")
}


func TestDifference( t *testing.T ) {
    result := Difference([]T{1, 2, 3}, []T{2, 30, 40})
    asserts.Equals( t, "takes the difference of two arrays", 
		fmt.Sprintf("%v",result), "[1 3]")

    result2 := Difference([]T{1, 2, 3, 4}, []T{2, 30, 40}, []T{1, 11, 111})
    asserts.Equals( t, "takes the difference of three arrays",
		fmt.Sprintf("%v",result2), "[3 4]")
}

func TestWithout( t *testing.T ) {

    list := []T{1, 2, 1, 0, 3, 1, 4}
    asserts.Equals( t, "can remove all instances of func args from an object",
		fmt.Sprintf("%v", Without(list, 0, 1)), "[2 3 4]")

    asserts.Equals( t, "can remove all instances from an array of args from an object",
		fmt.Sprintf("%v", Without(list, []T{0, 1})), "[2 3 4]")

    list2 := []T{1, 2, "1", 0, 3, 1, 4}
    asserts.Equals( t, "can remove all instances of func args from an object",
		fmt.Sprintf("%v", Without(list2, 0, 1)), "[2 1 3 4]")

    asserts.Equals( t, "can remove all instances from an array of args from an object",
		fmt.Sprintf("%v", Without(list2, []T{0, 1})), "[2 1 3 4]")

    asserts.Equals( t, "can remove all instances from an array of args from an object",
		fmt.Sprintf("%v", Without(list2, []T{0, "1"})), "[1 2 3 1 4]")

    asserts.Equals( t, "can remove all instances from an array of args from an object",
		fmt.Sprintf("%v", Without(list2, []T{0, "1", 1})), "[2 3 4]")

	/* TODO: fix ArrayOfMaps
   listm := []T{ {"one" : 1}, {"two" : 2}}
   asserts.Equals( t, "uses real object identity for comparisons.",
		fmt.Sprintf("%v",Without(listm, map[T]T{"one" : 1})), "barg") //).length == 2, 
   asserts.True("ditto", len(Without(listm, listm[0])) == 1)
	*/
}


func TestUniq( t *testing.T ) {
   list := []T{1, 2, 1, 3, 1, 4}
   asserts.Equals( t, "can find the unique values of an unsorted array",
		 fmt.Sprintf("%v", Uniq(list,false)), "[1 2 3 4]") 

   list2 := []T{1, 1, 1, 2, 2, 3}
   asserts.Equals( t, "can find the unique values of a sorted array faster",
		 fmt.Sprintf("%v", Uniq(list2,true)), "[1 2 3]") 

	list3 := []map[T]T{{ "name":"moe"}, { "name":"curly"}, { "name":"larry"}, {"name":"curly"}}
   iterator   := func(value T, key T, list T) T { return value.(map[T]T)["name"] }
   comparator := func(a T, b T) bool { return a.(map[T]T)["name"] == b.(map[T]T)["name"] }
   asserts.Equals( t, "can find the unique values of an array using a custom iterator",
		fmt.Sprintf("%v",Uniq(list3, false, iterator,comparator)), "[map[name:moe] map[name:curly] map[name:larry]]")
   asserts.Equals( t, "can find the unique values of an array using a custom iterator",
		fmt.Sprintf("%v",Map(Uniq(list3, false, iterator,comparator),iterator)), "[moe curly larry]")

}


func TestUnion( t *testing.T ) {
   result := Union([]T{1, 2, 3}, []T{2, 30, 1}, []T{1, 40})
   asserts.Equals( t, "takes the union of a list of arrays", fmt.Sprintf("%v",result), "[1 2 3 30 40]")

   result2 := Union([]T{1, 2, 3}, []T{2, 30, 1}, []T{1, 40, []T{1}})
   asserts.Equals( t, "takes the union of a list of nested arrays", fmt.Sprintf("%v",result2),"[1 2 3 30 40 [1]]")

   result3 := Union(nil, []T{1, 2, 3})
   asserts.Equals( t, "takes the union of nil and a list of arrays", fmt.Sprintf("%v",result3),"[<nil> 1 2 3]")
}

func TestIntersection( t *testing.T ) {
	strLessThan := func( this T, that T ) bool {
		return this.(string) < that.(string)
	}
	stooges := []T{"moe", "curly", "larry"}
	leaders := []T{"moe", "groucho"}
	asserts.Equals( t, "can take the set intersection of two arrays",
		fmt.Sprintf("%v",Intersection(strLessThan,stooges, leaders)), "[moe]")

	theSixStooges := []T{"moe", "moe", "curly", "curly", "larry", "larry"}
	asserts.Equals( t, "returns a duplicate-free array",
		fmt.Sprintf("%v",Intersection(strLessThan, theSixStooges, leaders)), "[moe]")
}

func TestIndexOf( t *testing.T ) {
	intLessThan := func( this T, that T ) bool {
		return this.(int) < that.(int)
	}
	numbers := []T{1, 2, 3}
	asserts.IntEquals( t, "can compute indexOf, even without the native function",
	IndexOf(numbers, 2,intLessThan), 1)

	asserts.IntEquals( t, "handles nulls properly",
	IndexOf(nil, 2,intLessThan), -1)

	num := 35;
	numbers2 := []T{10, 20, 30, 40, 50}
	asserts.IntEquals( t, "35 is not in the list", IndexOf(numbers2, num, intLessThan, true), -1 )
	num  = 40
	asserts.IntEquals( t, "40 is in the list", IndexOf(numbers2, num, intLessThan,true), 3)

	numbers3 := []T{1, 40, 40, 40, 40, 40, 40, 40, 50, 60, 70}
	num = 40
	asserts.IntEquals( t, "40 is in the list ", IndexOf(numbers3,num,intLessThan,true), 1)

}


func TestZip( t *testing.T ) {

   names := []T{"moe", "larry", "curly"}
	ages := []T{30, 40, 50}
	leaders := []T{true}

   stooges := Zip(names, ages, leaders)
   asserts.Equals( t, "zipped together arrays of different lengths",
		fmt.Sprintf("%v", stooges), "[[moe 30 true] [larry 40 <nil>] [curly 50 <nil>]]")

   stooges2 := Zip([]T{"moe",30, "stooge 1"},[]T{"larry",40, "stooge 2"},[]T{"curly",50, "stooge 3"})
	asserts.Equals( t, "zipped pairs",
    fmt.Sprintf("%v",stooges2), "[[moe larry curly] [30 40 50] [stooge 1 stooge 2 stooge 3]]")

    // In the case of difference lengths of the tuples undefineds
    // should be used as placeholder
    stooges3 := Zip([]T{"moe",30},[]T{"larry",40},[]T{"curly",50, "extra data"})
    asserts.Equals( t, "zipped pairs with empties", 
		fmt.Sprintf("%v",stooges3), "[[moe larry curly] [30 40 50] [<nil> <nil> extra data]]")
	empty := Zip([]T{})
	asserts.Equals( t, "unzipped empty", fmt.Sprintf("%v",empty),"[]")

	empty2 := Zip([]T{})
	asserts.Equals( t, "unzipped empty2", fmt.Sprintf("%v",empty2),"[]")
}


func TestObject( t *testing.T ) {
   result := Object([]T{"moe", "larry", "curly"}, []T{30, 40, 50})
   shouldBe := map[T]T{ "moe": 30, "larry": 40, "curly": 50}
   asserts.Equals( t, "two arrays zipped together into an object",
		fmt.Sprintf("%v", result), fmt.Sprintf("%v",shouldBe) )

	result2   := Object([]T{"one", 1, "two", 2, "three", 3})
	shouldBe2 := map[T]T{ "one": 1, "two": 2, "three": 3}
	asserts.Equals( t, "an array of pairs zipped together into an object",
		fmt.Sprintf("%v", result2), fmt.Sprintf("%v",shouldBe2) )

	result3   := Object([]T{[]T{"one", 1}, []T{"two", 2}, []T{"three", 3}})
	shouldBe3 := map[T]T{ "one": 1, "two": 2, "three": 3}
	asserts.Equals( t, "an array of pairs zipped together into an object",
		fmt.Sprintf("%v", result3), fmt.Sprintf("%v",shouldBe3) )

    asserts.Nil( t, "handles nils", Object(nil))

   stooges := map[T]T{ "moe": 30, "larry": 40, "curly": 50}
	asserts.Equals( t, "an object converted to pairs and back to an object",
		fmt.Sprintf("%v",stooges), fmt.Sprintf("%v",Object(Pairs(stooges))))
}


func TestPairs( t *testing.T ){
	asserts.Equals( t,  "can convert an object into pairs",
		fmt.Sprintf("%v",Pairs( map[T]T{"one": 1, "two": 2})), 
		fmt.Sprintf("%v",[]T{[]T{"one", 1}, []T{"two", 2}}))

	asserts.Equals( t, "... even when one of them is 'length'",
		fmt.Sprintf("%v",Pairs(map[T]T{"one": 1, "two": 2, "length": 3})), 
		fmt.Sprintf("%v",[]T{ []T{"one", 1}, []T{"two", 2}, []T{"length", 3}}))

}


func TestLastIndexOf( t *testing.T ) {
	numbers := []T{1, 0, 1}
	asserts.IntEquals( t, "lastindexof simple test", LastIndexOf(numbers, 1), 2)

	numbers2 := []T{1, 0, 1, 0, 0, 1, 0, 0, 0}
	asserts.IntEquals( t, "can compute lastIndexOf", LastIndexOf(numbers2, 1), 5)
	asserts.IntEquals( t, "lastIndexOf the other element", LastIndexOf(numbers2, 0), 8)

	asserts.IntEquals( t, "handles nulls properly", LastIndexOf(nil, 2), -1) 

	numbers3 := []T{1, 2, 3, 1, 2, 3, 1, 2, 3}
	index := LastIndexOf(numbers3, 2, 2)
	asserts.IntEquals( t, "supports the fromIndex argument",index, 1)

}

func TestRange( t *testing.T ) {
	asserts.Equals( t, "range with no arguments generates an empty array",
		fmt.Sprintf("%v",Range()),"[]")
	asserts.Equals( t, "range with 0 as a first argument generates an empty array",
		fmt.Sprintf("%v",Range(0)), "[]")
	asserts.Equals( t, "range with a single positive argument generates an array of elements 0,1,2,...,n-1",
		fmt.Sprintf("%v",Range(4)), "[0 1 2 3]") 
	asserts.Equals( t, "range with two arguments a &amp; b, a&lt;b generates an array of elements a,a+1,a+2,...,b-2,b-1",
		fmt.Sprintf("%v",Range(5, 8)), "[5 6 7]") 
	asserts.Equals( t, "range with two arguments a &amp; b, b&lt;a generates an empty array",
		fmt.Sprintf("%v",Range(8, 5)), "[]") 
	asserts.Equals( t, "range with three arguments a &amp; b &amp; c, c &lt; b-a, a &lt; b generates an array of elements a,a+c,a+2c,...,b - (multiplier of a) &lt; c",
		fmt.Sprintf("%v",Range(3, 10, 3)), "[3 6 9]") 
	asserts.Equals( t, "range with three arguments a &amp; b &amp; c, c &gt; b-a, a &lt; b generates an array with a single element, equal to a",
		fmt.Sprintf("%v",Range(3, 10, 15)), "[3]") 
	asserts.Equals( t,  "range with three arguments a &amp; b &amp; c, a &gt; b, c &lt; 0 generates an array of elements a,a-c,a-2c and ends with the number not less than b",
		fmt.Sprintf("%v",Range(12, 7, -2)), "[12 10 8]")
	asserts.Equals( t,  "final example in the Python docs",
		fmt.Sprintf("%v",Range(0, -10, -1)), "[0 -1 -2 -3 -4 -5 -6 -7 -8 -9]")

}


func TestPartial( t *testing.T ) {

	funk := func( args ...T) T { 
		return fmt.Sprintf("%v",args )
	}

	passAB := Partial(funk, "a", "b");

	asserts.Equals( t, "can partially apply", 
		fmt.Sprintf("%v", passAB("c", "d")), "[a b c d]")

	asserts.Equals( t, "can partially apply", 
		fmt.Sprintf("%v", passAB("e", "f")), "[a b e f]")

	asserts.Equals( t, "can partially apply", 
		fmt.Sprintf("%v", passAB("1", 2)), "[a b 1 2]")
}

func TestMemoize( t *testing.T ) {
	asserts.IntEquals( t, "a memoized version of fibonacci produces identical results",
		fib(10).(int), 55)
	fib2 := Memoize(fib, func(n ...T) T { return n[0]}) // Redefine `fib` for memoization
	asserts.IntEquals( t, "a memoized version of fibonacci produces identical results",
		fib2(10).(int), 55)

	o := func(str ...T) T { return str[0] }
	fastO := Memoize(o)
	asserts.Equals( t, "blah", o("toString").(string), "toString")
	asserts.Equals( t, "blah blah", fastO("toString").(string), "toString")
}


func TestOnce( t *testing.T ) {
	num := 0
	increment := Once(func(...T) T { 
		num += 1 
		return num 
	})
	increment()
	increment()
	asserts.IntEquals( t, "can increment once", num, 1)
}

func TestWrap( t *testing.T ) {
	// from http://play.golang.org/p/Ic5G5QEO93
	reverse := func(s string) string {
		n := len(s)
		runes := make([]rune, n)
		for _, rune := range s {
			n--
			runes[n] = rune
		}
		return string(runes[n:])
	}
	greet := func(name ...T) T { return fmt.Sprintf("hi: %v",name[0]) }
	backwards := Wrap(greet, func(args ...T)T {
		fn := args[0].(func(...T)T)
		name := args[1]
		return fmt.Sprintf("%v %v", fn(name) , 
			reverse(name.(string)))
	})
	asserts.Equals( t,"wrapped the salutation function", backwards("moe").(string), "hi: moe eom" )

	inner := func(...T) T{ return "Hello " }
	obj   := map[T]T{ "name" : "Moe"}
	hi    := Wrap(inner, func(args ...T)T {
		fn := args[0].(func(...T)T)
		this := args[1].(map[T]T)
		return fn().(string) + this["name"].(string)
	})
	asserts.Equals( t, "pass obj as arg ", hi(obj).(string), "Hello Moe")

	noop  := func(...T) T {return nil }
	wrapped := Wrap(noop, func( args ...T)T{ 
		return args
	})
	ret    := wrapped([]T{"whats", "your"}, "vector", "victor")
	_,ok := ret.([]T)[0].(func(...T)T)
	asserts.True( t, "noop test first arg is a func", ok)
	asserts.Equals( t, "noop test, rest of args", fmt.Sprintf("%v",ret.([]T)[1:]),
		"[[whats your] vector victor]")
}


func TestCompose( t *testing.T ) {
	greet    := func(name ...T)T{ return "hi: " + name[0].(string) }
	exclaim  := func(sentence ...T)T{ return sentence[0].(string) + "!" }
	pause    := func(midway ...T)T{ return midway[0].(string) + ", " }
	composed := Compose(exclaim, greet)

	asserts.Equals( t, "can compose a function that takes another",composed("moe").(string), "hi: moe!") 

	composed2 := Compose(greet, exclaim)
	asserts.Equals( t, "in this case, the functions are also commutative",composed2("moe").(string), "hi: moe!") 

	composed3 := Compose(greet, pause, exclaim)
	asserts.Equals( t, "in this case, the functions are not commutative",composed3("moe").(string), "hi: moe!, ") 

	composed4 := Compose( greet, exclaim, pause)
	asserts.Equals( t, "in this case, the functions are not commutative",composed4("moe").(string), "hi: moe, !") 
} 


func TestAfter( t *testing.T ) {
	testAfter := func(afterAmount , timesCalled int) int {
		afterCalled := 0
		after := After(afterAmount, func(...T)T {
			afterCalled += 1
			return nil
		})
		for timesCalled > 0 {
			timesCalled -= 1
			after()
		}
		return afterCalled
	}

	asserts.IntEquals( t, "after(N) should fire after being called N times", 
		testAfter(5, 5), 1)
	asserts.IntEquals( t, "after(N) should not fire unless called N times",  
		testAfter(5, 4), 0) 
	asserts.IntEquals( t, "after(0) should not fire immediately", 
		testAfter(0, 0), 0) 
	asserts.IntEquals( t, "after(0) should fire when first invoked", 
		testAfter(0, 1), 1)
}
