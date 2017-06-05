package emails

import (
	"net/http"
	"strings"

	"github.com/news-ai/tabulae/models"

	apiModels "github.com/news-ai/api/models"

	"google.golang.org/appengine"

	Gmail "github.com/news-ai/go-gmail"
)

func SendGmailEmail(r *http.Request, user apiModels.User, email models.Email, files []models.File) (string, string, error) {
	c := appengine.NewContext(r)

	userFullName := strings.Join([]string{user.FirstName, user.LastName}, " ")
	emailFullName := strings.Join([]string{email.FirstName, email.LastName}, " ")

	from := userFullName + "<" + user.Email + ">"

	if user.EmailAlias != "" {
		from = userFullName + "<" + user.EmailAlias + ">"
	}

	to := emailFullName + "<" + email.To + ">"

	gmail := Gmail.Gmail{}
	gmail.AccessToken = user.AccessToken

	if len(email.Attachments) > 0 && len(files) > 0 {
		return gmail.SendEmailWithAttachments(r, c, from, to, email.Subject, email.Body, email, files)
	}

	return gmail.SendEmail(c, from, to, email.Subject, email.Body, email)
}
