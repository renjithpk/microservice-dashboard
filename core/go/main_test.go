package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
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
	mergedData := map[string]interface{}{
		"GroupA": map[string]interface{}{
			"color": "red",
			"items": []map[string]interface{}{
				{
					"name": "ItemA1",
					"menu": []map[string]interface{}{
						{
							"name": "MenuA1",
							"link": "linkA1",
						},
					},
					"backends": "buttonA",
				},
				{
					"name": "ItemA2",
					"menu": []map[string]interface{}{
						{
							"name": "MenuB1",
							"link": "linkB1",
						},
					},
					"backends": "buttonB",
				},
			},
		},
		"GroupB": map[string]interface{}{
			"color": "blue",
			"items": []map[string]interface{}{
				{
					"name": "ItemB1",
					"menu": []map[string]interface{}{
						{
							"name": "MenuC1",
							"link": "linkC1",
						},
					},
				},
			},
		},
	}

	// Create a new Processor with test-specific template and output files
	processor := NewProcessor("./testdata/data-template.yaml", "temp_output.yaml")

	// Set a temporary output file for the test
	tempOutputFile := processor.GetOutputFile()
	defer os.Remove(tempOutputFile)

	// Test the ProcessTemplate function
	err := processor.ProcessTemplate(mergedData)
	if err != nil {
		t.Fatalf("Failed to process template: %v", err)
	}

	// Read the generated output file and compare its content with the expected template
	expectedOutput := `
- group: GroupA
  color: "red"
  items:
    - name: ItemA1
      menu:
        - name: MenuA1
          link: linkA1
        backends: buttonA
    - name: ItemA2
      menu:
        - name: MenuB1
          link: linkB1
        backends: buttonB
- group: GroupB
  color: "blue"
  items:
    - name: ItemB1
      menu:
        - name: MenuC1
          link: linkC1
`
	outputData, err := ioutil.ReadFile(tempOutputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if string(outputData) != expectedOutput {
		t.Errorf("Processed template output does not match expected output:\nExpected:\n#%s#\nGot:\n#%s#", expectedOutput, string(outputData))
	}
}
