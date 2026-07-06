package model

import "time"

type CopyrightCache struct {
	WorkID    string    `gorm:"primaryKey;size:36" json:"workID"`
	Title     string    `gorm:"size:256" json:"title"`
	Artist    string    `gorm:"size:128" json:"artist"`
	OwnerID   string    `gorm:"size:36;index:idx_owner" json:"ownerID"`
	FileHash  string    `gorm:"size:64" json:"fileHash"`
	Genre     string    `gorm:"size:64" json:"genre"`
	TxID      string    `gorm:"size:128" json:"txID"`
	CreatedAt time.Time `json:"createdAt"`
}
