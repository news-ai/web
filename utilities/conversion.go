package utilities

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

func NormalizeUrlToUsername(twitter string, url string) string {
	twitter = strings.Replace(twitter, "http:", "", -1)
	twitter = strings.Replace(twitter, "https:", "", -1)
	twitter = strings.Replace(twitter, "www.", "", -1)
	twitter = strings.Replace(twitter, "/", "", -1)
	twitter = strings.Replace(twitter, url, "", -1)
	twitter = strings.Replace(twitter, "@", "", -1)
	twitter = strings.Replace(twitter, ".", "", -1)
	return twitter
}

func NormalizeUrl(initialUrl string) (string, error) {
	u, err := url.Parse(initialUrl)
	if err != nil {
		return "", err
	}
	urlHost := strings.Replace(u.Host, "www.", "", 1)
	return u.Scheme + "://" + urlHost, nil
}

func UpdateIfNotBlank(initial *string, replace string) {
	if replace != "" {
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
