# Guide for builq

## Safety & Sanitization

If query args/params are passed via `%$` or `%?` they will not appear in query but will be appended to the arguments slice (2nd resulting param of `Builder.Build()` method).

`ExampleQuery1` in [examples_test.go](https://github.com/cristalhq/builq/blob/main/example_test.go) explicitly shows that arguments are not a part of `query` string but are in `args` slice.

Such usage prevents potential SQL-injection and/or any other harmful user inputs. However, `%s` should be used with a care and no user input should be passed through `%s`,

## String placeholder

To write just a string there is the `%s` formatting verb. Works the same as in the `fmt` package.

## Argument placeholder

`builq` supports only 2 popular formats:

* PostgreSQL via `%$` (`$1, $2, $3..`)
* MySQL/SQLite via `%?` (`?, ?, ?..`)

This should cover almost all available databases, if not - feel free to make an issue.

## Slice/batch modifiers

Both `%$` and `%?` formats can be extended with `+` or `#`:

* `%+$` will expand slice argument as `$1, $2, ... $(len(arg)-1)`
* `%#?` will expand slice of slices argument as `(?, ?), (?, ?), ... (?, ?)`

Argument must be a slice (for `+`) or a slice of slices (for `#`), otherwise the `.Build()` method returns an error.

## Speed

Even with `fmt` packge speed is very good. If case you want zero-allocation query builder consider to cache query and just use it's value (works only for static queries like in `BenchmarkBuildCached`)

```
goos: darwin
goarch: arm64
pkg: github.com/cristalhq/builq
BenchmarkBuildNaive
BenchmarkBuildNaive-10       	  299349	      4214 ns/op	     488 B/op	       8 allocs/op
BenchmarkBuildManyArgs
BenchmarkBuildManyArgs-10    	   85747	     14027 ns/op	    1056 B/op	      24 allocs/op
BenchmarkBuildCached
BenchmarkBuildCached-10      	1000000000	         0.6213 ns/op	       0 B/op	       0 allocs/op
PASS
```