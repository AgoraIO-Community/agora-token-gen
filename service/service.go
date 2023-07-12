package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Service is the backend service
type Service struct {
	Server *http.Server
	Sigint chan os.Signal
}

// Stop service safely, closing additional connections if needed.
func (s *Service) Stop() {
	// Will continue once an interrupt has occurred
	signal.Notify(s.Sigint, os.Interrupt)
	<-s.Sigint

	// cancel would be useful if we had to close third party connection first
	// Like connections to a db or cache
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	cancel()
	err := s.Server.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}
}

// Start runs the service by listening to the specified port
func (s *Service) Start() {
	log.Println("Listening to port " + s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func NewService() *Service {

	godotenv.Load()

	serverPort, serverPortExists := os.LookupEnv("SERVER_PORT")
	if !serverPortExists || len(serverPort) == 0 {
		// Check $PORT, this is used by Railway.
		port, portExists := os.LookupEnv("PORT")
		if portExists && len(port) > 0 {
			serverPort = port
		} else {
			serverPort = "8080"
		}
	}

	s := &Service{
		Sigint: make(chan os.Signal, 1),
		Server: &http.Server{
			Addr: fmt.Sprintf(":%s", serverPort),
		},
	}

	api := gin.Default()

	api.POST("rtc/", s.rtcToken)
	api.POST("rtm/", s.rtmToken)
	api.POST("chat/", s.chatToken)

	s.Server.Handler = api
	return s
}

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
