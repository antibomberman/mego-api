package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/antibomberman/mego-api/pkg/response"
	postPb "github.com/antibomberman/mego-protos/gen/go/post"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

func (s *Server) PostList(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		PageSize  int32  `json:"page_size" validate:"required"`
		PageToken string `json:"page_token"`
		Search    string `json:"search"`
		Sort      int    `json:"sort" validate:"gte=0,lte=3"`
	}
	var data Data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		response.Fail(w, "Ошибка при десериализации тела запроса")
		return
	}
	validate := validator.New()
	err = validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			response.Fail(w, fmt.Sprintf("Ошибка валидации поля: %s", err.Field()))
			return
		}
	}

	posts, err := s.postClient.FindPost(r.Context(), &postPb.FindPostRequest{
		PageSize:  data.PageSize,
		PageToken: data.PageToken,
		Search:    data.Search,
		SortOrder: postPb.SortOrder(data.Sort),
	})
	if err != nil {
		response.Fail(w, "Ошибка получения постов")
		return
	}
	fmt.Printf("Count: %d\n", len(posts.Posts))
	response.Success(w, "Посты успешно получены", posts)

}
func (s *Server) PostMyList(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		PageSize  int32  `json:"page_size" validate:"required,gte=10,lte=100"`
		PageToken string `json:"page_token"`
		Search    string `json:"search"`
		Sort      int    `json:"sort" validate:"gte=0,lte=3"`
	}
	var data Data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		response.Fail(w, "Ошибка при десериализации тела запроса")
		return
	}
	validate := validator.New()
	err = validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			response.Fail(w, fmt.Sprintf("Ошибка валидации поля: %s", err.Field()))
			return
		}
	}

	posts, err := s.postClient.GetByAuthor(r.Context(), &postPb.GetByAuthorRequest{
		AuthorId:  r.Context().Value("user_id").(string),
		PageSize:  data.PageSize,
		PageToken: data.PageToken,
		Search:    data.Search,
		SortOrder: postPb.SortOrder(data.Sort),
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
		Contents []ContentItem `json:"contents"`
	}

	var reqBody RequestBody
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Fail(w, "Ошибка при чтении тела запроса")
		return
	}

	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		response.Fail(w, "Ошибка при разборе JSON")
		return
	}

	for i, item := range reqBody.Contents {
		fmt.Printf("Элемент %d:\n", i)
		fmt.Printf("  Заголовок: %s\n", item.Title)
		fmt.Printf("  Содержание: %s\n", item.Content)

		for j, file := range item.Files {
			fmt.Printf("    Файл %d: %s\n", j, file.FileName)
			fileBytes, err := base64.StdEncoding.DecodeString(file.Data)
			if err != nil {
				fmt.Printf("Ошибка при декодировании файла: %v\n", err)
				continue
			}
			fmt.Printf("    Длина файла: %d байт\n", len(fileBytes))

			if err != nil {
				fmt.Printf("Ошибка при декодировании файла: %v\n", err)
				continue
			}

			contentType := http.DetectContentType(fileBytes)
			file.ContentType = contentType

			fmt.Printf("    ContentType: %s\n", contentType)

		}
	}

	response.Success(w, "Данные успешно обработаны", reqBody)
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
