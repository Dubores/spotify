package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

var (
	auth  *spotifyauth.Authenticator
	state = "abc123"
)

func getTopTracks(client *spotify.Client) (*spotify.FullTrackPage, error) {
	tracks, err := client.CurrentUsersTopTracks(context.Background(), spotify.Limit(10), spotify.Timerange("short_term"))
	if err != nil {
		return nil, err
	}
	return tracks, nil
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL("http://localhost:8080/callback"),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserReadEmail, spotifyauth.ScopeUserTopRead),
		spotifyauth.WithClientID(os.Getenv("SPOTIFY_ID")),
		spotifyauth.WithClientSecret(os.Getenv("SPOTIFY_SECRET")),
	)

	fmt.Println("SPOTIFY_ID:", os.Getenv("SPOTIFY_ID"))
	r := gin.Default()

	r.GET("/login", func(c *gin.Context) {
		url := auth.AuthURL(state)
		c.Redirect(http.StatusTemporaryRedirect, url)
	})

	r.GET("/callback", func(c *gin.Context) {
		token, err := auth.Token(c.Request.Context(), state, c.Request)
		if err != nil {
			c.String(http.StatusForbidden, "Couldn't get token")
			return
		}

		client := spotify.New(auth.Client(c.Request.Context(), token))
		topTracks, err := getTopTracks(client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get top tracks"})
			return
		}

		c.JSON(http.StatusOK, topTracks)
	})

	r.Run(":8080")
}
