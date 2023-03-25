package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	sim "github.com/fabian-z/fabrico-ledger/gcodesim"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var instructions []*sim.Instruction

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		instruction, err := sim.ParseInstruction(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		if instruction != nil {
			instructions = append(instructions, instruction)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Parsed %d instructions", len(instructions))

	printer := sim.NewPrinter()

	for _, inst := range instructions {
		//log.Println(inst)
		err = inst.InstructionCode.Simulate(printer, inst.Parameters)
		if err != nil {
			log.Println(inst)
			log.Fatal(err)
		}
	}

	log.Printf("%+v", printer)

	for z, l := range printer.Layers {
		if l.ExtrudeMovements > 0 {
			ioutil.WriteFile(fmt.Sprintf("%v.svg", z), l.SVG(), 0666)
		}
	}
}
