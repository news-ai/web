package emails

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"net"
	"net/smtp"
	"strings"

	"github.com/news-ai/tabulae/models"

	apiModels "github.com/news-ai/api/models"
)

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
