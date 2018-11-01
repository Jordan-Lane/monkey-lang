package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkeylang/evaluator"
	"monkeylang/lexer"
	"monkeylang/parser"
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
		parser := parser.New(lexer)

		program := parser.ParseProgram()
		if len(parser.Errors()) != 0 {
			printParserErrors(out, parser.Errors())
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, errorMsg := range errors {
		io.WriteString(out, "Oops we got an unexpected Parser Error: \n")
		io.WriteString(out, "\t"+errorMsg+"\n")
	}
}
