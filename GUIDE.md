# Guide for builq

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
