package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Instruction corresponds to a single machine instruction in a Marie machine.
type Instruction func(*Machine, Word)

// Each function defined in this file of type Instruction.

// instruction provides convenient access to Instruction functions.
var instruction map[Opcode]Instruction

func init() {
	instruction = map[Opcode]Instruction{
		OpJnS:      JnS,
		OpLoad:     Load,
		OpStore:    Store,
		OpAdd:      Add,
		OpSubt:     Subt,
		OpInput:    Input,
		OpOutput:   Output,
		OpHalt:     Halt,
		OpSkipcond: Skipcond,
		OpJump:     Jump,
		OpClear:    Clear,
		OpAddI:     AddI,
		OpJumpI:    JumpI,
	}
}

func Load(m *Machine, x Word) {
	m.MAR = x
	m.MBR = m.M[m.MAR]
	m.AC = m.MBR
}

func Store(m *Machine, x Word) {
	m.MAR = x
	m.MBR = m.AC
	m.M[m.MAR] = m.MBR
}

func Add(m *Machine, x Word) {
	m.MAR = x
	m.MBR = m.M[m.MAR]
	m.AC += m.MBR
}

func Subt(m *Machine, x Word) {
	m.MAR = x
	m.MBR = m.M[m.MAR]
	m.AC -= m.MBR
}

func Input(m *Machine, _ Word) {
	var x int
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		var err error
		hex := s.Text()
		x, err = strconv.Atoi(hex)
		if err == nil {
			break
		}
	}
	m.IN = Word(x)
	m.AC = m.IN
}

func Output(m *Machine, _ Word) {
	m.OUT = m.AC
	fmt.Printf("%x\n", m.OUT)
}

func Halt(m *Machine, _ Word) {
	os.Exit(0)
}

func Skipcond(m *Machine, x Word) {
	switch x >> 10 & 3 {
	case 0:
		if m.AC < 0 {
			m.PC++
		}
	case 1:
		if m.AC == 0 {
			m.PC++
		}
	case 2:
		if m.AC > 0 {
			m.PC++
		}
	case 3:
		fmt.Fprintln(os.Stderr, "bad instruction:", m.IR)
		os.Exit(1)
	}
}

func Jump(m *Machine, x Word) {
	m.PC = x
}

func JnS(m *Machine, x Word) {
	m.MAR = x
	m.MBR = m.PC
	m.M[m.MAR] = m.MBR
	m.MBR = x
	m.AC = 1
	m.AC += m.MBR
	m.PC = m.AC
}

func Clear(m *Machine, x Word) {
	m.AC = 0
}

func AddI(m *Machine, x Word) {
	m.MAR = x
	m.MBR = m.M[m.MAR]
	m.MAR = m.MBR
	m.MBR = m.M[m.MAR]
	m.AC += m.MBR
}

func JumpI(m *Machine, x Word) {
	m.MAR = x
	m.MBR = m.M[m.MAR]
	m.PC = m.MBR
}
