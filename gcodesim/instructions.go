package gcodesim

import (
	"errors"
)

type InstructionCode interface {
	String() string // With description
	Code() string   // Only G/M Code for reproduction
	Simulate(printer *Printer, parameters []Parameter) error
}

/* TODO
- Only absolute positioning is supported currently
- A, B, C axis are ignored
- Only one hotend (index 0) supported
- Following codes are missing
...
*/

var (
	ErrNotImplemented   = errors.New("not implemented")
	ErrInvalidParameter = errors.New("invalid parameter")
)

var InstructionCodeMap map[string]InstructionCode = map[string]InstructionCode{
	"G0":   new(G0),
	"G1":   new(G1),
	"G28":  new(G28),
	"G92":  new(G92),
	"M82":  new(M82),
	"M84":  new(M84),
	"M104": new(M104),
	"M106": new(M106),
	"M107": new(M107),
	"M109": new(M109),
	"M140": new(M140),
	"M190": new(M190),
}

type G0 struct{}

func (i *G0) String() string {
	return "G0: Linear Move (Non-Extrusion)"
}

func (i *G0) Code() string {
	return "G0"
}

func (i *G0) Simulate(printer *Printer, parameters []Parameter) error {
	// Possible parameters: F, X, Y, Z

	target := printer.HeadPosition

	for _, p := range parameters {
		switch p.Letter {
		case 'A', 'B', 'C':
			//ignore
		case 'F':
			printer.CurrentFeedRate = p.Value
		case 'X':
			target.X = p.Value + printer.Offset.X
		case 'Y':
			target.Y = p.Value + printer.Offset.Y
		case 'Z':
			target.Z = p.Value + printer.Offset.Z
		case 'S':
			// laser power, ignore
		default:
			return ErrInvalidParameter
		}
	}

	if _, ok := printer.Layers[target.Z]; target.Z != printer.HeadPosition.Z && !ok {
		printer.Layers[target.Z] = NewLayer()
	}

	// Move print head
	printer.HeadPosition = target

	return nil
}

type G1 struct{}

func (i *G1) String() string {
	return "G1: Linear Move (Extrusion)"
}

func (i *G1) Code() string {
	return "G1"
}

func (i *G1) Simulate(printer *Printer, parameters []Parameter) error {
	// Possible parameters: E, F, X, Y, Z

	target := printer.HeadPosition

	for _, p := range parameters {
		switch p.Letter {
		case 'A', 'B', 'C':
			//ignore
		case 'F':
			printer.CurrentFeedRate = p.Value
		case 'E':
			target.E = p.Value + printer.Offset.E
		case 'X':
			target.X = p.Value + printer.Offset.X
		case 'Y':
			target.Y = p.Value + printer.Offset.Y
		case 'Z':
			target.Z = p.Value + printer.Offset.Z
		case 'S':
			// laser power, ignore
		default:
			return ErrInvalidParameter
		}
	}

	if _, ok := printer.Layers[target.Z]; target.Z != printer.HeadPosition.Z && !ok {
		printer.Layers[target.Z] = NewLayer()
	}

	printer.Layers[target.Z].ExtrudeMovements++
	printer.Layers[target.Z].Lines = append(printer.Layers[target.Z].Lines,
		Line{
			x0: printer.HeadPosition.X,
			y0: printer.HeadPosition.Y,
			x1: target.X,
			y1: target.Y,
		})

	// Move print head
	printer.HeadPosition = target

	if printer.ExtrudedFilament < target.E {
		printer.ExtrudedFilament = target.E
	}

	return nil
}

type G28 struct{}

func (i *G28) String() string {
	return "G28: Auto Home"
}

func (i *G28) Code() string {
	return "G28"
}

func (i *G28) Simulate(printer *Printer, parameters []Parameter) error {

	if len(parameters) == 0 {
		printer.HeadPosition = Position{
			X: 0,
			Y: 0,
			Z: 0,
		}
		return nil
	}

	target := printer.HeadPosition

	for _, p := range parameters {
		switch p.Letter {
		case 'L', 'O', 'R':
			//ignore
		case 'A', 'B', 'C':
			//ignore
		case 'X':
			target.X = 0
		case 'Y':
			target.Y = 0
		case 'Z':
			target.Z = 0
		default:
			return ErrInvalidParameter
		}
	}

	printer.HeadPosition = target

	return nil
}

type G92 struct{}

func (i *G92) String() string {
	return "G92: Set Position"
}

func (i *G92) Code() string {
	return "G92"
}

func (i *G92) Simulate(printer *Printer, parameters []Parameter) error {
	targetOffset := printer.Offset

	/* Target Axis parameters become the current point without physical movement of pos
	by setting offset: */

	for _, p := range parameters {
		switch p.Letter {
		case 'A', 'B', 'C':
			//ignore
		case 'X':
			targetOffset.X = printer.HeadPosition.X - p.Value
		case 'Y':
			targetOffset.Y = printer.HeadPosition.Y - p.Value
		case 'Z':
			targetOffset.Z = printer.HeadPosition.Z - p.Value
		case 'E':
			targetOffset.E = printer.HeadPosition.E - p.Value
		default:
			return ErrInvalidParameter
		}
	}

	printer.Offset = targetOffset

	return nil

}

type M82 struct{}

func (i *M82) String() string {
	return "M82: E-Axis Absolute"
}

func (i *M82) Code() string {
	return "M82"
}

func (i *M82) Simulate(printer *Printer, parameters []Parameter) error {
	if len(parameters) > 0 {
		return ErrInvalidParameter
	}
	printer.RelativeExtruder = false
	return nil
}

type M84 struct{}

func (i *M84) String() string {
	return "M82: Disable / Idle Steppers (Until next move)"
}

func (i *M84) Code() string {
	return "M84"
}

func (i *M84) Simulate(printer *Printer, parameters []Parameter) error {
	// ignore for simulation purpose
	return nil
}

type M104 struct{}

func (i *M104) String() string {
	return "M104: Set Hotend Temperature"
}

func (i *M104) Code() string {
	return "M104"
}

func (i *M104) Simulate(printer *Printer, parameters []Parameter) error {
	if len(parameters) < 1 {
		return ErrInvalidParameter
	}

	for _, p := range parameters {
		switch p.Letter {
		case 'S':
			printer.ExtruderTemperature = p.Value
		case 'T':
			if p.Value != 0 {
				return ErrInvalidParameter
			}
		default:
			return ErrInvalidParameter
		}
	}

	return nil
}

type M106 struct{}

func (i *M106) String() string {
	return "M106: Set Fan Speed"
}

func (i *M106) Code() string {
	return "M106"
}

func (i *M106) Simulate(printer *Printer, parameters []Parameter) error {
	if len(parameters) < 1 {
		return ErrInvalidParameter
	}

	for _, p := range parameters {
		switch p.Letter {
		case 'S':
			printer.FanSpeed = uint8(p.Value)
		default:
			return ErrInvalidParameter
		}
	}

	return nil
}

type M107 struct{}

func (i *M107) String() string {
	return "M107: Turn off fan"
}

func (i *M107) Code() string {
	return "M107"
}

func (i *M107) Simulate(printer *Printer, parameters []Parameter) error {
	if len(parameters) > 0 {
		return ErrInvalidParameter
	}

	printer.FanSpeed = 0

	return nil
}

type M109 struct{}

func (i *M109) String() string {
	return "M109: Wait for Hotend Temperature"
}

func (i *M109) Code() string {
	return "M109"
}

func (i *M109) Simulate(printer *Printer, parameters []Parameter) error {
	// ignore for simulation purpose
	return nil
}

type M140 struct{}

func (i *M140) String() string {
	return "M140: Set Bed Temperature"
}

func (i *M140) Code() string {
	return "M140"
}

func (i *M140) Simulate(printer *Printer, parameters []Parameter) error {
	if len(parameters) < 1 {
		return ErrInvalidParameter
	}

	for _, p := range parameters {
		switch p.Letter {
		case 'S':
			printer.BedTemperature = p.Value
		default:
			return ErrInvalidParameter
		}
	}

	return nil
}

type M190 struct{}

func (i *M190) String() string {
	return "M190: Wait for Bed Temperature"
}

func (i *M190) Code() string {
	return "M190"
}

func (i *M190) Simulate(printer *Printer, parameters []Parameter) error {
	// ignore for simulation purpose
	return nil
}
