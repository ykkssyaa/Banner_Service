package server

import (
	"BannerService/internal/models"
	sErr "BannerService/pkg/serverError"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func UrlArgToInt32(w http.ResponseWriter, r *http.Request, arg string) (int32, error) {

	argStr := r.URL.Query().Get(arg)

	if argStr != "" {
		argInt, err := strconv.Atoi(argStr)

		if err != nil {
			sErr.ErrorResponse(w, sErr.ServerError{
				Message:    "Bad Request " + err.Error(),
				StatusCode: http.StatusBadRequest,
			})
			return 0, err
		}

		return int32(argInt), nil
	}

	return 0, nil
}

func (s *HttpServer) BannerGet(w http.ResponseWriter, r *http.Request) {

	featureId, err := UrlArgToInt32(w, r, "feature_id")
	tagId, err := UrlArgToInt32(w, r, "tag_id")
	limit, err := UrlArgToInt32(w, r, "limit")
	offset, err := UrlArgToInt32(w, r, "offset")

	if err != nil {
		return
	}

	banners, err := s.services.GetBanner(tagId, featureId, limit, offset)
	if err != nil {
		sErr.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(banners)
}

func (s *HttpServer) BannerIdDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *HttpServer) BannerIdPatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *HttpServer) BannerPost(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		sErr.ErrorResponse(w, sErr.ServerError{
			Message:    "Content Type is not application/json",
			StatusCode: http.StatusUnsupportedMediaType,
		})
		return
	}

	var banner models.Banner
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&banner)

	if err != nil {
		if errors.As(err, &unmarshalErr) {
			sErr.ErrorResponse(w, sErr.ServerError{
				Message:    "Bad Request. Wrong Type provided for field " + unmarshalErr.Field,
				StatusCode: http.StatusBadRequest,
			})
		} else {
			sErr.ErrorResponse(w, sErr.ServerError{
				Message:    "Bad Request " + err.Error(),
				StatusCode: http.StatusBadRequest,
			})
		}
		return
	}

	id, err := s.services.Banner.CreateBanner(banner)
	if err != nil {
		sErr.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(id)
}

func (s *HttpServer) UserBannerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
