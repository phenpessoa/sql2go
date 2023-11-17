package sql2go

import (
	"io"

	"github.com/phenpessoa/sql2go/internal/parser"
)

func Parse[T any](dst *T, r io.Reader) error {
	return parser.Parse(dst, r)
}
