package Auth

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// regular authentication
func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		header = header[len("Bearer "):]
		godotenv.Load(".env")
		token, err := jwt.Parse(header, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_G")), nil
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "JWT validation error.",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("id", claims["id"])
			c.Next()
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "JWT invalid.",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
	}
}

// google authentication
func GInit(c *gin.Context) {
	godotenv.Load(".env")
	var GoogleAuth *oauth2.Config = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SEC"),
		RedirectURL:  "http://localhost:8080/user/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	url := GoogleAuth.AuthCodeURL("state", oauth2.AccessTypeOnline)
	c.Redirect(302, url)
}

func GCallback(c *gin.Context) string {
	godotenv.Load(".env")
	var GoogleAuth *oauth2.Config = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SEC"),
		RedirectURL:  "http://intern-3.ap-southeast-1.elasticbeanstalk.com/user/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	code := c.Query("code")
	tok, err := GoogleAuth.Exchange(context.Background(), code)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		// return
	}
	client := GoogleAuth.Client(context.Background(), tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		// return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("body", string(body))

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "")
	if err != nil {
		log.Println("JSON parse error: ", err)
		// return
	}
	println(prettyJSON.String())
	return prettyJSON.String()
}
