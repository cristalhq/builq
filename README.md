# builq

[![build-img]][build-url]
[![pkg-img]][pkg-url]
[![reportcard-img]][reportcard-url]
[![coverage-img]][coverage-url]
[![version-img]][version-url]

Easily build queries in Go.

## Rationale

Simplest way to represent SQL query is a string. But the query arguments and their indexing (`$1`, `$2` etc) requires additional attention. This tiny library helps to build queries and handles parameter indexing.

## Features

* Simple and easy.
* Tested.
* Dependency-free.

See [GUIDE.md](https://github.com/cristalhq/builq/blob/main/GUIDE.md) for more details.

## Install

Go version 1.19+

```
go get github.com/cristalhq/builq
```

## Example

```go
cols := builq.Columns{"foo, bar"}

var b builq.Builder
b.Addf("SELECT %s FROM %s", cols, "users").
	Addf("WHERE active IS TRUE").
	Addf("AND user_id = %$ OR user = %$", 42, "root")

query, args, err := b.Build()
if err != nil {
	panic(err)
}

fmt.Printf("query:\n%v", query)
fmt.Printf("args:\n%v", args)

// Output:
//
// query:
// SELECT foo, bar FROM users
// WHERE active IS TRUE
// AND user_id = $1 OR user = $2
// args:
// [42 root]
```

See examples: [examples_test.go](https://github.com/cristalhq/builq/blob/main/example_test.go).

## Documentation

See [these docs][pkg-url].

## License

[MIT License](LICENSE).

[build-img]: https://github.com/cristalhq/builq/workflows/build/badge.svg
[build-url]: https://github.com/cristalhq/builq/actions
[pkg-img]: https://pkg.go.dev/badge/cristalhq/builq
[pkg-url]: https://pkg.go.dev/github.com/cristalhq/builq
[reportcard-img]: https://goreportcard.com/badge/cristalhq/builq
[reportcard-url]: https://goreportcard.com/report/cristalhq/builq
[coverage-img]: https://codecov.io/gh/cristalhq/builq/branch/main/graph/badge.svg
[coverage-url]: https://codecov.io/gh/cristalhq/builq
[version-img]: https://img.shields.io/github/v/release/cristalhq/builq
[version-url]: https://github.com/cristalhq/builq/releases
