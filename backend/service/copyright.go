package service

import (
	"encoding/json"
	"fmt"
	"strconv"

	"blockchain/backend/model"

	"github.com/google/uuid"
)

type CopyrightService struct {
	Fabric *FabricClient
}

type MusicWorkDTO struct {
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

type HistoryRecord struct {
	TxID      string        `json:"txID"`
	Timestamp string        `json:"timestamp"`
	IsDelete  bool          `json:"isDelete"`
	Value     *MusicWorkDTO `json:"value,omitempty"`
}

type SearchResult struct {
	Total int             `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
	List  []*MusicWorkDTO `json:"list"`
}

type VerifyHashResult struct {
	Match        bool   `json:"match"`
	OnChainHash  string `json:"onChainHash"`
	UploadedHash string `json:"uploadedHash"`
}

func (s *CopyrightService) RegisterWork(ownerID, title, artist, fileHash, description, genre string) (*MusicWorkDTO, error) {
	workID := uuid.New().String()

	_, err := s.Fabric.Submit("RegisterWork",
		workID, title, artist, ownerID, fileHash, description, genre)
	if err != nil {
		return nil, fmt.Errorf("chaincode RegisterWork: %w", err)
	}

	// Query the newly created work from chain
	work, err := s.QueryWork(workID)
	if err != nil {
		return nil, err
	}

	// Cache to MySQL for search acceleration
	model.DB.Create(&model.CopyrightCache{
		WorkID:   workID,
		Title:    title,
		Artist:   artist,
		OwnerID:  ownerID,
		FileHash: fileHash,
		Genre:    genre,
		TxID:     work.TxID,
	})

	return work, nil
}

func (s *CopyrightService) QueryWork(workID string) (*MusicWorkDTO, error) {
	data, err := s.Fabric.Evaluate("QueryWork", workID)
	if err != nil {
		return nil, fmt.Errorf("chaincode QueryWork: %w", err)
	}
	var work MusicWorkDTO
	if err := json.Unmarshal(data, &work); err != nil {
		return nil, err
	}
	return &work, nil
}

func (s *CopyrightService) QueryWorksByOwner(ownerID string) ([]*MusicWorkDTO, error) {
	data, err := s.Fabric.Evaluate("QueryWorksByOwner", ownerID)
	if err != nil {
		return nil, fmt.Errorf("chaincode QueryWorksByOwner: %w", err)
	}
	var works []*MusicWorkDTO
	if err := json.Unmarshal(data, &works); err != nil {
		return nil, err
	}
	if works == nil {
		works = []*MusicWorkDTO{}
	}
	return works, nil
}

func (s *CopyrightService) SearchWorks(keyword string, page, size int) (*SearchResult, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 10
	}
	skip := (page - 1) * size

	data, err := s.Fabric.Evaluate("SearchWorks", keyword, strconv.Itoa(size), strconv.Itoa(skip))
	if err != nil {
		return nil, fmt.Errorf("chaincode SearchWorks: %w", err)
	}
	var works []*MusicWorkDTO
	if err := json.Unmarshal(data, &works); err != nil {
		return nil, err
	}
	if works == nil {
		works = []*MusicWorkDTO{}
	}
	return &SearchResult{
		Total: len(works),
		Page:  page,
		Size:  size,
		List:  works,
	}, nil
}

func (s *CopyrightService) VerifyFileHash(workID, fileHash string) (*VerifyHashResult, error) {
	data, err := s.Fabric.Evaluate("VerifyFileHash", workID, fileHash)
	if err != nil {
		return nil, fmt.Errorf("chaincode VerifyFileHash: %w", err)
	}
	var match bool
	json.Unmarshal(data, &match)

	work, _ := s.QueryWork(workID)
	onChainHash := ""
	if work != nil {
		onChainHash = work.FileHash
	}
	return &VerifyHashResult{
		Match:        match,
		OnChainHash:  onChainHash,
		UploadedHash: fileHash,
	}, nil
}

func (s *CopyrightService) TransferCopyright(workID, fromID, toID string) (*MusicWorkDTO, error) {
	_, err := s.Fabric.Submit("TransferCopyright", workID, fromID, toID)
	if err != nil {
		return nil, fmt.Errorf("chaincode TransferCopyright: %w", err)
	}
	model.DB.Model(&model.CopyrightCache{}).Where("work_id = ?", workID).Update("owner_id", toID)
	return s.QueryWork(workID)
}

func (s *CopyrightService) GetHistory(workID string) ([]*HistoryRecord, error) {
	data, err := s.Fabric.Evaluate("GetHistory", workID)
	if err != nil {
		return nil, fmt.Errorf("chaincode GetHistory: %w", err)
	}
	var records []*HistoryRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}
	if records == nil {
		records = []*HistoryRecord{}
	}
	return records, nil
}
