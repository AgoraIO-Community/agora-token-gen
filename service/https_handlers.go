package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AgoraIO-Community/go-tokenbuilder/chatTokenBuilder"
	rtctokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	rtmtokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtmtokenbuilder"
	whiteboardtokenbuilder "github.com/netless-io/netless-token/golang"
)

type RtcTokenReq struct {
	AppId          string `json:"appId"`
	AppCertificate string `json:"certificate"`
	Channel        string `json:"channel"`
	Uid            string `json:"uid"`
	Role           string `json:"role,omitempty"`
	Expiration     int    `json:"expire,omitempty"`
}

type RtmTokenReq struct {
	AppId          string `json:"appId"`
	AppCertificate string `json:"certificate"`
	Channel        string `json:"channel"`
	Uid            string `json:"uid"`
	Expiration     int    `json:"expire,omitempty"`
}

type ChatTokenReq struct {
	AppId          string `json:"appId"`
	AppCertificate string `json:"certificate"`
	Uid            string `json:"uid"`
	Expiration     int    `json:"expire,omitempty"`
}

type WhiteboardSDKTokenReq struct {
	Role       string `json:"role"`
	AccessKey  string `json:"accessKey"`
	SecretKey  string `json:"SecretKey"`
	Expiration int    `json:"expire,omitempty"`
}

type WhiteboardRoomTokenReq struct {
	Role       string `json:"role"`
	AccessKey  string `json:"accessKey"`
	SecretKey  string `json:"SecretKey"`
	Expiration int    `json:"expire,omitempty"`
	RoomUuid   string `json:"roomuuid"`
}

type WhiteboardTaskTokenReq struct {
	Role       string `json:"role"`
	AccessKey  string `json:"accessKey"`
	SecretKey  string `json:"SecretKey"`
	Expiration int    `json:"expire,omitempty"`
	TaskUuid   string `json:"taskuuid"`
}

func RtcToken(w http.ResponseWriter, r *http.Request) {
	log.Println("Generating RTC token")
	var tokenRequest RtcTokenReq

	err := json.NewDecoder(r.Body).Decode(&tokenRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
		log.Println(tokenErr)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": http.StatusBadRequest,
			"error":  "Error generating RTC token: " + tokenErr.Error(),
		})
		return
	}
	log.Println("RTC Token generated")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"rtcToken": rtcToken,
	})
}

func RtmToken(w http.ResponseWriter, r *http.Request) {
	log.Println("Generating RTM token")
	var tokenRequest RtmTokenReq
	err := json.NewDecoder(r.Body).Decode(&tokenRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	expireTimestamp := expirationFromNow(tokenRequest.Expiration)

	rtmToken, tokenErr := rtmtokenbuilder2.BuildToken(
		tokenRequest.AppId,
		tokenRequest.AppCertificate,
		tokenRequest.Uid,
		expireTimestamp,
	)

	if tokenErr != nil {
		log.Println(tokenErr)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": http.StatusBadRequest,
			"error":  "Error generating RTC token: " + tokenErr.Error(),
		})
		return
	}

	log.Println("RTC Token generated")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"rtmToken": rtmToken,
	})
}

func ChatToken(w http.ResponseWriter, r *http.Request) {
	log.Println("Generating Chat token")
	var tokenRequest ChatTokenReq
	err := json.NewDecoder(r.Body).Decode(&tokenRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	expireTimestamp := expirationFromNow(tokenRequest.Expiration)

	chatToken, tokenErr := chatTokenBuilder.BuildChatUserToken(
		tokenRequest.AppId,
		tokenRequest.AppCertificate,
		tokenRequest.Uid,
		expireTimestamp,
	)

	if tokenErr != nil {
		log.Println(tokenErr)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": http.StatusBadRequest,
			"error":  "Error generating chat token: " + tokenErr.Error(),
		})
		return
	}

	log.Println("Chat token generated")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"chatToken": chatToken,
	})
}

// Create Whiteboard SDK token
func WhiteboardSDKToken(w http.ResponseWriter, r *http.Request) {
	log.Println("Generating Whiteboard SDK token")
	var tokenRequest WhiteboardSDKTokenReq
	err := json.NewDecoder(r.Body).Decode(&tokenRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	content := whiteboardtokenbuilder.SDKContent{
		Role: tokenRequest.Role,
	}

	sdkToken := whiteboardtokenbuilder.SDKToken(
		tokenRequest.AccessKey,
		tokenRequest.SecretKey,
		int64(tokenRequest.Expiration),
		&content,
	)

	log.Println("Whiteboard SDK token generated")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sdkToken": sdkToken,
	})
}

// Create Whiteboard Room token
func WhiteboardRoomToken(w http.ResponseWriter, r *http.Request) {
	log.Println("Generating Whiteboard Room token")
	var tokenRequest WhiteboardRoomTokenReq
	err := json.NewDecoder(r.Body).Decode(&tokenRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	content := whiteboardtokenbuilder.RoomContent{
		Role: tokenRequest.Role,
		Uuid: tokenRequest.RoomUuid,
	}

	roomToken := whiteboardtokenbuilder.RoomToken(
		tokenRequest.AccessKey,
		tokenRequest.SecretKey,
		int64(tokenRequest.Expiration),
		&content,
	)

	log.Println("Whiteboard room token generated")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"roomToken": roomToken,
	})
}

// Create Whiteboard Task token
func WhiteboardTaskToken(w http.ResponseWriter, r *http.Request) {
	log.Println("Generating Whiteboard task token")
	var tokenRequest WhiteboardTaskTokenReq
	err := json.NewDecoder(r.Body).Decode(&tokenRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	content := whiteboardtokenbuilder.TaskContent{
		Role: tokenRequest.Role,
		Uuid: tokenRequest.TaskUuid,
	}

	taskToken := whiteboardtokenbuilder.TaskToken(
		tokenRequest.AccessKey,
		tokenRequest.SecretKey,
		int64(tokenRequest.Expiration),
		&content,
	)

	log.Println("Whiteboard task token generated")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"taskToken": taskToken,
	})
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
