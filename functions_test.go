package underscore

// TODO, for the new debounce, throttle, delay tests, I'm using a DelayAndWait func added to underscore, 
//  but perhaps it would be better to append all delay channels to one channel slice, and block in one place
// in this file, instead of blocking in the DelayAndWait call.  

import (
	"github.com/markmontymark/asserts"
	"fmt"
	"sync"
	"testing"
	"time"
	_ "os"
)

var fib func(...T) T

var globalWg sync.WaitGroup

func init() {
	fib = func(n ...T) T {
		if n[0].(int) < 2 {
			return n[0]
		}
		return fib(n[0].(int)-1).(int) + fib(n[0].(int)-2).(int)
	}
}

// XXX: missing bind

func TestPartial(t *testing.T) {

	funk := func(args ...T) T {
		return fmt.Sprint(args)
	}

	passAB := Partial(funk, "a", "b")

	asserts.Equals(t, "can partially apply",
		fmt.Sprint(passAB("c", "d")), "[a b c d]")

	asserts.Equals(t, "can partially apply",
		fmt.Sprint(passAB("e", "f")), "[a b e f]")

	asserts.Equals(t, "can partially apply",
		fmt.Sprint(passAB("1", 2)), "[a b 1 2]")
}

// XXX: missing bindAll

func TestMemoize(t *testing.T) {
	asserts.IntEquals(t, "a memoized version of fibonacci produces identical results",
		fib(10).(int), 55)
	fib2 := Memoize(fib, func(n ...T) T { return n[0] }) // Redefine `fib` for memoization
	asserts.IntEquals(t, "a memoized version of fibonacci produces identical results",
		fib2(10).(int), 55)

	o := func(str ...T) T { return str[0] }
	fastO := Memoize(o)
	asserts.Equals(t, "blah", o("toString").(string), "toString")
	asserts.Equals(t, "blah blah", fastO("toString").(string), "toString")
}


func TestOnce(t *testing.T) {
	num := 0
	increment := Once(func(...T) T {
		num += 1
		return num
	})
	increment()
	increment()
	asserts.IntEquals(t, "can increment once", num, 1)
}

// TODO: missing recursive once

func TestWrap(t *testing.T) {
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
	greet := func(name ...T) T { return fmt.Sprintf("hi: %v", name[0]) }
	backwards := Wrap(greet, func(args ...T) T {
		fn := args[0].(func(...T) T)
		name := args[1]
		return fmt.Sprintf("%v %v", fn(name),
			reverse(name.(string)))
	})
	asserts.Equals(t, "wrapped the salutation function", backwards("moe").(string), "hi: moe eom")

	inner := func(...T) T { return "Hello " }
	obj := map[T]T{"name": "Moe"}
	hi := Wrap(inner, func(args ...T) T {
		fn := args[0].(func(...T) T)
		this := args[1].(map[T]T)
		return fn().(string) + this["name"].(string)
	})
	asserts.Equals(t, "pass obj as arg ", hi(obj).(string), "Hello Moe")

	noop := func(...T) T { return nil }
	wrapped := Wrap(noop, func(args ...T) T {
		return args
	})
	ret := wrapped([]T{"whats", "your"}, "vector", "victor")
	_, ok := ret.([]T)[0].(func(...T) T)
	asserts.True(t, "noop test first arg is a func", ok)
	asserts.Equals(t, "noop test, rest of args", fmt.Sprint(ret.([]T)[1:]),
		"[[whats your] vector victor]")
}

func TestCompose(t *testing.T) {
	greet := func(name ...T) T { return "hi: " + name[0].(string) }
	exclaim := func(sentence ...T) T { return sentence[0].(string) + "!" }
	pause := func(midway ...T) T { return midway[0].(string) + ", " }
	composed := Compose(exclaim, greet)

	asserts.Equals(t, "can compose a function that takes another", composed("moe").(string), "hi: moe!")

	composed2 := Compose(greet, exclaim)
	asserts.Equals(t, "in this case, the functions are also commutative", composed2("moe").(string), "hi: moe!")

	composed3 := Compose(greet, pause, exclaim)
	asserts.Equals(t, "in this case, the functions are not commutative", composed3("moe").(string), "hi: moe!, ")

	composed4 := Compose(greet, exclaim, pause)
	asserts.Equals(t, "in this case, the functions are not commutative", composed4("moe").(string), "hi: moe, !")
}

func TestAfter(t *testing.T) {
	testAfter := func(afterAmount, timesCalled int) int {
		afterCalled := 0
		after := After(afterAmount, func(...T) T {
			afterCalled += 1
			return nil
		})
		for timesCalled > 0 {
			timesCalled -= 1
			after()
		}
		return afterCalled
	}

	asserts.IntEquals(t, "after(N) should fire after being called N times",
		testAfter(5, 5), 1)
	asserts.IntEquals(t, "after(N) should not fire unless called N times",
		testAfter(5, 4), 0)
	asserts.IntEquals(t, "after(0) should not fire immediately",
		testAfter(0, 0), 0)
	asserts.IntEquals(t, "after(0) should fire when first invoked",
		testAfter(0, 1), 1)
}


func TestNow(t *testing.T) {
	diff := Now() - time.Now().UnixNano()
	asserts.True(t, "Produces the correct time in milliseconds", diff <= 0 && diff > -50 );//within 50ns
}

func TestDelay(t *testing.T) {
	delayed := false
	Delay(func() T { delayed = true; return delayed }, 100)
	Delay(func() T {
		asserts.False( t, "didn't delay the function quite yet", delayed)
		return delayed },
		50)
	Delay(func() T {
		asserts.True( t, "delayed the function", delayed)
		return delayed
	},
		150)

	// wait for Delay calls above to run their course
	select {
		case <-time.After(300 * time.Millisecond):
		break
	}
}


func TestDefer(t *testing.T) {
	deferred := false
	func(){
		defer func(boole bool){ deferred = boole }(true)
		Delay(func() T{ asserts.Ok(t, "deferred the function", deferred) ; return deferred }, 50)
	}()
	// wait for Delay calls above to run their course
	select {
		case <-time.After(100 * time.Millisecond):
		break
	}
}

// XXX: debounce, tests in progress
func TestDebounce(t *testing.T) {
	var counter int = 0
	incr := func() T { counter += 1
		return counter }
	debouncedIncr := Debounce(incr, 32)
	debouncedIncr()
	debouncedIncr()
	Delay(debouncedIncr, 32)
	//Delay(func(){ asserts.IntEquals(t, "incr was debounced", counter, 0) },
		//96)
	select {
		case <-time.After(5 * time.Millisecond):
			asserts.IntEquals(t, "incr was debounced", counter, 0)
		break
	}
	select {
		case <-time.After(90 * time.Millisecond):
			asserts.IntEquals(t, "incr was debounced", counter, 1)
		break
	}
}

func TestDebounceASAP(t *testing.T) {
	counter := 0
	incr := func() T { counter += 1; return counter }
	debouncedIncr := Debounce(incr, 64, true)
	a := debouncedIncr().(int)
	b := debouncedIncr().(int)
	asserts.IntEquals(t, "debounced immediate return a",a, 1)
	asserts.IntEquals(t, "debounced immediate return b",b, 1)
	asserts.IntEquals(t, "incr was called immediately",1,counter)
	Delay(debouncedIncr, 16)
	Delay(debouncedIncr, 32)
	Delay(debouncedIncr, 48)
	 globalWg.Add(1)
	/*w*/Delay(func () T {
			asserts.IntEquals(t, "incr was debounced", counter, 1)
		globalWg.Done()
		return counter
	},128)
}

func TestDebounceASAPRecursively(t *testing.T) {
	counter := 0
	var debouncedIncr func()T
	debouncedIncr = Debounce(func() T{
			counter += 1
			if counter < 10 {
				debouncedIncr()
			}
			return counter
		}, 32, true)
	debouncedIncr()
	asserts.IntEquals(t, "incr was called immediately", counter, 1)
	//Delay(func() T{ asserts.IntEquals( t, "Incr was debounced, recursively", counter, 1); return nil },
		//96)
	select {
		case <-time.After(96 * time.Millisecond):
			asserts.IntEquals(t, "incr was debounced, recursively", counter, 1)
		break
	}
}

/* not relevant to Go, can't redefine Now()
func TestDebounceAfterSystemTimeIsMuckedWith(t *testing.T) {
    counter := 0
    origNowFunc := Now
    debouncedIncr := Debounce(func()T{ counter += 1; return counter }, 100, true)
    debouncedIncr();
    asserts.IntEquals(t, "Incr called immediately",counter, 1)
    Now = func() int64 {
      return 201301111
    }

	select {
		case <-time.After(200 * time.Millisecond):
			debouncedIncr()
			asserts.IntEquals(t, "incr was debounced, successively", counter, 2)
			Now = origNowFunc
		break
	}
  }
*/

func TestThrottle(t *testing.T){
	counter := 0
	incr := func(...T) T { counter += 1; return counter}
	throttledIncr := Throttle(incr, 32)
	throttledIncr()
	throttledIncr()
   asserts.IntEquals(t, "incr was called immediately", counter, 1)
	 globalWg.Add(1)
	/*w*/Delay(func () T {
		asserts.IntEquals(t, "incr was throttled", counter, 2)
		globalWg.Done()
		return counter
	},64)
}

func TestThrottleWithArgs (t *testing.T){
	var value int = 0
	update := func(val ...T) T {
		if len(val) > 0 {
			value = (val[0]).(int)
		}
		return value
	}
	throttledUpdate := Throttle(update, 32)
	throttledUpdate(1)
	throttledUpdate(2)
	Delay(func() T {
		throttledUpdate(3)
		return value
	}, 85)
	asserts.IntEquals(t, "updated to first value", value, 1)
	 globalWg.Add(1)
	/*w*/Delay(func() T {
		asserts.IntEquals(t, "updated to latest value", value, 3)
		globalWg.Done()
		return nil }, 100)
}

func TestThrottleOnce(t *testing.T){
    counter := 0
    incr := func(...T) T{ counter += 1; return counter }
    throttledIncr := Throttle(incr, 2)
    result := throttledIncr().(int)
	 globalWg.Add(1)
    /*w*/Delay(func()T {
      asserts.IntEquals( t, "throttled functions return their value", result, 1)
      asserts.IntEquals( t, "Incr was called once ", counter , 1)
		globalWg.Done()
		return nil
    }, 6)
}

func TestThrottleTwice(t *testing.T){
    counter := 0
    incr := func(...T) T{ counter += 1; return counter }
    throttledIncr := Throttle(incr, 2)
    throttledIncr()
    throttledIncr()
	 globalWg.Add(1)
    /*w*/Delay(func()T {
      asserts.IntEquals( t, "Incr was called twice ", counter , 2)
		globalWg.Done()
		return nil
    }, 6)
}

func TestMoreThrottling(t *testing.T){
	counter := 0
	incr := func(...T) T{ counter += 1; return counter }
	throttledIncr := Throttle(incr, 2)
	throttledIncr()
	throttledIncr()
   asserts.True(t, "checking counter amidst throttling", counter == 1)
	 globalWg.Add(1)
   /*w*/Delay(func()T {
      asserts.Ok(t,"checking counter amidst throttling, 2",counter == 2)
      throttledIncr()
      asserts.Ok(t,"checking counter amidst throttling, 3",counter == 3)
		globalWg.Done()
		return nil
   }, 85)
}


func TestThrottleRepeatedlyWithResults( t *testing.T ) {
	counter := 0
	incr := func(...T) T{ counter += 1; return counter }
	throttledIncr := Throttle(incr, 100)
   results := []int{}
	saveResult := func() T { results = append(results,throttledIncr().(int)); return results }
	saveResult()
	saveResult()
   Delay(saveResult, 50)
   Delay(saveResult, 150)
   Delay(saveResult, 160)
   Delay(saveResult, 230)
	var localWg sync.WaitGroup
	localWg.Add(1)
   /*w*/Delay(func() T {
      asserts.IntEquals( t, "incr was called once", results[0], 1)
      asserts.IntEquals( t, "incr was throttled", results[1], 1)
      asserts.IntEquals( t, "incr was throttled", results[2], 1)
      asserts.IntEquals( t, "incr was called twice", results[3], 2)
      asserts.IntEquals( t, "incr was throttled", results[4], 2)
      asserts.IntEquals( t, "incr was called trailing", results[5], 3)
		localWg.Done()
		return nil
    }, 300)
	localWg.Wait()
}

func TestThrottleTriggersTrailingCallWhenInvokedRepeatedly( t *testing.T) {
    counter := 0
    counterRet := 0
    var limit int64 = 4800
    incr := func(...T) T{ counter += 1; return counter }
    throttledIncr := Throttle(incr, 32)

    stamp := Now()
    for (Now() - stamp) < limit {
      throttledIncr()
		counterRet += 1
    }
    lastCount := counter
    asserts.Ok( t, "Trailing test", counterRet > 1)
	 globalWg.Add(1)
    Delay(func() T {
		asserts.Ok( t, "Counter > lastCount", counter > lastCount)
		globalWg.Done()
		return nil
    }, 96)
}

func TestThrottleDoesNotTriggerLeadingCallWhenLeadingIsFalse (t *testing.T){
	counter := 0
	incr := func(...T) T{ counter += 1; return counter }
	throttledIncr := Throttle(incr, 60, map[string]bool{"leading":false,"name":true})
	throttledIncr()
	throttledIncr()
	asserts.IntEquals( t, "Throttle does not trigger leading call 1", counter , 0)
	var localWg sync.WaitGroup
	localWg.Add(1)
	Delay(func() T {
		asserts.IntEquals( t, "Throttle does not trigger leading call 2", counter , 1)
		localWg.Done()
		return nil
    }, 100)
	localWg.Wait()
}


func TestMoreThrottleDoesNotTriggerLeadingCallWhenLeadingIsFalse (t *testing.T){
	counter := new(int)
	incr := func(...T) T{ *counter += 1; return counter }
	throttledIncr := Throttle(incr, 100, map[string]bool{"leading":false})
	throttledIncr()
   Delay(func()T{ throttledIncr(); return *counter }, 50)
   Delay(func()T{ throttledIncr(); return *counter }, 60)
   Delay(func()T{ throttledIncr(); return *counter }, 200)
   Delay(func()T{ throttledIncr(); return *counter }, 350)

	var localWg sync.WaitGroup
	localWg.Add(2)
   Delay(func()T {
		asserts.IntEquals(t, "More throttle, first delay, counter should be 1", *counter , 1)
		localWg.Done()
		return nil
    }, 250)

   Delay(func() T{
		asserts.IntEquals(t, "More throttle, second delay, counter should be 2", *counter , 2)
		localWg.Done()
		return nil
    }, 350)
   asserts.IntEquals(t, "More throttle, counter still zero ", *counter , 0)
	localWg.Wait()

	/*select {
		case <-time.After(1000):
	}
	fmt.Printf("end counter %v, now %v\n",*counter,time.Now().UnixNano() / 1000000);
	*/
}

/*
  asyncTest('one more throttle with leading: false test', 2, function() {
    var counter = 0;
    var incr = function(){ counter++; };
    var throttledIncr = _.throttle(incr, 100, {leading: false});

    var time = new Date;
    while (new Date - time < 350) throttledIncr();
    ok(counter <= 3);

    _.delay(function() {
      ok(counter <= 4);
      start();
    }, 200);
  });

  asyncTest('throttle does not trigger trailing call when trailing is set to false', 4, function() {
    var counter = 0;
    var incr = function(){ counter++; };
    var throttledIncr = _.throttle(incr, 60, {trailing: false});

    throttledIncr(); throttledIncr(); throttledIncr();
    ok(counter === 1);

    _.delay(function() {
      ok(counter == 1);

      throttledIncr(); throttledIncr();
      ok(counter == 2);

      _.delay(function() {
        ok(counter == 2);
        start();
      }, 96);
    }, 96);
  });
*/
/*
  asyncTest('throttle continues to function after system time is set backwards', 2, function() {
    var counter = 0;
    var incr = function(){ counter++; };
    var throttledIncr = _.throttle(incr, 100);
    var origNowFunc = _.now;

    throttledIncr();
    ok(counter == 1);
    _.now = function () {
      return new Date(2013, 0, 1, 1, 1, 1);
    };

    _.delay(function() {
      throttledIncr();
      ok(counter == 2);
      start();
      _.now = origNowFunc;
    }, 200);
  });

*/

func TestWaitPlaceholder(t *testing.T){
	globalWg.Wait()
}
