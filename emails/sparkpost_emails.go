package emails

import (
	"encoding/base64"
	"net/http"
	"os"
	"strings"

	"github.com/news-ai/tabulae/attach"
	"github.com/news-ai/tabulae/models"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	sp "github.com/news-ai/gosparkpost"

	apiModels "github.com/news-ai/api/models"
)

func SendSparkPostEmail(r *http.Request, email models.Email, user apiModels.User, files []models.File) (bool, string, error) {
	c := appengine.NewContext(r)

	apiKey := os.Getenv("SPARKPOST_API_KEY")
	cfg := &sp.Config{
		BaseUrl:    "https://api.sparkpost.com",
		ApiKey:     apiKey,
		ApiVersion: 1,
	}
	var client sp.Client
	err := client.Init(cfg)
	if err != nil {
		return false, "", err
	}

	client.Client = urlfetch.Client(c)

	emailSubject := ""
	if email.Subject == "" {
		emailSubject = "(no subject)"
	} else {
		emailSubject = email.Subject
	}

	from := user.Email
	if user.EmailAlias != "" {
		from = user.EmailAlias
	}
	if email.FromEmail != "" {
		from = email.FromEmail
	}

	tx := &sp.Transmission{
		Recipients: []string{email.To},
	}

	content := sp.Content{
		From:    from,
		HTML:    email.Body,
		Subject: emailSubject,
	}

	headerTo := email.To

	// Attach CC and BCC
	if len(email.CC) > 0 {
		for _, c := range email.CC {
			tx.Recipients = append(tx.Recipients.([]sp.Recipient), sp.Recipient{
				Address: sp.Address{Email: c, HeaderTo: headerTo},
			})
		}
		// add cc header to content
		if content.Headers == nil {
			content.Headers = map[string]string{}
		}
		content.Headers["cc"] = strings.Join(email.CC, ",")
	}

	if len(email.BCC) > 0 {
		for _, b := range email.BCC {
			tx.Recipients = append(tx.Recipients.([]sp.Recipient), sp.Recipient{
				Address: sp.Address{Email: b, HeaderTo: headerTo},
			})
		}
	}

	if len(email.Attachments) > 0 {
		attachments := []sp.Attachment{}

		bytesArray, attachmentType, fileNames, err := attach.GetAttachmentsForEmail(r, email, files)
		log.Infof(c, "%v", bytesArray)
		if err == nil {
			for i := 0; i < len(bytesArray); i++ {
				str := base64.StdEncoding.EncodeToString(bytesArray[i])

				singleAttachment := sp.Attachment{}
				singleAttachment.Filename = fileNames[i]
				singleAttachment.MIMEType = attachmentType[i]
				singleAttachment.B64Data = str

				log.Infof(c, "%v", singleAttachment)

				attachments = append(attachments, singleAttachment)
			}
		}

		content.Attachments = attachments
	}

	tx.Content = content

	id, _, err := client.Send(tx)
	if err != nil {
		return false, "", err
	}

	return true, id, nil
}
