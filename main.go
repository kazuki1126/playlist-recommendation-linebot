package main

import (
	"fmt"

	line "github.com/kazuki1126/playlist-recommendation-linebot/pkg/LINE"
	spotify "github.com/kazuki1126/playlist-recommendation-linebot/pkg/Spotify"

	"log"
	"net/http"
	"os"

	"github.com/robfig/cron"
)

func main() {
	c := cron.New()
	c.AddFunc("*/40 * * * *", func() { spotify.UpdateToken(os.Getenv("SPOTIFY_SECRET_KEY")) })
	c.Start()
	http.HandleFunc("/callback", line.SendReply)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
