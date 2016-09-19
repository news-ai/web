package utilities_test

import (
	"fmt"

	"google.golang.org/appengine/aetest"

	"github.com/news-ai/web/utilities"
)

func ExampleGetTitleFromHTTPRequest() {
	c, done, err := aetest.NewContext()
	if err != nil {
		fmt.Println(err)
	}
	defer done()

	title, err := utilities.GetTitleFromHTTPRequest(c, "http://nypost.com")

	fmt.Println(title)
	// Output: New York Post
}
