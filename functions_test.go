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

// XXX: missing bind - wont add

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

// XXX: missing bindAll - wont add

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

// XXX: missing delay - wont add
// XXX: missing defer - wont add
// XXX: missing throttle - wont add
// XXX: missing debounce - wont add

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

// TODO: missing recursive once

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
