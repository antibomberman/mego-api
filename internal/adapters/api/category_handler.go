package api

import (
	"encoding/base64"
	"encoding/json"
	"github.com/antibomberman/mego-api/pkg/response"
	postPb "github.com/antibomberman/mego-protos/gen/go/post"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
)

func (s *Server) CategoryList(w http.ResponseWriter, r *http.Request) {
	categories, err := s.postClient.FindCategory(r.Context(), &postPb.FindCategoryRequest{})
	if err != nil {
		response.Fail(w, "Ошибка получения")
		return
	}
	response.Success(w, "Категории успешно получены", categories)

}

func (s *Server) CategoryCreate(w http.ResponseWriter, r *http.Request) {
	type File struct {
		FileName string `json:"file_name"`
		Data     string `json:"data"` // base64 encoded
	}
	type RequestBody struct {
		Name string `json:"name"`
		Icon File   `json:"icon"`
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
	file := postPb.FileCreateOrUpdate{}
	if reqBody.Icon.Data != "" {
		fileBytes, err := base64.StdEncoding.DecodeString(reqBody.Icon.Data)

		if err != nil {
			log.Printf("Error decoding image data: %v", err)
			response.Fail(w, "Ошибка при декодировании изображения: "+err.Error())
			return
		}
		file.Data = fileBytes
		file.FileName = reqBody.Icon.FileName
	}
	pbData := &postPb.CreateCategoryRequest{
		Name: reqBody.Name,
		Icon: &file,
	}

	createdPost, err := s.postClient.CreateCategory(r.Context(), pbData)
	if err != nil {
		response.Fail(w, "Ошибка создания поста: "+err.Error())
		return
	}

	response.Success(w, "Данные успешно обработаны", createdPost)
}

func (s *Server) CategoryUpdate(w http.ResponseWriter, r *http.Request) {
	type File struct {
		FileName string `json:"file_name"`
		Data     string `json:"data"` // base64 encoded
	}
	type RequestBody struct {
		Id   string `json:"id" valid:"required"`
		Name string `json:"name" valid:"required"`
		Icon File   `json:"icon"`
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
	file := postPb.FileCreateOrUpdate{}

	if reqBody.Icon.Data != "" {
		fileBytes, err := base64.StdEncoding.DecodeString(reqBody.Icon.Data)

		if err != nil {
			log.Printf("Error decoding image data: %v", err)
			response.Fail(w, "Ошибка при декодировании изображения: "+err.Error())
			return
		}
		file.Data = fileBytes
		file.FileName = reqBody.Icon.FileName
	}

	createdPost, err := s.postClient.UpdateCategory(r.Context(), &postPb.UpdateCategoryRequest{
		Id:   reqBody.Id,
		Name: reqBody.Name,
		Icon: &file,
	})
	if err != nil {
		response.Fail(w, "Ошибка создания поста: "+err.Error())
		return
	}

	response.Success(w, "Данные успешно обработаны", createdPost)
}

func (s *Server) CategoryDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Fail(w, "Не указан идентификатор поста")
		return
	}
	_, err := s.postClient.DeleteCategory(r.Context(), &postPb.DeleteCategoryRequest{Id: id})
	if err != nil {
		response.Fail(w, "Ошибка удаления поста")
		return
	}
	response.Success(w, "Пост успешно удален", nil)
}
