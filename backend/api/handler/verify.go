package handler

import (
	"blockchain/backend/service"
	"blockchain/backend/utils"

	"github.com/gin-gonic/gin"
)

type VerifyHandler struct {
	VerifySvc *service.VerifyService
}

func (h *VerifyHandler) VerifyLicense(c *gin.Context) {
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
