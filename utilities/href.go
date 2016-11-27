package utilities

import (
	"bytes"

	"golang.org/x/net/context"

	"golang.org/x/net/html"

	"google.golang.org/appengine/log"
)

func AppendHrefWithLink(c context.Context, body string, emailId string, hrefAppend string) string {
	r := bytes.NewReader([]byte(body))
	z := html.NewTokenizer(r)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			buf := new(bytes.Buffer)
			buf.ReadFrom(r)
			log.Infof(c, "%v", buf.String())
			return buf.String()
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			for i, a := range t.Attr {
				if a.Key == "href" {
					log.Infof(c, "%v", t.Attr[i])
					t.Attr[i].Val = hrefAppend + "/?id=" + emailId + "&url=" + a.Val
					log.Infof(c, "%v", t.Attr[i])
				}
			}

		}

	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	log.Infof(c, "%v", buf.String())
	return buf.String()
}
