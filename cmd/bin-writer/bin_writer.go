package main

import (
	"encoding/hex"
	"log"
	"os"
	"strings"
)

func main() {
	// Define the input hex dump as a string
	hexDump := `00 C1 F0 F2 F0 F0 F2 2A 44 01 29 80 90 00 00 00
	00 04 04 00 00 00 F1 F6 F4 F3 F7 F5 F0 F7 F0 F0
	F0 F1 F4 F2 F3 F9 F5 F5 F3 F0 F1 F0 F0 F0 F0 F0
	F0 F0 F0 F0 F0 F0 F0 F0 F0 F0 F0 F4 F0 F2 F0 F1
	F3 F5 F3 F3 F9 F0 F7 F8 F5 F3 F0 F4 F0 F2 F0 F4
	F0 F2 F6 F0 F1 F1 F0 F2 F1 F1 F0 F0 F0 F0 F0 F0
	F0 F9 F9 F9 F0 F3 F7 F4 F3 F7 F5 F0 F7 F0 F0 F0
	F1 F4 F2 F3 F9 F5 F5 7E F0 F4 F1 F2 F5 F0 F1 F1
	F1 F2 F3 F4 F1 F2 F3 F0 F0 F0 F0 F0 F0 F0 F0 F0
	F0 F0 F6 F4 F3 F7 F6 F2 C2 D5 E3 F0 F0 F0 F0 F0
	F0 F0 F1 F6 F0 F8 19 F3 7F 50 B5 57 35 5E F1 F0
	00 00 00 00 00 F1 F2 F0 F0 F7 F5 F6 F8 F0 F0 F3
	F4 F7 F5`

	// Remove leading/trailing whitespace and split into individual hex values
	hexDump = strings.ReplaceAll(hexDump, "\n", "")
	hexDump = strings.ReplaceAll(hexDump, "\t", "")
	hexDump = strings.ReplaceAll(hexDump, " ", "")

	hexValues, err := hex.DecodeString(hexDump)
	if err != nil {
		log.Fatal(err)
	}

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// Write the binary data to a binary file
	if err := os.WriteFile(dirname+"/output.bin", hexValues, 0644); err != nil {
		log.Fatal(err)
	}

	log.Println("Binary file 'output.bin' has been created.")
}
