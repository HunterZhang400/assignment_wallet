package wallet

import (
	"assigement_wallet/basedata"
	"assigement_wallet/config"
	"assigement_wallet/http_core"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// LoginController to simplify demo, just simulate login by user id, without password check design
func LoginController(ctx *gin.Context) {
	userID := ctx.PostForm("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, basedata.NewErrorResponse("user_id not allowed to be empty"))
		return
	}
	userExist, err := CheckUserExist(userID)
	if err != nil {
		log.Println("LoginController:", err.Error())
		ctx.JSON(http.StatusInternalServerError, basedata.NewErrorResponse(basedata.ServerUnavailable))
		return
	}
	if userExist {
		token, err := http_core.EncodeJWT(http_core.UserInfo{UserID: userID}, []byte(config.ServerConfigs.Server.JWTKey))
		if err != nil {
			log.Println("LoginController:", err.Error())
			ctx.JSON(http.StatusInternalServerError, basedata.NewErrorResponse(basedata.ServerUnavailable))
			return
		}
		ctx.SetCookie(http_core.SessionName, token, 3600, "/", "127.0.0.1", false, false)
		ctx.JSON(http.StatusOK, basedata.NewResponse(basedata.Success))
	} else {
		ctx.JSON(http.StatusBadRequest, basedata.NewErrorResponse(fmt.Sprintf("not found %s", userID)))
	}
}

func LogoutController(ctx *gin.Context) {
	ctx.SetCookie(http_core.SessionName, "", 3600, "/", "127.0.0.1", false, false)
	ctx.JSON(http.StatusOK, basedata.NewResponse(basedata.Success))
}

func QueryBalanceController(ctx *gin.Context) {
	userID := getUserIDFromSession(ctx)
	if userID == "" {
		ctx.JSON(http.StatusForbidden, basedata.NewErrorResponse(basedata.NotLogIn))
		return
	}
	balance, err := QueryBalance(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, basedata.NewErrorResponse(basedata.ServerUnavailable))
		return
	}
	ctx.JSON(http.StatusOK, basedata.NewResponse(balance))
}

type depositResponse struct {
	//indicate deposit success or fail
	DepositFlag string `json:"deposit_flag"`
	//indicate query balance success or fail
	QueryBalanceFlag string `json:"query_balance_flag"`
	CurrentBalance   int64  `json:"current_balance"`
}

func DepositController(ctx *gin.Context) {
	userID := getUserIDFromSession(ctx)
	if userID == "" {
		ctx.JSON(http.StatusForbidden, basedata.NewErrorResponse(basedata.NotLogIn))
		return
	}
	amountString := ctx.PostForm("amount")
	if amountString == "" {
		ctx.JSON(http.StatusBadRequest, basedata.NewErrorResponse("amount required"))
		return
	}
	amount, err := strconv.ParseInt(amountString, 10, 64)
	if err != nil {
		log.Println("DepositController:", fmt.Sprintf("wrong amount:%s", amountString))
		ctx.JSON(http.StatusBadRequest, basedata.NewErrorResponse("invalid amount"))
		return
	}
	if amount <= 0 {
		ctx.JSON(http.StatusBadRequest, basedata.NewErrorResponse("invalid amount"))
		return
	}
	err = Deposit(userID, amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, basedata.NewErrorResponse(basedata.ServerUnavailable))
		return
	}
	resp := depositResponse{
		DepositFlag:      basedata.Success,
		QueryBalanceFlag: basedata.Success,
		CurrentBalance:   0,
	}
	resp.CurrentBalance, err = QueryBalance(userID)
	if err != nil {
		//only query current balance fail, but the deposit succeeded, avoid make user confused
		resp.QueryBalanceFlag = basedata.Fail
		ctx.JSON(http.StatusOK, basedata.NewResponse(resp))
		return
	}
	ctx.JSON(http.StatusOK, basedata.NewResponse(resp))
}

type withdrawResponse struct {
	//indicate withdraw success or fail
	WithdrawFlag string `json:"withdraw_flag"`
	//indicate query balance success or fail
	QueryBalanceFlag string `json:"query_balance_flag"`
	CurrentBalance   int64  `json:"current_balance"`
}

func WithdrawController(ctx *gin.Context) {
	userID := getUserIDFromSession(ctx)
	if userID == "" {
		ctx.JSON(http.StatusForbidden, basedata.NewErrorResponse(basedata.NotLogIn))
		return
	}
	amountString := ctx.PostForm("amount")
	if amountString == "" {
		ctx.JSON(http.StatusBadRequest, basedata.NewErrorResponse("amount required"))
		return
	}
	amount, err := strconv.ParseInt(amountString, 10, 64)
	if err != nil {
		log.Println("WithdrawController:", fmt.Sprintf("wrong amount:%s", amountString))
		ctx.JSON(http.StatusBadRequest, basedata.NewErrorResponse("invalid amount"))
		return
	}
	if amount <= 0 {
		ctx.JSON(http.StatusBadRequest, basedata.NewErrorResponse("invalid amount"))
		return
	}
	err = Withdraw(userID, amount)
	if err != nil {
		log.Println("WithdrawController Withdraw:", err.Error())
		ctx.JSON(http.StatusInternalServerError, basedata.NewErrorResponse(basedata.ServerUnavailable))
		return
	}
	resp := withdrawResponse{
		WithdrawFlag:     basedata.Success,
		QueryBalanceFlag: basedata.Success,
		CurrentBalance:   0,
	}
	resp.CurrentBalance, err = QueryBalance(userID)
	if err != nil {
		//only query current balance fail, but the withdrawal succeeded, avoid make user confused
		resp.QueryBalanceFlag = basedata.Fail
		ctx.JSON(http.StatusOK, basedata.NewResponse(resp))
		return
	}
	ctx.JSON(http.StatusOK, basedata.NewResponse(resp))
}

type transferParameter struct {
	ToUserID string `json:"to_user_id" form:"to_user_id"`
	Amount   int64  `json:"amount"  form:"amount"`
}

type transferResponse struct {
	//indicate withdraw success or fail
	TransferFlag string `json:"transfer_flag"`
	//indicate query balance success or fail
	QueryBalanceFlag string `json:"query_balance_flag"`
	CurrentBalance   int64  `json:"current_balance"`
}

func TransferController(ctx *gin.Context) {
	fromUserID := getUserIDFromSession(ctx)
	if fromUserID == "" {
		ctx.JSON(http.StatusForbidden, basedata.NewErrorResponse(basedata.NotLogIn))
		return
	}
	var param transferParameter
	ctx.Bind(&param)
	if param.ToUserID == "" {
		ctx.JSON(http.StatusBadRequest, basedata.NewErrorResponse("to_user_id required"))
		return
	}
	if param.Amount <= 0 {
		ctx.JSON(http.StatusBadRequest, basedata.NewErrorResponse("invalid amount"))
		return
	}
	err := Transfer(fromUserID, param.ToUserID, param.Amount)
	if err != nil {
		log.Println("TransferController Transfer:", err.Error())
		ctx.JSON(http.StatusInternalServerError, basedata.NewErrorResponse(basedata.ServerUnavailable))
		return
	}
	resp := transferResponse{
		TransferFlag:     basedata.Success,
		QueryBalanceFlag: basedata.Success,
		CurrentBalance:   0,
	}
	resp.CurrentBalance, err = QueryBalance(fromUserID)
	if err != nil {
		//only query current balance fail, but the transfer succeeded, avoid make user confused
		resp.QueryBalanceFlag = basedata.Fail
		ctx.JSON(http.StatusOK, basedata.NewResponse(resp))
		return
	}
	ctx.JSON(http.StatusOK, basedata.NewResponse(resp))
}

func QueryHistoryController(ctx *gin.Context) {
	userID := getUserIDFromSession(ctx)
	if userID == "" {
		ctx.JSON(http.StatusForbidden, basedata.NewErrorResponse(basedata.NotLogIn))
		return
	}
	sizeString := ctx.Query("size")
	if sizeString == "" {
		//use default size
		sizeString = "100"
	}
	size, err := strconv.ParseInt(sizeString, 10, 64)
	if err != nil {
		log.Println("WithdrawController:", fmt.Sprintf("wrong size:%s", sizeString))
		ctx.JSON(http.StatusBadRequest, basedata.NewErrorResponse("invalid size"))
		return
	}
	details, err := QueryHistory(userID, size)
	if err != nil {
		log.Println("QueryHistoryController QueryHistory:", err.Error())
		ctx.JSON(http.StatusInternalServerError, basedata.NewErrorResponse(basedata.ServerUnavailable))
		return
	}
	ctx.JSON(http.StatusOK, basedata.NewResponse(details))
}

func getUserIDFromSession(ctx *gin.Context) string {
	v, exist := ctx.Get(http_core.ContextUserKey)
	if exist == false {
		return ""
	}
	userID, ok := v.(string)
	if !ok {
		log.Println("getUserIDFromSession:", fmt.Sprintf("invalid userID:%+v", v))
		return ""
	}
	return userID
}
