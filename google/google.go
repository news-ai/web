package google

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/news-ai/tabulae/models"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"golang.org/x/net/context"
)

const BASEURL = "https://www.googleapis.com/"

type User struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	Hd            string `json:"hd"`
}

type RefreshTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	GrantType    string `json:"grant_type"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// Check if access token is valid
func ValidateAccessToken(r *http.Request, user models.User) error {
	c := appengine.NewContext(r)

	client := urlfetch.Client(c)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo?alt=json&access_token=" + user.AccessToken)
	if err != nil {
		log.Errorf(c, "%v", err)
		return err
	}

	// Decode JSON from Google
	decoder := json.NewDecoder(resp.Body)
	var googleUser User
	err = decoder.Decode(&googleUser)
	if err != nil {
		log.Errorf(c, "%v", err)
		return err
	}

	if googleUser.ID != "" {
		return nil
	}

	log.Infof(c, "%v", googleUser)

	return errors.New("Access token expired")
}

func RefreshAccessToken(r *http.Request, user models.User) error {
	c := appengine.NewContext(r)

	contextWithTimeout, _ := context.WithTimeout(c, time.Second*15)
	client := urlfetch.Client(contextWithTimeout)

	if user.RefreshToken == "" {
		return errors.New("User does not have a refresh token")
	}

	refreshtokenRequest := RefreshTokenRequest{}
	refreshtokenRequest.ClientID = os.Getenv("GOOGLEAUTHKEY")
	refreshtokenRequest.ClientSecret = os.Getenv("GOOGLEAUTHSECRET")
	refreshtokenRequest.RefreshToken = user.RefreshToken
	refreshtokenRequest.GrantType = "refresh_token"

	messageJson, err := json.Marshal(refreshtokenRequest)
	if err != nil {
		log.Errorf(c, "%v", err)
		return err
	}

	messageQuery := bytes.NewReader(messageJson)

	URL := BASEURL + "oauth2/v4/token"
	req, _ := http.NewRequest("POST", URL, messageQuery)

	response, err := client.Do(req)
	if err != nil {
		log.Errorf(c, "%v", err)
		return err
	}

	defer response.Body.Close()

	toReturn := RefreshTokenResponse{}
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf(c, "%v", err)
		return err
	}

	err = json.Unmarshal(contents, &toReturn)
	if err != nil {
		log.Errorf(c, "%v", err)
		return err
	}

	log.Infof(c, "%v", response.Body)
	log.Infof(c, "%v", contents)

	timeToAdd := time.Duration(time.Duration(toReturn.ExpiresIn) * time.Second)

	user.AccessToken = toReturn.AccessToken
	user.GoogleExpiresIn = time.Now().Add(timeToAdd)
	user.TokenType = toReturn.TokenType
	user.Save(c)

	return nil
}
