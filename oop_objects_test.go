package underscore

import (
	"github.com/markmontymark/asserts"
	"fmt"
	"math"
	"testing"
	"time"
)

func TestKeysOOP(t *testing.T) {
	data := New(map[T]T{"a": 1, "b": 4, "c": 6}).Chain().Keys().Value()
	asserts.Equals(t, "keys of a map", "[a b c]", fmt.Sprint(data))

	nodata := New(map[T]T{}).Chain().Keys().Value()
	asserts.Equals(t, "keys of an empty map", "[]", fmt.Sprint(nodata))

	var nildata map[T]T = nil
	nilretval := New(nildata).Chain().Keys().Value()
	asserts.Equals(t, "keys of an empty map", "[]", fmt.Sprint(nilretval))
}

func TestKeysOOPChain(t *testing.T) {
	data := New(map[T]T{"a": 1, "b": 4, "c": 6}).Chain().Keys().Value()
	asserts.Equals(t, "keys of a map", "[a b c]", fmt.Sprint(data))

	nodata := New(map[T]T{}).Chain().Keys().Value()
	asserts.Equals(t, "keys of an empty map", "[]", fmt.Sprint(nodata))

	var nildata map[T]T = nil
	nilretval := New(nildata).Chain().Keys().Value()
	asserts.Equals(t, "keys of an empty map", "[]", fmt.Sprint(nilretval))
}

func TestValuesOOP(t *testing.T) {
	data := New(map[T]T{"a": 1, "b": 4, "c": 6}).Chain().Values().Value()
	asserts.Equals(t, "values of a map", "[1 4 6]", fmt.Sprint(data))

	data2 := New(map[T]T{"a": 1, "b": 1, "c": 6}).Chain().Values().Value()
	asserts.Equals(t, "values of a map", "[1 1 6]", fmt.Sprint(data2))

	nodata := New(map[T]T{}).Chain().Values().Value()
	asserts.Equals(t, "values of an empty map", "[]", fmt.Sprint(nodata))

	var nildata map[T]T = nil
	nilretval := New(nildata).Chain().Values().Value()
	asserts.Equals(t, "values of an empty map", "[]", fmt.Sprint(nilretval))
}

func TestPairsOOP(t *testing.T) {
	asserts.Equals(t, "can convert an object into pairs",
		fmt.Sprint(New(map[T]T{"one": 1, "two": 2}).Chain().Pairs().Value()),
		fmt.Sprint([]T{[]T{"one", 1}, []T{"two", 2}}))

	asserts.Equals(t, "... even when one of them is length",
		fmt.Sprint(New(map[T]T{"one": 1, "two": 2, "length": 3}).Chain().Pairs().Value()),
		fmt.Sprint([]T{[]T{"one", 1}, []T{"two", 2}, []T{"length", 3}}))

}

func TestInvertOOP(t *testing.T) {
	obj := map[T]T{"first": "Moe", "second": "Larry", "third": "Curly"}
	asserts.Equals(t, "can invert an object", fmt.Sprint(New(obj).Chain().Invert().Keys().Value()), "[Moe Larry Curly]")
	asserts.Equals(t, "two inverts gets you back where you started",
		fmt.Sprint(New(obj).Chain().Invert().Invert().Value()), fmt.Sprint(obj))
}

func TestExtendOOP(t *testing.T) {

	asserts.Equals(t, "can extend an object with the attributes of another",
		New(map[T]T{}).Chain().Extend(map[T]T{"a": "b"}).Value().(map[T]T)["a"].(string), "b")

	asserts.Equals(t, "properties in source override destination",
		New(map[T]T{"a": "x"}).Chain().Extend(map[T]T{"a": "b"}).Value().(map[T]T)["a"].(string), "b")

	asserts.Equals(t, "properties not in source dont get overriden",
		New(map[T]T{"x": "x"}).Chain().Extend(map[T]T{"a": "b"}).Value().(map[T]T)["x"].(string), "x")

	result := New(map[T]T{"x": "x"}).Chain().Extend(map[T]T{"a": "a"}, map[T]T{"b": "b"}).Value()
	asserts.Equals(t, "can extend from multiple source objects",
		fmt.Sprint(result), fmt.Sprint(map[T]T{"x": "x", "a": "a", "b": "b"}))

	result2 := New(map[T]T{"x": "x"}).Chain().Extend(map[T]T{"a": "a", "x": 2}, map[T]T{"a": "b"}).Value()
	asserts.Equals(t, "extending from multiple source objects last property trumps",
		fmt.Sprint(result2), fmt.Sprint(map[T]T{"x": 2, "a": "b"}))

	result3 := New(map[T]T{}).Chain().Extend(map[T]T{"a": 0, "b": nil}).Keys().Value()
	asserts.Equals(t, "extend copies undefined values",
		fmt.Sprint(result3), "[a b]")

	result4 := New(map[T]T{})
	result5 := result4.Extend(nil, 0, map[T]T{"a": 1}) //.Value()
	asserts.IntEquals(t, "should not error on `null` or `undefined` sources",
		result5.Value().(map[T]T)["a"].(int), 1)

}

func TestPickOOP(t *testing.T) {
	result := New(map[T]T{"a": 1, "b": 2, "c": 3}).Chain().Pick("a", "c").Value()
	asserts.Equals(t, "can restrict properties to those named",
		fmt.Sprint(result), fmt.Sprint(map[T]T{"a": 1, "c": 3}))

	result2 := New(map[T]T{"a": 1, "b": 2, "c": 3}).Chain().Pick([]T{"b", "c"}).Value()
	asserts.Equals(t, "can restrict properties to those named in an array",
		fmt.Sprint(result2), fmt.Sprint(map[T]T{"b": 2, "c": 3}))

	result3 := New(map[T]T{"a": 1, "b": 2, "c": 3}).Chain().Pick([]T{"a"}, "b").Value()
	asserts.Equals(t, "can restrict properties to those named in mixed args",
		fmt.Sprint(result3), fmt.Sprint(map[T]T{"a": 1, "b": 2}))
}

func TestOmitOOP(t *testing.T) {
	result := New(map[T]T{"a": 1, "b": 2, "c": 3}).Chain().Omit("b").Value()
	asserts.Equals(t, "can omit a single named property",
		fmt.Sprint(result), fmt.Sprint(map[T]T{"a": 1, "c": 3}))
	result2 := New(map[T]T{"a": 1, "b": 2, "c": 3}).Chain().Omit("a", "c").Value()
	asserts.Equals(t, "can omit several named properties",
		fmt.Sprint(result2), fmt.Sprint(map[T]T{"b": 2}))
	result3 := New(map[T]T{"a": 1, "b": 2, "c": 3}).Chain().Omit([]T{"b", "c"}).Value()
	asserts.Equals(t, "can omit properties named in an array",
		fmt.Sprint(result3), fmt.Sprint(map[T]T{"a": 1}))
}

func TestDefaultsOOP(t *testing.T) {

	options := New(map[T]T{"zero": 0, "one": 1, "empty": "", "nan": math.NaN(), "nothing": nil})
	options2 := options.Defaults(map[T]T{"zero": 1, "one": 10, "twenty": 20, "nothing": "str"}).Value().(map[T]T)

	asserts.IntEquals(t, "value exists", options2["zero"].(int), 0)
	asserts.IntEquals(t, "value exists", options2["one"].(int), 1)
	asserts.IntEquals(t, "default applied", options2["twenty"].(int), 20)
	asserts.Nil(t, "null isnt overridden", options2["nothing"])

	options3 := options.Defaults(map[T]T{"empty": "full"}, map[T]T{"nan": "nan"}, map[T]T{"word": "word"}, map[T]T{"word": "dog"}).Value().(map[T]T)
	asserts.Equals(t, "value exists", options3["empty"].(string), "")
	asserts.True(t, "NaN isnt overridden", math.IsNaN(options3["nan"].(float64)))
	asserts.Equals(t, "new value is added, first one wins", options3["word"].(string), "word")

	options4 := New(map[T]T{}).Defaults(nil, map[T]T{"a": 1}).Value().(map[T]T)

	asserts.IntEquals(t, "should not error on `null` or `undefined` sources", options4["a"].(int), 1)
}

func TestCloneOOP(t *testing.T) {
	moe := map[T]T{"name": "moe", "lucky": []T{13, 27, 34}}
	var clone map[T]T
	clone = New(moe).Chain().Clone().Value().(map[T]T)
	asserts.Equals(t, "the clone as the attributes of the original",
		clone["name"].(string), "moe")

	clone["name"] = "curly"
	asserts.True(t, "clones can change shallow attributes without affecting the original",
		clone["name"].(string) == "curly" && moe["name"].(string) == "moe")

	clone["lucky"] = append(clone["lucky"].([]T), 101)
	// CHANGE FROM ORIG Underscore.js, changes to deep attributes are not shared!!!
	// TODO:  Fix???
	asserts.IntEquals(t, "changes to deep attributes are NOT shared with the original",
		New(moe["lucky"]).Chain().Last().Value().([]T)[0].(int), 34) // was 101 in Underscore.js
	asserts.IntEquals(t, "non objects should not be changed by clone", New(1).Clone().Value().(int), 1)
	asserts.Nil(t, "non objects should not be changed by clone", New(nil).Clone().Value())

	var cloneArray []T
	cloneArray = New(moe["lucky"]).Clone().Value().([]T)
	asserts.Equals(t, "clone an array", fmt.Sprint(cloneArray), fmt.Sprint(moe["lucky"]))
	cloneArray = append(cloneArray, 10101)
	asserts.NotEquals(t, "clone an array is shallow?", fmt.Sprint(cloneArray), fmt.Sprint(moe["lucky"]))
}

// TODO: missing TestIsEqual from objects.js
// TODO: missing TestIsEmpty from objects.js
// XXX: missing TestIsElement from objects.js - wont add
// XXX: missing TestIsArguments from objects.js - wont add
// TODO: missing TestIsObject from objects.js

func TestIsArrayWithArrayOOP(t *testing.T) {
	list := New([]T{"name", "moe", "age", 30}).Chain()
	asserts.True(t, "Testing IsArray", list.IsArray().Value().(bool))
}

func TestIsArrayWithStringOOP(t *testing.T) {
	scalar := New("name").Chain()
	asserts.False(t, "Testing IsArray", scalar.IsArray().Value().(bool))
}

func TestIsArrayWithMapOOP(t *testing.T) {
	mapp := New(make(map[string]int, 0)).Chain()
	asserts.False(t, "Testing IsArray", mapp.IsArray().Value().(bool))
}

func TestIsStringWithArrayOOP(t *testing.T) {
	list := New([]T{"name", "moe", "age", 30}).Chain()
	asserts.False(t, "Testing IsString", list.IsString().Value().(bool))
}

func TestIsStringWithStringOOP(t *testing.T) {
	scalar := New("name").Chain()
	asserts.True(t, "Testing IsString function", scalar.IsString().Value().(bool))
	asserts.False(t, "Testing New([]T{scalar}).IsString method()", New([]T{"name"}).Chain().IsString().Value().(bool))
}

func TestIsStringWithMapOOP(t *testing.T) {
	mapp := New(make(map[string]int, 0)).Chain()
	asserts.False(t, "Testing IsString", mapp.IsString().Value().(bool))
}

func TestIsMapWithArrayOOP(t *testing.T) {
	list := New([]T{"name", "moe", "age", 30}).Chain()
	asserts.False(t, "Testing IsMap with array ", list.IsMap().Value().(bool))
}

func TestIsMapWithStringOOP(t *testing.T) {
	scalar := New("name").Chain()
	asserts.False(t, "Testing IsMap with string", scalar.IsMap().Value().(bool))
}

func TestIsMapWithMapOOP(t *testing.T) {
	mapp := New(make(map[T]T, 0)).Chain()
	asserts.True(t, "Testing IsMap", mapp.IsMap().Value().(bool))
}

//XXX: missing TestIsNumber from objects.js
//XXX: missing TestIsBool from objects.js
//XXX: missing TestIsFunction from objects.js
//XXX: missing TestIsDate from objects.js
//XXX: missing TestIsRegExp from objects.js
//XXX: missing TestIsFinite from objects.js
//XXX: missing TestIsNaN from objects.js
//XXX: missing TestIsNil from objects.js

func TestTapOOP(t *testing.T) {
	var intercepted T
	interceptor := func(args ...T) T {
		obj := args[0]
		intercepted = obj
		return nil
	}
	returned := Tap(1, interceptor)
	asserts.IntEquals(t, "passes tapped object to interceptor", intercepted.(int), 1)
	asserts.IntEquals(t, "returns tapped object", returned.(int), 1)

	returned2 := New([]T{1, 2, 3}).Chain().Map(func(val, idx, list T) T { return val.(int) * 2 }).Max(func(a T, b T) bool { return a.(int) < b.(int) }).Tap(interceptor) //.Value()

	asserts.True(t, "can use tapped objects in a chain",
		returned2.Value().(int) == 6 && intercepted.(int) == 6)

	returned3 := New([]T{1}).Chain().Map(func(val, idx, list T) T { return val.(int) * 2 }).Max(func(a T, b T) bool { return a.(int) < b.(int) }).Tap(interceptor) //.Value()
	asserts.True(t, "can use tapped scalar in a chain",
		returned3.Value().(int) == 2 && intercepted.(int) == 2)

	returned4 := New(map[T]T{"a": 1}).Chain().Map(func(val, idx, list T) T { return val.(int) * 2 }).Max(func(a T, b T) bool { return a.(int) < b.(int) }).Tap(interceptor) //.Value()
	asserts.True(t, "can use tapped scalar in a chain",
		returned4.Value().(int) == 2 && intercepted.(int) == 2)
}

func TestNowOOP(t *testing.T) {
	diff := New().Chain().Now().Value().(int64) - time.Now().UnixNano()
	asserts.True(t, "Now called in OOP style, in a chain",
		diff <= 0 && diff >= -500)
}

// TODO:  missing TestHas from objects.js "has"
