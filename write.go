package builq

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (b *Builder) write(sb *strings.Builder, resArgs *[]any, s string, args ...any) error {
	for argID := 0; ; argID++ {
		n := strings.IndexByte(s, '%')
		if n == -1 {
			if argID != len(args) {
				b.setErr(errTooManyArguments)
			}

			sb.WriteString(s)
			sb.WriteByte(b.sep)
			return nil
		}

		sb.WriteString(s[:n])

		if argID >= len(args) {
			return errTooFewArguments
		}

		arg := args[argID]

		s = s[n+1:] // skip '%'
		switch verb := s[0]; verb {
		case '$', '?', 's', 'd':
			s = s[1:]
			b.writeArg(sb, resArgs, verb, arg)

		case '+', '#':
			isBatch := verb == '#'
			s = s[1:]
			if len(s) < 1 {
				b.setErr(errIncorrectVerb)
				continue
			}

			switch verb := s[0]; verb {
			case '$', '?':
				s = s[1:]

				if isBatch {
					b.writeBatch(sb, resArgs, verb, arg)
				} else {
					b.writeSlice(sb, resArgs, verb, arg)
				}
			default:
				b.setErr(errUnsupportedVerb)
			}
		default:
			b.setErr(errUnsupportedVerb)
		}
	}
}

func (b *Builder) writeBatch(sb *strings.Builder, resArgs *[]any, verb byte, arg any) {
	for i, arg := range b.asSlice(arg) {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteByte('(')
		b.writeSlice(sb, resArgs, verb, arg)
		sb.WriteByte(')')
	}
}

func (b *Builder) writeSlice(sb *strings.Builder, resArgs *[]any, verb byte, arg any) {
	for i, arg := range b.asSlice(arg) {
		if i > 0 {
			sb.WriteString(", ")
		}
		b.writeArg(sb, resArgs, verb, arg)
	}
}

func (b *Builder) writeArg(sb *strings.Builder, resArgs *[]any, verb byte, arg any) {
	if b.debug {
		b.writeDebug(sb, arg)
		return
	}

	var isSimple bool

	switch verb {
	case '$':
		b.counter++
		sb.WriteByte('$')
		sb.WriteString(strconv.Itoa(b.counter))
		*resArgs = append(*resArgs, arg)
	case '?':
		sb.WriteByte('?')
		*resArgs = append(*resArgs, arg)
	case 's':
		isSimple = true
		switch arg := arg.(type) {
		case string:
			sb.WriteString(arg)
		case fmt.Stringer:
			sb.WriteString(arg.String())
		default:
			fmt.Fprint(sb, arg)
		}
	case 'd':
		isSimple = true
		b.assertNumber(arg)
		fmt.Fprint(sb, arg)
	}

	// ok to have many simple placeholders
	if isSimple {
		return
	}

	// store the first placeholder used in the query
	// to check for mixed placeholders later
	if b.placeholder == 0 {
		b.placeholder = verb
	}
	if b.placeholder != verb {
		b.setErr(errMixedPlaceholders)
	}
}

func (b *Builder) writeDebug(sb *strings.Builder, arg any) {
	switch arg := arg.(type) {
	case Columns:
		sb.WriteString(arg.String())
	case time.Time:
		sb.WriteByte('\'')
		sb.WriteString(arg.UTC().Format("2006-01-02 15:04:05:999999"))
		sb.WriteByte('\'')
	case fmt.Stringer:
		sb.WriteByte('\'')
		sb.WriteString(arg.String())
		sb.WriteByte('\'')
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		fmt.Fprint(sb, arg)
	case string:
		sb.WriteByte('\'')
		sb.WriteString(arg)
		sb.WriteByte('\'')
	default:
		sb.WriteByte('\'')
		fmt.Fprint(sb, arg)
		sb.WriteByte('\'')
	}
}

func (b *Builder) assertNumber(v any) {
	switch v.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
	default:
		b.setErr(errNonNumericArg)
	}
}
