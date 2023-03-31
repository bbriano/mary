package main

import (
	"bufio"
	"fmt"
	"os"
)

// Opcode is the 4-bit operation code of an instruction.
type Opcode int

// opcode maps operation string literals to opcode values.
// It is used to parse Marie assembly code in Machine.Load.
var opcode map[string]Opcode = map[string]Opcode{
	"JnS":      OpJnS,
	"Load":     OpLoad,
	"Store":    OpStore,
	"Add":      OpAdd,
	"Subt":     OpSubt,
	"Input":    OpInput,
	"Output":   OpOutput,
	"Halt":     OpHalt,
	"Skipcond": OpSkipcond,
	"Jump":     OpJump,
	"Clear":    OpClear,
	"AddI":     OpAddI,
	"JumpI":    OpJumpI,
	"LoadI":    OpLoadI,
	"StoreI":   OpStoreI,
	"Dump":     OpDump,
}

// Instruction encodes the execute operation of an instruction.
type Instruction func(*Machine, Word)

// instruction maps opcode to Instruction functions.
// It is used to decode the machine code in Machine.Run.
var instruction map[Opcode]Instruction = map[Opcode]Instruction{
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
	OpLoadI:    LoadI,
	OpStoreI:   StoreI,
	OpDump:     Dump,
}

const (
	OpJnS Opcode = iota
	OpLoad
	OpStore
	OpAdd
	OpSubt
	OpInput
	OpOutput
	OpHalt
	OpSkipcond
	OpJump
	OpClear
	OpAddI
	OpJumpI
	OpLoadI
	OpStoreI
	OpDump
)

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
	var x Word
	s := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for s.Scan() {
		var err error
		hex := s.Text()
		x, err = parseWord(hex, 16)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			fmt.Print("> ")
			continue
		}
		break
	}
	m.IN = x
	m.AC = m.IN
}

func Output(m *Machine, _ Word) {
	m.OUT = m.AC
	fmt.Printf("%04x\n", m.OUT)
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

func LoadI(m *Machine, x Word) {
	m.MAR = x
	m.MBR = m.M[m.MAR]
	m.MAR = m.MBR
	m.MBR = m.M[m.MAR]
	m.AC = m.MBR
}

func StoreI(m *Machine, x Word) {
	m.MAR = x
	m.MBR = m.M[m.MAR]
	m.MAR = m.MBR
	m.MBR = m.AC
	m.M[m.MAR] = m.MBR
}

func Dump(m *Machine, x Word) {
	fmt.Printf("AC=%d PC=%d MAR=%d MBR=%d IR=%d IN=%d OUT=%d\n",
		m.AC, m.PC, m.MAR, m.MBR, m.IR, m.IN, m.OUT)
	rows := int((x-1)/16) + 1
	for i := 0; i < rows; i++ {
		fmt.Printf("%04X:", i*16)
		for j := 0; j < 16; j++ {
			if i*16+j == int(x) {
				break
			}
			fmt.Printf(" %04X", m.M[i*16+j])
		}
		fmt.Println()
	}
}
