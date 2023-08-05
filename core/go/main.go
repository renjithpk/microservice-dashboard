package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Define flags for templateFile, outputFile, and inputFiles
	var templateFile string
	var outputFile string
	var valuesFile string
	var showHelp bool

	flag.StringVar(&templateFile, "template", "./data-template.yaml", "Path to the template YAML file")
	flag.StringVar(&outputFile, "output", "./data.yaml", "Path to the output YAML file")
	flag.StringVar(&valuesFile, "values", "", "Path to comma separated additional values YAML file")
	flag.BoolVar(&showHelp, "help", false, "Show usage information")

	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	fmt.Println("template: " + templateFile)
	fmt.Println("output: " + outputFile)
	if valuesFile != "" {
		fmt.Println("values: " + valuesFile)
		valuesFile = "./values.yaml," + valuesFile
	}

	// Get the comma-separated values from the inputFilesStr
	inputFiles := strings.Split(valuesFile, ",")
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
