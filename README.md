# qder

[![build-img]][build-url]
[![pkg-img]][pkg-url]
[![reportcard-img]][reportcard-url]
[![coverage-img]][coverage-url]
[![version-img]][version-url]

Query builder for Go.

## Rationale

Often simple string representing SQL query is what we are looking for. But the parameters and their indexing (`$1`, `$2` etc) requires additional attention. This tiny library helps to build queries and handles parameter indexing.

## Features

* Simple.
* Tested.
* Dependency-free.

## Install

Go version 1.17+

```
go get github.com/cristalhq/qder
```

## Example

```go
q := qder.Newf("SELECT %s FROM %s", "foo, bar", "users")
q.Append("WHERE")
q.Add("active = ", true)
q.Add("AND user_id = ", 42)
q.Append("ORDER BY created_at")
q.Append("LIMIT 100;")

doQuery(q.Query(), q.Args()...)

// Output:
//
// SELECT foo, bar FROM users
// WHERE
// active = $1
// AND user_id = $2
// ORDER BY created_at
// LIMIT 100;
```

## Documentation

See [these docs][pkg-url].

## License

[MIT License](LICENSE).

[build-img]: https://github.com/cristalhq/qder/workflows/build/badge.svg
[build-url]: https://github.com/cristalhq/qder/actions
[pkg-img]: https://pkg.go.dev/badge/cristalhq/qder
[pkg-url]: https://pkg.go.dev/github.com/cristalhq/qder
[reportcard-img]: https://goreportcard.com/badge/cristalhq/qder
[reportcard-url]: https://goreportcard.com/report/cristalhq/qder
[coverage-img]: https://codecov.io/gh/cristalhq/qder/branch/main/graph/badge.svg
[coverage-url]: https://codecov.io/gh/cristalhq/qder
[version-img]: https://img.shields.io/github/v/release/cristalhq/qder
[version-url]: https://github.com/cristalhq/qder/releases
