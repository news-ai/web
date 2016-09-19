package utilities

import (
	"errors"
	"io"

	"golang.org/x/net/context"
	"golang.org/x/net/html"

	"google.golang.org/appengine/urlfetch"
)

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, bool, error) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true, nil
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok, _ := traverse(c)
		if ok {
			return result, ok, nil
		}
	}

	return "", false, errors.New("Could not find a title")
}

func getHtmlTitle(r io.Reader) (string, bool, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", false, err
	}

	return traverse(doc)
}

func GetTitleFromHTTPRequest(c context.Context, url string) (string, error) {
	client := urlfetch.Client(c)
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	title, ok, err := getHtmlTitle(resp.Body)
	if !ok {
		return "", err
	}

	return title, nil
}
