This tool aims to provide regex based refactoring of Go code.

## Warnings:
* This is currently an experiment/itch scratch
* There are likely many many corner cases or even common use-cases not yet covered

## Examples

### Replace any number of params using 1 wildcard
`$ refex -b 'statsd.Count1($1$)' -a 'stats.D.Count1($1$)' mycode.go`
 
Before:

	package mypackage
	
	func something() {
		statsd.Count1("call")
		statsd.Count1("me")
		statsd.Count1("baby")
	}

After:

	package mypackage
	
	func something() {
		stats.D.Count1("call")
		stats.D.Count1("me")
		stats.D.Count1("baby")
	}

### Replace 2 params using 2 wildcards
`$ refex -b 'statsd.Count1($1$, $2$)' -a 'stats.D.Count1($1$, $2$)' mycode.go`
 
Before:

	package mypackage
	
	func something() {
		statsd.Count1("don't", "call")
		statsd.Count1("me", "baby")
	}

After:

	package mypackage
	
	func something() {
		stats.D.Count1("don't", "call")
		stats.D.Count1("me", "baby")
	}

### Replace 2 params using 2 wildcards but change the order
`$ refex -b 'statsd.Count1($1$, $2$)' -a 'stats.D.Count1($2$, $1$)' mycode.go`
 
Before:

	package mypackage
	
	func something() {
		statsd.Count1("don't", "call")
		statsd.Count1("me", "baby")
	}

After:

	package mypackage
	
	func something() {
		stats.D.Count1("call", "don't")
		stats.D.Count1("baby", "me")
	}


### Provided examples

* `$ refex -d -b 'rand.Intn($1$)' -a 'rand.In63n($1$)' test-data/example1.go`
* `$ refex -d -b 'fmt.Print($1$)' -a 'fmt.Fprintf(os.Stderr, $1$)' test-data/example2.go`
