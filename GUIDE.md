# Guide for builq

## Safety & Sanitization

Query arguments should be passed via `%$` or `%?`, this way they won't appear in the query string but instead will be appended to the arguments slice (2nd return value of the `Builder.Build()` method), thus preventing potential SQL injections.

Examples in [example_test.go](example_test.go) explicitly show that query arguments are not a part of the `query` string but returned separately as the `args` slice.

The `%s` verb should be used with an extra care, no user input should be passed through it.

## String placeholder

To write just a string there is the `%s` formatting verb. Works the same as in the `fmt` package.

Please note that unlike `fmt`, `builq` does not support width and explicit argument indexes.

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

Even with the `fmt` package speed is very good. If case you want zero-allocation query builder consider to cache query and just use it's value (works only for static queries like in `BenchmarkBuildCached`)

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