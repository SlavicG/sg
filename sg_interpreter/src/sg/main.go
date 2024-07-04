//Interpreter for StarGust SG

package main

import (
	"fmt"
	"os"
	"os/user"
	"sg_interpreter/src/sg/repl"
)

func main() {
	_, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Type in your commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
