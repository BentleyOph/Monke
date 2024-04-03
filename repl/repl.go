package repl

import (
	"bufio"
	"fmt"
	"github.com/BentleyOph/monke/lexer"
	"github.com/BentleyOph/monke/token"
	"io"
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
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}

}
