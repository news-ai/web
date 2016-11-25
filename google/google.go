package google

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/news-ai/tabulae/models"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

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
	return nil
}
