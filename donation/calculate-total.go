package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type DonationData struct {
	Amount int    `json:"amount"`
	Status string `json:"status"`
}

func main() {
	data, err := os.ReadFile("donations.json")
	if err != nil {
		fmt.Println("0")
		return
	}

	var donations []DonationData
	if err := json.Unmarshal(data, &donations); err != nil {
		fmt.Println("0")
		return
	}

	total := 0
	for _, d := range donations {
		if d.Status == "PAID" {
			total += d.Amount
		}
	}

	fmt.Println(total)
} 