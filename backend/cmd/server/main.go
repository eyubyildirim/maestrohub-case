package main

import (
	"log"
	"net/http"
	"tower-defence-engine/internal/config"
	"tower-defence-engine/internal/game"
	"tower-defence-engine/internal/websocket"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("FATAL: Could not load config.json: %v", err)
	}
	log.Println("Configuration loaded successfully.")

	hub := websocket.NewHub()
	go hub.Run()

	gameEngine := game.NewEngine(hub, cfg)
	go gameEngine.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWs(w, r, gameEngine.CommandCh)
	})

	log.Println("Command Center starting on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("ListenAndServe failed: %v", err)
	}
}
