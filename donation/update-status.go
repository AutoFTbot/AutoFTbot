package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type DonationData struct {
	Amount      int    `json:"amount"`
	Reference   string `json:"reference"`
	Status      string `json:"status"`
	Date        string `json:"date"`
	DonorName   string `json:"donor_name,omitempty"`
	Message     string `json:"message,omitempty"`
}

func main() {
	if len(os.Args) != 3 {
		return
	}

	reference := os.Args[1]
	newStatus := os.Args[2]

	scriptDir, _ := os.Getwd()
	rootDir := filepath.Dir(scriptDir)
	donationsPath := filepath.Join(rootDir, "donations.json")

	// Read donations file
	data, err := os.ReadFile(donationsPath)
	if err != nil {
		return
	}

	var donations []DonationData
	if err := json.Unmarshal(data, &donations); err != nil {
		return
	}

	// Update status
	for i, d := range donations {
		if d.Reference == reference {
			donations[i].Status = newStatus
			break
		}
	}

	// Save updated donations
	data, err = json.MarshalIndent(donations, "", "  ")
	if err != nil {
		return
	}

	if err := os.WriteFile(donationsPath, data, 0644); err != nil {
		return
	}

	// Print updated donation as JSON (for notification)
	for _, d := range donations {
		if d.Reference == reference {
			json, _ := json.Marshal(d)
			fmt.Println(string(json))
			break
		}
	}
} 
