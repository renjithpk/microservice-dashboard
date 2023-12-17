package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMergeValueFiles(t *testing.T) {
	// Create temporary test files
	testData := []string{
		`a: 1
b: 2`,
		`b: 3
c: 4`,
		`a: 5
d: 6`,
	}

	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	files := []string{}
	for i, data := range testData {
		filename := fmt.Sprintf("values%d.yaml", i+1)
		filePath := filepath.Join(tempDir, filename)
		err := ioutil.WriteFile(filePath, []byte(data), 0644)
		if err != nil {
			t.Fatalf("Failed to write test data: %v", err)
		}
		files = append(files, filePath)
	}

	// Create a new Processor with test-specific template and output files
	processor := NewProcessor("../testdata/data-template.yaml", "temp_output.yaml")

	// Test the MergeValueFiles function
	mergedData, err := processor.MergeValueFiles(files)
	if err != nil {
		t.Fatalf("Failed to merge values: %v", err)
	}

	// Verify the merged data
	expectedMergedData := map[string]interface{}{
		"a": 5,
		"b": 3,
		"c": 4,
		"d": 6,
	}

	for key, val := range expectedMergedData {
		if mergedData[key] != val {
			t.Errorf("Expected mergedData[%s] to be %v, but got %v", key, val, mergedData[key])
		}
	}
}

func TestProcessTemplate(t *testing.T) {
	// Template data directly in the test case
	templateData := `
groups:
{{- range .groups }}
  {{- $groupMetadata := .metadata}}
  - group: {{ .group }}
    tiles:
  {{- range .categories }}
    {{- $categoryMetadata := .metadata}}
    {{- range .apis}}
      - name: {{ .name }}
      {{- $metadata := mergeValues $groupMetadata $categoryMetadata }}
      {{- if $metadata.jenkin_badge}}
        jenkin_badge: true
      {{- end}}
    {{- end}}
  {{- end}}
{{- end }}
`

	// Define your merged data (replace this with your actual merged data)
	valuesInYamlForm := `
groups:
  - group: group-name
    metadata:
      git_badge: true
      jenkin_badge: true
    categories:
      - type: type
        metadata:
          git_badge: false
        apis:
          - name: api1
            jenkin_badge: true
          - name: api2
            jenkin_badge: true
`

	// Read the generated output file and compare its content with the expected template
	expectedOutput := `
groups:
  - group: group-name
    tiles:
      - name: api1
        jenkin_badge: true
      - name: api2
        jenkin_badge: true
`

	// Create temporary test directory
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a temporary file for the template
	templateFile := filepath.Join(tempDir, "data-template.yaml")
	err = ioutil.WriteFile(templateFile, []byte(templateData), 0644)
	if err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// Output file in the temporary directory
	outputFile := filepath.Join(tempDir, "temp_output.yaml")

	// Your test logic using the temporary template and output files...
	processor := NewProcessor(templateFile, outputFile)

	// TODO add code to populate mergedData
	// mergedData = convert valuesInYamlForm
	var mergedData map[string]interface{}
	err = yaml.Unmarshal([]byte(valuesInYamlForm), &mergedData)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML data: %v", err)
	}

	// Test the ProcessTemplate function
	err = processor.ProcessTemplate(mergedData)
	if err != nil {
		t.Fatalf("Failed to process template: %v", err)
	}

	outputData, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if string(outputData) != expectedOutput {
		t.Errorf("Processed template output does not match expected output:\nExpected:\n#%s#\nGot:\n#%s#", expectedOutput, string(outputData))
	}
}
