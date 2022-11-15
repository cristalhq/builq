# Guide for builq

## String placeholder

To write just a string there is `%s` formatting verb. Works the same as in `fmt` package.

## Placeholder formats

`builq` supports only 2 popular formats: PostgreSQL via `%$` (`$1, $2, $3..`) and MySQL/SQLite via `%?` (`?, ?, ?..`).

This should cover almost all avaliable databases, if not - feel free to make an issue.

## Slice/batch modifiers

Both `%$` and `%?` formats can be extended with `+` or `#`.

`%+$` will expand slice argument as `$1, $2, ... $(len(arg)-1)` 

If argument isn't a slice - error will be printed and returned from `.Build()` method.

`%#?` will expand slice of slices argument as `(?1, ?2), (?3,..),... ?(len(arg)-2, ?len(arg)-1)` 

If argument isn't a slice of slices - error will be printed and returned from `.Build()` method.
