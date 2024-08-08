package api

import (
	"github.com/antibomberman/mego-api/pkg/response"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (s *Server) PostCreate(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) PostShow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Fail(w, http.StatusBadRequest, "Не указан идентификатор поста")
		return
	}
}

func (s *Server) PostUpdate(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) PostDelete(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) PostHide(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) PostList(w http.ResponseWriter, r *http.Request) {
}
