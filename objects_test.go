package underscore

import (
	"./lib/asserts"
	"fmt"
	"math"
	"testing"
)

func TestKeys(t *testing.T) {
	data := Keys(map[T]T{"a": 1, "b": 4, "c": 6})
	asserts.Equals(t, "keys of a map", "[a b c]", fmt.Sprintf("%v", data))

	nodata := Keys(map[T]T{})
	asserts.Equals(t, "keys of an empty map", "[]", fmt.Sprintf("%v", nodata))

	var nildata map[T]T = nil
	nilretval := Keys(nildata)
	asserts.Equals(t, "keys of an empty map", "[]", fmt.Sprintf("%v", nilretval))
}

func TestValues(t *testing.T) {
	data := Values(map[T]T{"a": 1, "b": 4, "c": 6})
	asserts.Equals(t, "values of a map", "[1 4 6]", fmt.Sprintf("%v", data))

	data2 := Values(map[T]T{"a": 1, "b": 1, "c": 6})
	asserts.Equals(t, "values of a map", "[1 1 6]", fmt.Sprintf("%v", data2))

	nodata := Values(map[T]T{})
	asserts.Equals(t, "values of an empty map", "[]", fmt.Sprintf("%v", nodata))

	var nildata map[T]T = nil
	nilretval := Values(nildata)
	asserts.Equals(t, "values of an empty map", "[]", fmt.Sprintf("%v", nilretval))
}

func TestPairs(t *testing.T) {
	asserts.Equals(t, "can convert an object into pairs",
		fmt.Sprintf("%v", Pairs(map[T]T{"one": 1, "two": 2})),
		fmt.Sprintf("%v", []T{[]T{"one", 1}, []T{"two", 2}}))

	asserts.Equals(t, "... even when one of them is length",
		fmt.Sprintf("%v", Pairs(map[T]T{"one": 1, "two": 2, "length": 3})),
		fmt.Sprintf("%v", []T{[]T{"one", 1}, []T{"two", 2}, []T{"length", 3}}))

}

func TestInvert(t *testing.T) {
	obj := map[T]T{"first": "Moe", "second": "Larry", "third": "Curly"}
	asserts.Equals(t, "can invert an object", fmt.Sprintf("%v", Keys(Invert(obj))), "[Moe Larry Curly]")
	asserts.Equals(t, "two inverts gets you back where you started",
		fmt.Sprintf("%v", Invert(Invert(obj))), fmt.Sprintf("%v", obj))
}

func TestExtend(t *testing.T) {

	asserts.Equals(t, "can extend an object with the attributes of another",
		Extend(map[T]T{}, map[T]T{"a": "b"})["a"].(string), "b")

	asserts.Equals(t, "properties in source override destination",
		Extend(map[T]T{"a": "x"}, map[T]T{"a": "b"})["a"].(string), "b")

	asserts.Equals(t, "properties not in source dont get overriden",
		Extend(map[T]T{"x": "x"}, map[T]T{"a": "b"})["x"].(string), "x")

	result := Extend(map[T]T{"x": "x"}, map[T]T{"a": "a"}, map[T]T{"b": "b"})
	asserts.Equals(t, "can extend from multiple source objects",
		fmt.Sprintf("%v", result), fmt.Sprintf("%v", map[T]T{"x": "x", "a": "a", "b": "b"}))

	result2 := Extend(map[T]T{"x": "x"}, map[T]T{"a": "a", "x": 2}, map[T]T{"a": "b"})
	asserts.Equals(t, "extending from multiple source objects last property trumps",
		fmt.Sprintf("%v", result2), fmt.Sprintf("%v", map[T]T{"x": 2, "a": "b"}))

	result3 := Extend(map[T]T{}, map[T]T{"a": 0, "b": nil})
	asserts.Equals(t, "extend copies undefined values",
		fmt.Sprintf("%v", Keys(result3)), "[a b]")

	result4 := map[T]T{}
	result5 := Extend(result4, nil, 0, map[T]T{"a": 1})
	asserts.IntEquals(t, "should not error on `null` or `undefined` sources",
		result5["a"].(int), 1)

}

func TestPick(t *testing.T) {
	result := Pick(map[T]T{"a": 1, "b": 2, "c": 3}, "a", "c")
	asserts.Equals(t, "can restrict properties to those named",
		fmt.Sprintf("%v", result), fmt.Sprintf("%v", map[T]T{"a": 1, "c": 3}))

	result2 := Pick(map[T]T{"a": 1, "b": 2, "c": 3}, []T{"b", "c"})
	asserts.Equals(t, "can restrict properties to those named in an array",
		fmt.Sprintf("%v", result2), fmt.Sprintf("%v", map[T]T{"b": 2, "c": 3}))

	result3 := Pick(map[T]T{"a": 1, "b": 2, "c": 3}, []T{"a"}, "b")
	asserts.Equals(t, "can restrict properties to those named in mixed args",
		fmt.Sprintf("%v", result3), fmt.Sprintf("%v", map[T]T{"a": 1, "b": 2}))
}

func TestOmit(t *testing.T) {
	result := Omit(map[T]T{"a": 1, "b": 2, "c": 3}, "b")
	asserts.Equals(t, "can omit a single named property",
		fmt.Sprintf("%v", result), fmt.Sprintf("%v", map[T]T{"a": 1, "c": 3}))
	result2 := Omit(map[T]T{"a": 1, "b": 2, "c": 3}, "a", "c")
	asserts.Equals(t, "can omit several named properties",
		fmt.Sprintf("%v", result2), fmt.Sprintf("%v", map[T]T{"b": 2}))
	result3 := Omit(map[T]T{"a": 1, "b": 2, "c": 3}, []T{"b", "c"})
	asserts.Equals(t, "can omit properties named in an array",
		fmt.Sprintf("%v", result3), fmt.Sprintf("%v", map[T]T{"a": 1}))
}

func TestDefaults(t *testing.T) {

	options := map[T]T{"zero": 0, "one": 1, "empty": "", "nan": math.NaN(), "nothing": nil}
	Defaults(options, map[T]T{"zero": 1, "one": 10, "twenty": 20, "nothing": "str"})

	asserts.IntEquals(t, "value exists", options["zero"].(int), 0)
	asserts.IntEquals(t, "value exists", options["one"].(int), 1)
	asserts.IntEquals(t, "default applied", options["twenty"].(int), 20)
	asserts.Nil(t, "null isnt overridden", options["nothing"])

	Defaults(options, map[T]T{"empty": "full"}, map[T]T{"nan": "nan"}, map[T]T{"word": "word"}, map[T]T{"word": "dog"})
	asserts.Equals(t, "value exists", options["empty"].(string), "")
	asserts.True(t, "NaN isnt overridden", math.IsNaN(options["nan"].(float64)))
	asserts.Equals(t, "new value is added, first one wins", options["word"].(string), "word")

	options2 := map[T]T{}
	Defaults(options2, nil, map[T]T{"a": 1})

	asserts.IntEquals(t, "should not error on `null` or `undefined` sources", options2["a"].(int), 1)
}

func TestClone(t *testing.T) {
	moe := map[T]T{"name": "moe", "lucky": []T{13, 27, 34}}
	var clone map[T]T
	clone = Clone(moe).(map[T]T)
	asserts.Equals(t, "the clone as the attributes of the original",
		clone["name"].(string), "moe")

	clone["name"] = "curly"
	asserts.True(t, "clones can change shallow attributes without affecting the original",
		clone["name"].(string) == "curly" && moe["name"].(string) == "moe")

	clone["lucky"] = append(clone["lucky"].([]T), 101)
	// CHANGE FROM ORIG Underscore.js, changes to deep attributes are not shared!!!
	// TODO:  Fix???
	asserts.IntEquals(t, "changes to deep attributes are NOT shared with the original",
		Last(moe["lucky"].([]T))[0].(int), 34) // was 101 in Underscore.js
	asserts.IntEquals(t, "non objects should not be changed by clone", Clone(1).(int), 1)
	asserts.Nil(t, "non objects should not be changed by clone", Clone(nil))

	var cloneArray []T
	cloneArray = Clone(moe["lucky"]).([]T)
	asserts.Equals(t, "clone an array", fmt.Sprintf("%v", cloneArray), fmt.Sprintf("%v", moe["lucky"]))
	cloneArray = append(cloneArray, 10101)
	asserts.NotEquals(t, "clone an array is shallow?", fmt.Sprintf("%v", cloneArray), fmt.Sprintf("%v", moe["lucky"]))
}

// TODO: missing TestIsEqual from objects.js
// TODO: missing TestIsEmpty from objects.js
// XXX: missing TestIsElement from objects.js - wont add
// XXX: missing TestIsArguments from objects.js - wont add
// TODO: missing TestIsObject from objects.js

func TestIsArrayWithArray(t *testing.T) {
	list := []T{"name", "moe", "age", 30}
	asserts.True(t, "Testing IsArray", IsArray(list))
}

func TestIsArrayWithString(t *testing.T) {
	scalar := "name"
	asserts.False(t, "Testing IsArray", IsArray(scalar))
}

func TestIsArrayWithMap(t *testing.T) {
	mapp := make(map[string]int, 0)
	asserts.False(t, "Testing IsArray", IsArray(mapp))
}

func TestIsStringWithArray(t *testing.T) {
	list := []T{"name", "moe", "age", 30}
	asserts.False(t, "Testing IsString", IsString(list))
}

func TestIsStringWithString(t *testing.T) {
	scalar := "name"
	asserts.True(t, "Testing IsString function", IsString(scalar))
	asserts.True(t, "Testing New(scalar).IsString method()", New(scalar).IsString().Value().(bool))
	asserts.False(t, "Testing New([]T{scalar}).IsString method()", New([]T{scalar}).IsString().Value().(bool))
}

func TestIsStringWithMap(t *testing.T) {
	mapp := make(map[string]int, 0)
	asserts.False(t, "Testing IsString", IsString(mapp))
}

func TestIsMapWithArray(t *testing.T) {
	list := []T{"name", "moe", "age", 30}
	asserts.False(t, "Testing IsMap with array ", IsMap(list))
}

func TestIsMapWithString(t *testing.T) {
	scalar := "name"
	asserts.False(t, "Testing IsMap with string", IsMap(scalar))
}

func TestIsMapWithMap(t *testing.T) {
	mapp := make(map[T]T, 0)
	asserts.True(t, "Testing IsMap", IsMap(mapp))
}

//XXX: missing TestIsNumber from objects.js
//XXX: missing TestIsBool from objects.js
//XXX: missing TestIsFunction from objects.js
//XXX: missing TestIsDate from objects.js
//XXX: missing TestIsRegExp from objects.js
//XXX: missing TestIsFinite from objects.js
//XXX: missing TestIsNaN from objects.js
//XXX: missing TestIsNil from objects.js

func TestTap(t *testing.T) {
	var intercepted T
	interceptor := func(args ...T) T {
		obj := args[0]
		intercepted = obj
		return nil
	}
	returned := Tap(1, interceptor)
	asserts.IntEquals(t, "passes tapped object to interceptor", intercepted.(int), 1)
	asserts.IntEquals(t, "returns tapped object", returned.(int), 1)

	returned2 := New([]T{1, 2, 3}).Chain().Map(func(val, idx, list T) T { return val.(int) * 2 }).Max(func(a T, b T) bool { return a.(int) < b.(int) }).Tap(interceptor).Value()

	asserts.True(t, "can use tapped objects in a chain",
		returned2.(int) == 6 && intercepted.(int) == 6)

	returned3 := New([]T{1}).Chain().Map(func(val, idx, list T) T { return val.(int) * 2 }).Max(func(a T, b T) bool { return a.(int) < b.(int) }).Tap(interceptor).Value()
	asserts.True(t, "can use tapped scalar in a chain",
		returned3.(int) == 2 && intercepted.(int) == 2)

	returned4 := New(map[T]T{"a": 1}).Chain().Map(func(val, idx, list T) T { return val.(int) * 2 }).Max(func(a T, b T) bool { return a.(int) < b.(int) }).Tap(interceptor).Value()
	asserts.True(t, "can use tapped scalar in a chain",
		returned4.(int) == 2 && intercepted.(int) == 2)
}

// TODO:  missing TestHas from objects.js "has"
