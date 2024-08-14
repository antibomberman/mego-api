package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/antibomberman/mego-api/pkg/response"
	postPb "github.com/antibomberman/mego-protos/gen/go/post"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
)

func (s *Server) PostList(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query().Get("page_size"))
	if r.URL.Query().Get("page_size") == "" {
		response.Fail(w, "Не указан размер страницы")
		return
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		response.Fail(w, "Неверный размер страницы")
		return
	}
	pageToken := r.URL.Query().Get("page_token")
	search := r.URL.Query().Get("search")
	sortStr := r.URL.Query().Get("sort")
	sort := 0
	if sortStr != "" {
		if sortStr != "0" && sortStr != "1" {
			response.Fail(w, "Неверный порядок сорти")
		} else {
			sort, _ = strconv.Atoi(sortStr)
		}
	}
	posts, err := s.postClient.FindPost(r.Context(), &postPb.FindPostRequest{
		PageSize:  int32(pageSize),
		PageToken: pageToken,
		Search:    search,
		SortOrder: postPb.SortOrder(sort),
	})
	if err != nil {
		response.Fail(w, "Ошибка получения постов")
		return
	}
	fmt.Printf("Count: %d\n", len(posts.Posts))
	response.Success(w, "Посты успешно получены", posts)

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
		response.Fail(w, "Ошибка получения постов")
		return
	}
	response.Success(w, "Посты успешно получены", posts)
}

func (s *Server) PostCreate(w http.ResponseWriter, r *http.Request) {
	type File struct {
		FileName    string `json:"file_name"`
		Data        string `json:"data"` // base64 encoded
		ContentType string `json:"content_type"`
	}
	type ContentItem struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Files   []File `json:"files"`
	}
	type RequestBody struct {
		Title    string        `json:"title"`
		Type     int32         `json:"type"`
		Contents []ContentItem `json:"contents"`
	}
	var reqBody RequestBody
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Fail(w, "Ошибка при чтении тела запроса: "+err.Error())
		return
	}

	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		response.Fail(w, "Ошибка при разборе JSON: "+err.Error())
		return
	}
	pbData := &postPb.CreatePostRequest{
		Title:    reqBody.Title,
		AuthorId: r.Context().Value("user_id").(string),
		Type:     postPb.PostType(reqBody.Type),
		Contents: make([]*postPb.PostContentCreateOrUpdate, len(reqBody.Contents)),
	}
	for i, item := range reqBody.Contents {
		pbData.Contents[i] = &postPb.PostContentCreateOrUpdate{
			Title:   item.Title,
			Content: item.Content,
			Files:   make([]*postPb.PostContentFileCreateOrUpdate, len(item.Files)),
		}
		for j, file := range item.Files {
			fileBytes, err := base64.StdEncoding.DecodeString(file.Data)
			if err != nil {
				fmt.Printf("Ошибка при декодировании файла: %v\n", err)
				continue
			}
			contentType := http.DetectContentType(fileBytes)
			file.ContentType = contentType
			pbData.Contents[i].Files[j] = &postPb.PostContentFileCreateOrUpdate{
				FileName:    file.FileName,
				ContentType: contentType,
				Data:        fileBytes,
			}
		}
	}

	createdPost, err := s.postClient.CreatePost(r.Context(), pbData)
	if err != nil {
		response.Fail(w, "Ошибка создания поста: "+err.Error())
		return
	}

	response.Success(w, "Данные успешно обработаны", createdPost)
}

func (s *Server) PostShow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Fail(w, "Не указан идентификатор поста")
		return
	}
	postDetail, err := s.postClient.GetById(r.Context(), &postPb.GetByIdRequest{Id: id})
	if err != nil {
		response.Fail(w, "Ошибка получения поста")
		return
	}
	response.Success(w, "Пост успешно получен", postDetail)
}

func (s *Server) PostUpdate(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) PostDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Fail(w, "Не указан идентификатор поста")
		return
	}
	_, err := s.postClient.DeletePost(r.Context(), &postPb.DeletePostRequest{Id: id, AuthorId: r.Context().Value("user_id").(string)})
	if err != nil {
		response.Fail(w, "Ошибка удаления поста")
		return
	}
	response.Success(w, "Пост успешно удален", nil)
}

func (s *Server) PostHide(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Fail(w, "Не указан идентификатор поста")
		return
	}
	_, err := s.postClient.HidePost(r.Context(), &postPb.HidePostRequest{Id: id, AuthorId: r.Context().Value("user_id").(string)})
	if err != nil {
		response.Fail(w, "Ошибка скрытия поста")
		return
	}
	response.Success(w, "Пост успешно скрыт", nil)
}
