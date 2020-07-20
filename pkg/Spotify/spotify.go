package spotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type reqToSpotify struct {
	Country string
	Limit   int
	Token   string
}

var spotifyRequest = reqToSpotify{
	Country: "US",
	Limit:   3,
	Token:   fmt.Sprintf("Bearer %s", os.Getenv("SPOTIFY_TOKEN")),
}

const (
	ChillMusic = "chill"
	PartyMusic = "party"
	Hiphop     = "hiphop"
	EDM        = "edm_dance"
	OldMusic   = "decades"
	AnyMusic   = "anything"
)

var PlaylistCategories = []string{"toplists", "summer", "hiphop", "pop", "country", "workout", "latin", "mood", "rock", "rnb"}

type parsedPlaylists struct {
	Playlists struct {
		Items []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
		} `json:"items"`
	} `json:"playlists"`
}

func GetPlayLists(musicCategory string) ([]string, error) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	offset := r.Perm(8)[0]
	url := fmt.Sprintf("https://api.spotify.com/v1/browse/categories/%s/playlists?country=%s&limit=%d&offset=%d",
		musicCategory, spotifyRequest.Country, spotifyRequest.Limit, offset)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", spotifyRequest.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("StatusNotOK")
	}
	var parsedData = parsedPlaylists{}
	if err := json.NewDecoder(resp.Body).Decode(&parsedData); err != nil {
		return nil, err
	}
	recommendedPlaylists := []string{}
	for _, item := range parsedData.Playlists.Items {
		recommendedPlaylists = append(recommendedPlaylists, item.ExternalUrls.Spotify)
	}
	return recommendedPlaylists, nil
}

func GetRandomCategory(playlistCategories []string) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	num := r.Perm(10)[0]
	randomCategory := playlistCategories[num]
	return randomCategory
}

func UpdateToken(spotifySecretKey string) error {
	apiURL := "https://accounts.spotify.com/api/token"
	body := url.Values{}
	body.Set("grant_type", "client_credentials")
	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", spotifySecretKey))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var spotifyRespReceiver struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&spotifyRespReceiver); err != nil {
		return err
	}
	var configVar = struct {
		SPOTIFY_TOKEN string
	}{
		spotifyRespReceiver.AccessToken,
	}
	bsBody, err := json.Marshal(configVar)
	if err != nil {
		return err
	}
	req, err = http.NewRequest(http.MethodPatch, fmt.Sprintf("https://api.heroku.com/apps/%s/config-vars",
		os.Getenv("APP_NAME")), bytes.NewReader(bsBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("HEROKU_TOKEN")))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("StatusNotOK %d", resp.StatusCode))
	}
	spotifyRequest.Token = os.Getenv("SPOTIFY_TOKEN")
	log.Println("Changed the Spotify token")
	return nil
}
