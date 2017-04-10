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

func GenerateEmail(r *http.Request, user models.User, email models.Email, files []models.File) (string, error) {
	userFullName := strings.Join([]string{user.FirstName, user.LastName}, " ")
	emailFullName := strings.Join([]string{email.FirstName, email.LastName}, " ")

	from := userFullName + " <" + user.SMTPUsername + ">"
	to := emailFullName + " <" + email.To + ">"

	if len(email.Attachments) > 0 && len(files) > 0 {
		return GenerateEmailWithAttachments(r, from, to, email.Subject, email.Body, email, files)
	}

	return GenerateEmailWithoutAttachments(from, to, email.Subject, email.Body, email)
}

func GenerateEmailWithoutAttachments(from string, to string, subject string, body string, email models.Email) (string, error) {
	CC := ""
	BCC := ""

	if len(email.CC) > 0 {
		CC = "Cc: " + strings.Join(email.CC, ",") + "\r\n"
	}

	if len(email.BCC) > 0 {
		BCC = "Bcc: " + strings.Join(email.BCC, ",") + "\r\n"
	}

	temp := "From: " + from + "\r\n" +
		CC +
		BCC +
		"reply-to: " + from + "\r\n" +
		"Content-type: text/html;charset=iso-8859-1\r\n" +
		"MIME-Version: 1.0\r\n" +
		"To:  " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body

	return temp, nil
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

	return string(temp[:]), nil
}

func SendSMTPEmail(servername string, email string, password string, to string, subject string, body string) error {
	headers := make(map[string]string)
	headers["From"] = email
	headers["To"] = to
	headers["Subject"] = subject

	host, _, _ := net.SplitHostPort(servername)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	if servername == "smtp.office365.com:587" {
		auth := LoginAuth(email, password)

		conn, err := net.Dial("tcp", servername)
		if err != nil {
			return err
		}

		smtpC, err := smtp.NewClient(conn, host)
		if err != nil {
			return err
		}

		smtpC.StartTLS(tlsconfig)

		if err = smtpC.Auth(auth); err != nil {
			return err
		}

		// To && From
		if err = smtpC.Mail(email); err != nil {
			return err
		}

		if err = smtpC.Rcpt(to); err != nil {
			return err
		}

		// Data
		w, err := smtpC.Data()
		if err != nil {
			return err
		}

		_, err = w.Write([]byte(body))
		if err != nil {
			return err
		}

		err = w.Close()
		if err != nil {
			return err
		}

		smtpC.Quit()
	} else {
		auth := smtp.PlainAuth("", email, password, host)

		// Here is the key, you need to call tls.Dial instead of smtp.Dial
		// for smtp servers running on 465 that require an ssl connection
		// from the very beginning (no starttls)
		conn, err := tls.Dial("tcp", servername, tlsconfig)
		if err != nil {
			return err
		}

		c, err := smtp.NewClient(conn, host)
		if err != nil {
			return err
		}

		// Auth
		if err = c.Auth(auth); err != nil {
			return err
		}

		// To && From
		if err = c.Mail(email); err != nil {
			return err
		}

		if err = c.Rcpt(to); err != nil {
			return err
		}

		// Data
		w, err := c.Data()
		if err != nil {
			return err
		}

		_, err = w.Write([]byte(body))
		if err != nil {
			return err
		}

		err = w.Close()
		if err != nil {
			return err
		}

		c.Quit()
	}

	return nil
}

func VerifySMTP(servername string, email string, password string) error {
	host, _, _ := net.SplitHostPort(servername)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	if servername == "smtp.office365.com:587" {
		auth := LoginAuth(email, password)

		conn, err := net.Dial("tcp", servername)
		if err != nil {
			return errors.New("Error Dialing")
		}

		smtpC, err := smtp.NewClient(conn, host)
		if err != nil {
			return errors.New("Could not connect to your SMTP host")
		}

		smtpC.StartTLS(tlsconfig)

		if err = smtpC.Auth(auth); err != nil {
			return errors.New("Your email or password is invalid " + err.Error())
		}
	} else {
		auth := smtp.PlainAuth("", email, password, host)

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
	}

	return nil
}
