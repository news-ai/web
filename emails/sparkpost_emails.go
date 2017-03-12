package emails

import (
	"net/http"
	"os"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"github.com/news-ai/tabulae/models"

	sp "github.com/news-ai/gosparkpost"
)

func SendSparkPostEmail(r *http.Request, email models.Email, user models.User, files []models.File) (bool, string, error) {
	c := appengine.NewContext(r)

	apiKey := os.Getenv("SPARKPOST_API_KEY")
	cfg := &sp.Config{
		BaseUrl:    "https://api.sparkpost.com",
		ApiKey:     apiKey,
		ApiVersion: 1,
	}
	var client sp.Client
	err := client.Init(cfg)
	if err != nil {
		return false, "", err
	}

	client.Client = urlfetch.Client(c)

	return false, "", nil
}
