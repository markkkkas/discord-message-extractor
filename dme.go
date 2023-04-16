package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := flag.String("token", "", "discord bot auth token")
	flag.Parse()

	if *token == "" {
		log.Fatalf("missing token")
	}

	dg, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatalf("failed to create discord client: %v\n", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

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
			respondError(w, "Internal server error", http.StatusBadRequest)
			return
		}

		json, err := json.Marshal(messages)
		if err != nil {
			log.Printf("error while marshalling json: %v\n", err)
			respondError(w, "Internal server error", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(json))
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func respondError(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "{\"error\": \"%s\"}", msg)
}
