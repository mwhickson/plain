package main

import (
	"bufio"
	"fmt"
	"os"
	"plain/pkg/compiler"
	"plain/pkg/parser"
	"plain/pkg/vm"
	"strconv"
	"strings"
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
	fmt.Println("Type 'help' for a list of commands, or 'exit' to quit.")
	fmt.Println("------------------------------------------------------------")

	var lines []string
	scanner := bufio.NewScanner(os.Stdin)

	executeSession := func() bool {
		source := strings.Join(lines, "\n")
		if strings.TrimSpace(source) == "" {
			return true
		}

		l := parser.NewLexer(source)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			for _, err := range p.Errors() {
				fmt.Println("Parser error:", err)
			}
			return false
		}

		comp := compiler.New()
		if err := comp.Compile(program); err != nil {
			fmt.Println("Compiler error:", err)
			return false
		}

		chunk := comp.Bytecode()
		machine := vm.NewVM(chunk)
		if err := machine.Run(); err != nil {
			fmt.Println("VM runtime error:", err)
			return false
		}

		if val, err := machine.Pop(); err == nil {
			fmt.Println(val.String())
		}
		return true
	}

	for {
		fmt.Print("plain> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		trimmed := strings.TrimSpace(input)

		if trimmed == "exit" || trimmed == "quit" {
			break
		}
		if trimmed == "" {
			continue
		}

		parts := strings.Fields(trimmed)
		cmd := strings.ToLower(parts[0])

		switch cmd {
		case "help":
			fmt.Println("Available REPL Commands:")
			fmt.Println("  help           - Show this help menu")
			fmt.Println("  new            - Clear the current session")
			fmt.Println("  list           - Show code lines in memory")
			fmt.Println("  run            - Execute code currently in memory")
			fmt.Println("  save <file>    - Save session lines to a file")
			fmt.Println("  load <file>    - Load session lines from a file")
			fmt.Println("  edit <line>    - Edit a specific line number (blank to delete)")
			fmt.Println("  exit / quit    - Exit the REPL")

		case "new":
			lines = []string{}
			fmt.Println("Session cleared.")

		case "list":
			if len(lines) == 0 {
				fmt.Println("No code in session.")
			} else {
				for i, l := range lines {
					fmt.Printf("%3d: %s\n", i+1, l)
				}
			}

		case "run":
			fmt.Println("Running session...")
			executeSession()

		case "save":
			if len(parts) < 2 {
				fmt.Println("Usage: save <filename>")
				continue
			}
			filename := parts[1]
			data := strings.Join(lines, "\n") + "\n"
			err := os.WriteFile(filename, []byte(data), 0644)
			if err != nil {
				fmt.Println("Error saving file:", err)
			} else {
				fmt.Printf("Saved %d lines to %s\n", len(lines), filename)
			}

		case "load":
			if len(parts) < 2 {
				fmt.Println("Usage: load <filename>")
				continue
			}
			filename := parts[1]
			data, err := os.ReadFile(filename)
			if err != nil {
				fmt.Println("Error reading file:", err)
				continue
			}
			fileLines := strings.Split(string(data), "\n")
			lines = []string{}
			for _, fl := range fileLines {
				if trimmedLine := strings.TrimSpace(fl); trimmedLine != "" {
					lines = append(lines, trimmedLine)
				}
			}
			fmt.Printf("Loaded %d lines from %s. Running session...\n", len(lines), filename)
			executeSession()

		case "edit":
			if len(parts) < 2 {
				fmt.Println("Usage: edit <line number>")
				continue
			}
			lineNum, err := strconv.Atoi(parts[1])
			if err != nil || lineNum < 1 || lineNum > len(lines) {
				fmt.Printf("Invalid line number. Current session has %d lines.\n", len(lines))
				continue
			}

			idx := lineNum - 1
			fmt.Printf("Current [%d]: %s\n", lineNum, lines[idx])
			fmt.Printf("edit[%d]> ", lineNum)
			if scanner.Scan() {
				newLine := scanner.Text()
				if strings.TrimSpace(newLine) == "" {
					lines = append(lines[:idx], lines[idx+1:]...)
					fmt.Printf("Deleted line %d.\n", lineNum)
				} else {
					lines[idx] = newLine
					fmt.Printf("Updated line %d.\n", lineNum)
				}
				executeSession()
			}

		default:
			// Check if it's a Direct Mode bang-command (e.g. !10 + 10)
			if strings.HasPrefix(input, "!") {
				directSource := strings.TrimSpace(input[1:])
				if directSource == "" {
					continue
				}

				l := parser.NewLexer(directSource)
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

				if val, err := machine.Pop(); err == nil {
					fmt.Println(val.String())
				}
			} else {
				// Program Mode: perform a silent syntax check before buffering
				l := parser.NewLexer(input)
				p := parser.NewParser(l)
				_ = p.ParseProgram()
				if len(p.Errors()) > 0 {
					for _, err := range p.Errors() {
						fmt.Println("Parser error:", err)
					}
					continue // Reject invalid lines from buffer
				}

				lines = append(lines, input)
				fmt.Printf("[%d lines in buffer]\n", len(lines))
			}
		}
	}
}
