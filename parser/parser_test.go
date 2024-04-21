package parser // Package declaration for the parser package

import (
	"testing" // Importing the testing package

	"github.com/BentleyOph/monke/ast" // Importing the ast package
	"github.com/BentleyOph/monke/lexer" // Importing the lexer package
	
)

func TestLetStatements(t *testing.T){ // Defining a test function named TestLetStatements
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`
	l := lexer.New(input) // Creating a new lexer instance with the input string

	P := New(l)
	program := P.ParseProgram() // Parsing the program using the parser instance

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
