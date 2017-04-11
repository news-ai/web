package outlook

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/news-ai/tabulae/models"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"golang.org/x/net/context"
)

type OutlookRefreshTokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	ExpiresOn    string `json:"expires_on"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	PwdExp       string `json:"pwd_exp"`
	PwdURL       string `json:"pwd_url"`
}

// Check if access token is valid
func ValidateAccessToken(r *http.Request, user models.User) error {
	c := appengine.NewContext(r)

	client := urlfetch.Client(c)
	req, _ := http.NewRequest("GET", "https://outlook.office.com/api/v2.0/me", nil)
	req.Header.Add("Authorization", "Bearer "+user.OutlookAccessToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf(c, "%v", "there was an issue getting your token "+err.Error())
		return err
	}

	if resp.StatusCode == 401 {
		return errors.New("Access token expired")
	}

	return nil
}

func RefreshAccessToken(r *http.Request, user models.User) (models.User, error) {
	c := appengine.NewContext(r)

	contextWithTimeout, _ := context.WithTimeout(c, time.Second*15)
	client := urlfetch.Client(contextWithTimeout)

	if user.OutlookRefreshToken == "" {
		return user, errors.New("User does not have a refresh token")
	}

	URL := "https://login.microsoftonline.com/common/oauth2/v2.0/token"

	form := url.Values{}
	form.Add("client_id", os.Getenv("OUTLOOKAUTHKEY"))
	form.Add("client_secret", os.Getenv("OUTLOOKAUTHSECRET"))
	form.Add("refresh_token", user.OutlookRefreshToken)
	form.Add("grant_type", "refresh_token")

	response, err := client.PostForm(URL, form)
	if err != nil {
		log.Errorf(c, "%v", err)
		return user, err
	}

	// Decode JSON from Google
	decoder := json.NewDecoder(response.Body)
	var refreshtoken OutlookRefreshTokenResponse
	err = decoder.Decode(&refreshtoken)
	if err != nil {
		log.Errorf(c, "%v", err)
		return user, err
	}

	if refreshtoken.AccessToken != "" {
		refreshTime, err := strconv.Atoi(refreshtoken.ExpiresIn)
		if err != nil {
			refreshTime = 3600
		}

		timeToAdd := time.Duration(time.Duration(refreshTime) * time.Second)
		user.OutlookAccessToken = refreshtoken.AccessToken
		user.OutlookExpiresIn = time.Now().Add(timeToAdd)
		user.OutlookTokenType = refreshtoken.TokenType
		if refreshtoken.RefreshToken != "" {
			user.OutlookRefreshToken = refreshtoken.RefreshToken
		}
		user.Save(c)
		return user, nil
	}

	return user, errors.New("Access token not present")
}
