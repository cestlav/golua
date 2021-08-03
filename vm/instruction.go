package vm

type Instruction uint32

const MAXARG_Bx = 1 << 18 - 1
const MAXARG_sBx = MAXARG_Bx >> 1

func (i Instruction) OpCode() int {
	return int(i & 0x3F)
}

func (i Instruction) ABC() (int, int, int) {
	a := int(i >> 6 & 0xFF)
	b := int(i >> 14 & 0x1FF)
	c := int(i >> 23 & 0x1FF)
	return a, b, c
}

func (i Instruction) ABx() (int, int) {
	a := int(i >> 6 & 0xFF)
	b := int(i >> 14)
	return a, b
}

func (i Instruction) AsBx() (int, int) {
	a, bx := i.ABx()
	return a, bx - MAXARG_sBx
}

func (i Instruction) Ax() int {
	return int(i >> 6)
}

func (i Instruction) OpName() string {
	return opcodes[i.OpCode()].name
}

func (i Instruction) OpMode() byte {
	return opcodes[i.OpCode()].opMode
}

func (i Instruction) ArgBMode() byte {
	return opcodes[i.OpCode()].argBMode
}

func (i Instruction) ArgCMode() byte {
	return opcodes[i.OpCode()].argCMode
}