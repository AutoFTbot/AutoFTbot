package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type DonationData struct {
	Amount    int    `json:"amount"`
	Reference string `json:"reference"`
	Status    string `json:"status"`
}

func main() {
	data, err := os.ReadFile("donations.json")
	if err != nil {
		return
	}

	var donations []DonationData
	if err := json.Unmarshal(data, &donations); err != nil {
		return
	}

	// Filter pending donations
	for _, d := range donations {
		if d.Status == "PENDING" {
			// Output as JSON for easy parsing
			json, _ := json.Marshal(d)
			fmt.Println(string(json))
		}
	}
} 