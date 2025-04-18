/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package shortener

import (
	"testing"
)

func TestGenerateUnique(t *testing.T) {
	// Test unique ID generation
	testCases := []struct {
		name              string
		existenceChecker  func(string) bool
		expectedUniqueIDs int
	}{
		{
			name: "No Existing IDs",
			existenceChecker: func(id string) bool {
				return false
			},
			expectedUniqueIDs: 5,
		},
		{
			name: "Some Existing IDs",
			existenceChecker: func(id string) bool {
				// Let's assume there is a conflict at every 3rd ID
				return len(id)%3 == 0
			},
			expectedUniqueIDs: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			generatedIDs := make(map[string]bool)

			for i := 0; i < tc.expectedUniqueIDs; i++ {
				shortID, err := GenerateUnique(tc.existenceChecker)

				// Error check
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Uniqueness check
				if generatedIDs[shortID] {
					t.Errorf("Duplicate ID generated: %s", shortID)
				}

				generatedIDs[shortID] = true
			}

			// Control of the number of IDs generated
			if len(generatedIDs) != tc.expectedUniqueIDs {
				t.Errorf("Expected %d unique IDs, got %d", tc.expectedUniqueIDs, len(generatedIDs))
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	// Simple ID generation test
	shortID, err := Generate()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// ID length control
	if len(shortID) != 8 { // default length
		t.Errorf("Expected ID length 8, got %d", len(shortID))
	}
}

func TestCustomGenerator(t *testing.T) {
	// special character set and length test
	generator := NewIDGenerator(10)
	shortID, err := generator.Generate()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Id length control
	if len(shortID) != 10 {
		t.Errorf("Expected ID length 10, got %d", len(shortID))
	}
}

// performance test for ID generation
func BenchmarkGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Generate()
	}
}

func BenchmarkGenerateUnique(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateUnique(func(id string) bool {
			return false
		})
	}
}
