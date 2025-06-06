package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"

	"github.com/AutoFTbot/OrderKuota-go/qris"
)

type DonationData struct {
	Amount      int    `json:"amount"`
	Reference   string `json:"reference"`
	Status      string `json:"status"`
	Date        string `json:"date"`
	QRString    string `json:"qr_string"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	amount := rand.Intn(99) + 1 // 1-99
	reference := "QRIS-" + time.Now().Format("20060102150405")

	config := qris.QRISConfig{
		MerchantID:   os.Getenv("MERCHANT_ID"),
		APIKey:       os.Getenv("API_KEY"),
		BaseQrString: os.Getenv("BASE_QR_STRING"),
	}
	qr, _ := qris.NewQRIS(config)
	data := qris.QRISData{
		Amount:      int64(amount),
		TransactionID: reference,
	}
	qrString, _ := qr.GetQRISString(data)

	donation := DonationData{
		Amount:    amount,
		Reference: reference,
		Status:    "PENDING",
		Date:      time.Now().Format(time.RFC3339),
		QRString:  qrString,
	}

	// Simpan ke donations.json (replace atau append, sesuai kebutuhan)
	file, _ := os.Create("donations.json")
	defer file.Close()
	json.NewEncoder(file).Encode([]DonationData{donation})
}
