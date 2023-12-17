package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fatih/color"
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

	// Define the custom function map with the added withMergedValues function
	funcMap := template.FuncMap{
		"mergeValues": p.MergeValues,
		"log":         p.DebugPrettyPrint,
		"override":    p.override,
	}

	// Create a new template and parse the template data
	tmpl, err := template.New("data-template").Funcs(funcMap).Funcs(sprig.FuncMap()).Parse(string(templateData))
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

// MergeMaps merges two maps recursively, with the values from "overwrite" map taking precedence.
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

// MergeValues merges multiple maps and returns the result.
func (p *Processor) MergeValues(maps ...map[string]interface{}) map[string]interface{} {
	// Initialize an empty base map to merge with
	base := make(map[string]interface{})

	// Merge each map into the base
	for _, m := range maps {
		base = p.MergeMaps(base, m)
	}

	return base
}

func (p *Processor) override(args ...interface{}) interface{} {
	var prevType reflect.Type
	var base interface{}

	// Check for nil base outside the loop
	for _, arg := range args {
		if arg == nil {
			continue
		}
		// Get the type of the current argument
		currType := reflect.TypeOf(arg)
		if base == nil {
			base = arg
			prevType = currType
			continue
		}

		// If the previous type is not nil, compare types
		if currType != prevType {
			panic(fmt.Sprintf("Error: Inconsistent types in override. Expected %v, got %v", prevType, currType))
		}

		// If the type is map, merge maps
		if currType.Kind() == reflect.Map {
			base = p.MergeMaps(base.(map[string]interface{}), arg.(map[string]interface{}))
		} else {
			base = arg
		}
	}
	return base
}

func (p *Processor) DebugPrettyPrint(args ...interface{}) string {
	fmt.Print("\nDebug Outputs:\n")
	for _, val := range args {
		yamlBytes, err := yaml.Marshal(val)
		if err != nil {
			color.Red(fmt.Sprintf("%s %v\n", "formatting value:", err))
			continue
		}
		// result += string(yamlBytes) + "\n"
		fmt.Printf("\n---\n")
		color.Yellow(string(yamlBytes))
	}
	fmt.Print("\n\n")
	return ""
}

// GetOutputFile returns the path of the output file.
func (p *Processor) GetOutputFile() string {
	return p.OutputFile
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			color.Red(fmt.Sprintf("Failed: %v", r))
			os.Exit(1)
		}
	}()

	// Define flags for templateFile, outputFile, and inputFiles
	var templateFile string
	var outputFile string
	var valuesFiles string
	var showHelp bool

	flag.StringVar(&templateFile, "template", "./data-template.yaml", "Path to the template YAML file")
	flag.StringVar(&outputFile, "output", "./data.yaml", "Path to the output YAML file")
	flag.StringVar(&valuesFiles, "values", "./values.yaml", "Path to comma-separated values YAML file")
	flag.BoolVar(&showHelp, "help", false, "Show usage information")

	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	fmt.Println("    template: " + templateFile)
	fmt.Println("      output: " + outputFile)
	fmt.Println("values files:")
	if valuesFiles != "" {
		inputFiles := strings.Split(valuesFiles, ",")
		for i, file := range inputFiles {
			valueFile := strings.TrimSpace(file)
			inputFiles[i] = valueFile
			fmt.Printf(" %d)  %s\n", i+1, valueFile)
		}
	} else {
		color.Red("Failed no values files given")
		os.Exit(1)
	}

	// Get the comma-separated values from the inputFilesStr
	inputFiles := strings.Split(valuesFiles, ",")
	// Merge the values using the processor
	processor := NewProcessor(templateFile, outputFile)
	mergedData, err := processor.MergeValueFiles(inputFiles)
	if err != nil {
		color.Red("Failed to merge values")
		fmt.Printf("%v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	// Process the template with the merged values
	err = processor.ProcessTemplate(mergedData)
	if err != nil {
		color.Red("Failed to execute template:")
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	color.Green(fmt.Sprintln("Template processing completed. Output file:", processor.GetOutputFile()))
}
