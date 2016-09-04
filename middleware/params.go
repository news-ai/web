package middleware

import (
	"net/http"
	"strconv"

	gcontext "github.com/gorilla/context"

	"github.com/news-ai/web/utilities"
)

func GetPagination(r *http.Request) (int, int, string) {
	limit := 20
	offset := 0

	queryLimit := r.URL.Query().Get("limit")
	queryOffset := r.URL.Query().Get("offset")
	queryAfter := r.URL.Query().Get("after")

	// check if query exists
	if len(queryLimit) != 0 {
		limit, _ = strconv.Atoi(queryLimit)
	}

	// check if offset exists
	if len(queryOffset) != 0 {
		offset, _ = strconv.Atoi(queryOffset)
	}

	// Boundary checks
	max_limit := 50
	if limit > max_limit {
		limit = max_limit
	}

	return limit, offset, queryAfter
}

func GetParams(r *http.Request) (string, string, string) {
	url := utilities.StripQueryString(r.URL.String())
	searchQuery := r.URL.Query().Get("q")
	order := r.URL.Query().Get("order")
	return url, searchQuery, order
}

func ConstructNext(r *http.Request, limit int, offset int, query string, order string) string {
	url := r.URL
	q := r.URL.Query()
	q.Set("limit", strconv.Itoa(limit))
	q.Set("offset", strconv.Itoa(offset+limit))

	if query != "" {
		q.Set("q", query)
	}

	if order != "" {
		q.Set("order", order)
	}

	url.RawQuery = q.Encode()
	return url.String()
}

func AttachParameters(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	limit, offset, after := GetPagination(r)
	url, query, order := GetParams(r)
	nextUrl := ConstructNext(r, limit, offset, query, order)
	gcontext.Set(r, "q", query)
	gcontext.Set(r, "url", url)
	gcontext.Set(r, "order", order)
	gcontext.Set(r, "limit", limit)
	gcontext.Set(r, "offset", offset)
	gcontext.Set(r, "after", after)
	gcontext.Set(r, "next", nextUrl)
	next(w, r)
}
