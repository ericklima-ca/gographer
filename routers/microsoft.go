package routers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/ericklima-ca/gographer/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var microsoftConfig *oauth2.Config

func setupMicrosoftLogin() (*oauth2.Config, string) {
	microsoftConfig = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes: []string{
			"https://graph.microsoft.com/User.Read",
		},
		Endpoint: microsoft.AzureADEndpoint(os.Getenv("TENANT_ID")),
	}

	state := xid.New().String()
	return microsoftConfig, state
}

func MicrosoftLogin(c *gin.Context) {
	urlToRedirect := c.Query("source")
	c.SetCookie("sourceURL", urlToRedirect, 0, "", "", true, true)
	config, state := setupMicrosoftLogin()
	redirectURL := config.AuthCodeURL(state)
	c.Redirect(http.StatusSeeOther, redirectURL)
}

func MicrosoftCallback(c *gin.Context) {
	code := c.Query("code")
	sourceURL, _ := c.Cookie("sourceURL")
	token, err := microsoftConfig.Exchange(context.Background(), code)
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	client := microsoftConfig.Client(context.TODO(), token)
	userInfo, err := client.Get("https://graph.microsoft.com/v1.0/me")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer userInfo.Body.Close()

	info, err := ioutil.ReadAll(userInfo.Body)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var user models.User
	err = json.Unmarshal(info, &user)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	newURL, _ := url.Parse(sourceURL)
	query, _ := url.ParseQuery(newURL.RawQuery)
	query.Add("token", getJWTToken(user))
	newURL.RawQuery = query.Encode()
	c.Redirect(http.StatusSeeOther, newURL.String())
}

func getJWTToken(user models.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":          user.Email,
		"displayName":    user.DisplayName,
		"jobTitle":       user.JobTitle,
		"id":             user.Id,
		"mail":           user.Email,
		"officeLocation": user.OfficeLocation,
	})

	secret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, _ := token.SignedString(secret)

	return tokenString
}
