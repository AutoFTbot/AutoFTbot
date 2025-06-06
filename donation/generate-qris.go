package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/AutoFTbot/OrderKuota-go/qris"
)

// DonationData struct untuk menyimpan data donasi
type DonationData struct {
	Amount     int    `json:"amount"`
	Reference  string `json:"reference"`
	Status     string `json:"status"`
	Date       string `json:"date"`
	QRString   string `json:"qr_string"`
}

func main() {
	// Ambil konfigurasi dari environment variable
	config := qris.QRISConfig{
		MerchantID:   os.Getenv("MERCHANT_ID"),
		APIKey:       os.Getenv("API_KEY"),
		BaseQrString: os.Getenv("BASE_QR_STRING"),
	}

	// Cek error inisialisasi QRIS
	qr, err := qris.NewQRIS(config)
	if err != nil {
		fmt.Println("Failed to initialize QRIS:", err)
		os.Exit(1)
	}

	// Generate amount random 1-99
	rand.Seed(time.Now().UnixNano())
	amount := rand.Intn(99) + 1
	reference := fmt.Sprintf("QRIS-%d", time.Now().Unix())

	data := qris.QRISData{
		Amount:        int64(amount),
		TransactionID: reference,
	}

	qrString, err := qr.GetQRISString(data)
	if err != nil {
		fmt.Println("Failed to generate QRIS string:", err)
		os.Exit(1)
	}

	// Buat data donasi baru
	donation := DonationData{
		Amount:    amount,
		Reference: reference,
		Status:    "PENDING",
		Date:      time.Now().Format(time.RFC3339),
		QRString:  qrString,
	}

	// Simpan ke donations.json (replace seluruh isi dengan donasi baru)
	file, err := os.Create("donations.json")
	if err != nil {
		fmt.Println("Failed to create donations.json:", err)
		os.Exit(1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode([]DonationData{donation}); err != nil {
		fmt.Println("Failed to write donations.json:", err)
		os.Exit(1)
	}

	fmt.Println("QRIS generated successfully!")
}
