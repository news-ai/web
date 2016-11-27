package utilities

import (
	"regexp"
	"strings"

	"golang.org/x/net/context"
)

func AppendHrefWithLink(c context.Context, body string, emailId string, hrefAppend string) string {
	re := regexp.MustCompile(`href=\"(.+?)\"`)
	reLinks := regexp.MustCompile(`\"(.+?)\"`)

	htmlALinks := re.FindAllString(body, -1)

	for i := 0; i < len(htmlALinks); i++ {
		original := htmlALinks[i]

		htmlALink := reLinks.FindString(htmlALinks[i])
		htmlALink = strings.Replace(htmlALink, "\"", "", -1)
		newLink := hrefAppend + "?id=" + emailId + "&url=" + htmlALink
		htmlALinks[i] = strings.Replace(htmlALinks[i], htmlALink, newLink, 1)

		body = strings.Replace(body, original, htmlALinks[i], 1)
	}

	return body
}
