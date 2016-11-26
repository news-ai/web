package emails

import (
	"net/http"
	"strings"

	"github.com/news-ai/tabulae/models"

	"google.golang.org/appengine"

	Gmail "github.com/news-ai/go-gmail"
)

func SendGmailEmail(r *http.Request, user models.User, email models.Email) (string, string, error) {
	c := appengine.NewContext(r)

	userFullName := strings.Join([]string{user.FirstName, user.LastName}, " ")
	emailFullName := strings.Join([]string{email.FirstName, email.LastName}, " ")

	from := userFullName + "<" + user.Email + ">"
	to := emailFullName + "<" + email.To + ">"

	gmail := Gmail.Gmail{}
	gmail.AccessToken = user.AccessToken
	return gmail.SendEmail(c, from, to, email.Subject, email.Body)
}
