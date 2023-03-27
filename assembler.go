package main

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Assemble assembles src.
func Assemble(src io.Reader) ([]Word, error) {
	raw, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(raw), "\n")

	// symtab is mapping identifier to address of identifier label.
	symtab := make(map[string]Word)

	// first pass; fill symtab.
	var addr Word
	for i, line := range lines {
		lineNo := i + 1
		tokens, err := tokenize(line)
		if err != nil {
			return nil, SyntaxError{lineNo, line, err}
		}
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

	// second pass; write to out.
	var out []Word
	for i, line := range lines {
		lineNo := i + 1
		tokens, _ := tokenize(line) // error already checked in first pass
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
				return nil, SyntaxError{lineNo, line, nil}
			}
			out = append(out, Word(opcode[instruction]<<12))
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
				return nil, SyntaxError{lineNo, line, nil}
			}
			out = append(out, Word(opcode[instruction]<<12))
			out[len(out)-1] |= symtab[identifier] & 0xFFF
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
				return nil, SyntaxError{lineNo, line, nil}
			}
			out = append(out, Word(opcode[instruction]<<12))
			n, err := strconv.ParseInt(number, 16, 0)
			if err != nil || n < minWordInt || n > maxWordInt {
				return nil, SyntaxError{lineNo, line, err}
			}
			out[len(out)-1] |= Word(n & 0xFFF)
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
			n, err := strconv.ParseInt(number, base, 0)
			if err != nil || n < minWordInt || n > maxWordInt {
				return nil, SyntaxError{lineNo, line, err}
			}
			out = append(out, Word(n))
		default:
			return nil, SyntaxError{lineNo, line, nil}
		}
	}
	return out, nil
}

type SyntaxError struct {
	lineNo int
	line   string
	error
}

func (s SyntaxError) Error() string {
	if s.error == nil {
		return fmt.Sprintf("syntax: line %d: %s", s.lineNo, s.line)
	}
	return fmt.Sprintf("syntax: line %d: %s: %s", s.lineNo, s.error, s.line)
}

// Token is the smallest sub-string unit of the src.
type Token struct {
	typ TokenType
	str string
}

// TokenType is a function that returns true if the string is a TokenType. It is used to classify Token.
type TokenType func(string) bool

// TokenInstruction is a TokenType for instructions. eg., "Load" or "Add".
func TokenInstruction(s string) bool {
	_, ok := opcode[s]
	return ok
}

// TokenDirective is a TokenType for directives. eg., "DEC" or "HEX".
func TokenDirective(s string) bool {
	return regexp.MustCompile(`^(DEC|HEX)$`).FindStringIndex(s) != nil
}

// TokenNumber is a TokenType for numbers. eg., "15" or "0xF".
func TokenNumber(s string) bool {
	return regexp.MustCompile(`^[-+]?[0-9A-Fa-f]+$`).FindStringIndex(s) != nil
}

// TokenIdentifier is a TokenType for identifiers. eg., "var" or "x1".
func TokenIdentifier(s string) bool {
	return regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*$`).FindStringIndex(s) != nil
}

// TokenComma is a TokenType for commas. eg., ",".
func TokenComma(s string) bool {
	return s == ","
}

func tokenize(line string) ([]Token, error) {
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
			return nil, fmt.Errorf("bad token: %q", s)
		}
	}
	return out, nil
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
