package gcodesim

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var (
	commentMatcher        = regexp.MustCompile(`;.*$|\(.*\)`)
	whiteSpaceMatcher     = regexp.MustCompile(`\s+`)
	ErrUnknownInstruction = errors.New("unknown instruction")
)

type Parameter struct {
	Letter rune // e.g. X / Y / Z / E axis, etc.
	Value  float64
}

func (p Parameter) String() string {
	return fmt.Sprintf("%s%v", string(p.Letter), p.Value)
}

// GCode represents a single line in GCode format
type Instruction struct {
	InstructionCode InstructionCode
	Parameters      []Parameter
}

func (i *Instruction) Format() string {
	res := i.InstructionCode.Code()
	for _, p := range i.Parameters {
		res = res + " " + p.String()
	}
	return res
}

func ParseInstruction(s string) (*Instruction, error) {
	// general format: G0
	var inst Instruction

	// remove commments, normalize whitespace
	s = commentMatcher.ReplaceAllString(s, "")
	s = whiteSpaceMatcher.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)

	if len(s) == 0 {
		return nil, nil
	}

	components := strings.Split(s, " ")

	// Expect instruction code first
	code, err := ParseInstructionCode(components[0])

	if err != nil {
		log.Printf("%v unknown", components[0])
		return nil, ErrUnknownInstruction
	}

	inst.InstructionCode = code

	if len(components) < 2 {
		return &inst, nil
	}

	components = components[1:]
	var c string

	for len(components) > 0 {
		var par Parameter
		c, components = components[0], components[1:]

		//log.Println("Component: ", c)

		if strings.HasPrefix(c, ";") {
			// EOL comment, can contain spaces
			break
		}

		if c[0] >= 65 && c[0] <= 90 {
			// First character is upper case A-Z
			par.Letter = rune(c[0])

			if len(c[1:]) > 0 {
				// Values may be optional, e.g. for G28
				par.Value, err = strconv.ParseFloat(c[1:], 64)
				if err != nil {
					return nil, fmt.Errorf("error parsing parameter value: %w", err)
				}
			}

			inst.Parameters = append(inst.Parameters, par)
		}

	}

	return &inst, nil
}

func ParseInstructionCode(s string) (InstructionCode, error) {
	impl, ok := InstructionCodeMap[s]

	if !ok {
		return nil, errors.New("unknown instruction code")
	}

	return impl, nil
}
