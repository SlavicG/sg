package repl

import (
	"bufio"
	"fmt"
	"io"
	"sg_interpreter/src/sg/Item"
	"sg_interpreter/src/sg/evaluator"
	"sg_interpreter/src/sg/lexer"
	"sg_interpreter/src/sg/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := Item.NewScope()
	fmt.Printf(PROMPT)

	line := ""
	for {
		scanned := scanner.Scan()
		line += scanner.Text()
		if !scanned {
			l := lexer.New(line)
			p := parser.New(l)
			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil {
				//check stuff
			}
			return
		}

	}
}

const ERROR_MESSAGE = `
VLADISLAV FOUND BUG
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, ERROR_MESSAGE)
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
