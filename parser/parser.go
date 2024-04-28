package parser

import (
	"fmt"

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
	prefixParseFns map[token.TokenType]prefixParseFn //map of functions that parse prefix expressions
	infixParseFns map[token.TokenType]infixParseFn //map of functions that parse infix expressions

	errors    []string
}
type (
	prefixParseFn func() ast.Expression //parse functions for prefix expressions
	infixParseFn func(ast.Expression) ast.Expression //parse functions for infix expressions
)

const(
	_ int = iota //assigns the zero value to the first constant in the group
	LOWEST //lowest precedence
	EQUALS // ==
	LESSGREATER // > or <
	SUM // +
	PRODUCT // *
	PREFIX // -X or !X
	CALL // myFunction(X)
)



func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
		errors: []string{},
	}
	//Read two tokens so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	p.prefixParseFns = make (map[token.TokenType]prefixParseFn) //initialize the map
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program { //returns the root node of our AST
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken() //repeatedly call nextToken to advance the curToken and peekToken pointers

	}
	return program

}

func (p *Parser) parseIdentifier() ast.Expression{
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.ParseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	// enforce that the next token is an identifier
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}
func (p *Parser) expectPeek(t token.TokenType) bool { //enforce the correctness of the order of tokens by checking the type of the next token
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}

}

func (p *Parser) Errors() [] string {
	return p.errors
}
func (p *Parser) peekError(t token.TokenType){
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",t ,p.peekToken.Type)
	p.errors = append(p.errors, msg)

}

func (p *Parser)ParseReturnStatement() *ast.ReturnStatement{
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	//TODO: we're skipping the expression until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	return stmt
}


func (p *Parser) registerPrefix (tokenType token.TokenType, fn prefixParseFn){
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn){
	p.infixParseFns[tokenType] = fn
}



func(p *Parser) parseExpressionStatement() *ast.ExpressionStatement{
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	return stmt
}
func (p *Parser) parseExpression(precedence int) ast.Expression{
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil  //no prefix parse function found
	}
	leftExp := prefix()
	return leftExp
}