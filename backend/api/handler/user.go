package handler

import (
	"blockchain/backend/model"
	"blockchain/backend/service"
	"blockchain/backend/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *service.UserService
}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		RealName string `json:"realName"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 1001, "参数错误: "+err.Error())
		return
	}

	user, err := h.Service.Register(req.Username, req.Password, req.RealName)
	if err != nil {
		utils.Error(c, 1001, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"realName": user.RealName,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 1001, "参数错误: "+err.Error())
		return
	}

	token, user, err := h.Service.Login(req.Username, req.Password)
	if err != nil {
		utils.Error(c, 1001, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"realName": user.RealName,
		},
	})
}

func (h *UserHandler) Profile(c *gin.Context) {
	userID := c.GetString("userID")
	user, err := model.FindUserByID(userID)
	if err != nil {
		utils.Error(c, 1002, "用户不存在")
		return
	}
	utils.Success(c, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"realName":  user.RealName,
		"createdAt": user.CreatedAt,
	})
}
