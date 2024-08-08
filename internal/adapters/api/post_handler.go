package api

import (
	"github.com/antibomberman/mego-api/pkg/response"
	postPb "github.com/antibomberman/mego-protos/gen/go/post"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (s *Server) PostList(w http.ResponseWriter, r *http.Request) {
	pageSize, err := strconv.Atoi(chi.URLParam(r, "page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}
	pageToken := chi.URLParam(r, "page_token")
	search := chi.URLParam(r, "search")
	sort, err := strconv.Atoi(chi.URLParam(r, "sort"))
	if err != nil {
		sort = 0
	}
	posts, err := s.postClient.FindPost(r.Context(), &postPb.FindPostRequest{
		PageSize:  int32(pageSize),
		PageToken: pageToken,
		Search:    search,
		SortOrder: postPb.SortOrder(sort),
	})
	if err != nil {
		response.Fail(w, 400, "Ошибка получения постов")
		return
	}
	response.Success(w, 200, "Посты успешно получены", posts)

}
func (s *Server) PostMyList(w http.ResponseWriter, r *http.Request) {
	pageSize, err := strconv.Atoi(chi.URLParam(r, "page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}
	pageToken := chi.URLParam(r, "page_token")
	search := chi.URLParam(r, "search")
	sort, err := strconv.Atoi(chi.URLParam(r, "sort"))
	if err != nil {
		sort = 0
	}
	posts, err := s.postClient.GetByAuthor(r.Context(), &postPb.GetByAuthorRequest{
		AuthorId:  r.Context().Value("user_id").(string),
		PageSize:  int32(pageSize),
		PageToken: pageToken,
		Search:    search,
		SortOrder: postPb.SortOrder(sort),
	})
	if err != nil {
		response.Fail(w, 400, "Ошибка получения постов")
		return
	}
	response.Success(w, 200, "Посты успешно получены", posts)
}

func (s *Server) PostCreate(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) PostShow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Fail(w, http.StatusBadRequest, "Не указан идентификатор поста")
		return
	}
	postDetail, err := s.postClient.GetById(r.Context(), &postPb.GetByIdRequest{Id: id})
	if err != nil {
		response.Fail(w, 400, "Ошибка получения поста")
		return
	}
	response.Success(w, 200, "Пост успешно получен", postDetail)
}

func (s *Server) PostUpdate(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) PostDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Fail(w, http.StatusBadRequest, "Не указан идентификатор поста")
		return
	}
	_, err := s.postClient.DeletePost(r.Context(), &postPb.DeletePostRequest{Id: id, AuthorId: r.Context().Value("user_id").(string)})
	if err != nil {
		response.Fail(w, 400, "Ошибка удаления поста")
		return
	}
	response.Success(w, 200, "Пост успешно удален", nil)
}

func (s *Server) PostHide(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Fail(w, http.StatusBadRequest, "Не указан идентификатор поста")
		return
	}
	_, err := s.postClient.HidePost(r.Context(), &postPb.HidePostRequest{Id: id, AuthorId: r.Context().Value("user_id").(string)})
	if err != nil {
		response.Fail(w, 400, "Ошибка скрытия поста")
		return
	}
	response.Success(w, 200, "Пост успешно скрыт", nil)
}
