package frontend

import (
	"reflect"
	"testing"
)

func TestReadAll(t *testing.T) {
	// Test case 1: Valid input
	validInput := []byte("Name,Age,City\nAlice,25,New York\nBob,30,San Francisco\n")
	expectedValidOutput := [][]string{
		{"Name", "Age", "City"},
		{"Alice", "25", "New York"},
		{"Bob", "30", "San Francisco"},
	}

	validOutput, err := ReadAll(validInput)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(validOutput, expectedValidOutput) {
		t.Errorf("ReadAll(validInput) = %v, expected %v", validOutput, expectedValidOutput)
	}

	// Test case 2: Invalid CSV data
	invalidInput := []byte("Name,Age,City\nAlice,25\nBob,30,San Francisco\n")
	_, invalidErr := ReadAll(invalidInput)
	if invalidErr == nil {
		t.Error("Expected error for invalid CSV data, but got nil")
	} else if invalidErr.Error() != "failed to read CSV: record on line 2: wrong number of fields" {
		t.Errorf("Unexpected error message for invalid CSV data: %v", invalidErr)
	}

	// Test case 3: Empty input
	emptyInput := []byte("")
	emptyOutput, emptyErr := ReadAll(emptyInput)
	if emptyErr != nil {
		t.Errorf("Unexpected error for empty input: %v", emptyErr)
	}
	if len(emptyOutput) != 0 {
		t.Errorf("ReadAll(emptyInput) = %v, expected empty slice", emptyOutput)
	}
}
