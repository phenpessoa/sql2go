package parser

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

func Parse[T any](dst *T, r io.Reader) error {
	v := reflect.ValueOf(dst).Elem()

	data, _ := io.ReadAll(r)
	input := string(data)
	l := newLexer(input)
	p := newParser(l)
	tree := p.parse()

	i := 0
	for i < len(tree.nodes) {
		n := tree.nodes[i]
		switch t := n.(type) {
		case nodeName:
			if !t.valid {
				return errors.New("sql2go: found an empty name")
			}

			field := v.FieldByName(t.val)
			if !field.IsValid() || !field.CanSet() || !field.CanInterface() {
				return fmt.Errorf(
					"sql2go: field not found or invalid in dst struct: %s",
					t.val,
				)
			}

			if _, ok := field.Interface().(string); !ok {
				return fmt.Errorf(
					"sql2go: field %s is not of type string", t.val,
				)
			}

			var (
				query    strings.Builder
				lastByte byte
			)
			i++
			for _, nn := range tree.nodes[i:] {
				switch t := nn.(type) {
				case nodeEnfOfQuery:
					query.Grow(1)
					query.WriteByte(';')
					lastByte = ';'
				case nodeName:
					goto out
				case nodeQuery:
					val := strings.TrimSpace(t.val)
					if lastByte == '\'' || lastByte == '"' || lastByte == '`' {
						query.Grow(len(val) + 1)
						query.WriteRune(' ')
					} else {
						query.Grow(len(val))
					}
					query.WriteString(val)
					lastByte = val[len(val)-1]
				case nodeStringLiteral:
					if lastByte != '\n' && lastByte != ' ' {
						query.Grow(len(t.val) + 3)
						query.WriteByte(' ')
					} else {
						query.Grow(len(t.val) + 2)
					}
					query.WriteByte('\'')
					query.WriteString(t.val)
					query.WriteByte('\'')
					lastByte = '\''
				case nodeIdentifier:
					if lastByte != '\n' && lastByte != ' ' {
						query.Grow(len(t.val) + 3)
						query.WriteByte(' ')
					} else {
						query.Grow(len(t.val) + 2)
					}
					query.WriteString(t.tok.literal)
					query.WriteString(t.val)
					query.WriteString(t.tok.literal)
					lastByte = t.tok.literal[0]
				case nodeNewLine:
					if lastByte != '\n' {
						query.Grow(1)
						query.WriteByte('\n')
					}
					lastByte = '\n'
				}

				i++
			}

		out:
			field.Set(reflect.ValueOf(strings.TrimSpace(query.String())))
		default:
			i++
		}
	}

	return nil
}
