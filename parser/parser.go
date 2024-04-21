package parser

import (
	"github.com/BentleyOph/monke/ast"
	"github.com/BentleyOph/monke/lexer"
	"github.com/BentleyOph/monke/token"
)

// constructing our AST that begins with *ast.Program at the top
// Parser is a struct that holds the lexer and the current token
type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token // points to the current token
	peekToken token.Token // points to the next token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	//Read two tokens so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
	}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()

	}
	return program

}

func (p *Parser)parseStatement() ast.Statement{
	switch p.curToken.Type{
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser)parseLetStatement() *ast.LetStatement{
	stmt := &ast.LetStatement{Token:p.curToken}
	if !p.expectPeek(token.IDENT){
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN){
		return nil 
	}
	for !p.curTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool{
	return p.curToken.Type == t 
}

func (p * Parser)peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}
func(p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t){
		p.nextToken()
		return true
	} else {
		return false
	}

}


