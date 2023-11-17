package parser

import (
	"strings"
)

type node interface{}

type nodeComment struct {
	valid bool
	val   string
}

type nodeName struct {
	valid bool
	val   string
}

type nodeStringLiteral struct {
	valid bool
	val   string
}

type nodeIdentifier struct {
	valid bool
	val   string
	tok   token
}

type nodeQuery struct {
	valid bool
	val   string
}

type nodeEnfOfQuery struct{}

type nodeNewLine struct{}

type ast struct {
	nodes []node
}

type parser struct {
	l       *lexer
	lastTok token
	curTok  token
	peekTok token
}

func newParser(l *lexer) *parser {
	p := &parser{l: l}
	return p
}

func (p *parser) nextToken() {
	p.lastTok = p.curTok
	p.curTok = p.peekTok
	p.peekTok = p.l.nextToken()
}

func (p *parser) parse() *ast {
	ast := new(ast)
	p.findFirstName()
	for p.curTok.typ != tokenTypeEOF {
		stmt := p.parseStatement()
		if stmt != nil {
			ast.nodes = append(ast.nodes, stmt)
		}
	}
	return ast
}

func (p *parser) findFirstName() {
	for {
		if p.lastTok.typ != tokenTypeStringLiteral &&
			p.lastTok.typ != tokenTypeIdentifier &&
			p.curTok.typ == tokenTypeName &&
			p.peekTok.typ == tokenTypeRawInput {
			return
		}
		p.nextToken()
		if p.curTok.typ == tokenTypeEOF {
			return
		}
	}
}

func (p *parser) parseStatement() node {
_switch:
	switch p.curTok.typ {
	case tokenTypeName:
		return p.parseNameStatement()
	case tokenTypeComment:
		return p.parseCommentStatement()
	case tokenTypeStringLiteral:
		return p.parseStringLiteralStatement()
	case tokenTypeIdentifier:
		return p.parseIdentifierStatement()
	case tokenTypeRawInput:
		return p.parseQueryStatement()
	case tokenTypeEOF:
		return nil
	case tokenTypeSemicolon:
		return p.parseEndOfQueryStatement()
	case tokenTypeNewLine:
		return p.parseNewLineStatement()
	default:
		p.nextToken()
		goto _switch
	}
}

func (p *parser) parseNameStatement() nodeName {
	if p.peekTok.typ == tokenTypeRawInput {
		node := nodeName{
			valid: true,
			val:   strings.TrimSpace(p.peekTok.literal),
		}
		p.nextToken()
		p.nextToken()
		return node
	}
	return nodeName{}
}

func (p *parser) parseCommentStatement() nodeComment {
	var b strings.Builder

	p.nextToken()

	for {
		if p.curTok.typ == tokenTypeNewLine || p.curTok.typ == tokenTypeEOF {
			break
		}
		b.Grow(len(p.curTok.literal))
		b.WriteString(p.curTok.literal)
		p.nextToken()
	}

	return nodeComment{
		valid: true,
		val:   b.String(),
	}
}

func (p *parser) parseStringLiteralStatement() nodeStringLiteral {
	var b strings.Builder

	p.nextToken()

	for {
		if p.curTok.typ == tokenTypeStringLiteral {
			p.nextToken()
			break
		}

		if p.curTok.typ == tokenTypeEOF {
			break
		}

		b.Grow(len(p.curTok.literal))
		b.WriteString(p.curTok.literal)
		p.nextToken()
	}

	return nodeStringLiteral{
		valid: true,
		val:   b.String(),
	}
}

func (p *parser) parseIdentifierStatement() nodeIdentifier {
	seeking := p.curTok.literal
	tok := p.curTok

	var b strings.Builder

	p.nextToken()

	for {
		if p.curTok.typ == tokenTypeIdentifier && p.curTok.literal == seeking {
			p.nextToken()
			break
		}

		if p.curTok.typ == tokenTypeEOF {
			break
		}

		b.Grow(len(p.curTok.literal))
		b.WriteString(p.curTok.literal)
		p.nextToken()
	}

	return nodeIdentifier{
		valid: true,
		val:   b.String(),
		tok:   tok,
	}
}

func (p *parser) parseQueryStatement() nodeQuery {
	var b strings.Builder

	for {
		if p.curTok.typ == tokenTypeEOF {
			break
		}

		if p.curTok.typ != tokenTypeRawInput {
			break
		}

		b.Grow(len(p.curTok.literal))
		b.WriteString(p.curTok.literal)
		p.nextToken()
	}

	return nodeQuery{
		valid: true,
		val:   b.String(),
	}
}

func (p *parser) parseEndOfQueryStatement() nodeEnfOfQuery {
	p.nextToken()
	return nodeEnfOfQuery{}
}

func (p *parser) parseNewLineStatement() nodeNewLine {
	p.nextToken()
	return nodeNewLine{}
}
