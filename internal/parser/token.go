package parser

type tokenType uint8

const (
	tokenTypeUndefined tokenType = iota
	tokenTypeEOF

	// -- comment
	tokenTypeComment

	// ;
	tokenTypeSemicolon

	// -- name:
	tokenTypeName

	// '
	tokenTypeStringLiteral

	// " or `
	tokenTypeIdentifier

	// \n
	tokenTypeNewLine

	// anything else
	tokenTypeRawInput
)

type token struct {
	typ     tokenType
	literal string
}
