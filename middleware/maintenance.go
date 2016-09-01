package middleware

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/capability"

	"github.com/news-ai/web/errors"
)

func AppEngineCheck(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := appengine.NewContext(r)
	if !capability.Enabled(ctx, "datastore_v3", "*") {
		w.Header().Set("Content-Type", "application/json")
		errors.ReturnError(w, http.StatusInternalServerError, "Datastore offline", "Please try again later")
		return
	}
	next(w, r)
}
