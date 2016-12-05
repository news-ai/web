package emails

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"net"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/news-ai/tabulae/attach"
	"github.com/news-ai/tabulae/models"
)

func GenerateEmail(from string, to string, subject string, body string, email models.Email) (string, error) {
	CC := ""
	BCC := ""

	if len(email.CC) > 0 {
		CC = "Cc: " + strings.Join(email.CC, ",") + "\r\n"
	}

	if len(email.BCC) > 0 {
		BCC = "Bcc: " + strings.Join(email.BCC, ",") + "\r\n"
	}

	temp := []byte("From: " + from + "\r\n" +
		CC +
		BCC +
		"reply-to: " + from + "\r\n" +
		"Content-type: text/html;charset=iso-8859-1\r\n" +
		"MIME-Version: 1.0\r\n" +
		"To:  " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body)

	raw := base64.StdEncoding.EncodeToString(temp)
	raw = strings.Replace(raw, "/", "_", -1)
	raw = strings.Replace(raw, "+", "-", -1)
	raw = strings.Replace(raw, "=", "", -1)

	return raw, nil
}

func GenerateEmailWithAttachments(r *http.Request, from string, to string, subject string, body string, email models.Email, files []models.File) (string, error) {
	nl := "\r\n" // newline
	boundary := "__newsai_tabulae__"

	CC := ""
	BCC := ""

	if len(email.CC) > 0 {
		CC = "Cc: " + strings.Join(email.CC, ",") + nl
	}

	if len(email.BCC) > 0 {
		BCC = "Bcc: " + strings.Join(email.BCC, ",") + nl
	}

	temp := []byte("MIME-Version: 1.0" + nl +
		"To:  " + to + nl +
		CC +
		BCC +
		"From: " + from + nl +
		"reply-to: " + from + nl +
		"Subject: " + subject + nl +

		"Content-Type: multipart/mixed; boundary=\"" + boundary + "\"" + nl + nl +

		// Boundary one is email itself
		"--" + boundary + nl +

		"Content-Type: text/html; charset=UTF-8" + nl +
		"MIME-Version: 1.0" + nl +
		"Content-Transfer-Encoding: base64" + nl + nl +

		// Body itself
		body + nl + nl)

	for i := 0; i < len(files); i++ {
		bytesArray, attachmentType, fileNames, err := attach.GetAttachmentsForEmail(r, email, files)
		if err == nil {
			for i := 0; i < len(bytesArray); i++ {
				str := base64.StdEncoding.EncodeToString(bytesArray[i])

				attachment := []byte(
					"--" + boundary + nl +
						"Content-Type: " + attachmentType[i] + nl +
						"MIME-Version: 1.0" + nl +
						"Content-Disposition: attachment; filename=\"" + fileNames[i] + "\"" + nl +
						"Content-Transfer-Encoding: base64" + nl + nl +
						str + nl + nl,
				)

				temp = append(temp, attachment...)
			}
		}
	}

	finalBoundry := []byte(
		"--" + boundary + "--",
	)

	temp = append(temp, finalBoundry...)

	raw := base64.StdEncoding.EncodeToString(temp)
	raw = strings.Replace(raw, "/", "_", -1)
	raw = strings.Replace(raw, "+", "-", -1)
	raw = strings.Replace(raw, "=", "", -1)

	return raw, nil
}

func VerifySMTP(servername string, email string, password string) error {
	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", email, password, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return errors.New("The servername you have entered for your SMTP connection is invalid")
	}

	smtpC, err := smtp.NewClient(conn, host)
	if err != nil {
		return errors.New("Could not connect to your SMTP host")
	}

	if err = smtpC.Auth(auth); err != nil {
		return errors.New("Your email or password is invalid")
	}

	return nil
}
