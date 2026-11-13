package main

import (
	"bufio"
	"fmt"
	"os"
	"plain/pkg/compiler"
	"plain/pkg/parser"
	"plain/pkg/vm"
)

func main() {
	if len(os.Args) > 1 {
		filename := os.Args[1]
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read file %s: %v\n", filename, err)
			os.Exit(1)
		}
		runSource(string(data))
	} else {
		startREPL()
	}
}

func runSource(source string) {
	l := parser.NewLexer(source)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			fmt.Fprintln(os.Stderr, "Parser error:", err)
		}
		os.Exit(1)
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		fmt.Fprintln(os.Stderr, "Compiler error:", err)
		os.Exit(1)
	}

	chunk := comp.Bytecode()
	machine := vm.NewVM(chunk)
	if err := machine.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "VM runtime error:", err)
		os.Exit(1)
	}
}

func startREPL() {
	fmt.Println("PLAIN Programming Language REPL")
	fmt.Println("Type 'exit' or 'quit' to exit.")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("plain> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if line == "exit" || line == "quit" {
			break
		}
		if line == "" {
			continue
		}

		l := parser.NewLexer(line)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			for _, err := range p.Errors() {
				fmt.Println("Parser error:", err)
			}
			continue
		}

		comp := compiler.New()
		if err := comp.Compile(program); err != nil {
			fmt.Println("Compiler error:", err)
			continue
		}

		chunk := comp.Bytecode()
		machine := vm.NewVM(chunk)
		if err := machine.Run(); err != nil {
			fmt.Println("VM runtime error:", err)
			continue
		}

		// If an expression leaves a value on the stack, pop and print it
		if val, err := machine.Pop(); err == nil {
			fmt.Println(val.String())
		}
	}
}