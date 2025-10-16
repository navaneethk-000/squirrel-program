package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"slices"
)

type JournalEntry struct {
	Events   []string `json:"events"`
	Squirrel bool     `json:"squirrel"`
}

type Counts struct {
	n00 uint
	n01 uint
	n10 uint
	n11 uint
}

type Minmax struct {
	maxValue     float64
	mostPosEvent string
	minValue     float64
	mostNegEvent string
}

func phi(counts Counts) float64 {

	n00 := float64(counts.n00)
	n01 := float64(counts.n01)
	n10 := float64(counts.n10)
	n11 := float64(counts.n11)

	num := float64((n11 * n00) - (n10 * n01))
	n1_ := float64(n11 + n10)
	n0_ := float64(n01 + n00)
	n_1 := float64(n11 + n01)
	n_0 := float64(n10 + n00)

	den := math.Sqrt(float64((n1_ * n0_ * n_1 * n_0)))

	if den == 0 {
		return 0
	}
	return num / den
}

func getCount(entries []JournalEntry, event string) Counts {

	var n11, n01, n10, n00 uint

	for _, entry := range entries {
		if slices.Contains(entry.Events, event) {
			if entry.Squirrel {
				n11++
			} else {
				n10++
			}

		} else {
			if entry.Squirrel {
				n01++
			} else {
				n00++
			}

		}
	}

	counts := Counts{
		n00: n00,
		n11: n11,
		n01: n01,
		n10: n10,
	}

	return counts

}

func getCorrelations(journal []JournalEntry) map[string]float64 {
	correlations := make(map[string]float64)
	for _, entry := range journal {
		for _, e := range entry.Events {
			c := getCount(journal, e)
			correlations[e] = phi(c)

		}
	}

	return correlations
}

func getValues(journal []JournalEntry) Minmax {
	correlations := getCorrelations(journal)

	var mostPosEvent string

	var mostNegEvent string

	maxValue := -1.0

	minValue := 1.0

	for key, value := range correlations {

		if maxValue < value {

			maxValue = value

			mostPosEvent = key

		} else if minValue > value {

			minValue = value

			mostNegEvent = key

		}

	}
	result := Minmax{
		maxValue:     maxValue,
		minValue:     minValue,
		mostPosEvent: mostPosEvent,
		mostNegEvent: mostNegEvent,
	}
	return result

}

func preprocess(journalEntries []JournalEntry) []JournalEntry {
	var journal []JournalEntry
	for _, entry := range journalEntries {
		hasPeanut := slices.Contains(entry.Events, "peanuts")
		notBrushedTeeth := !slices.Contains(entry.Events, "brushed teeth")

		if hasPeanut && notBrushedTeeth {
			entry.Events = append(entry.Events, "dirty teeth")
			journal = append(journal, entry)
		} else {
			journal = append(journal, entry)
		}
	}
	return journal
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

	journal = preprocess(journal)
	minMax := getValues(journal)

	fmt.Println("Most Positive value : ", minMax.maxValue, minMax.mostPosEvent)
	fmt.Println("Most Negative value : ", minMax.minValue, minMax.mostNegEvent)

}
