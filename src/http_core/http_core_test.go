package http_core

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestEncodeJWT(t *testing.T) {
	secret := []byte("1234")
	userID := "a"
	ss, err := EncodeJWT(UserInfo{UserID: userID}, secret)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	//fmt.Println(ss)
	userInfo, err := DecodeJWT(ss, secret)
	if err != nil {
		log.Fatalf(err.Error())
	}
	assert.Equal(t, userID, userInfo.UserID, "decode got wrong user data")
}
