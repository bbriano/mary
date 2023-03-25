package main

// Opcodes as defined by the instruction set.
type Opcode int

// opcode maps opcode string literals to Opcode.
// It is used to parse Marie assembly code in Machine.Load.
var opcode map[string]Opcode

func init() {
	opcode = map[string]Opcode{
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
	}
}

const (
	// Do not re-order these constants.
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
)
