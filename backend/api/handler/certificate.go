package handler

import (
	"fmt"

	"blockchain/backend/service"
	"blockchain/backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CertificateHandler struct {
	CertSvc    *service.CertificateService
	Fabric     *service.FabricClient
	CopyrightSvc *service.CopyrightService
}

func (h *CertificateHandler) Download(c *gin.Context) {
	workID := c.Param("workID")

	pdfBytes, err := h.CertSvc.GenerateCertificate(workID)
	if err != nil {
		utils.Error(c, 5001, "证书生成失败: "+err.Error())
		return
	}

	filename := fmt.Sprintf("copyright_certificate_%s.pdf", workID[:8])
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/pdf")
	c.Data(200, "application/pdf", pdfBytes)
}

func (h *CertificateHandler) FileDispute(c *gin.Context) {
	var req struct {
		WorkID   string `json:"workID" binding:"required"`
		Evidence string `json:"evidence" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 1001, "参数错误: "+err.Error())
		return
	}

	disputeID := "disp-" + uuid.New().String()[:8]
	claimantID := c.GetString("userID")

	_, err := h.Fabric.Submit("FileDispute", disputeID, req.WorkID, claimantID, req.Evidence)
	if err != nil {
		utils.Error(c, 5001, "提交争议失败: "+err.Error())
		return
	}
	utils.Success(c, gin.H{
		"disputeID":  disputeID,
		"workID":     req.WorkID,
		"claimantID": claimantID,
		"status":     "PENDING",
	})
}

func (h *CertificateHandler) QueryDisputes(c *gin.Context) {
	workID := c.Param("workID")

	data, err := h.Fabric.Evaluate("QueryDisputesByWork", workID)
	if err != nil {
		utils.Error(c, 5001, "查询争议失败: "+err.Error())
		return
	}
	utils.RawJSON(c, 200, data)
}
