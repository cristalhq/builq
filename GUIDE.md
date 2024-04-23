# Guide for builq

## Safety & Sanitization

Query arguments should be passed via `%$`, `%?`, or `%@`, this way they won't appear in the query string but instead will be appended to the arguments slice (2nd return value of the `Builder.Build()` method), thus preventing potential SQL injections.

Examples in [example_test.go](example_test.go) explicitly show that query arguments are not a part of the `query` string but returned separately as the `args` slice.

The `%s` verb should be used with an extra care, no user input should be passed through it.

## SQL injections and `%s` usage

`builq` is language-agnostic query builder that doesn't differentiate Postgres SQL syntax from MySQL. The `%s` verb was introduced to give flexibility to the library users.

The following query is valid for `builq` but isn't valid for Postgres:

```go
const tableName = "my_table"
user := "admin"

q := builq.New()
q("SELECT * FROM %$ WHERE username = %$", tableName, user)

// will generate query: SELECT * FROM $1 WHERE username = $2
```

The query above is correct for `builq` and is incorrect for Postgres (error `SQLSTATE 42601`). Exactly for such cases, `%s` was added:

```go
q("SELECT * FROM %s WHERE username = %$", tableName, user)

// will generate query: SELECT * FROM my_table WHERE username = $1
```

Remember that `%s` should be used with care, and as mentioned in the section above, no user input should be passed via `%s`.

## Compile-time queries

To enforce compile-time queries `builq.Builder` accepts only constant strings:

```go
var sb builq.Builder
sb.Addf("SELECT %s FROM %s", cols, "table")
sb.Addf("WHERE id = %$", 123)

// this WILL NOT complile, orClause isn't const
// var orClause = "OR id = %$"
// sb.Addf(orClause, 42)

// WILL compile, orClause2 is known at compile-time
const orClause2 = "OR id = %$"
sb.Addf(orClause2, 42)
```

The reason behing this API is to improve security and to prevent bad runtime queries.
Also, some projects require constant queries due to security policies (precise definition might be different but you get the idea).

## String placeholder

To write just a string there is the `%s` formatting verb. Works the same as in the `fmt` package.

Please note that unlike `fmt`, `builq` does not support width and explicit argument indexes.

## Argument placeholder

`builq` supports 3 formats:

* PostgreSQL via `%$` (`$1, $2, $3..`)
* MySQL/SQLite via `%?` (`?, ?, ?..`)
* MSSQL via `%@` (`@p1, @p2, @p3..`)

This should cover almost all available databases, if not - feel free to make an issue.

## Slice/batch modifiers

All formats can be extended with `+` or `#`:

* `%+$` will expand slice argument as `$1, $2, ... $(len(arg)-1)`
* `%#?` will expand slice of slices argument as `(?, ?), (?, ?), ... (?, ?)`

Argument must be a slice (for `+`) or a slice of slices (for `#`), otherwise the `.Build()` method returns an error.

## Debug

The convenience `DebugBuild` method can be used to debug queries.
Unlike `Build` it returns a complete query with all the placeholders being replaced with their arguments.
The query can then be copy-pasted and executed directly from DB.
While handy during development this method could lead to SQL injections, so be careful and avoid it in production code.

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
