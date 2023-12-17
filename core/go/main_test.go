package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
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

func TestOverride(t *testing.T) {
	processor := NewProcessor("", "") // Assuming you have proper initialization

	tests := []struct {
		name          string
		args          []interface{}
		expected      interface{}
		expectedPanic string
	}{
		{
			name:     "Merge Maps",
			args:     []interface{}{map[string]interface{}{"a": 1}, map[string]interface{}{"b": 2}},
			expected: map[string]interface{}{"a": 1, "b": 2},
		},
		{
			name:     "Replace String",
			args:     []interface{}{"abc", "xyz"},
			expected: "xyz",
		},
		{
			name:     "string and nil",
			args:     []interface{}{"abc", nil},
			expected: "abc",
		},
		{
			name:     "nil and String",
			args:     []interface{}{nil, "abc"},
			expected: "abc",
		},
		{
			name:     "Nil Arguments",
			args:     []interface{}{nil, nil},
			expected: nil,
		},
		{
			name:     "Map and Nil",
			args:     []interface{}{map[string]interface{}{"a": 1}, nil},
			expected: map[string]interface{}{"a": 1},
		},
		{
			name:     "nil and Map",
			args:     []interface{}{nil, map[string]interface{}{"a": 1}},
			expected: map[string]interface{}{"a": 1},
		},
		{
			name:          "Inconsistent Types",
			args:          []interface{}{map[string]interface{}{"a": 1}, "xyz", 123, nil},
			expected:      nil,
			expectedPanic: "Error: Inconsistent types in override. Expected map[string]interface {}, got string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Recover from panics and compare with the expected panic value
			defer func() {
				if r := recover(); r != nil {
					if got, want := fmt.Sprintf("%v", r), test.expectedPanic; got != want {
						t.Errorf("Expected panic: %v, got: %v", want, got)
					}
				}
			}()

			result := processor.override(test.args...)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}
func TestProcessTemplate(t *testing.T) {
	// Template data directly in the test case
	templateData := `
groups:
{{- range .groups }}
  {{- log . }}
  {{- $groupMetadata := .metadata}}
  - group: {{ upper .group }}
    tiles:
  {{- range .categories }}
    {{- $categoryMetadata := .metadata}}
    {{- range .apis}}
      - name: {{ .name }}
      {{- $metadata := override $groupMetadata $categoryMetadata }}
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
  - group: GROUP-NAME
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
