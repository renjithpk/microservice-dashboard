package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Processor represents the template processor.
type Processor struct {
	OutputFile   string
	TemplatePath string // Add a new field for the template path
}

// NewProcessor creates a new instance of Processor with the provided template file and output file.
func NewProcessor(templatePath, outputFile string) *Processor {
	return &Processor{
		TemplatePath: templatePath,
		OutputFile:   outputFile,
	}
}

// MergeValueFiles takes a slice of file paths as input and returns the merged map of values from those files.
func (p *Processor) MergeValueFiles(files []string) (map[string]interface{}, error) {
	var mergedData map[string]interface{}

	for _, file := range files {
		if file == "" {
			continue
		}
		// Read the input values from each file
		valuesData, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %v", file, err)
		}

		var values map[string]interface{}
		err = yaml.Unmarshal(valuesData, &values)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal %s: %v", file, err)
		}

		// Merge the current file's content with the previous content (last file has higher precedence)
		mergedData = p.MergeMaps(mergedData, values)
	}

	return mergedData, nil
}

// ProcessTemplate takes the merged data and the template path, processes the template,
// and writes the output to the specified OutputFile.
func (p *Processor) ProcessTemplate(mergedData map[string]interface{}) error {
	// Read the template file
	templateData, err := ioutil.ReadFile(p.TemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %v", p.TemplatePath, err)
	}

	// Create a new template and parse the template data
	tmpl, err := template.New("data-template").Parse(string(templateData))
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	// Process the template with the merged values
	outputFileHandle, err := os.Create(p.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %v", p.OutputFile, err)
	}
	defer outputFileHandle.Close()

	err = tmpl.Execute(outputFileHandle, mergedData)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
}

// mergeMaps merges two maps recursively, with the values from "overwrite" map taking precedence.
func (p *Processor) MergeMaps(base, overwrite map[string]interface{}) map[string]interface{} {
	if base == nil {
		base = make(map[string]interface{})
	}

	for key, val := range overwrite {
		baseVal, found := base[key]
		if found {
			switch val.(type) {
			case map[string]interface{}:
				// If both base and overwrite have nested maps, merge them recursively.
				base[key] = p.MergeMaps(baseVal.(map[string]interface{}), val.(map[string]interface{}))
			default:
				// Otherwise, use the value from overwrite, as it takes precedence.
				base[key] = val
			}
		} else {
			base[key] = val
		}
	}

	return base
}

// GetOutputFile returns the path of the output file.
func (p *Processor) GetOutputFile() string {
	return p.OutputFile
}
