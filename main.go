package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

// MarkovChain stores the transitions for a character-level Markov chain.
type MarkovChain struct {
	// transitions maps a state (string) to all possible next runes.
	transitions map[string][]rune
	order       int
}

// NewMarkovChain initializes a MarkovChain of the specified order.
func NewMarkovChain(order int) *MarkovChain {
	return &MarkovChain{
		transitions: make(map[string][]rune),
		order:       order,
	}
}

// AddText processes the given text to populate the transitions map.
func (mc *MarkovChain) AddText(text string) {
	// If the text is shorter than the order, nothing to process
	if len(text) <= mc.order {
		return
	}

	// Build transitions by sliding over the text
	for i := 0; i < len(text)-mc.order; i++ {
		// Current state is the substring of length 'order'
		state := text[i : i+mc.order]
		// The next character after this state
		nextChar := rune(text[i+mc.order])
		mc.transitions[state] = append(mc.transitions[state], nextChar)
	}
}

// Generate produces 'length' characters of text using the Markov chain,
// optionally starting with a given 'starter' string. If the starter is
// longer than 'length', it will be truncated to fit. The total output
// will always be exactly 'length' characters (if enough transitions exist).
func (mc *MarkovChain) Generate(length int, seed int64, starter string) string {
	if length <= 0 {
		return ""
	}

	// Seed the random number generator
	if seed < 0 {
		rand.Seed(time.Now().UnixNano())
	} else {
		rand.Seed(seed)
	}

	// If we have no transitions, there's nothing to generate.
	if len(mc.transitions) == 0 {
		// Return just the truncated starter, if any.
		if len(starter) > length {
			return starter[:length]
		}
		return starter
	}

	// Prepare a builder for the final output
	var result strings.Builder

	// If the starter text is already >= length, just truncate and return it.
	if len(starter) >= length {
		return starter[:length]
	}

	// Otherwise, we add the entire starter to the result
	result.WriteString(starter)

	// We'll generate enough characters to reach 'length' total
	needed := length - len(starter)

	// Compute the initial state from the starter, if possible
	var currentState string
	if len(starter) >= mc.order {
		// Use the last 'order' characters of starter
		currentState = starter[len(starter)-mc.order:]
	} else {
		// If not enough characters in the starter, pick a random state
		var states []string
		for state := range mc.transitions {
			states = append(states, state)
		}
		currentState = states[rand.Intn(len(states))]
		// Also append currentState if we don't have a starter,
		// but that would count toward the result. For simplicity,
		// we won't add it in the result right now, because we
		// are continuing from the starter. We'll just treat
		// "missing characters" as if they never existed.
	}

	// Now generate the remaining characters
	for i := 0; i < needed; i++ {
		// Possible next runes from currentState
		nextRunes := mc.transitions[currentState]
		if len(nextRunes) == 0 {
			// No known transitions from this state, pick a random new one
			var states []string
			for s := range mc.transitions {
				states = append(states, s)
			}
			currentState = states[rand.Intn(len(states))]
			// Write currentState to continue generation
			// but we only want to write one character to the result, not the entire state.
			// We'll pick a single random nextChar from that new state's transitions, if possible.
			nextRunes = mc.transitions[currentState]
			if len(nextRunes) == 0 {
				// If even this new state has no transitions, we're stuck
				break
			}
		}
		nextChar := nextRunes[rand.Intn(len(nextRunes))]
		result.WriteRune(nextChar)

		// Update currentState by dropping the first character and adding the new one
		if mc.order > 1 {
			if len(currentState) > 0 {
				currentState = currentState[1:] + string(nextChar)
			} else {
				// If for some reason currentState is empty, just set to new char
				currentState = string(nextChar)
			}
		} else {
			currentState = string(nextChar)
		}
	}

	return result.String()
}

func main() {
	// Define command-line flags
	k := flag.Int("k", 1, "Order of the Markov chain")
	l := flag.Int("l", 100, "Number of characters to generate (total output length)")
	inputFile := flag.String("i", "", "Input file (optional, reads from stdin if not provided)")
	seedFlag := flag.Int64("seed", -1, "Random seed (optional, defaults to current time if not provided)")
	starter := flag.String("starter", "", "Starter text to prepend to the output")
	flag.Parse()

	// Read the input text from file or stdin
	var reader io.Reader
	if *inputFile != "" {
		f, err := os.Open(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", *inputFile, err)
			os.Exit(1)
		}
		defer f.Close()
		reader = f
	} else {
		// Read from stdin
		reader = os.Stdin
	}

	// Capture the entire text
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanBytes) // we want to scan by character
	var builder strings.Builder

	for scanner.Scan() {
		builder.WriteString(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	text := builder.String()

	// Build the Markov Chain
	mc := NewMarkovChain(*k)
	mc.AddText(text)

	// Generate the output
	output := mc.Generate(*l, *seedFlag, *starter)
	fmt.Println(output)
}
