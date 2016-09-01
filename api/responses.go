package api

import (
	"net/http"

	gcontext "github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"

	"github.com/news-ai/web/errors"
)

type BaseResponse struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Data     interface{} `json:"data"`
	Included interface{} `json:"included"`
}

type BaseSingleResponse struct {
	Data     interface{} `json:"data"`
	Included interface{} `json:"included"`
}

func BaseResponseHandler(val interface{}, included interface{}, count int, err error, r *http.Request) (BaseResponse, error) {
	response := BaseResponse{}
	response.Data = val
	response.Included = included
	response.Count = count
	response.Next = gcontext.Get(r, "next").(string)
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
