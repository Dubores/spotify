package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

var (
	auth          *spotifyauth.Authenticator
	secretKey     []byte
	userTokens    = make(map[string]string)
	userAPITokens = make(map[string]*oauth2.Token)
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

	secretKey := []byte(os.Getenv("SPOTIFY_ID"))

	if secretKey == nil {
		// Cheating the compiler
	}

	r := gin.Default()

	// Use the Cors middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	r.Use(cors.New(config))

	r.GET("/login", spotifyAuthHandler)

	// r.GET("/isLoggedIn", isLoggedInHandler)

	r.GET("/callback", spotifyCallbackHandler)

	apiGroup := r.Group("/api")
	{
		// API Endpoints group for user info
		userGroup := apiGroup.Group("/user")
		{
			userGroup.GET("/profile", ensureAuthenticated(), spotifyUserHandler)
		}
	}
	r.Run(":8080")
}

// Handlers

// func isLoggedInHandler(c *gin.Context){

// }

func spotifyUserHandler(c *gin.Context) {
	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert userID to string"})
		return
	}

	client := spotify.New(auth.Client(c.Request.Context(), userAPITokens[userIDStr]))

	user, err := client.CurrentUser(c)
	if err != nil {
		c.String(http.StatusNotFound, "Couldn't get current user")
		return
	}

	fmt.Println(user)

	c.JSON(http.StatusOK, gin.H{"userName": user.DisplayName, "url": user.Images[0].URL})
}

func spotifyAuthHandler(c *gin.Context) {
	state, _ := generateRandomString(16)
	url := auth.AuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func spotifyCallbackHandler(c *gin.Context) {
	state := c.Query("state") // Retrieve the state from the query parameters
	//code := c.Query("code")

	token, err := auth.Token(c.Request.Context(), state, c.Request)
	if err != nil {
		c.String(http.StatusNotFound, "Couldn't get token")
		return
	}

	// Create a client using the specified token
	client := spotify.New(auth.Client(c.Request.Context(), token))

	// The client can now be used to make authenticated requests
	//Get the current user's profile
	// Generate and store the JWT token
	user, err := client.CurrentUser(c)
	if err != nil {
		c.String(http.StatusNotFound, "Couldn't get current user")
		return
	}

	jwtToken, err := generateJWT(user.ID)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	// Save the tokens in database
	userTokens[user.ID] = jwtToken
	userAPITokens[user.ID] = token

	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to get user profile")
		return
	}

	// Set the token in a HTTP-only cookie

	c.SetCookie("jwtToken", jwtToken, 3600, "/", "localhost", false, true)

	frontendURL := "http://localhost:3000" // Include the userID as a query parameter
	c.Redirect(http.StatusFound, frontendURL)
}

func ensureAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Ensuring auth")
		tokenString, err := c.Cookie("jwtToken")
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization missing"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token claims"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iss": "https://kagsaboys.com",               // Replace with your issuer
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expiration time
		"iat": time.Now().Unix(),                     // Issued at time
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
