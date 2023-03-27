package main

import (
	"fmt"
	"os"
)

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
	M [machineMemory]Word
}

// Word is the machine's 16 bit data bus.
type Word int

// minWordInt is the minimum integer that can be represented with a Word.
const minWordInt = -1 << 15 // -32768

// maxWordInt is the maximum integer that can be represented with a Word.
const maxWordInt = 0xFFFF // 65535

// machineMemory is the number of words in the machine's 12-bit addressed memory.
const machineMemory = 1 << 12 // 4096

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
func (m *Machine) Load(f *os.File) {
	program, err := Assemble(f)
	switch err := err.(type) {
	case nil:
	case SyntaxError:
		if err.error == nil {
			fmt.Fprintf(os.Stderr, "syntax: %s:%d: %s\n", f.Name(), err.lineNo, err.line)
		} else {
			fmt.Fprintf(os.Stderr, "syntax: %s: %s:%d: %s\n", err.error, f.Name(), err.lineNo, err.line)
		}
		os.Exit(1)
	default:
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(program) >= machineMemory {
		fmt.Fprintln(os.Stderr, "program too long:", len(program), "instructions")
	}
	for i, w := range program {
		m.M[i] = w
	}
}
