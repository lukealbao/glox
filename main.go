/*
Package glox provides a golang port of the lox language.

It has two basic usage patterns: evaluating a lox script and an interactive REPL.

Evaluating a script

Just like any other scripting language runtime.
    glox somescript.lox

Interactive REPL

A very basic REPL can be opened by running the interpreter without any
arguments.

    glox

It currently has only the most basic functionality, which is to say that it
interprets each line of input at a time. It does not have full readline support.
*/
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox <script>")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		RunFile(os.Args[1])
	} else {
		RunPrompt()
	}
}

func RunPrompt() {
	prompt := bufio.NewScanner(os.Stdin)

	fmt.Printf("> ")

	for prompt.Scan() {
		err := prompt.Err()
		check(err)

		Run([]byte(prompt.Text()))
		fmt.Printf("> ")
	}

	check(prompt.Err())
}

func RunFile(filename string) {
	src, err := ioutil.ReadFile(filename)
	check(err)

	Run(src)
}

func Run(src []byte) {
	scanner := NewScanner(string(src))
	scanner.ScanTokens()

	for _, token := range scanner.tokens {
		fmt.Printf("Token: `%s', Type: %s\n", token.Lexeme, token.Type)
	}
}
