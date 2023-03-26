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

// minWordInt is the minimum integer that can be represented with a Word (-32768).
const minWordInt = -1 << 15

// maxWordInt is the maximum integer that can be represented with a Word (-65535).
const maxWordInt = 1<<16 - 1

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

// Load loads a Marie assembly program and assembles it to the machine's memory.
func (m *Machine) Load(r io.Reader) {
	raw, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(raw), "\n")

	// Construct symbolic table mapping identifier to address of identifier label.
	symtab := make(map[string]Word)
	var addr Word
	for _, line := range lines {
		tokens := tokenize(line)
		switch len(tokens) {
		case 0:
			continue
		case 1:
			addr++
			continue
		}
		switch hashTokens(tokens[:2]) {
		case hashTokenTypes(TokenIdentifier, TokenComma):
			identifier := tokens[0].str
			symtab[identifier] = addr
			addr++
		default:
			addr++
		}
	}

	// second pass. fill m.M
	addr = 0
	for _, line := range lines {
		tokens := tokenize(line)
		if len(tokens) >= 2 {
			switch hashTokens(tokens[:2]) {
			case hashTokenTypes(TokenIdentifier, TokenComma):
				tokens = tokens[2:]
			}
		}
		switch hashTokens(tokens) {
		case hashTokenTypes(): // empty (or comment) lines
		case hashTokenTypes(TokenInstruction):
			instruction := tokens[0].str
			switch opcode[instruction] {
			case OpInput:
			case OpOutput:
			case OpHalt:
			case OpClear:
			default:
				syntaxerr(line)
			}
			m.M[addr] = Word(opcode[instruction] << 12)
			addr++
		case hashTokenTypes(TokenInstruction, TokenIdentifier):
			instruction := tokens[0].str
			identifier := tokens[1].str
			switch opcode[instruction] {
			case OpJnS:
			case OpLoad:
			case OpStore:
			case OpAdd:
			case OpSubt:
			case OpSkipcond:
			case OpJump:
			case OpAddI:
			case OpJumpI:
			case OpLoadI:
			case OpStoreI:
			default:
				syntaxerr(line)
			}
			m.M[addr] = Word(opcode[instruction] << 12)
			m.M[addr] |= symtab[identifier] & 0xFFF
			addr++
		case hashTokenTypes(TokenInstruction, TokenNumber):
			instruction := tokens[0].str
			number := tokens[1].str
			switch opcode[instruction] {
			case OpJnS:
			case OpLoad:
			case OpStore:
			case OpAdd:
			case OpSubt:
			case OpSkipcond:
			case OpJump:
			case OpAddI:
			case OpJumpI:
			case OpLoadI:
			case OpStoreI:
			default:
				syntaxerr(line)
			}
			m.M[addr] = Word(opcode[instruction] << 12)
			n, err := strconv.Atoi(number)
			fmt.Println(n, minWordInt, maxWordInt)
			if err != nil || n < minWordInt || n > maxWordInt {
				syntaxerr(line)
			}
			m.M[addr] |= Word(n & 0xFFF)
			addr++
		case hashTokenTypes(TokenDirective, TokenNumber):
			directive := tokens[0].str
			number := tokens[1].str
			var base int
			switch directive {
			case "HEX":
				base = 16
			case "DEC":
				base = 10
			default:
				panic("unreachable")
			}
			n, err := strconv.ParseInt(number, base, 17)
			if err != nil || n < minWordInt || n > maxWordInt {
				syntaxerr(line)
			}
			m.M[addr] |= Word(n)
			addr++
		default:
			syntaxerr(line)
		}
	}
}

type Token struct {
	typ TokenType
	str string
}

type TokenType func(string) bool

func TokenInstruction(s string) bool {
	_, ok := opcode[s]
	return ok
}

func TokenDirective(s string) bool {
	return regexp.MustCompile(`^(DEC|HEX)$`).FindStringIndex(s) != nil
}

func TokenNumber(s string) bool {
	return regexp.MustCompile(`^[-+]?[0-9A-Fa-f]+$`).FindStringIndex(s) != nil
}

func TokenIdentifier(s string) bool {
	return regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*$`).FindStringIndex(s) != nil
}

func TokenComma(s string) bool {
	return s == ","
}

func tokenize(line string) []Token {
	var out []Token
	line = strings.Split(line, "/")[0]
	line = strings.ReplaceAll(line, ",", " , ")
	line = regexp.MustCompile(`[ \t\n]+`).ReplaceAllString(line, " ")
	line = strings.Trim(line, " ")
	for _, s := range strings.Split(line, " ") {
		if s == "" {
			continue
		}
		switch {
		case TokenInstruction(s):
			out = append(out, Token{TokenInstruction, s})
		case TokenDirective(s):
			out = append(out, Token{TokenDirective, s})
		case TokenNumber(s):
			out = append(out, Token{TokenNumber, s})
		case TokenIdentifier(s):
			out = append(out, Token{TokenIdentifier, s})
		case TokenComma(s):
			out = append(out, Token{TokenComma, s})
		default:
			syntaxerr(line)
		}
	}
	return out
}

func hashTokens(tokens []Token) string {
	var ttypes []TokenType
	for _, t := range tokens {
		ttypes = append(ttypes, t.typ)
	}
	return hashTokenTypes(ttypes...)
}

func hashTokenTypes(ttypes ...TokenType) string {
	return fmt.Sprint(ttypes)
}

func syntaxerr(line string) {
	fmt.Fprintln(os.Stderr, "syntax error:", line)
	panic("syntaxerr")
	os.Exit(1)
}
