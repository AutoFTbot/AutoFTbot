package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type DonationData struct {
	Amount int    `json:"amount"`
	Status string `json:"status"`
}

func main() {
	// Pastikan path ke donations.json benar (di root repo)
	scriptDir, _ := os.Getwd()
	rootDir := filepath.Dir(scriptDir)
	donationsPath := filepath.Join(rootDir, "donations.json")

	data, err := os.ReadFile(donationsPath)
	if err != nil {
		fmt.Print("0")
		return
	}

	var donations []DonationData
	if err := json.Unmarshal(data, &donations); err != nil {
		fmt.Print("0")
		return
	}

	total := 0
	for _, d := range donations {
		if d.Status == "PAID" {
			total += d.Amount
		}
	}

	fmt.Print(total)
} 
