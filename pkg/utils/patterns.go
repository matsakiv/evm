package utils

import "io"

type Pattern interface {
    Match(inst *Instruction) bool
}

type IsOpcode struct {
    Opcode byte
}

func (p IsOpcode) Match(inst *Instruction) bool {
    return p.Opcode == inst.Opcode
}

type IsPushX struct {
}

func (p IsPushX) Match(inst *Instruction) bool {
    return inst.Opcode >= PUSH1 && inst.Opcode <= PUSH32
}

func FindPattern(reader io.ByteReader, pattern []Pattern, prevInstr *OrderedInstruction) ([]*OrderedInstruction, error) {
    index := 0
    address := 0

    if prevInstr != nil {
        index = prevInstr.Index
        address = prevInstr.Address
    }

    lps := computeLps(pattern)
    patternLength := len(pattern)
    patternCandidate := make(chan *OrderedInstruction, patternLength)

    j := 0

    for {
        instr, err := Read(reader)

        if err == io.EOF {
            return nil, err
        }

        if len(patternCandidate) == patternLength {
            <-patternCandidate
        }

        patternCandidate <- &OrderedInstruction{Instruction: instr, Index: index, Address: address}

        for {
            if pattern[j].Match(instr) {
                j++

                if j == patternLength {
                    // found pattern at: index - j + 1
                    foundPattern := make([]*OrderedInstruction, patternLength)

                    for k := 0; k < patternLength; k++ {
                        foundPattern[k] = <-patternCandidate
                    }

                    return foundPattern, nil
                    //j = lps[j-1]
                }

                break
            }

            if j > 0 {
                j = lps[j-1]
            } else {
                break
            }
        }

        if instr.Data != nil {
            address += len(instr.Data)
        }

        index++
        address++
    }
}

func computeLps(pattern []Pattern) (lps []int) {
    patternLength := len(pattern)
    lps = make([]int, patternLength)

    l := 0

    for i := 1; i < patternLength; i++ {
        for {
            if pattern[i] == pattern[l] {
                l++
                break
            }

            if l == 0 {
                break
            }

            l = lps[l-1]
        }

        lps[i] = l
    }

    return lps
}
