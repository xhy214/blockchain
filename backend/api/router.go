package api

import (
	"blockchain/backend/api/handler"
	"blockchain/backend/api/middleware"
	"blockchain/backend/config"
	"blockchain/backend/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, fabric *service.FabricClient) *gin.Engine {
	service.JWTSecret = cfg.JWT.Secret
	service.JWTExpireHours = cfg.JWT.ExpireHours

	userSvc := &service.UserService{}
	copyrightSvc := &service.CopyrightService{Fabric: fabric}
	licenseSvc := &service.LicenseService{Fabric: fabric}
	verifySvc := &service.VerifyService{Fabric: fabric}
	certSvc := &service.CertificateService{Fabric: fabric}

	userH := &handler.UserHandler{Service: userSvc}
	copyrightH := &handler.CopyrightHandler{Service: copyrightSvc}
	licenseH := &handler.LicenseHandler{LicenseSvc: licenseSvc, VerifySvc: verifySvc}
	certH := &handler.CertificateHandler{
		CertSvc:      certSvc,
		Fabric:       fabric,
		CopyrightSvc: copyrightSvc,
	}

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api/v1")
	{
		// 用户模块（无需认证）
		api.POST("/user/register", userH.Register)
		api.POST("/user/login", userH.Login)

		auth := api.Group("", middleware.AuthMiddleware())
		{
			auth.GET("/user/profile", userH.Profile)

			// 版权模块
			auth.POST("/copyright/register", copyrightH.Register)
			auth.GET("/copyright/:workID", copyrightH.Query)
			auth.GET("/copyright/my/list", copyrightH.MyList)
			auth.GET("/copyright/search", copyrightH.Search)
			auth.GET("/copyright/:workID/history", copyrightH.History)
			auth.POST("/copyright/verify-hash", copyrightH.VerifyHash)
			auth.POST("/copyright/transfer", copyrightH.Transfer)

			// 授权模块
			auth.POST("/license/grant", licenseH.Grant)
			auth.GET("/license/verify", licenseH.Verify)
			auth.GET("/license/my", licenseH.My)
			auth.POST("/license/revoke", licenseH.Revoke)
			auth.POST("/license/record-usage", licenseH.RecordUsage)

			// 创新功能
			auth.GET("/copyright/:workID/certificate", certH.Download)
			auth.POST("/dispute/file", certH.FileDispute)
			auth.GET("/dispute/:workID", certH.QueryDisputes)
		}
	}

	return r
}
