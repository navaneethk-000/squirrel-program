package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type JournalEntry struct {
	Events   []string `json:"events"`
	Squirrel bool     `json:"squirrel"`
}

func phi(x, y []bool) float64 {
	if len(x) != len(y) {
		return 0
	}
	var n11, n10, n01, n00 float64
	for i := range x {
		switch {
		case x[i] && y[i]:
			n11++
		case x[i] && !y[i]:
			n10++
		case !x[i] && y[i]:
			n01++
		default:
			n00++
		}
	}
	numerator := (n11 * n00) - (n10 * n01)
	denominator := math.Sqrt((n11 + n10) * (n01 + n00) * (n11 + n01) * (n10 + n00))
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

func main() {
	data, err := os.ReadFile("journal.json")
	if err != nil {
		fmt.Println("Error reading json file:", err)
		return
	}

	var journal []JournalEntry
	if err := json.Unmarshal(data, &journal); err != nil {
		fmt.Println("Error during unmarshal: ", err)
		return
	}

	n := len(journal)
	if n == 0 {
		fmt.Println("No entries found.")
		return
	}

	events := map[string]bool{}
	for _, entry := range journal {
		for _, e := range entry.Events {
			events[e] = true
		}
	}

	y := make([]bool, n)
	for i := range journal {
		y[i] = journal[i].Squirrel
	}

	var mostPosEvent, mostNegEvent string
	mostPosVal := -1.0
	mostNegVal := 1.0

	for ev := range events {

		x := make([]bool, n)
		for i := range journal {
			found := false
			for _, e := range journal[i].Events {
				if e == ev {
					found = true
					break
				}
			}
			x[i] = found
		}

		c := phi(x, y)

		if c > mostPosVal {
			mostPosVal = c
			mostPosEvent = ev
		}
		if c < mostNegVal {
			mostNegVal = c
			mostNegEvent = ev
		}
	}

	fmt.Println("\nMost positively correlated event:")
	fmt.Printf("  %s (%.4f)\n", mostPosEvent, mostPosVal)

	fmt.Println("\nMost negatively correlated event:")
	fmt.Printf("  %s (%.4f)\n", mostNegEvent, mostNegVal)
}
