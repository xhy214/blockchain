package handler

import (
	"blockchain/backend/service"
	"blockchain/backend/utils"

	"github.com/gin-gonic/gin"
)

type LicenseHandler struct {
	LicenseSvc *service.LicenseService
	VerifySvc  *service.VerifyService
}

func (h *LicenseHandler) Grant(c *gin.Context) {
	var req struct {
		WorkID      string `json:"workID" binding:"required"`
		LicenseeID  string `json:"licenseeID" binding:"required"`
		LicenseType string `json:"licenseType" binding:"required"`
		StartDate   string `json:"startDate" binding:"required"`
		EndDate     string `json:"endDate" binding:"required"`
		MaxUsage    int    `json:"maxUsage"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 1001, "参数错误: "+err.Error())
		return
	}

	grantorID := c.GetString("userID")

	lic, err := h.LicenseSvc.GrantLicense(
		req.WorkID, grantorID, req.LicenseeID, req.LicenseType,
		req.StartDate, req.EndDate, req.MaxUsage)
	if err != nil {
		utils.Error(c, 1003, "授权失败: "+err.Error())
		return
	}
	utils.Success(c, lic)
}

func (h *LicenseHandler) Verify(c *gin.Context) {
	workID := c.Query("workID")
	licenseeID := c.Query("licenseeID")
	if workID == "" || licenseeID == "" {
		utils.Error(c, 1001, "workID 和 licenseeID 为必填参数")
		return
	}

	result, err := h.VerifySvc.VerifyLicense(workID, licenseeID)
	if err != nil {
		utils.Error(c, 5001, "核验失败: "+err.Error())
		return
	}
	utils.Success(c, result)
}

func (h *LicenseHandler) My(c *gin.Context) {
	userID := c.GetString("userID")
	licenses, err := h.LicenseSvc.QueryLicensesByLicensee(userID)
	if err != nil {
		utils.Error(c, 5001, "查询失败: "+err.Error())
		return
	}
	utils.Success(c, licenses)
}

func (h *LicenseHandler) Revoke(c *gin.Context) {
	var req struct {
		LicenseID string `json:"licenseID" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 1001, "参数错误: "+err.Error())
		return
	}

	revokerID := c.GetString("userID")

	if err := h.LicenseSvc.RevokeLicense(req.LicenseID, revokerID); err != nil {
		utils.Error(c, 1003, "撤销失败: "+err.Error())
		return
	}
	utils.Success(c, gin.H{"licenseID": req.LicenseID, "status": "REVOKED"})
}

func (h *LicenseHandler) RecordUsage(c *gin.Context) {
	var req struct {
		LicenseID string `json:"licenseID" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 1001, "参数错误: "+err.Error())
		return
	}

	if err := h.LicenseSvc.RecordUsage(req.LicenseID); err != nil {
		utils.Error(c, 3003, "记录使用失败: "+err.Error())
		return
	}
	lic, _ := h.LicenseSvc.QueryLicense(req.LicenseID)
	utils.Success(c, gin.H{
		"licenseID": req.LicenseID,
		"usedCount": func() int {
			if lic != nil {
				return lic.UsedCount
			}
			return 0
		}(),
	})
}
