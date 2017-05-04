package api

import (
	"net/http"

	gcontext "github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"

	"github.com/news-ai/web/errors"
)

type BasePagingCursors struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

type BasePagingResponse struct {
	Cursors BasePagingCursors `json:"cursors"`
	Next    string            `json:"next"`
}

type BaseSummaryResponse struct {
	Total int `json:"total"`
}

type BaseResponse struct {
	Count    int                 `json:"count"`
	Data     interface{}         `json:"data"`
	Included interface{}         `json:"included"`
	Paging   BasePagingResponse  `json:"paging"`
	Summary  BaseSummaryResponse `json:"summary"`
}

type BaseSingleResponse struct {
	Data     interface{} `json:"data"`
	Included interface{} `json:"included"`
}

func BaseResponseHandler(val interface{}, included interface{}, count int, total int, err error, r *http.Request) (BaseResponse, error) {
	response := BaseResponse{}
	response.Data = val
	response.Included = included
	response.Count = count

	basePagingResponse := BasePagingResponse{}
	basePagingResponse.Next = gcontext.Get(r, "next").(string)
	response.Paging = basePagingResponse

	baseSummaryResponse := BaseSummaryResponse{}
	baseSummaryResponse.Total = total
	response.Summary = baseSummaryResponse

	return response, err
}

func BaseSingleResponseHandler(val interface{}, included interface{}, err error) (BaseSingleResponse, error) {
	response := BaseSingleResponse{}
	response.Data = val
	response.Included = included
	return response, err
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	errors.ReturnError(w, http.StatusNotFound, "An unknown error occurred while trying to process this request.", "Not Found")
	return
}

// Handler for when there is a key present after /users/<id> route.
func NotFoundHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	errors.ReturnError(w, http.StatusNotFound, "An unknown error occurred while trying to process this request.", "Not Found")
	return
}
