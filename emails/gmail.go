package emails

import (
	"net/http"

	"github.com/news-ai/tabulae/models"

	"google.golang.org/appengine"

	Gmail "github.com/news-ai/go-gmail"
)

func SendGmailEmail(r *http.Request, user models.User, email models.Email) error {
	c := appengine.NewContext(r)

	gmail := Gmail.Gmail{}
	gmail.AccessToken = user.AccessToken
	err := gmail.SendEmail(c, user.Email, email.To, email.Subject, email.Body)

	return err
}
