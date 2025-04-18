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
	"crypto/rand"
	"errors"
	"math/big"
)

var (
	// DefaultGenerator is the default ID generator
	DefaultGenerator *IDGenerator

	// ErrUniqueIDGenerationFailed occurs when unique ID generation fails
	ErrUniqueIDGenerationFailed = errors.New("unique ID generation failed after multiple attempts")
)

const (
	defaultAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defaultLength   = 8
)

// init initializes the default generator
func init() {
	DefaultGenerator = NewIDGenerator(defaultLength) // Default 6 character length
}

// IDGenerator creates short unique IDs
type IDGenerator struct {
	alphabet string // Character set for ID generation
	length   int    // Length of generated ID
}

// NewIDGenerator creates a new ID generator
func NewIDGenerator(length int) *IDGenerator {
	return &IDGenerator{
		alphabet: defaultAlphabet,
		length:   length,
	}
}

// Generate creates a random short ID
func (g *IDGenerator) Generate() (string, error) {
	shortID := make([]byte, g.length)
	alphabetLength := big.NewInt(int64(len(g.alphabet)))

	for i := 0; i < g.length; i++ {
		// Cryptographically secure random index
		randIndex, err := rand.Int(rand.Reader, alphabetLength)
		if err != nil {
			return "", err
		}
		shortID[i] = g.alphabet[randIndex.Int64()]
	}

	return string(shortID), nil
}

// GenerateUnique creates a unique ID using an existence checker
func (g *IDGenerator) GenerateUnique(existenceChecker func(string) bool) (string, error) {
	maxAttempts := 10 // Prevent infinite loop

	for attempt := 0; attempt < maxAttempts; attempt++ {
		shortID, err := g.Generate()
		if err != nil {
			return "", err
		}

		// Return ID if it doesn't exist
		if !existenceChecker(shortID) {
			return shortID, nil
		}
	}

	return "", ErrUniqueIDGenerationFailed
}

// CustomAlphabet allows setting a custom character set
func (g *IDGenerator) CustomAlphabet(alphabet string) *IDGenerator {
	g.alphabet = alphabet
	return g
}

// Generate uses the default generator to create an ID
func Generate() (string, error) {
	return DefaultGenerator.Generate()
}

// GenerateUnique uses the default generator to create a unique ID
func GenerateUnique(existenceChecker func(string) bool) (string, error) {
	return DefaultGenerator.GenerateUnique(existenceChecker)
}
