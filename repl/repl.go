package repl

import (
	"bufio"
	"fmt"
	"github.com/BentleyOph/monke/parser"
	"io"

	"github.com/BentleyOph/monke/lexer"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in) // initialize scanner to read from input
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan() // scan the input
		if !scanned {
			return
		}
		line := scanner.Text() // get the input
		l := lexer.New(line)   // initialize lexer
		p := parser.New(l)     // initialize parser

		program := p.ParseProgram() // parse the program
		if len(p.Errors()) != 0 {   // check for errors
			printParserErrors(out, p.Errors())
			continue
		}
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")

	}

}

// const MONKEY_FACE = `
// 			__,__
//    .--.  .-"     "-.  .--.
//   / .. \/  .-. .-.  \/ .. \
//  | |  '|  /   Y   \  |'  | |
//  | \   \  \ 0 | 0 /  /   / |
//   \ '- ,\.-"""""""-./, -' /
//    ''-' /_   ^ ^   _\ '-''
// 	   |  \._   _./  |
// 	   \   \ '~' /   /
// 		'._ '-=-' _.'
// 		   '-----'
// `

func printParserErrors(out io.Writer, errors []string) {
	// io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
