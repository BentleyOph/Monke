package parser

import (
	"fmt"
	"strconv"

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


var precedences = map[token.TokenType]int{ //map of precedences for each token type
	token.EQ : EQUALS,
	token.NOT_EQ : EQUALS,
	token.LT : LESSGREATER,
	token.GT : LESSGREATER,
	token.PLUS: SUM,
	token.MINUS: SUM,
	token.SLASH: PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN: CALL,
}



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
	p.registerPrefix(token.INT,p.parseIntegerLiteral)
	p.registerPrefix(token.BANG,p.parsePrefixExpression)
	p.registerPrefix(token.MINUS,p.parsePrefixExpression)
	p.registerPrefix(token.TRUE,p.parseBoolean)
	p.registerPrefix(token.FALSE,p.parseBoolean)
	p.registerPrefix(token.LPAREN,p.parseGroupedExpression)
	p.registerPrefix(token.IF,p.parseIfExpression)
	p.registerPrefix(token.FUNCTION,p.parseFunctionLiteral)
	p.registerPrefix(token.STRING,p.parseStringLiteral)
	p.infixParseFns = make (map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS,p.parseInfixExpression)
	p.registerInfix(token.MINUS,p.parseInfixExpression)
	p.registerInfix(token.SLASH,p.parseInfixExpression)
	p.registerInfix(token.ASTERISK,p.parseInfixExpression)
	p.registerInfix(token.EQ,p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ,p.parseInfixExpression)
	p.registerInfix(token.LT,p.parseInfixExpression)
	p.registerInfix(token.GT,p.parseInfixExpression)
	p.registerInfix(token.LPAREN,p.parseCallExpression)
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
			program.Statements = append(program.Statements, stmt)
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

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON){
		p.nextToken()
	
	}
	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// check if the next token is of a certain type
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
	stmt.ReturnValue = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	return stmt
}

//register a prefix parse function for a token type
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn){ 
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn){
	p.infixParseFns[tokenType] = fn
}



func(p *Parser) parseExpressionStatement() *ast.ExpressionStatement{
	// defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	return stmt
}
// parseExpression parses an expression based on the given precedence.
// It uses prefix and infix parse functions to build the expression tree.
// If a prefix parse function is not found for the current token type, it returns an error.
// It continues parsing infix expressions as long as the precedence of the next token is higher.
// It returns the resulting expression.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// defer untrace(trace("parseExpression"))
	prefix := p.prefixParseFns[p.curToken.Type] // check if a prefix parse function exists for the current token type
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {

		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}


func (p *Parser) parsePrefixExpression() ast.Expression{
	// defer untrace(trace("parsePrefixExpression"))
	expression := &ast.PrefixExpression{
		Token: p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}





func (p *Parser) noPrefixParseFnError(t token.TokenType){ //error message for when no prefix parse function is found
	msg := fmt.Sprintf("no prefix parse function for %s found",t)
	p.errors = append(p.errors, msg)
}




func (p *Parser) parseIntegerLiteral() ast.Expression {
	// defer untrace(trace("parseIntegerLiteral"))
	lit := &ast.IntegerLiteral{Token: p.curToken}
	
	value, err := strconv.ParseInt(p.curToken.Literal,0,64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer",p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil 
	}
	lit.Value = value
	return lit 

}


func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p 
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type];ok{
		return p
	}
	return LOWEST
}



func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression{
	// defer untrace(trace("parseInfixExpression"))
	expression := &ast.InfixExpression{
		Token: p.curToken,
		Operator: p.curToken.Literal,
		Left: left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression 
}


func (p *Parser) parseBoolean() ast.Expression{
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression{
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN){
		return nil
	}
	return exp
}


func (p *Parser) parseIfExpression() ast.Expression{
	expression := &ast.IfExpression{Token:p.curToken}

	if !p.expectPeek(token.LPAREN){
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN){
		return nil
	}
	if !p.expectPeek(token.LBRACE){
		return nil
	}
	expression.Consequence = p.parseBlockStatement()
	if p.peekTokenIs(token.ELSE){
		p.nextToken()
		if !p.expectPeek(token.LBRACE){
			return nil 
		}
		expression.Alternative = p.parseBlockStatement()

	}
	return expression

}


func (p *Parser) parseBlockStatement() *ast.BlockStatement{
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}
	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF){
		stmt := p.parseStatement()
			block.Statements = append(block.Statements, stmt)
		p.nextToken()
	}
	return block
}


func (p *Parser) parseFunctionLiteral() ast.Expression{
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(token.LPAREN){
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE){
		return nil
	}
	lit.Body = p.parseBlockStatement()

	return lit
}

func(p *Parser) parseFunctionParameters() []*ast.Identifier{
	identifiers:= []*ast.Identifier{}
	if p.peekTokenIs(token.RPAREN){
		p.nextToken()
		return identifiers
	}
	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA){
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}
	if !p.expectPeek(token.RPAREN){
		return nil
	}
	return identifiers
}

func(p *Parser) parseCallExpression(function ast.Expression) ast.Expression{ // receives the already parsed function literal and uses it to create a call expression
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}


func (p *Parser) parseCallArguments() []ast.Expression{
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN){
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA){
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN){
		return nil
	}

	return args 
}

func (p *Parser) parseStringLiteral() ast.Expression{
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}