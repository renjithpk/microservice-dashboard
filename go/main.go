package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

func main() {
	// Read the input values from value.yaml file
	valuesFile := "./values.yaml"
	valuesData, err := ioutil.ReadFile(valuesFile)
	if err != nil {
		fmt.Printf("Failed to read %s: %v\n", valuesFile, err)
		os.Exit(1)
	}

	var values interface{}
	err = yaml.Unmarshal(valuesData, &values)
	if err != nil {
		fmt.Printf("Failed to unmarshal %s: %v\n", valuesFile, err)
		os.Exit(1)
	}

	// Read the template file
	templateFile := "./data-template.yaml"
	templateData, err := ioutil.ReadFile(templateFile)
	if err != nil {
		fmt.Printf("Failed to read %s: %v\n", templateFile, err)
		os.Exit(1)
	}

	// Create a new template and parse the template data
	tmpl, err := template.New("data-template").Parse(string(templateData))
	if err != nil {
		fmt.Printf("Failed to parse template: %v\n", err)
		os.Exit(1)
	}

	// Process the template with the input values
	// var resultData []byte
	outputFile := "./data.yaml"
	outputFileHandle, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Failed to create %s: %v\n", outputFile, err)
		os.Exit(1)
	}
	defer outputFileHandle.Close()

	err = tmpl.Execute(outputFileHandle, values)
	if err != nil {
		fmt.Printf("Failed to execute template: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Template processing completed. Output file:", outputFile)
}
