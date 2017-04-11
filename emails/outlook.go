package emails

import (
	"net/http"
	"strings"

	"github.com/news-ai/tabulae/models"

	"google.golang.org/appengine"

	Outlook "github.com/news-ai/go-outlook"
)

func SendOutlookEmail(r *http.Request, user models.User, email models.Email, files []models.File) error {
	c := appengine.NewContext(r)

	userFullName := strings.Join([]string{user.FirstName, user.LastName}, " ")
	emailFullName := strings.Join([]string{email.FirstName, email.LastName}, " ")

	from := userFullName + "<" + user.OutlookEmail + ">"
	to := emailFullName + "<" + email.To + ">"

	outlook := Outlook.Outlook{}
	outlook.AccessToken = user.OutlookAccessToken

	if len(email.Attachments) > 0 && len(files) > 0 {
		return outlook.SendEmailWithAttachments(r, c, from, to, email.Subject, email.Body, email, files)
	}

	return outlook.SendEmail(c, from, to, email.Subject, email.Body, email)
}
