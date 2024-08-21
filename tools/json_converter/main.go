package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: json_converter <input_file> <output_file>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Open input file.
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return
	}
	defer file.Close()

	var jsonArray []json.RawMessage

	// Reads each line of the file and parses it into a JSON object.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		jsonArray = append(jsonArray, json.RawMessage(line))
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	// Writes an array of JSON objects to the output file.
	outputFileHandle, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outputFileHandle.Close()

	// Manually build an array of JSON with indentation.
	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	for i, jsonObj := range jsonArray {
		var indented bytes.Buffer
		if err := json.Indent(&indented, jsonObj, "  ", "  "); err != nil {
			fmt.Printf("Error indenting JSON: %v\n", err)
			return
		}
		buffer.WriteString("  ")
		buffer.Write(indented.Bytes())
		if i < len(jsonArray)-1 {
			buffer.WriteString(",")
		}
		buffer.WriteString("\n")
	}
	buffer.WriteString("]")

	_, err = outputFileHandle.Write(buffer.Bytes())
	if err != nil {
		fmt.Printf("Error writing JSON to output file: %v\n", err)
		return
	}

	fmt.Println("Successfully converted JSON objects to JSON array with indentation.")
}
