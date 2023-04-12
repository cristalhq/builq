package builq

import (
	"fmt"
	"strconv"
	"strings"
)

func (b *Builder) write(s string, args ...any) error {
	for argID := 0; ; argID++ {
		n := strings.IndexByte(s, '%')
		if n == -1 {
			if argID != len(args) {
				b.setErr(errTooManyArguments)
			}

			b.query.WriteString(s)
			b.query.WriteByte(b.sep)
			return nil
		}

		b.query.WriteString(s[:n])

		if argID >= len(args) {
			return errTooFewArguments
		}

		arg := args[argID]

		s = s[n+1:] // skip '%'
		switch verb := s[0]; verb {
		case 's', '$', '?':
			s = s[1:]
			b.writeArg(verb, arg)

		case '+', '#':
			isBatch := verb == '#'
			s = s[1:]
			if len(s) < 1 {
				// TODO: Error
				continue
			}

			switch verb := (s[0]); verb {
			case '$', '?':
				s = s[1:]

				if isBatch {
					b.writeBatch(verb, arg)
				} else {
					b.writeSlice(verb, arg)
				}
			default:
				b.setErr(errUnsupportedVerb)
			}
		default:
			b.setErr(errUnsupportedVerb)
		}
	}
}

func (b *Builder) writeBatch(verb byte, arg any) {
	for i, arg := range b.asSlice(arg) {
		if i > 0 {
			b.query.WriteString(", ")
		}
		b.query.WriteByte('(')
		b.writeSlice(verb, arg)
		b.query.WriteByte(')')
	}
}

func (b *Builder) writeSlice(verb byte, arg any) {
	for i, arg := range b.asSlice(arg) {
		if i > 0 {
			b.query.WriteString(", ")
		}
		b.writeArg(verb, arg)
	}
}

func (b *Builder) writeArg(verb byte, arg any) {
	switch verb {
	case 's':
		switch arg := arg.(type) {
		case string:
			b.query.WriteString(arg)
		case fmt.Stringer:
			b.query.WriteString(arg.String())
		default:
			fmt.Fprint(&b.query, arg)
		}
	case '$':
		b.counter++
		b.query.WriteByte('$')
		b.query.WriteString(strconv.Itoa(b.counter))
		b.args = append(b.args, arg)
	case '?':
		b.query.WriteByte('?')
		b.args = append(b.args, arg)
	}

	// store the first placeholder used in the query
	// to check for mixed placeholders later
	// ignore 's' because it's ok to have in any case
	if b.placeholder == 0 && verb != 's' {
		b.placeholder = rune(verb)
	}
	if verb != 's' && b.placeholder != rune(verb) {
		b.setErr(errMixedPlaceholders)
	}
}
