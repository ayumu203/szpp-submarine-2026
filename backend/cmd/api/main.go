package main

import (
	"backend/application"
	"backend/infrastructure"
	"backend/presentation"
	"log"
	"net/http"
	"os"
)

func main() {
	gameRepo := infrastructure.NewMemoryGameRepository()
	turnLogRepo := infrastructure.NewMemoryTurnLogRepository()
	predictionRepo := infrastructure.NewMemoryPredictionRepository()
	playerGamesRepo := infrastructure.NewMemoryPlayerGamesIndexRepository()

	cpuAnalysis := application.NewCpuAnalysisService(turnLogRepo, predictionRepo)
	cpuDecision := application.NewCpuDecisionService()
	gameService := application.NewGameService(
		gameRepo,
		turnLogRepo,
		predictionRepo,
		playerGamesRepo,
		cpuDecision,
		cpuAnalysis,
	)

	handler := presentation.NewGameHandler(gameService)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	server := withCORS(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("api server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatal(err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
