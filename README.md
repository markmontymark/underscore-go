# underscore-go

A Go port of Underscore.js

# setup

## Install gvm

	git clone git@github.com:moovweb/gvm.git
	cd gvm/
	vi binscripts/gvm-installer 
	./binscripts/gvm-installer 
	source $HOME/.gvm/scripts/gvm

## Install go

	gvm listall
	gvm install go1.1.2
	gvm use go1.1.2


# testing

	git clone git@github.com/markmontymark/underscore-go.git
	cd underscore-go/
	go test

# status

	Most of the Underscore.js functionality is ported, finishing porting tests


# todos

 -	Finish porting tests from chaining.js, speed.js
 - Add the rest of the public api as methods on a \*Underscore object
 - Port Underscore.js functions like defer, throttle, debounce using setTimeout with go-routines and Go's Defer language support

		

