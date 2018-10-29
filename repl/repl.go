package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkeylang/lexer"
	"monkeylang/token"
)

const PROMPT = ">> "

// Start - Start the MonkeyLang REPL
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		lexer := lexer.New(line)

		for tok := lexer.NextToken(); tok.Type != token.EOF; tok = lexer.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
