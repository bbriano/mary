// Maria is a simulation of the Marie machine described in chapter 4 of
// "Computer Organization and Architecture" by Linda Null and Julia Lobur.
package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: maria file")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer f.Close()

	m := &Machine{}
	m.Load(f)
	m.Run()
}

// Machine simulates a Marie machine. Most of the registers are not needed for the simulation,
// but they are added to illustrate the Marie machine described in the book.
type Machine struct {
	// Accumulator
	AC Word
	// Program Counter
	PC Word
	// Memory Access Register
	MAR Word
	// Program Buffer Register
	MBR Word
	// Instruction Register
	IR Word
	// Input Register
	IN Word
	// Output Register
	OUT Word
	// Memory. Used to store both program and data.
	M [1 << 12]Word // 12-bit addressing
}

// Word is the width of Marie's data bus.
type Word uint16

var symbol = regexp.MustCompile("^[A-Za-z][A-Za-z0-9]*$")
var directive = regexp.MustCompile("^(DEC|HEX)$")
var number = regexp.MustCompile("^[-+]?[0-9A-Fa-f]+$")
var white = regexp.MustCompile("[ \t\n]+")

// Load loads a Marie assembly program and assembles it to the machine's memory.
func (m *Machine) Load(r io.Reader) {
	b, err := io.ReadAll(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	lines := strings.Split(string(b), "\n")
	symtab := make(map[string]Word)

	// first pass. fill symbol table
	for _, line := range lines {
		tokens := tokenize(line)
		switch len(tokens) {
		case 0:
			continue
		case 4:
		default:
			m.PC++
			continue
		}
		sym := symbol.FindString(tokens[0])
		dir := directive.FindString(tokens[2])
		num := number.FindString(tokens[3])
		if tokens[1] != "," || sym == "" || dir == "" || num == "" {
			fmt.Fprintln(os.Stderr, "syntax error:", line)
			os.Exit(1)
		}
		symtab[sym] = m.PC
		m.PC++
	}
	m.PC = 0

	// second pass. fill m.M
	for _, line := range lines {
		tokens := tokenize(line)
		switch len(tokens) {
		case 0: // empty (or comment) lines
		case 1:
			switch opcode[tokens[0]] {
			case OpInput:
			case OpOutput:
			case OpHalt:
			case OpSkipcond:
			case OpClear:
			default:
				fmt.Fprintln(os.Stderr, "syntax error:", line)
				os.Exit(1)
			}
			m.M[m.PC] = Word(opcode[tokens[0]] << 12)
			m.PC++
		case 2:
			switch opcode[tokens[0]] {
			case OpJnS:
			case OpLoad:
			case OpStore:
			case OpAdd:
			case OpSubt:
			case OpJump:
			case OpAddI:
			case OpJumpI:
			default:
				fmt.Fprintln(os.Stderr, "syntax error:", line)
				os.Exit(1)
			}
			operand, ok := symtab[tokens[1]]
			if !ok {
				n, err := strconv.Atoi(tokens[1])
				if err != nil {
					fmt.Fprintln(os.Stderr, "syntax error:", line)
					os.Exit(1)
				}
				operand = Word(n)
			}
			m.M[m.PC] = Word(opcode[tokens[0]] << 12)
			m.M[m.PC] |= operand & 0xFFF
			m.PC++
		case 4:
			if tokens[1] != "," {
				fmt.Fprintln(os.Stderr, "syntax error:", line)
				os.Exit(1)
			}
			n, err := parseint(tokens[2], tokens[3])
			if err != nil {
				fmt.Fprintln(os.Stderr, "num too big:", line)
				os.Exit(1)
			}
			m.M[m.PC] = n
			m.PC++
		default:
			fmt.Fprintln(os.Stderr, "syntax error:", line)
			os.Exit(1)
		}
	}
	m.PC = 0
}

func tokenize(line string) []string {
	line = strings.Split(line, "/")[0]
	line = strings.ReplaceAll(line, ",", " , ")
	line = white.ReplaceAllString(line, " ")
	line = strings.Trim(line, " ")
	var out []string
	for _, s := range strings.Split(line, " ") {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

func parseint(directive, num string) (Word, error) {
	var base int
	switch directive {
	case "HEX":
		base = 16
	case "DEC":
		base = 10
	default:
		panic("unreachable")
	}
	n, err := strconv.ParseInt(num, base, 17)
	if err != nil {
		return 0, err
	}
	return Word(n), nil
}

// Run starts execution of the program stored in the machine's memory.
func (m *Machine) Run() {
	for {
		m.MAR = m.PC
		m.MBR = m.M[m.PC]
		m.IR = m.MBR
		m.PC++
		opcode := Opcode(m.IR >> 12)
		operand := m.IR & 0xFFF
		instruction[opcode](m, operand)
	}
}
