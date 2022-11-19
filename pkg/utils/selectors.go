package utils

import (
    "encoding/hex"
    "io"

    "golang.org/x/crypto/sha3"
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

func GetSelector(signature string) string {
    keccak256 := sha3.NewLegacyKeccak256()

    keccak256.Write([]byte(signature))

    return hex.EncodeToString(keccak256.Sum(nil)[0:4])
}
