package server

import (
	"BannerService/internal/models"
	sErr "BannerService/pkg/serverError"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func UrlArgToInt32(w http.ResponseWriter, argStr string) (int32, error) {

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

	featureId, err := UrlArgToInt32(w, r.URL.Query().Get("feature_id"))
	if err != nil {
		return
	}
	tagId, err := UrlArgToInt32(w, r.URL.Query().Get("tag_id"))
	if err != nil {
		return
	}
	limit, err := UrlArgToInt32(w, r.URL.Query().Get("limit"))
	if err != nil {
		return
	}
	offset, err := UrlArgToInt32(w, r.URL.Query().Get("offset"))
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

	id, err := UrlArgToInt32(w, mux.Vars(r)["id"])
	if err != nil {
		return
	}

	if err := s.services.DeleteBanner(id); err != nil {
		sErr.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)
}

func (s *HttpServer) BannerIdPatch(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		sErr.ErrorResponse(w, sErr.ServerError{
			Message:    "Content Type is not application/json",
			StatusCode: http.StatusUnsupportedMediaType,
		})
		return
	}

	id, err := UrlArgToInt32(w, mux.Vars(r)["id"])
	if err != nil {
		return
	}

	var banner models.Banner
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&banner)

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

	banner.Id = id

	if err := s.services.PatchBanner(banner); err != nil {
		sErr.ErrorResponse(w, err)
		return
	}

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
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(id)
}

func (s *HttpServer) UserBannerGet(w http.ResponseWriter, r *http.Request) {

	featureId, err := UrlArgToInt32(w, r.URL.Query().Get("feature_id"))
	if err != nil {
		return
	}
	tagId, err := UrlArgToInt32(w, r.URL.Query().Get("tag_id"))
	if err != nil {
		return
	}

	useLastRevision := r.URL.Query().Get("use_last_revision") == "true"
	role := r.Context().Value("role").(string)

	banner, err := s.services.GetUserBanner(tagId, featureId, role, useLastRevision)
	if err != nil {
		sErr.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(banner)
}

func (s *HttpServer) GetBannerVersions(w http.ResponseWriter, r *http.Request) {

	id, err := UrlArgToInt32(w, mux.Vars(r)["id"])
	if err != nil {
		return
	}

	banners, err := s.services.GetBannerVersions(id)
	if err != nil {
		sErr.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(banners)
}

func (s *HttpServer) SetBannerVersion(w http.ResponseWriter, r *http.Request) {

	id, err := UrlArgToInt32(w, mux.Vars(r)["id"])
	if err != nil {
		return
	}
	version, err := UrlArgToInt32(w, r.URL.Query().Get("version"))
	if err != nil {
		return
	}

	err = s.services.SetBannerVersion(id, version)
	if err != nil {
		sErr.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
