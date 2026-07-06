package handler

import (
	"strconv"

	"blockchain/backend/service"
	"blockchain/backend/utils"

	"github.com/gin-gonic/gin"
)

type CopyrightHandler struct {
	Service *service.CopyrightService
}

func (h *CopyrightHandler) Register(c *gin.Context) {
	title := c.PostForm("title")
	artist := c.PostForm("artist")
	genre := c.PostForm("genre")
	description := c.PostForm("description")

	if title == "" || artist == "" {
		utils.Error(c, 1001, "title 和 artist 为必填项")
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		utils.Error(c, 1001, "请上传音频文件")
		return
	}
	defer file.Close()

	fileHash, err := utils.SHA256File(file)
	if err != nil {
		utils.Error(c, 1001, "文件哈希计算失败: "+err.Error())
		return
	}

	ownerID := c.GetString("userID")

	work, err := h.Service.RegisterWork(ownerID, title, artist, fileHash, description, genre)
	if err != nil {
		utils.Error(c, 5001, "存证失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"workID":     work.WorkID,
		"fileHash":   work.FileHash,
		"txID":       work.TxID,
		"registerAt": work.RegisterAt,
	})
}

func (h *CopyrightHandler) Query(c *gin.Context) {
	workID := c.Param("workID")
	work, err := h.Service.QueryWork(workID)
	if err != nil {
		utils.Error(c, 2001, "作品不存在")
		return
	}
	utils.Success(c, work)
}

func (h *CopyrightHandler) MyList(c *gin.Context) {
	userID := c.GetString("userID")
	works, err := h.Service.QueryWorksByOwner(userID)
	if err != nil {
		utils.Error(c, 5001, "查询失败: "+err.Error())
		return
	}
	utils.Success(c, works)
}

func (h *CopyrightHandler) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		utils.Error(c, 1001, "keyword 为必填参数")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	result, err := h.Service.SearchWorks(keyword, page, size)
	if err != nil {
		utils.Error(c, 5001, "搜索失败: "+err.Error())
		return
	}
	utils.Success(c, result)
}

func (h *CopyrightHandler) History(c *gin.Context) {
	workID := c.Param("workID")
	records, err := h.Service.GetHistory(workID)
	if err != nil {
		utils.Error(c, 5001, "查询历史失败: "+err.Error())
		return
	}
	utils.Success(c, records)
}

func (h *CopyrightHandler) VerifyHash(c *gin.Context) {
	workID := c.PostForm("workID")
	if workID == "" {
		utils.Error(c, 1001, "workID 为必填参数")
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		utils.Error(c, 1001, "请上传文件")
		return
	}
	defer file.Close()

	fileHash, err := utils.SHA256File(file)
	if err != nil {
		utils.Error(c, 1001, "文件哈希计算失败: "+err.Error())
		return
	}

	result, err := h.Service.VerifyFileHash(workID, fileHash)
	if err != nil {
		utils.Error(c, 5001, "验证失败: "+err.Error())
		return
	}

	if !result.Match {
		utils.ErrorWithData(c, 2003, "文件哈希不匹配", result)
		return
	}
	utils.Success(c, result)
}

func (h *CopyrightHandler) Transfer(c *gin.Context) {
	var req struct {
		WorkID string `json:"workID" binding:"required"`
		ToID   string `json:"toID" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 1001, "参数错误: "+err.Error())
		return
	}

	fromID := c.GetString("userID")

	work, err := h.Service.TransferCopyright(req.WorkID, fromID, req.ToID)
	if err != nil {
		utils.Error(c, 1003, "转让失败: "+err.Error())
		return
	}
	utils.Success(c, work)
}
