package service

import (
	"encoding/json"
	"fmt"
)

type VerifyService struct {
	Fabric *FabricClient
}

type VerifyResultDTO struct {
	Valid   bool        `json:"valid"`
	Reason  string      `json:"reason,omitempty"`
	License *LicenseDTO `json:"license,omitempty"`
}

func (s *VerifyService) VerifyLicense(workID, licenseeID string) (*VerifyResultDTO, error) {
	data, err := s.Fabric.Evaluate("VerifyLicense", workID, licenseeID)
	if err != nil {
		return nil, fmt.Errorf("chaincode VerifyLicense: %w", err)
	}
	var result VerifyResultDTO
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
