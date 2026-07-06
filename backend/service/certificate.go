package service

import (
	"encoding/json"
	"fmt"

	"blockchain/backend/utils"
)

type CertificateService struct {
	Fabric *FabricClient
}

func (s *CertificateService) GenerateCertificate(workID string) ([]byte, error) {
	raw, err := s.Fabric.Evaluate("QueryWork", workID)
	if err != nil {
		return nil, fmt.Errorf("query work for certificate: %w", err)
	}

	var work struct {
		WorkID      string `json:"workID"`
		Title       string `json:"title"`
		Artist      string `json:"artist"`
		OwnerID     string `json:"ownerID"`
		FileHash    string `json:"fileHash"`
		RegisterAt  string `json:"registerAt"`
		Description string `json:"description"`
		Genre       string `json:"genre"`
		Status      string `json:"status"`
		TxID        string `json:"txID"`
	}
	if err := json.Unmarshal(raw, &work); err != nil {
		return nil, fmt.Errorf("unmarshal work: %w", err)
	}

	return utils.GenerateCertificatePDF(&utils.CertData{
		WorkID:     work.WorkID,
		Title:      work.Title,
		Artist:     work.Artist,
		OwnerName:  work.OwnerID,
		FileHash:   work.FileHash,
		RegisterAt: work.RegisterAt,
		TxID:       work.TxID,
	})
}
