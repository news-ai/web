package emails

import (
	"crypto/tls"
	"net"
	"net/smtp"

	"golang.org/x/net/context"

	"google.golang.org/appengine/log"
)

func VerifySMTP(c context.Context, servername string, email string, password string) error {
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
		log.Errorf(c, "%v", err)
		return err
	}

	smtpC, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Errorf(c, "%v", err)
		return err
	}

	if err = smtpC.Auth(auth); err != nil {
		log.Errorf(c, "%v", err)
		return err
	}

	return nil
}
