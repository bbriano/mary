package main

import (
	"fmt"
	"os"
)

// Word is the machine's 16 bit data bus.
type Word int

// minWordInt is the minimum integer that can be represented with a Word.
const minWordInt = -1 << 15 // -32768

// maxWordInt is the maximum integer that can be represented with a Word.
const maxWordInt = 0xFFFF // 65535

// machineMemory is the number of words in the machine's 12-bit addressed memory.
const machineMemory = 1 << 12 // 4096

// Machine simulates a Marie machine. Most of the registers are not needed for the simulation,
// but they are added to illustrate the Marie machine described in the book.
type Machine struct {
	AC Word
	PC Word
	MAR Word
	MBR Word
	IR Word
	IN Word
	OUT Word
	M [machineMemory]Word
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

// Load loads f to the machine's memory.
func (m *Machine) Load(f *os.File) error {
	program, err := Assemble(f)
	switch err := err.(type) {
	case nil:
	case SyntaxError:
		return fmt.Errorf("syntax: %s:%d: %s\n", f.Name(), err.lineNo, err.line)
	default:
		return fmt.Errorf("%v", err)
	}
	if len(program) >= machineMemory {
		return fmt.Errorf("program too long: %d/%d instructions", len(program), machineMemory)
	}
	for i, w := range program {
		m.M[i] = w
	}
	return nil
}
