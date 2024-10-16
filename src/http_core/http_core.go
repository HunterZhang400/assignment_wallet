package http_core

import (
	"assigement_wallet/src/basedata"
	"assigement_wallet/src/config"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
)

const SessionName = "wallet_session"
const ContextUserKey = "user_id"

type UserInfo struct {
	UserID string `json:"id"`
}

type WalletClaims struct {
	UserData UserInfo `json:"user_data"`
	jwt.RegisteredClaims
}

func EncodeJWT(userData UserInfo, secret []byte) (string, error) {
	claims := WalletClaims{
		UserData: userData,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secret)
	return ss, err
}

func DecodeJWT(tokenString string, secret []byte) (*UserInfo, error) {
	token, err := jwt.ParseWithClaims(tokenString, &WalletClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if userClaims, ok := token.Claims.(*WalletClaims); ok && token.Valid {
		return &userClaims.UserData, nil
	} else {
		return nil, errors.New("invalid token")
	}
}

var whitePathList = map[string]struct{}{"/api/wallet/v1/login": {}}

func SessionControl(ctx *gin.Context) {
	if _, e := whitePathList[ctx.Request.URL.Path]; e {
		return
	}
	ss, err := ctx.Cookie(SessionName)
	if err != nil {
		log.Println("SessionControl:", err.Error())
		ctx.AbortWithStatusJSON(http.StatusForbidden, basedata.NewErrorResponse(basedata.NotLogIn))
		return
	}
	userIndo, err := DecodeJWT(ss, []byte(config.ServerConfigs.Server.JWTKey))
	if err != nil {
		log.Println("SessionControl DecodeJWT:", err.Error())
		ctx.AbortWithStatusJSON(http.StatusForbidden, basedata.NewErrorResponse(basedata.NotLogIn))
		return
	}
	if userIndo == nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, basedata.NewErrorResponse(basedata.NotLogIn))
		return
	}
	ctx.Set(ContextUserKey, userIndo.UserID)
}
