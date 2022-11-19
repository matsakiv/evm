package utils

import (
    "fmt"
    "io"
)

type Instruction struct {
    Opcode byte
    Data   []byte
}

type OrderedInstruction struct {
    *Instruction
    Index   int
    Address int
}

func Read(reader io.ByteReader) (*Instruction, error) {
    op, err := reader.ReadByte()

    if err == io.EOF {
        return nil, io.EOF
    }

    if op >= PUSH1 && op <= PUSH32 {
        dataLength := int(op-PUSH1) + 1
        data := make([]byte, dataLength)

        for i := 0; i < dataLength; i++ {
            dataByte, dataErr := reader.ReadByte()

            if dataErr == io.EOF {
                return nil, io.EOF
            }

            data[i] = dataByte
        }

        return &Instruction{op, data}, nil
    } else {
        return &Instruction{op, nil}, nil
    }
}

func PrintDisassembly(reader io.ByteReader) {
    address := 0

    for {
        instr, err := Read(reader)

        if err == io.EOF {
            return
        }

        if instr.Data != nil {
            fmt.Printf("[0x%04X] %s 0x%x\n", address, opcode2String[instr.Opcode], instr.Data)
            address += len(instr.Data)
        } else {
            fmt.Printf("[0x%04X] %s\n", address, opcode2String[instr.Opcode])
        }

        address++
    }
}

func PrintInstruction(instr *Instruction) {
    if instr.Data != nil {
        fmt.Printf("%s 0x%x\n", opcode2String[instr.Opcode], instr.Data)
    } else {
        fmt.Printf("%s\n", opcode2String[instr.Opcode])
    }
}

func PrintOrderedInstruction(instr *OrderedInstruction) {
    if instr.Data != nil {
        fmt.Printf("[0x%04X] %s 0x%x\n", instr.Address, opcode2String[instr.Opcode], instr.Data)
    } else {
        fmt.Printf("[0x%04X] %s\n", instr.Address, opcode2String[instr.Opcode])
    }
}
