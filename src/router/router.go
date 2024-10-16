package router

import (
	"assigement_wallet/src/wallet"
	"github.com/gin-gonic/gin"
)

func Register(g *gin.Engine) {
	walletV1 := g.Group("/api/wallet/v1")
	{
		walletV1.POST("/login", wallet.LoginController)
		walletV1.POST("/logout", wallet.LogoutController)
		walletV1.POST("/deposit", wallet.DepositController)
		walletV1.POST("/withdraw", wallet.WithdrawController)
		walletV1.GET("/query/balance", wallet.QueryBalanceController)
		walletV1.GET("/query/history", wallet.QueryHistoryController)
		walletV1.POST("/transfer", wallet.TransferController)
	}
}
