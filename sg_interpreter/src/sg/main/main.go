package main

import (
	"fmt"
	"os"
	"sg_interpreter/src/sg/repl"
)

func main() {
	// Open the input file
	file, err := os.Open(`/Users/slavicg/work/StarGust/sg/sg_interpreter/src/sg/main/input.sg`)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Start the REPL, reading from the file and writing to stdout
	//file or os.Stdin
	repl.Start(file, os.Stdout)
}
