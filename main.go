package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := flag.String("token", "", "discord bot auth token")
	flag.Parse()

	if *token == "" {
		log.Fatalf("missing token")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dg, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatalf("failed to create discord client: %v\n", err)
	}

	dg.Identify.Intents = discordgo.IntentGuildMessages

	err = dg.Open()
	if err != nil {
		log.Fatalf("failed to create websocket connection: %v\n", err)
	}

	log.Printf("connected as %s\n", dg.State.User.Username)

	server := &http.Server{Addr: ":8080", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		channelID := r.URL.Query().Get("channelId")
		if channelID == "" {
			respondError(w, "missing query param: channelId", http.StatusBadRequest)
			return
		}

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			respondError(w, "missing or malformed query param: limit", http.StatusBadRequest)
			return
		}

		messages, err := dg.ChannelMessages(channelID, limit, "", "", "")
		if err != nil {
			log.Printf("error while getting channel messages: %v\n", err)
			respondError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		respondJSON(w, messages, http.StatusOK)
	})}

	go func() {
		log.Println("listening on :8080")

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("server error: %v\n", err)
			stop()
		}
	}()

	<-ctx.Done()

	log.Println("shutting down...")

	shutdownCtx, shutdownStop := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownStop()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v\n", err)
	}

	log.Println("shutting down client...")

	if err := dg.Close(); err != nil {
		log.Printf("client shutdown error: %v\n", err)
	}

	log.Println("shutdown complete")
}

func respondError(w http.ResponseWriter, msg string, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, "{\"error\": \"%s\"}", msg)
}

func respondJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Add("Content-Type", "application/json")

	json, err := json.Marshal(data)
	if err != nil {
		log.Printf("error while marshalling json: %v\n", err)
		respondError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	fmt.Fprint(w, string(json))
}
