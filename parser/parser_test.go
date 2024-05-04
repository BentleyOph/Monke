package parser // Package declaration for the parser package

import (
	"fmt"
	"testing" // Importing the testing package

	"github.com/BentleyOph/monke/ast"   // Importing the ast package
	"github.com/BentleyOph/monke/lexer" // Importing the lexer package
)

func TestLetStatements(t *testing.T){ // Defining a test function named TestLetStatements
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`
	l := lexer.New(input) // Creating a new lexer instance with the input string

	p := New(l)
	program := p.ParseProgram() // Parsing the program using the parser instance
	checkParserErrors(t,p) // Checking for parser errors

	if program == nil { // Checking if the parsed program is nil
		t.Fatalf("ParseProgram() returned nil") // Failing the test with an error message
	}
	if len(program.Statements) != 3 { // Checking if the program contains 3 statements
		t.Fatalf("program.Statements does not contain 3 statements. got= %d",len(program.Statements)) // Failing the test with an error message
	}
	tests := []struct{
		expectedIdentifier string // Declaring a struct field named expectedIdentifier of type string
	}{
		{"x"}, // Initializing a struct with expectedIdentifier set to "x"
		{"y"}, // Initializing a struct with expectedIdentifier set to "y"
		{"foobar"}, // Initializing a struct with expectedIdentifier set to "foobar"
	}
	for i, tt := range tests { // Looping over the tests slice
		stmt := program.Statements[i] // Getting the statement at index i
		if !testLetStatement(t,stmt,tt.expectedIdentifier){ // Calling the testLetStatement function with the statement and expected identifier
			return // Returning from the function
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool { // Defining a helper function named testLetStatement
	if s.TokenLiteral() != "let"{ // Checking if the statement's token literal is "let"
		t.Errorf("s.TokenLiteral not 'let'. got= %q",s.TokenLiteral()) // Failing the test with an error message
		return false 
	}
	letStmt, ok := s.(*ast.LetStatement) // Type asserting the statement to *ast.LetStatement
	if !ok { // Checking if the type assertion was successful
		t.Errorf("s not *ast.LetStatement. got= %T",s) // Failing the test with an error message
		return false // Returning false
	}
	if letStmt.Name.Value != name { // Checking if the let statement's name value matches the expected identifier
		t.Errorf("letStmt.Name.Value not '%s' got= %s",name,letStmt.Name.Value) // Failing the test with an error message
		return false // Returning false
	}
	if letStmt.Name.TokenLiteral() != name { 
		t.Errorf("s.Name not '%s' . got = %s",name,letStmt.Name) 
		return false // Returning false
	}
	return true // Returning true
}

func checkParserErrors(t *testing.T,p *Parser){
	errors := p.Errors() // Getting the parser's errors

	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors",len(errors))
	for _,msg :=range errors{
		t.Errorf("parser error: %q",msg)
	}
	t.FailNow()
}


func TestReturnStatements(t *testing.T){
	input := `
	return 5;
	return 10;
	return 993322;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t,p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got= %d",len(program.Statements))
	}
	for _,stmt := range program.Statements{
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got = %T",stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return"{
			t.Errorf("returnStmt.TokenLiteral not 'return' got %q ",returnStmt.TokenLiteral())
		}
	}

}


func TestIdentifierExpression(t *testing.T){
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t,p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got= %d",len(program.Statements)) 
	}
	_, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got= %T",program.Statements[0])
	}
	ident , ok := program.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.Identifier) // Type asserting the expression to *ast.Identifier
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got= %T",program.Statements[0])
	}
	if ident.Value != "foobar"{
		t.Errorf("ident.Value not %s. got = %s", "foobar",ident.Value)

	}
	if ident.TokenLiteral() != "foobar"{
		t.Errorf("ident.TokenLiteral not %s. got = %s","foobar",ident.TokenLiteral())
	}


}


func TestIntegerLiteralExpression(t *testing.T){
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t,p)

	if len(program.Statements) != 1{
		t.Fatalf("program has not enough statements. got = %d ",len(program.Statements))

	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement) // Type asserting the statement to *ast.ExpressionStatement
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got= %T",program.Statements[0])
	}

	literal,ok := stmt.Expression.(*ast.IntegerLiteral) // Type asserting the expression to *ast.IntegerLiteral
	if !ok {
		t.Fatalf("exp not *ast.IntegerLIteral. got = %T",stmt.Expression)
	}
	if literal.Value != 5{
		t.Errorf("literal.Value not %d. got = %d",5,literal.Value)
	}
	if literal.TokenLiteral()!= "5" {
		t.Errorf("literal.TokenLiteral not %d. got = %s",5,literal.TokenLiteral())
	}
}


func TestParsingPrefixExpressions(t *testing.T){
	prefixTests := []struct {
		input string
		operator string
		value interface {}
	}{
		{"!5","!",5},
		{"-15","-",15},
		{"!true;","!",true},
		{"!false","!",false},

	}

	for _,tt := range prefixTests{
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t,p)

		if len(program.Statements) != 1{
			t.Fatalf("program has not enough statements. got = %d",len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got = %T",program.Statements[0])
		}
		exp ,ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got = %T",stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got = %s",tt.operator,exp.Operator)
		}
		if !testLiteralExpression (t, exp.Right, tt.value){
			return 
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool{
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got = %T",il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got = %d",value,integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d",value){
		t.Errorf("integ.TokenLiteral not %d. got = %s",value,integ.TokenLiteral())
		return false
	}
	return true
}



func TestParsingInfixExpressions(t *testing.T){
	infixTests := []struct {
		input string
		leftValue interface {}
		operator string
		rightValue interface{}
	}{
		{"5+5",5,"+",5},
		{"5-5",5,"-",5},
		{"5*5",5,"*",5},
		{"5/5",5,"/",5},
		{"5>5",5,">",5},
		{"5<5",5,"<",5},
		{"5==5",5,"==",5},
		{"5!=5",5,"!=",5},
		{"true == true",true,"==",true},
		{"true != false",true,"!=",false},
		{"false == false",false,"==",false},
	}

	for _,tt := range infixTests{
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t,p)

		if len(program.Statements) != 1 {
			t.Fatalf("program does not contain %d statements. got = %d",1,len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement.got = %T",program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got = %T",stmt.Expression)
		}
		if !testLiteralExpression(t,exp.Left,tt.leftValue){
			return
		}
		if exp.Operator != tt.operator{
			t.Fatalf("exp.Operator is not '%s'. got = %s",tt.operator,exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right,tt.rightValue){
			return 
		}
	}
}



func TestOperatorPrecedenceParsing(t *testing.T){
	tests := []struct {
		input string
		expected string
	}{
		// {"-a * b","((-a) * b)"},
		// {"!-a","(!(-a))"},
		// {"a+b+c","((a + b) + c)"},
		// {"a+b-c","((a + b) - c)"},
		// {"a*b*c","((a * b) * c)"},
		// {"a*b/c","((a * b) / c)"},
		// {"a/b*c","((a / b) * c)"},
		// {"a+b/c","(a + (b / c))"},
		// {"a+b*c+d","((a + (b * c)) + d)"},
		// {"3+4; -5*5","(3 + 4)((-5) * 5)"},
		// {"5>4==3<4","((5 > 4) == (3 < 4))"},
		// {"5<4!=3>4","((5 < 4) != (3 > 4))"},
		// {"3+4*5==3*1+4*5","((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		// {"true","true"},
		// {"false","false"},
		// {"3>5 == false","((3 > 5) == false)"},
		// {"3<5 == true","((3 < 5) == true)"},
		// {"1 + (2 + 3) + 4","((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2","((5 + 5) * 2)"},
		// {"-(5 + 5)","(-(5 + 5))"},i
	}

	for _,tt := range tests{
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t,p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected = %q, got = %q",tt.expected,actual)
		}
	}
}


func testIdentifier(t *testing.T,exp ast.Expression,value string) bool {
	ident ,ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got = %T",exp)
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got = %s",value,ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got = %s",value,ident.TokenLiteral())
		return false
	}
	return true
}

// testLiteralExpression is a helper function used to test literal expressions in the parser.
// It takes a testing.T object, an ast.Expression, and an expected value as input.
// The function performs type assertions based on the expected value and calls the appropriate test function.
// It returns a boolean value indicating whether the test passed or failed.
func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type)  {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t ,exp ,v)
	}
	t.Errorf("type of exp not handled. got = %T", exp)
	return false
}


func  testInfixExpression(t *testing.T,exp ast.Expression,left interface{},operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got = %T",exp)
		return false
	}
	if !testLiteralExpression(t,opExp.Left,left){
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got = %s",operator,opExp.Operator)
		return false
	}
	if !testLiteralExpression(t,opExp.Right,right){
		return false
	}

	return true
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}
		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not %t. got=%t", tt.expectedBoolean,
				boolean.Value)
		}
	}
}


func testBooleanLiteral(t *testing.T,exp ast.Expression, value bool) bool {
	boolean , ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got = %T",exp)
		return false
	}	
	if boolean.Value != value {
		t.Errorf("boolean.Value not %t . got = %t",value,boolean.Value)
		return false
	}
	if boolean.TokenLiteral() != fmt.Sprintf("%t",value){
		t.Errorf("boolean.TokenLiteral not %t. got = %s",value,boolean.TokenLiteral())
		return false
	}
	return true
}