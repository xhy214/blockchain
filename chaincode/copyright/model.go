package main

type MusicWork struct {
	WorkID      string `json:"workID"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	OwnerID     string `json:"ownerID"`
	FileHash    string `json:"fileHash"`
	RegisterAt  string `json:"registerAt"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
	Status      string `json:"status"` // ACTIVE / TRANSFERRED / DISPUTED
	TxID        string `json:"txID"`
}

type License struct {
	LicenseID   string `json:"licenseID"`
	WorkID      string `json:"workID"`
	GrantorID   string `json:"grantorID"`
	LicenseeID  string `json:"licenseeID"`
	LicenseType string `json:"licenseType"` // COMMERCIAL / NON_COMMERCIAL / EXCLUSIVE
	StartDate   string `json:"startDate"`   // RFC3339
	EndDate     string `json:"endDate"`     // RFC3339
	GrantedAt   string `json:"grantedAt"`
	Status      string `json:"status"` // ACTIVE / REVOKED
	MaxUsage    int    `json:"maxUsage"`
	UsedCount   int    `json:"usedCount"`
}

type DisputeRecord struct {
	DisputeID  string `json:"disputeID"`
	WorkID     string `json:"workID"`
	ClaimantID string `json:"claimantID"`
	Evidence   string `json:"evidence"`
	FiledAt    string `json:"filedAt"`
	Status     string `json:"status"` // PENDING / RESOLVED
}

type VerifyResult struct {
	Valid   bool     `json:"valid"`
	Reason  string   `json:"reason,omitempty"`
	License *License `json:"license,omitempty"`
}
