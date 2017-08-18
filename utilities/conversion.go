package utilities

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

func NormalizeUrlToUsername(hostUrl string, socialNetworkUrl string) string {
	if strings.Contains(hostUrl, "/") {
		if !strings.Contains(hostUrl, socialNetworkUrl) {
			return ""
		}
	}

	hostUrl = strings.Replace(hostUrl, "http:", "", -1)
	hostUrl = strings.Replace(hostUrl, "https:", "", -1)
	hostUrl = strings.Replace(hostUrl, "www.", "", -1)
	hostUrl = strings.Replace(hostUrl, "/", "", -1)
	hostUrl = strings.Replace(hostUrl, socialNetworkUrl, "", -1)
	hostUrl = strings.Replace(hostUrl, "@", "", -1)
	hostUrl = strings.Replace(hostUrl, ".", "", -1)
	return hostUrl
}

func NormalizeUrl(initialUrl string) (string, error) {
	u, err := url.Parse(initialUrl)
	if err != nil {
		return "", err
	}
	urlHost := strings.Replace(u.Host, "www.", "", 1)
	return u.Scheme + "://" + urlHost, nil
}

func GetDomainName(initialUrl string) (string, error) {
	u, err := url.Parse(initialUrl)
	if err != nil {
		return "", err
	}
	urlHost := strings.Split(u.Host, ".")
	if len(urlHost) == 0 {
		return u.Host, nil
	}
	return urlHost[0], nil
}

func UpdateIfNotBlank(initial *string, replace string) {
	if replace != "" {
		*initial = replace
	}

	if replace == "" && *initial != "" {
		*initial = replace
	}
}

func StringIdToInt(id string) (int64, error) {
	currentId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, err
	}
	return currentId, nil
}

func IntIdToString(id int64) string {
	return strconv.FormatInt(id, 10)
}

func ExtractEmailExtension(email string) (string, error) {
	splitEmail := strings.Split(email, "@")
	if len(splitEmail) > 1 {
		return splitEmail[1], nil
	}
	return "", errors.New("Email is invalid")
}

func ExtractNameFromEmail(email string) (string, error) {
	splitEmail := strings.Split(email, ".")
	if len(splitEmail) > 1 {
		return strings.Title(splitEmail[0]), nil
	}
	return "", errors.New("Name is invalid")
}

func StripQueryString(inputUrl string) string {
	u, err := url.Parse(inputUrl)
	if err != nil {
		return inputUrl
	}
	if u.Scheme == "http" {
		u.Scheme = "https"
	}
	if u.Host != "" && !strings.Contains(u.Host, "www.") {
		u.Host = "www." + u.Host
	}
	u.RawQuery = ""
	return u.String()
}

func StripQueryStringForWebsite(inputUrl string) string {
	u, err := url.Parse(inputUrl)
	if err != nil {
		return inputUrl
	}
	if u.Host != "" {
		u.Host = strings.ToLower(u.Host)
	}
	return u.String()
}
