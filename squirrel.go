package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"slices"
)

type JournalEntries struct {
	Events   []string `json:"events"`
	Squirrel bool     `json:"squirrel"`
}

type Counts struct {
	n00 uint
	n01 uint
	n10 uint
	n11 uint
}

func phi(counts Counts) float64 {

	n00 := float64(counts.n00)
	n01 := float64(counts.n01)
	n10 := float64(counts.n10)
	n11 := float64(counts.n11)

	num := float64((n11 * n00) - (n10 * n01))
	n1 := float64(n11 + n10)
	n0 := float64(n01 + n00)
	n1_ := float64(n11 + n01)
	n0_ := float64(n10 + n00)

	den := math.Sqrt(float64((n1 * n0 * n1_ * n0_)))

	if den == 0 {
		return 0
	}
	return num / den
}

func getCount(entries []JournalEntries, event string) Counts {

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

func main() {

	data, err := os.ReadFile("journal.json")
	if err != nil {
		fmt.Println("Error reading json file:", err)
		return
	}

	var journal []JournalEntries
	if err := json.Unmarshal(data, &journal); err != nil {
		fmt.Println("Error during unmarshal: ", err)
		return
	}

	c := getCount(journal, "carrot")

	correlation := make(map[string]float64)
	for _, entry := range journal {
		for _, e := range entry.Events {
			c = getCount(journal, e)
			correlation[e] = phi(c)

		}
	}

	var mostPosEvent string
	var mostNegEvent string
	var mostNegVal float64
	var mostPosVal float64

	for key, value := range correlation {
		if mostPosVal < value {
			mostPosVal = value
			mostPosEvent = key
		} else if mostNegVal > value {
			mostNegVal = value
			mostNegEvent = key

		}

	}
	fmt.Println("Most Positive value : ", mostNegVal, mostPosEvent)
	fmt.Println("Most Negative value : ", mostPosVal, mostNegEvent)

}
