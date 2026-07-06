package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type CopyrightContract struct {
	contractapi.Contract
}

// ─── 版权存证 ───────────────────────────────────────────────

func (c *CopyrightContract) RegisterWork(ctx contractapi.TransactionContextInterface,
	workID, title, artist, ownerID, fileHash, description, genre string) error {

	existing, err := ctx.GetStub().GetState("WORK_" + workID)
	if err != nil {
		return fmt.Errorf("failed to read state: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("work %s already exists", workID)
	}

	ts, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get tx timestamp: %w", err)
	}
	registerAt := time.Unix(ts.Seconds, 0).UTC().Format(time.RFC3339)

	work := MusicWork{
		WorkID:      workID,
		Title:       title,
		Artist:      artist,
		OwnerID:     ownerID,
		FileHash:    fileHash,
		RegisterAt:  registerAt,
		Description: description,
		Genre:       genre,
		Status:      "ACTIVE",
		TxID:        ctx.GetStub().GetTxID(),
	}

	data, err := json.Marshal(work)
	if err != nil {
		return err
	}
	if err := ctx.GetStub().PutState("WORK_"+workID, data); err != nil {
		return err
	}

	// 复合键：owner~work，支持按 ownerID 范围查询
	ck, err := ctx.GetStub().CreateCompositeKey("owner~work", []string{ownerID, workID})
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(ck, []byte{0x00})
}

func (c *CopyrightContract) QueryWork(ctx contractapi.TransactionContextInterface, workID string) (*MusicWork, error) {
	data, err := ctx.GetStub().GetState("WORK_" + workID)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("work %s not found", workID)
	}
	var work MusicWork
	return &work, json.Unmarshal(data, &work)
}

// QueryWorksByOwner 用复合键查询某用户的全部作品（不依赖 CouchDB 富查询）
func (c *CopyrightContract) QueryWorksByOwner(ctx contractapi.TransactionContextInterface, ownerID string) ([]*MusicWork, error) {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey("owner~work", []string{ownerID})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var works []*MusicWork
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			return nil, err
		}
		_, parts, err := ctx.GetStub().SplitCompositeKey(item.Key)
		if err != nil || len(parts) < 2 {
			continue
		}
		workID := parts[1]
		work, err := c.QueryWork(ctx, workID)
		if err != nil {
			continue
		}
		works = append(works, work)
	}
	return works, nil
}

// SearchWorks CouchDB 富查询（按 title/artist 模糊匹配）
// limit 和 skip 作为字符串参数传入，链码内部转换
func (c *CopyrightContract) SearchWorks(ctx contractapi.TransactionContextInterface,
	keyword, limitStr, skipStr string) ([]*MusicWork, error) {

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	skip, err := strconv.Atoi(skipStr)
	if err != nil || skip < 0 {
		skip = 0
	}

	query := fmt.Sprintf(`{
		"selector": {
			"$or": [
				{"title":  {"$regex": "(?i)%s"}},
				{"artist": {"$regex": "(?i)%s"}}
			]
		},
		"limit": %d,
		"skip":  %d
	}`, keyword, keyword, limit, skip)

	iter, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var works []*MusicWork
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			return nil, err
		}
		var w MusicWork
		if err := json.Unmarshal(item.Value, &w); err != nil {
			continue
		}
		works = append(works, &w)
	}
	return works, nil
}

// VerifyFileHash 验证文件哈希与链上记录是否一致
func (c *CopyrightContract) VerifyFileHash(ctx contractapi.TransactionContextInterface,
	workID, fileHash string) (bool, error) {

	work, err := c.QueryWork(ctx, workID)
	if err != nil {
		return false, err
	}
	return work.FileHash == fileHash, nil
}

// ─── 版权转让 ───────────────────────────────────────────────

func (c *CopyrightContract) TransferCopyright(ctx contractapi.TransactionContextInterface,
	workID, fromID, toID string) error {

	work, err := c.QueryWork(ctx, workID)
	if err != nil {
		return err
	}
	if work.OwnerID != fromID {
		return fmt.Errorf("permission denied: %s is not the owner", fromID)
	}

	// 删除旧复合键，建新复合键
	oldCK, _ := ctx.GetStub().CreateCompositeKey("owner~work", []string{fromID, workID})
	ctx.GetStub().DelState(oldCK)
	newCK, _ := ctx.GetStub().CreateCompositeKey("owner~work", []string{toID, workID})
	ctx.GetStub().PutState(newCK, []byte{0x00})

	work.OwnerID = toID
	work.Status = "TRANSFERRED"

	data, err := json.Marshal(work)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("WORK_"+workID, data)
}

// ─── 授权管理 ───────────────────────────────────────────────

func (c *CopyrightContract) GrantLicense(ctx contractapi.TransactionContextInterface,
	licenseID, workID, grantorID, licenseeID, licenseType, startDate, endDate, maxUsageStr string) error {

	work, err := c.QueryWork(ctx, workID)
	if err != nil {
		return err
	}
	if work.OwnerID != grantorID {
		return fmt.Errorf("permission denied: %s is not the owner", grantorID)
	}

	existing, _ := ctx.GetStub().GetState("LICENSE_" + licenseID)
	if existing != nil {
		return fmt.Errorf("license %s already exists", licenseID)
	}

	maxUsage, err := strconv.Atoi(maxUsageStr)
	if err != nil {
		maxUsage = 0
	}

	ts, _ := ctx.GetStub().GetTxTimestamp()
	grantedAt := time.Unix(ts.Seconds, 0).UTC().Format(time.RFC3339)

	lic := License{
		LicenseID:   licenseID,
		WorkID:      workID,
		GrantorID:   grantorID,
		LicenseeID:  licenseeID,
		LicenseType: licenseType,
		StartDate:   startDate,
		EndDate:     endDate,
		GrantedAt:   grantedAt,
		Status:      "ACTIVE",
		MaxUsage:    maxUsage,
		UsedCount:   0,
	}

	data, err := json.Marshal(lic)
	if err != nil {
		return err
	}
	if err := ctx.GetStub().PutState("LICENSE_"+licenseID, data); err != nil {
		return err
	}

	// 复合键：work~license 和 licensee~license
	wck, _ := ctx.GetStub().CreateCompositeKey("work~license", []string{workID, licenseID})
	ctx.GetStub().PutState(wck, []byte{0x00})
	lck, _ := ctx.GetStub().CreateCompositeKey("licensee~license", []string{licenseeID, licenseID})
	return ctx.GetStub().PutState(lck, []byte{0x00})
}

// VerifyLicense 核验某用户对某作品是否持有当前有效授权
func (c *CopyrightContract) VerifyLicense(ctx contractapi.TransactionContextInterface,
	workID, licenseeID string) (*VerifyResult, error) {

	iter, err := ctx.GetStub().GetStateByPartialCompositeKey("work~license", []string{workID})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	now := time.Now().UTC()

	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			continue
		}
		_, parts, err := ctx.GetStub().SplitCompositeKey(item.Key)
		if err != nil || len(parts) < 2 {
			continue
		}
		licenseID := parts[1]

		data, err := ctx.GetStub().GetState("LICENSE_" + licenseID)
		if err != nil || data == nil {
			continue
		}
		var lic License
		if err := json.Unmarshal(data, &lic); err != nil {
			continue
		}

		if lic.LicenseeID != licenseeID {
			continue
		}
		if lic.Status != "ACTIVE" {
			continue
		}

		start, err := time.Parse(time.RFC3339, lic.StartDate)
		if err != nil {
			continue
		}
		end, err := time.Parse(time.RFC3339, lic.EndDate)
		if err != nil {
			continue
		}
		if now.Before(start) {
			return &VerifyResult{Valid: false, Reason: "授权尚未生效"}, nil
		}
		if now.After(end) {
			return &VerifyResult{Valid: false, Reason: "授权已过期"}, nil
		}
		if lic.MaxUsage > 0 && lic.UsedCount >= lic.MaxUsage {
			return &VerifyResult{Valid: false, Reason: "使用次数已达上限"}, nil
		}

		return &VerifyResult{Valid: true, License: &lic}, nil
	}

	return &VerifyResult{Valid: false, Reason: "无有效授权"}, nil
}

func (c *CopyrightContract) QueryLicense(ctx contractapi.TransactionContextInterface, licenseID string) (*License, error) {
	data, err := ctx.GetStub().GetState("LICENSE_" + licenseID)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("license %s not found", licenseID)
	}
	var lic License
	return &lic, json.Unmarshal(data, &lic)
}

func (c *CopyrightContract) QueryLicensesByWork(ctx contractapi.TransactionContextInterface, workID string) ([]*License, error) {
	return c.queryLicensesByIndex(ctx, "work~license", workID)
}

func (c *CopyrightContract) QueryLicensesByLicensee(ctx contractapi.TransactionContextInterface, licenseeID string) ([]*License, error) {
	return c.queryLicensesByIndex(ctx, "licensee~license", licenseeID)
}

func (c *CopyrightContract) queryLicensesByIndex(ctx contractapi.TransactionContextInterface, indexName, key string) ([]*License, error) {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(indexName, []string{key})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var licenses []*License
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			return nil, err
		}
		_, parts, err := ctx.GetStub().SplitCompositeKey(item.Key)
		if err != nil || len(parts) < 2 {
			continue
		}
		lic, err := c.QueryLicense(ctx, parts[1])
		if err != nil {
			continue
		}
		licenses = append(licenses, lic)
	}
	return licenses, nil
}

func (c *CopyrightContract) RevokeLicense(ctx contractapi.TransactionContextInterface,
	licenseID, revokerID string) error {

	lic, err := c.QueryLicense(ctx, licenseID)
	if err != nil {
		return err
	}
	if lic.GrantorID != revokerID {
		return fmt.Errorf("permission denied: only grantor can revoke")
	}
	lic.Status = "REVOKED"
	data, err := json.Marshal(lic)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("LICENSE_"+licenseID, data)
}

func (c *CopyrightContract) RecordUsage(ctx contractapi.TransactionContextInterface, licenseID string) error {
	lic, err := c.QueryLicense(ctx, licenseID)
	if err != nil {
		return err
	}
	if lic.Status != "ACTIVE" {
		return fmt.Errorf("license is not active")
	}
	if lic.MaxUsage > 0 && lic.UsedCount >= lic.MaxUsage {
		return fmt.Errorf("usage limit reached")
	}
	lic.UsedCount++
	data, err := json.Marshal(lic)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("LICENSE_"+licenseID, data)
}

// ─── 版权争议 ───────────────────────────────────────────────

func (c *CopyrightContract) FileDispute(ctx contractapi.TransactionContextInterface,
	disputeID, workID, claimantID, evidence string) error {

	if _, err := c.QueryWork(ctx, workID); err != nil {
		return err
	}

	ts, _ := ctx.GetStub().GetTxTimestamp()
	filedAt := time.Unix(ts.Seconds, 0).UTC().Format(time.RFC3339)

	dispute := DisputeRecord{
		DisputeID:  disputeID,
		WorkID:     workID,
		ClaimantID: claimantID,
		Evidence:   evidence,
		FiledAt:    filedAt,
		Status:     "PENDING",
	}
	data, err := json.Marshal(dispute)
	if err != nil {
		return err
	}
	if err := ctx.GetStub().PutState("DISPUTE_"+disputeID, data); err != nil {
		return err
	}

	// 复合键：work~dispute
	ck, _ := ctx.GetStub().CreateCompositeKey("work~dispute", []string{workID, disputeID})
	return ctx.GetStub().PutState(ck, []byte{0x00})
}

func (c *CopyrightContract) QueryDisputesByWork(ctx contractapi.TransactionContextInterface, workID string) ([]*DisputeRecord, error) {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey("work~dispute", []string{workID})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var disputes []*DisputeRecord
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			return nil, err
		}
		_, parts, _ := ctx.GetStub().SplitCompositeKey(item.Key)
		if len(parts) < 2 {
			continue
		}
		data, err := ctx.GetStub().GetState("DISPUTE_" + parts[1])
		if err != nil || data == nil {
			continue
		}
		var d DisputeRecord
		if err := json.Unmarshal(data, &d); err != nil {
			continue
		}
		disputes = append(disputes, &d)
	}
	return disputes, nil
}

// ─── 历史记录 ───────────────────────────────────────────────

type HistoryRecord struct {
	TxID      string     `json:"txID"`
	Timestamp string     `json:"timestamp"`
	IsDelete  bool       `json:"isDelete"`
	Value     *MusicWork `json:"value,omitempty"`
}

func (c *CopyrightContract) GetHistory(ctx contractapi.TransactionContextInterface, workID string) ([]*HistoryRecord, error) {
	iter, err := ctx.GetStub().GetHistoryForKey("WORK_" + workID)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var records []*HistoryRecord
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			return nil, err
		}
		rec := &HistoryRecord{
			TxID:      item.TxId,
			Timestamp: time.Unix(item.Timestamp.Seconds, 0).UTC().Format(time.RFC3339),
			IsDelete:  item.IsDelete,
		}
		if !item.IsDelete && len(item.Value) > 0 {
			var w MusicWork
			if err := json.Unmarshal(item.Value, &w); err == nil {
				rec.Value = &w
			}
		}
		records = append(records, rec)
	}
	return records, nil
}
