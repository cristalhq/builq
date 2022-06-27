# builq

[![build-img]][build-url]
[![pkg-img]][pkg-url]
[![reportcard-img]][reportcard-url]
[![coverage-img]][coverage-url]
[![version-img]][version-url]

Easily build queries in Go.

## Rationale

Often simple string representing SQL query is what we are looking for. But the parameters and their indexing (`$1`, `$2` etc) requires additional attention. This tiny library helps to build queries and handles parameter indexing.

## Features

* Simple.
* Tested.
* Dependency-free.

## Install

Go version 1.17+

```
go get github.com/cristalhq/builq
```

## Example

```go
q := builq.Newf("SELECT %s FROM %s", "foo, bar", "users")
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
