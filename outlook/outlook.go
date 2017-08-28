package outlook

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/news-ai/api/models"

	"golang.org/x/net/context"
)

type OutlookRefreshTokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExpiresOn    int    `json:"ext_expires_in"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	PwdExp       string `json:"pwd_exp"`
	PwdURL       string `json:"pwd_url"`
	IdToken      string `json:"id_token"`
}

// Check if access token is valid
func ValidateAccessToken(c context.Context, user models.User) error {
	req, _ := http.NewRequest("GET", "https://outlook.office.com/api/v2.0/me", nil)
	req.Header.Add("Authorization", "Bearer "+user.OutlookAccessToken)
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("%v", "there was an issue getting your token "+err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return errors.New("Access token expired")
	}

	return nil
}

func RefreshAccessToken(c context.Context, user models.User) (models.User, error) {
	if user.OutlookRefreshToken == "" {
		return user, errors.New("User does not have a refresh token")
	}

	URL := "https://login.microsoftonline.com/common/oauth2/v2.0/token"

	form := url.Values{}
	form.Add("client_id", os.Getenv("OUTLOOKAUTHKEY"))
	form.Add("client_secret", os.Getenv("OUTLOOKAUTHSECRET"))
	form.Add("refresh_token", user.OutlookRefreshToken)
	form.Add("grant_type", "refresh_token")

	response, err := http.PostForm(URL, form)
	if err != nil {
		log.Printf("%v", err)
		return user, err
	}
	defer response.Body.Close()

	// Decode JSON from Google
	decoder := json.NewDecoder(response.Body)
	var refreshtoken OutlookRefreshTokenResponse
	err = decoder.Decode(&refreshtoken)
	if err != nil {
		log.Printf("%v", err)
		return user, err
	}

	if refreshtoken.AccessToken != "" {
		timeToAdd := time.Duration(time.Duration(refreshtoken.ExpiresIn) * time.Second)
		user.OutlookAccessToken = refreshtoken.AccessToken
		user.OutlookExpiresIn = time.Now().Add(timeToAdd)
		user.OutlookTokenType = refreshtoken.TokenType
		if refreshtoken.RefreshToken != "" {
			user.OutlookRefreshToken = refreshtoken.RefreshToken
		}
		return user, nil
	}

	return user, errors.New("Access token not present")
}
