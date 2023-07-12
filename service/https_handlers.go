package service

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/AgoraIO-Community/go-tokenbuilder/chatTokenBuilder"
	rtctokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	rtmtokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtmtokenbuilder"
	"github.com/gin-gonic/gin"
)

func (s *Service) rtcToken(c *gin.Context) {
	log.Println("Generating RTC token")
	var tokenRequest RtcTokenReq
	json.NewDecoder(c.Request.Body).Decode(&tokenRequest)

	expireTimestamp := expirationFromNow(tokenRequest.Expiration)

	var userRole rtctokenbuilder2.Role
	if tokenRequest.Role == "publisher" {
		userRole = rtctokenbuilder2.RolePublisher
	} else {
		userRole = rtctokenbuilder2.RoleSubscriber
	}

	uid64, parseErr := strconv.ParseUint(tokenRequest.Uid, 10, 64)
	var rtcToken string
	var tokenErr error
	// check if conversion fails
	if parseErr != nil {
		rtcToken, tokenErr = rtctokenbuilder2.BuildTokenWithAccount(
			tokenRequest.AppId, tokenRequest.AppCertificate, tokenRequest.Channel,
			tokenRequest.Uid, userRole, expireTimestamp,
		)
	} else {
		rtcToken, tokenErr = rtctokenbuilder2.BuildTokenWithUid(
			tokenRequest.AppId, tokenRequest.AppCertificate, tokenRequest.Channel,
			uint32(uid64), userRole, expireTimestamp,
		)
	}

	if tokenErr != nil {
		log.Println(tokenErr) // token failed to generate
		c.Error(tokenErr)
		errMsg := "Error Generating RTC token - " + tokenErr.Error()
		c.AbortWithStatusJSON(400, gin.H{
			"status": 400,
			"error":  errMsg,
		})
	} else {
		log.Println("RTC Token generated")
		c.JSON(200, gin.H{
			"rtcToken": rtcToken,
		})
	}
}

func (s *Service) rtmToken(c *gin.Context) {
	log.Println("Generating RTC token")
	var tokenRequest RtmTokenReq
	json.NewDecoder(c.Request.Body).Decode(&tokenRequest)

	expireTimestamp := expirationFromNow(tokenRequest.Expiration)

	rtmToken, tokenErr := rtmtokenbuilder2.BuildToken(tokenRequest.AppId, tokenRequest.AppCertificate, tokenRequest.Uid, expireTimestamp)

	if tokenErr != nil {
		log.Println(tokenErr) // token failed to generate
		c.Error(tokenErr)
		errMsg := "Error Generating RTC token - " + tokenErr.Error()
		c.AbortWithStatusJSON(400, gin.H{
			"status": 400,
			"error":  errMsg,
		})
	} else {
		log.Println("RTC Token generated")
		c.JSON(200, gin.H{
			"rtmToken": rtmToken,
		})
	}
}

func (s *Service) chatToken(c *gin.Context) {
	log.Println("Generating Chat token")
	var tokenRequest ChatTokenReq
	json.NewDecoder(c.Request.Body).Decode(&tokenRequest)

	expireTimestamp := expirationFromNow(tokenRequest.Expiration)

	rtmToken, tokenErr := chatTokenBuilder.BuildChatUserToken(tokenRequest.AppId, tokenRequest.AppCertificate, tokenRequest.Uid, expireTimestamp)

	if tokenErr != nil {
		log.Println(tokenErr) // token failed to generate
		c.Error(tokenErr)
		errMsg := "Error Generating Chat token - " + tokenErr.Error()
		c.AbortWithStatusJSON(400, gin.H{
			"status": 400,
			"error":  errMsg,
		})
	} else {
		log.Println("RTC Token generated")
		c.JSON(200, gin.H{
			"chatToken": rtmToken,
		})
	}
}

func expirationFromNow(expiration int) uint32 {
	if expiration == 0 {
		expiration = 86400
	}
	expireTimeInSeconds := uint32(expiration)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp := currentTimestamp + expireTimeInSeconds
	return expireTimestamp
}
