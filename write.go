package builq

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (b *Builder) write(sb *strings.Builder, resArgs *[]any, s string, args ...any) error {
	for argID := 0; ; argID++ {
		idx := strings.IndexByte(s, '%')
		if idx == -1 {
			var err error
			if len(args) != argID {
				err = fmt.Errorf("%w: have %d args, expected %d", errTooManyArguments, len(args), argID)
			}

			sb.WriteString(s)
			sb.WriteByte(b.sep)
			return err
		}

		sb.WriteString(s[:idx])

		s = s[idx+1:] // skip '%'
		if len(s) == 0 {
			return errLonelyModifier
		}

		switch verb := s[0]; verb {
		case '$', '?', 's', 'd':
			if argID >= len(args) {
				return fmt.Errorf("%w: have %d args, want %d", errTooFewArguments, len(args), argID+1)
			}

			arg := args[argID]
			s = s[1:]
			b.writeArg(sb, resArgs, verb, arg)

		case '+', '#':
			isBatch := verb == '#'
			s = s[1:]
			if len(s) < 1 || s[0] == ' ' {
				return fmt.Errorf("%w: '%c' requires additional '$' or '?'", errIncorrectVerb, verb)
			}

			switch verb := s[0]; verb {
			case '$', '?':
				if argID >= len(args) {
					return fmt.Errorf("%w: have %d args, want %d", errTooFewArguments, len(args), argID+1)
				}

				arg := args[argID]
				s = s[1:]

				if isBatch {
					b.writeBatch(sb, resArgs, verb, arg)
				} else {
					b.writeSlice(sb, resArgs, verb, arg)
				}

			default:
				return fmt.Errorf("%w: '%c' is not supported", errUnsupportedVerb, verb)
			}

		case '%':
			argID--
			s = s[1:]
			sb.WriteByte('%')

		case ' ':
			return errLonelyModifier

		default:
			return fmt.Errorf("%w: '%c' is not supported", errUnsupportedVerb, verb)
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

	switch {
	case b.placeholder == 0:
		b.placeholder = verb
	case b.placeholder != verb:
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

func (b *Builder) asSlice(v any) []any {
	value := reflect.ValueOf(v)

	if value.Kind() != reflect.Slice {
		b.setErr(errNonSliceArgument)
		return nil
	}

	res := make([]any, value.Len())
	for i := 0; i < value.Len(); i++ {
		res[i] = value.Index(i).Interface()
	}
	return res
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

func (b *Builder) setErr(err error) {
	if b.err == nil {
		b.err = err
	}
}
