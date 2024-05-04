package ast

import (
	"github.com/BentleyOph/monke/token"
	"bytes"
)

type Node interface {
	TokenLiteral() string 
	String () string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface { 
	Node
	expressionNode()
}

// Every valid Monkey program is a series of statements which are stored in the Program.Statements slice.
type Program struct {
	Statements []Statement

}
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string {
	var out bytes.Buffer // bytes.Buffer is a buffer of bytes with a Read and Write method
	for _, s := range p.Statements{ // iterate over each statement in the program
		out.WriteString(s.String()) // write the string representation of the statement to the buffer
	}
	return out.String() // return the buffer as a string
}


type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}



type Identifier struct {
	Token token.Token // token.IDENT
	Value string 
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string {
	return i.Value
}

type ReturnStatement struct{
	Token token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral()string{
	return rs.Token.Literal
}
func (rs *ReturnStatement) String()string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}




type ExpressionStatement struct {
	Token token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode(){}
func(es *ExpressionStatement)TokenLiteral()string{
	return es.Token.Literal
}
func(es *ExpressionStatement)String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}



type IntegerLiteral struct {
	Token token.Token
	Value int64
}
func (il *IntegerLiteral) expressionNode(){}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}



type PrefixExpression struct {
	Token token.Token 
	Operator string
	Right Expression
}

func (pe *PrefixExpression) expressionNode(){}
func (pe *PrefixExpression) TokenLiteral() string{
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()

}




type InfixExpression struct{
	Token token.Token // The operator token, e.g. '+'
	Left Expression
	Operator string
	Right Expression
}
func (oe *InfixExpression) expressionNode(){}
func (oe *InfixExpression)TokenLiteral() string{
	return oe.Token.Literal
}
func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String()) // String() will recursively call the String() method of the left expression
	out.WriteString(" " + oe.Operator + " ") // print the operator and a space
	out.WriteString(oe.Right.String()) // recursively call the String() method of the right expression
	out.WriteString(")")
	return out.String()
}



type Boolean struct {
	Token token.Token
	Value bool 
}
func (b *Boolean) expressionNode(){}
func(b *Boolean) TokenLiteral() string{
	return b.Token.Literal
}
func (b *Boolean) String() string {
	return b.Token.Literal
}