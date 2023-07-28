package main

import (
	"fmt"
	"os"
)

func main() {
	templateFile := "./data-template.yaml"
	outputFile := "./data.yaml"

	// List of YAML files to merge
	inputFiles := []string{"./values.yaml"}

	// Merge the values using the processor
	processor := NewProcessor(templateFile, outputFile)
	mergedData, err := processor.MergeValueFiles(inputFiles)
	if err != nil {
		fmt.Printf("Failed to merge values: %v\n", err)
		os.Exit(1)
	}

	// Process the template with the merged values
	err = processor.ProcessTemplate(mergedData)
	if err != nil {
		fmt.Printf("Failed to execute template: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Template processing completed. Output file:", processor.GetOutputFile())
}
