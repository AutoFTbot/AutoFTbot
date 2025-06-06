package donation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AutoFTbot/OrderKuota-go/qris"
)

// DonationData represents a donation record
type DonationData struct {
	Amount      int       `json:"amount"`
	Reference   string    `json:"reference"`
	Status      string    `json:"status"`
	Date        time.Time `json:"date"`
	DonorName   string    `json:"donor_name,omitempty"`
	Message     string    `json:"message,omitempty"`
}

// Manager handles donation operations
type Manager struct {
	qr            *qris.QRIS
	donationsFile string
}

// NewManager creates a new donation manager
func NewManager(qr *qris.QRIS, donationsFile string) *Manager {
	return &Manager{
		qr:            qr,
		donationsFile: donationsFile,
	}
}

// GenerateQR generates a QR code for donation
func (m *Manager) GenerateQR(amount int, reference string) ([]byte, error) {
	data := qris.QRISData{
		Amount:        int64(amount),
		TransactionID: reference,
	}
	return m.qr.GenerateQRCode(data)
}

// CheckStatus checks if a donation has been paid
func (m *Manager) CheckStatus(reference string, amount int) (bool, error) {
	status, err := m.qr.CheckPaymentStatus(reference, int64(amount))
	if err != nil {
		return false, err
	}
	return status.Status == "PAID", nil
}

// Save saves a donation record
func (m *Manager) Save(donation DonationData) error {
	// Create donations directory if it doesn't exist
	dir := filepath.Dir(m.donationsFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create donations directory: %v", err)
	}

	// Read existing donations
	var donations []DonationData
	if data, err := os.ReadFile(m.donationsFile); err == nil {
		if err := json.Unmarshal(data, &donations); err != nil {
			return fmt.Errorf("failed to parse donations file: %v", err)
		}
	}

	// Add new donation
	donations = append(donations, donation)

	// Save updated donations
	data, err := json.MarshalIndent(donations, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal donations: %v", err)
	}

	if err := os.WriteFile(m.donationsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to save donations: %v", err)
	}

	return nil
}

// GetAll retrieves all donations
func (m *Manager) GetAll() ([]DonationData, error) {
	data, err := os.ReadFile(m.donationsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []DonationData{}, nil
		}
		return nil, fmt.Errorf("failed to read donations file: %v", err)
	}

	var donations []DonationData
	if err := json.Unmarshal(data, &donations); err != nil {
		return nil, fmt.Errorf("failed to parse donations file: %v", err)
	}

	return donations, nil
}

// GetTotal calculates total amount of successful donations
func (m *Manager) GetTotal() (int, error) {
	donations, err := m.GetAll()
	if err != nil {
		return 0, err
	}

	total := 0
	for _, d := range donations {
		if d.Status == "PAID" {
			total += d.Amount
		}
	}

	return total, nil
} 