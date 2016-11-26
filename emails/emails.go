package emails

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/news-ai/tabulae/attach"
	"github.com/news-ai/tabulae/models"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailSubstitute struct {
	Name string
	Code string
}

type EmailCancel struct {
	BatchId string `json:"batch_id"`
	Status  string `json:"status"`
}

type EmailBatchResponse struct {
	BatchId string `json:"batch_id"`
}

func getBatchId(r *http.Request) string {
	c := appengine.NewContext(r)
	sendgrid.DefaultClient.HTTPClient = urlfetch.Client(c)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_BATCH_API_KEY"), "/v3/mail/batch", "https://api.sendgrid.com")
	request.Method = "POST"

	response, err := sendgrid.API(request)
	if err != nil {
		log.Errorf(c, "error: %v", err)
		return ""
	}

	var batch EmailBatchResponse

	if err := json.Unmarshal([]byte(response.Body), &batch); err != nil {
		log.Errorf(c, "%v", err)
		return ""
	}

	return batch.BatchId
}

func CancelEmail(r *http.Request, email models.Email) error {
	c := appengine.NewContext(r)
	sendgrid.DefaultClient.HTTPClient = urlfetch.Client(c)

	cancelEmail := EmailCancel{}
	cancelEmail.BatchId = email.BatchId
	cancelEmail.Status = "cancel"

	b, err := json.Marshal(cancelEmail)
	if err != nil {
		log.Errorf(c, "%v", err)
		return err
	}

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_SCH_API_KEY"), "/v3/user/scheduled_sends", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = b

	// Send the actual mail here
	response, err := sendgrid.API(request)
	if err != nil {
		log.Errorf(c, "error: %v", err)
		return err
	}

	if response.StatusCode == 201 || response.StatusCode == 200 {
		return nil
	}

	return errors.New("Invalid response from Sendgrid")
}

// Send an email confirmation to a new user
func SendEmail(r *http.Request, email models.Email, user models.User, files []models.File) (bool, string, string, error) {
	c := appengine.NewContext(r)

	sendgrid.DefaultClient.HTTPClient = urlfetch.Client(c)

	userFullName := strings.Join([]string{user.FirstName, user.LastName}, " ")
	emailFullName := strings.Join([]string{email.FirstName, email.LastName}, " ")

	from := mail.NewEmail(userFullName, user.Email)
	to := mail.NewEmail(emailFullName, email.To)
	content := mail.NewContent("text/html", email.Body)

	m := mail.NewV3Mail()

	// Set from
	m.SetFrom(from)
	m.Content = []*mail.Content{
		content,
	}

	// Adding a personalization for the email
	p := mail.NewPersonalization()
	p.Subject = email.Subject

	// Adding who we are sending the email to
	tos := []*mail.Email{
		to,
	}
	p.AddTos(tos...)

	batchId := ""

	if !email.SendAt.IsZero() {
		var timeInt int
		var unixTime int64
		unixTime = email.SendAt.Unix()
		timeInt = int(unixTime)
		p.SetSendAt(timeInt)

		batchId = getBatchId(r)
		if batchId != "" {
			m.SetBatchID(batchId)
		}
		log.Infof(c, "%v", batchId)
	}

	// Add personalization
	m.AddPersonalizations(p)

	// Add attachments
	if len(email.Attachments) > 0 {
		bytesArray, attachmentType, fileNames, err := attach.GetAttachmentsForEmail(r, email, files)
		log.Infof(c, "%v", bytesArray)
		if err == nil {
			for i := 0; i < len(bytesArray); i++ {
				a := mail.NewAttachment()

				str := base64.StdEncoding.EncodeToString(bytesArray[i])

				a.SetContent(str)
				a.SetType(attachmentType[i])
				a.SetFilename(fileNames[i])
				a.SetDisposition("attachment")

				log.Infof(c, "%v", a)

				m.AddAttachment(a)
			}
		}
	}

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)

	// Send the actual mail here
	response, err := sendgrid.API(request)
	if err != nil {
		log.Errorf(c, "error: %v", err)
		return false, "", "", err
	}

	emailId := ""
	if len(response.Headers["X-Message-Id"]) > 0 {
		emailId = response.Headers["X-Message-Id"][0]
	}
	return true, emailId, batchId, nil
}

// Send an email confirmation to a new user
// Someday convert this to a batch so we can send multiple confirmation emails at once
func SendInternalEmail(r *http.Request, email models.Email, templateId string, subject string, emailSubstitutes []EmailSubstitute, sendAt int) (bool, string, error) {
	c := appengine.NewContext(r)
	sendgrid.DefaultClient.HTTPClient = urlfetch.Client(c)

	m := mail.NewV3Mail()
	m.SetTemplateID(templateId)

	// Default from from a NewsAI account
	from := mail.NewEmail("Abhi from NewsAI", "abhi@newsai.org")
	m.SetFrom(from)

	// Adding a personalization for the email
	p := mail.NewPersonalization()
	p.Subject = subject

	// Adding who we are sending the email to
	emailFullName := ""
	if email.FirstName != "" && email.LastName != "" {
		emailFullName = strings.Join([]string{email.FirstName, email.LastName}, " ")
	}
	tos := []*mail.Email{
		mail.NewEmail(emailFullName, email.To),
	}
	p.AddTos(tos...)

	for i := 0; i < len(emailSubstitutes); i++ {
		p.SetSubstitution(emailSubstitutes[i].Name, emailSubstitutes[i].Code)
	}

	if !email.SendAt.IsZero() {
		var timeInt int
		var unixTime int64
		unixTime = email.SendAt.Unix()
		timeInt = int(unixTime)
		p.SetSendAt(timeInt)
	}

	// Add personalization
	m.AddPersonalizations(p)

	// Send the email
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)

	// Send the actual mail here
	response, err := sendgrid.API(request)
	if err != nil {
		log.Errorf(c, "error: %v", err)
		return false, "", err
	}

	emailId := ""
	if len(response.Headers["X-Message-Id"]) > 0 {
		emailId = response.Headers["X-Message-Id"][0]
	}
	return true, emailId, nil
}
