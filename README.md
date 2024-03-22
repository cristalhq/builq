# builq

[![build-img]][build-url]
[![pkg-img]][pkg-url]
[![reportcard-img]][reportcard-url]
[![coverage-img]][coverage-url]
[![version-img]][version-url]

Easily build queries in Go.

## Rationale

The simplest way to represent SQL query is a string.
But the query arguments and their indexing (`$1`, `$2` etc) require additional attention.
This tiny library helps to build queries and handles parameter indexing.

## Features

* Simple and easy.
* Safe and fast.
* Tested.
* Language-agnostic.
* Dependency-free.

See [docs][pkg-url] or [GUIDE.md](GUIDE.md) for more details.

## Install

Go version 1.19+

```
go get github.com/cristalhq/builq
```

## Example

```go
cols := builq.Columns{"foo, bar"}

q := builq.New()
q("SELECT %s FROM %s", cols, "users")
q("WHERE active IS TRUE")
q("AND user_id = %$ OR user = %$", 42, "root")

query, args, err := q.Build()
if err != nil {
	panic(err)
}

debug := q.DebugBuild()

fmt.Println("query:")
fmt.Println(query)
fmt.Println("\nargs:")
fmt.Println(args)
fmt.Println("\ndebug:")
fmt.Println(debug)

// query:
// SELECT foo, bar FROM users
// WHERE active IS TRUE
// AND user_id = $1 OR user = $2
//
// args:
// [42 root]
//
// debug:
// SELECT foo, bar FROM 'users'
// WHERE active IS TRUE
// AND user_id = 42 OR user = 'root'
```

See examples: [example_test.go](example_test.go).

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
