package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

var (
	auth        *spotifyauth.Authenticator
	state       = "abc124"
	clientStore = make(map[string]*spotify.Client)
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL("http://localhost:8080/callback"),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserReadEmail, spotifyauth.ScopeUserTopRead),
		spotifyauth.WithClientID(os.Getenv("SPOTIFY_ID")),
		spotifyauth.WithClientSecret(os.Getenv("SPOTIFY_SECRET")),
	)

	r := gin.Default()

	// Use the Cors middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	r.Use(cors.New(config))

	r.GET("/login", func(c *gin.Context) {
		url := auth.AuthURL(state)
		c.Redirect(http.StatusTemporaryRedirect, url)
	})

	r.GET("/user", func(c *gin.Context) {
		userID := c.Param("userId") // Retrieve the user ID from the request parameters

		// Retrieve the user's information based on the user ID
		// Example: You might retrieve the user's name from your user database
		user, err := clientStore[userID].CurrentUser(c)

		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to get user profile")
			return
		}

		userName := user.DisplayName

		// Return the user's information as JSON to the frontend
		userInfo := gin.H{
			"userId":   userID,
			"userName": userName,
			// Add any other necessary user information here
		}

		c.JSON(http.StatusOK, userInfo)
	})

	r.GET("/spotify/auth-url", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"authURL": auth.AuthURL(state)})
	})

	r.GET("/callback", func(c *gin.Context) {

		state := c.Query("state") // Retrieve the state from the query parameters

		// Use the same state string here that you used to generate the URL
		token, err := auth.Token(c.Request.Context(), state, c.Request)
		if err != nil {
			c.String(http.StatusNotFound, "Couldn't get token")
			return
		}

		// Create a client using the specified token
		client := spotify.New(auth.Client(c.Request.Context(), token))

		// The client can now be used to make authenticated requests
		// Example: Get the current user's profile
		user, err := client.CurrentUser(c)
		fmt.Println(user)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to get user profile")
			return
		}

		clientStore[user.ID] = client

		frontendURL := "http://localhost:3000?userId=" + user.ID // Include the userID as a query parameter
		c.Redirect(http.StatusFound, frontendURL)
	})

	r.Run(":8080")
}
