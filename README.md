# underscore-go

A Go port of Underscore.js

# setup



## install Go with gvm

	git clone git@github.com:moovweb/gvm.git
	cd gvm/
	vi binscripts/gvm-installer 
	./binscripts/gvm-installer 
	source $HOME/.gvm/scripts/gvm

## Install go

	gvm listall
	gvm install go1.2.1
	gvm use go1.2.1


# Testing

	git clone git@github.com/markmontymark/asserts.git
	cd asserts/
	go install
	git clone git@github.com/markmontymark/underscore-go.git
	cd underscore-go/
	go test

# Benchmarking

	cd underscore-go/
	go test -bench=.

 - so far, only a few Array functions are benchmarked, but it's already interesting to see what's fast and what's slow

# TODO

 -	Port speed.js to speed\_test.go
 - Add Bench\*() functions for benchmarking
 - Add Example\*() functions for showing how to use this code
 - Export an Underscore Interface interface, to allow different types to provide enumeration, possibly something like C#'s IEnumerable?

# Example usage

 - adapted from arrays\_test.go

		package main

		import "fmt"
		import "underscore"

		func main () {
			num := 35
			numbers2 := []underscore.T{10, 20, 30, 40, 50}
			if v := underscore.IndexOf(numbers2, num, func (a underscore.T, b underscore.T) bool { return a.(int) < b.(int)}, true); v == -1 {
				fmt.Println("35 is not in the list")
			}
		}

	- Copy the above code in a file, runme.go, and then on the command line, `go run runme.go`

