package utils

import (
    "encoding/hex"
    "io"
)

func FindSelectors(reader io.ByteReader) []string {
    pattern := []Pattern{
        IsOpcode{Opcode: DUP1},
        IsOpcode{Opcode: PUSH4},
        IsOpcode{Opcode: EQ},
        IsPushX{},
        IsOpcode{Opcode: JUMPI}}

    selectors := make([]string, 0)

    var lastInstr *OrderedInstruction = nil

    for {
        foundPattern, err := FindPattern(reader, pattern, lastInstr)

        if err == io.EOF {
            return selectors
        }

        selector := hex.EncodeToString(foundPattern[1].Data)
        selectors = append(selectors, selector)
    }
}
