package utilities

import (
	"net/mail"
	"strings"
)

func ValidateEmailFormat(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	emailAddressAndDomain := strings.Split(email, "@")
	if len(emailAddressAndDomain) == 1 {
		return false
	}

	emailUrlExtension := strings.Split(emailAddressAndDomain[1], ".")
	if len(emailUrlExtension) == 1 {
		return false
	}

	return true
}
