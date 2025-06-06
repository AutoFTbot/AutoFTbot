package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AutoFTbot/OrderKuota-go/qris"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("false")
		return
	}

	reference := os.Args[1]
	amount, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("false")
		return
	}

	// Initialize QRIS
	config := qris.QRISConfig{
		MerchantID:   os.Getenv("MERCHANT_ID"),
		APIKey:       os.Getenv("API_KEY"),
		BaseQrString: os.Getenv("BASE_QR_STRING"),
	}

	qr, err := qris.NewQRIS(config)
	if err != nil {
		fmt.Println("false")
		return
	}

	// Check payment status
	status, err := qr.CheckPaymentStatus(reference, int64(amount))
	if err != nil {
		fmt.Println("false")
		return
	}

	if status.Status == "PAID" {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
} 