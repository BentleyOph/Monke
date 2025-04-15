# Monke Programming Language

Monke is a simple, interpreted programming language written in Go. It is designed as an educational project to demonstrate the principles of lexing, parsing, and interpreting code. Inspired by similar projects such as [Writing An Interpreter In Go](https://interpreterbook.com/), Monke offers a lightweight environment to experiment with language concepts while learning Go.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

Monke is a toy programming language that implements core concepts of programming language design, including:

- **Lexical Analysis:** Tokenizes input strings into meaningful symbols.
- **Parsing:** Converts tokens into an abstract syntax tree (AST) to represent the code’s structure.
- **Evaluation:** Although not fully featured in the current code base, it lays the groundwork for interpreting expressions and statements.
- **REPL (Read–Eval–Print Loop):** Provides an interactive command-line interface for users to input commands and view generated parsed output.

Monke supports common language constructs, including:

- **Let Statements:** Declare variables.
- **Return Statements:** Return values from functions.
- **Expressions:** Integer, Boolean, and String literals, as well as infix and prefix expressions.
- **Conditionals:** `if` and `if-else` expressions.
- **Functions:** Function literals and call expressions.

The project is an excellent resource for learning about compiler and interpreter design while enjoying a playful, monkey-themed environment.

---

## Features

- **Interactive REPL:** An interactive prompt (>>), where users can type commands to see parsed output.
- **Lexer:** Breaks code into tokens such as identifiers, literals, and operators.
- **Parser:** Builds an abstract syntax tree (AST) from tokens while handling operator precedence.
- **AST Nodes:** Detailed implementation of language constructs like expressions, statements, functions, and conditionals.
- **Error Reporting:** Provides basic error reporting for syntax errors during parsing.
- **Extensible Design:** The code base is organized into separate packages (lexer, parser, AST, token, repl) for easy understanding and future expansion.

---

## Installation

Monke requires **Go 1.22.0** or later.

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/BentleyOph/monke.git
   cd monke
   ```

2. **Download Dependencies:**

   Monke uses Go modules, so use the following command to ensure all dependencies are available:

   ```bash
   go mod tidy
   ```

3. **Build or Run:**

   To run the REPL directly, execute:

   ```bash
   go run main.go
   ```

   Alternatively, build the application with:

   ```bash
   go build -o monke
   ./monke
   ```

---

## Usage

### Running the REPL

Once you run the application, you will see a prompt (>>). Here you can type Monkey-language commands. For example:

```monke
>> let x = 5;
>> return x;
>> fn(a, b) { a + b; }
```

The REPL reads the user input, passes it through the lexer and parser, and then displays the parsed output. If there are syntax errors or issues, the REPL will output error messages to help with debugging.

### Example Interaction

```plaintext
Hello <YourUsername>! This is the Monke programming language!
Type in any commands to see generated parsed code
>> let myVar = anotherVar;
let myVar = anotherVar;
>> if (x < y) { x } else { y }
if(x < y) { x }else{ y }
```

---

## Project Structure

The project is organized into several key packages:

- **`main.go`**  
  The application entry point. It greets the user using the current system username and starts the REPL.

- **`ast/ast.go`**  
  Contains the definitions of AST nodes including program, statements, and expressions. It also provides methods for converting nodes back into string representations.

- **`lexer/lexer.go`**  
  Implements the lexical analyzer (lexer) that reads the source code and converts it into tokens such as keywords, identifiers, literals, and operators.

- **`parser/parser.go`**  
  Contains the logic to parse tokens into an AST, handling operator precedence, prefix and infix expressions, function literals, and conditional expressions.

- **`repl/repl.go`**  
  Defines the interactive Read–Eval–Print Loop (REPL) for user input and output. It provides helpful prompts and error handling during code input.

- **`token/token.go`**  
  Defines the token types and includes helper functions to look up keywords.

- **`tests`**  
  Throughout the project, there are various test files (e.g., `ast_test.go`, `lexer_test.go`, `parser_test.go`) that verify the correctness of each component.

---

## Testing

Unit tests are provided to ensure that various components of Monke work as expected. To run all tests, execute the following command in the project root:

```bash
go test ./...
```

The tests cover:

- Lexical analysis for various tokens and literals.
- Parsing of different expressions, including let statements, return statements, prefix and infix expressions.
- Construction and stringification of the abstract syntax tree (AST).

---
