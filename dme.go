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
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "{\"error\": \"missing query param: channelId\"}")
			return
		}

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "{\"error\": \"missing or malformed query param: limit\"}")
			return
		}

		messages, err := dg.ChannelMessages(channelID, limit, "", "", "")
		if err != nil {
			log.Printf("error while getting channel messages: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "{\"error\": \"Internal server error\"}")
			return
		}

		json, err := json.Marshal(messages)
		if err != nil {
			log.Printf("error while marshalling json: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "{\"error\": \"Internal server error\"}")
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(json))
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
