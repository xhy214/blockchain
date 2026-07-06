package service

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type LicenseService struct {
	Fabric *FabricClient
}

type LicenseDTO struct {
	LicenseID   string `json:"licenseID"`
	WorkID      string `json:"workID"`
	GrantorID   string `json:"grantorID"`
	LicenseeID  string `json:"licenseeID"`
	LicenseType string `json:"licenseType"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	GrantedAt   string `json:"grantedAt"`
	Status      string `json:"status"`
	MaxUsage    int    `json:"maxUsage"`
	UsedCount   int    `json:"usedCount"`
}

func (s *LicenseService) GrantLicense(workID, grantorID, licenseeID, licenseType, startDate, endDate string, maxUsage int) (*LicenseDTO, error) {
	licenseID := "lic-" + uuid.New().String()[:8]

	_, err := s.Fabric.Submit("GrantLicense",
		licenseID, workID, grantorID, licenseeID, licenseType,
		startDate, endDate, fmt.Sprintf("%d", maxUsage))
	if err != nil {
		return nil, fmt.Errorf("chaincode GrantLicense: %w", err)
	}
	return s.QueryLicense(licenseID)
}

func (s *LicenseService) QueryLicense(licenseID string) (*LicenseDTO, error) {
	data, err := s.Fabric.Evaluate("QueryLicense", licenseID)
	if err != nil {
		return nil, fmt.Errorf("chaincode QueryLicense: %w", err)
	}
	var lic LicenseDTO
	if err := json.Unmarshal(data, &lic); err != nil {
		return nil, err
	}
	return &lic, nil
}

func (s *LicenseService) QueryLicensesByWork(workID string) ([]*LicenseDTO, error) {
	data, err := s.Fabric.Evaluate("QueryLicensesByWork", workID)
	if err != nil {
		return nil, fmt.Errorf("chaincode QueryLicensesByWork: %w", err)
	}
	var licenses []*LicenseDTO
	if err := json.Unmarshal(data, &licenses); err != nil {
		return nil, err
	}
	if licenses == nil {
		licenses = []*LicenseDTO{}
	}
	return licenses, nil
}

func (s *LicenseService) QueryLicensesByLicensee(licenseeID string) ([]*LicenseDTO, error) {
	data, err := s.Fabric.Evaluate("QueryLicensesByLicensee", licenseeID)
	if err != nil {
		return nil, fmt.Errorf("chaincode QueryLicensesByLicensee: %w", err)
	}
	var licenses []*LicenseDTO
	if err := json.Unmarshal(data, &licenses); err != nil {
		return nil, err
	}
	if licenses == nil {
		licenses = []*LicenseDTO{}
	}
	return licenses, nil
}

func (s *LicenseService) RevokeLicense(licenseID, revokerID string) error {
	_, err := s.Fabric.Submit("RevokeLicense", licenseID, revokerID)
	if err != nil {
		return fmt.Errorf("chaincode RevokeLicense: %w", err)
	}
	return nil
}

func (s *LicenseService) RecordUsage(licenseID string) error {
	_, err := s.Fabric.Submit("RecordUsage", licenseID)
	if err != nil {
		return fmt.Errorf("chaincode RecordUsage: %w", err)
	}
	return nil
}
