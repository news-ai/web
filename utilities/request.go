package utilities

import (
	"errors"
	"io"
	"net/http"

	"golang.org/x/net/html"
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

func GetTitleFromHTTPRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	title, ok, err := getHtmlTitle(resp.Body)
	if !ok {
		return "", err
	}

	return title, nil
}
