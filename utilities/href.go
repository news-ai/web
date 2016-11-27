package utilities

import (
	"bytes"

	"golang.org/x/net/html"
)

func AppendHrefWithLink(body string, emailId string, hrefAppend string) string {
	r := bytes.NewReader([]byte(body))
	z := html.NewTokenizer(r)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return string(z.Raw())
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			for _, a := range t.Attr {
				if a.Key == "href" {
					a.Val = hrefAppend + "/?id=" + emailId + "&url=" + a.Val
				}
			}

		}

	}

	return string(z.Raw())
}
