package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
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
	merchantID := os.Getenv("MERCHANT_ID")
	apiKey := os.Getenv("API_KEY")
	baseQrString := os.Getenv("BASE_QR_STRING")

	if merchantID == "" || apiKey == "" || baseQrString == "" {
		fmt.Println("[ERROR] MERCHANT_ID, API_KEY, dan BASE_QR_STRING harus di-set di environment variable!")
		fmt.Printf("MERCHANT_ID: %q\nAPI_KEY: %q\nBASE_QR_STRING: %q\n", merchantID, apiKey, baseQrString)
		os.Exit(1)
	}

	config := qris.QRISConfig{
		MerchantID:   merchantID,
		APIKey:       apiKey,
		BaseQrString: baseQrString,
	}

	fmt.Printf("[INFO] Config: %+v\n", config)

	// Cek error inisialisasi QRIS
	qr, err := qris.NewQRIS(config)
	if err != nil {
		fmt.Println("[ERROR] Failed to initialize QRIS:", err)
		os.Exit(1)
	}

	// Generate amount 5000 + random 1-99
	rand.Seed(time.Now().UnixNano())
	amount := 5000 + rand.Intn(99) + 1
	reference := fmt.Sprintf("QRIS-%d", time.Now().Unix())

	data := qris.QRISData{
		Amount:        int64(amount),
		TransactionID: reference,
	}

	qrString, err := qr.GetQRISString(data)
	if err != nil {
		fmt.Println("[ERROR] Failed to generate QRIS string:", err)
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

	// Pastikan file selalu di root repo
	scriptDir, _ := os.Getwd()
	rootDir := filepath.Dir(scriptDir)
	donationsPath := filepath.Join(rootDir, "donations.json")
	qrisPngPath := filepath.Join(rootDir, "qris.png")

	file, err := os.Create(donationsPath)
	if err != nil {
		fmt.Println("[ERROR] Failed to create donations.json:", err)
		os.Exit(1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode([]DonationData{donation}); err != nil {
		fmt.Println("[ERROR] Failed to write donations.json:", err)
		os.Exit(1)
	}
	fmt.Println("[SUCCESS] donations.json written at", donationsPath)

	// Generate QR code PNG di root
	qrCode, err := qr.GenerateQRCode(data)
	if err != nil {
		fmt.Println("[ERROR] Failed to generate QR code image:", err)
		os.Exit(1)
	}
	if err := qrCode.WriteFile(256, qrisPngPath); err != nil {
		fmt.Println("[ERROR] Failed to save QR code image:", err)
		os.Exit(1)
	}
	fmt.Println("[SUCCESS] QR code image saved as", qrisPngPath)
}
