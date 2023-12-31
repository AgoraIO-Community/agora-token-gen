package main

import "github.com/AgoraIO-Community/agora-token-gen/service"

func main() {
	s := service.NewService()
	// Stop is called on another thread, but waits for an interrupt
	go s.Stop()
	s.Start()
}
