package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/antibomberman/mego-api/pkg/response"
	userPb "github.com/antibomberman/mego-protos/gen/go/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

func (s *Server) UserShow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Fail(w, "Не указан идентификатор пользователя")
		return
	}
	userDetail, err := s.userClient.GetById(r.Context(), &userPb.Id{Id: chi.URLParam(r, "id")})
	if err != nil {
		response.Fail(w, "Ошибка получения пользовательской информации")
		return
	}
	response.Success(w, "Пользователь успешно получен", userDetail)
}
func (s *Server) UserMe(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id").(string)
	if id == "" {
		response.Fail(w, "Не указан идентификатор пользователя")
		return
	}
	userDetail, err := s.userClient.GetById(r.Context(), &userPb.Id{Id: id})
	fmt.Println(userDetail)
	if err != nil {
		response.Fail(w, "Ошибка получения пользовательской информации")
		return
	}
	response.Success(w, "Пользователь успешно получен", userDetail)
}
func (s *Server) UserUpdateProfile(w http.ResponseWriter, r *http.Request) {
	type AvatarRequest struct {
		FileName string `json:"file_name"`
		Data     string `json:"data"` // base64 encoded
	}
	type Data struct {
		FirstName  string         `json:"first_name" valid:"required"`
		MiddleName string         `json:"middle_name"`
		LastName   string         `json:"last_name"`
		About      string         `json:"about"`
		Avatar     *AvatarRequest `json:"avatar"`
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
	pbAvatar := &userPb.NewAvatar{}

	if data.Avatar != nil {
		fileBytes, err := base64.StdEncoding.DecodeString(data.Avatar.Data)
		if err != nil {
			response.Fail(w, "Ошибка при декодировании файла аватара")
			return
		}
		contentType := http.DetectContentType(fileBytes)
		pbAvatar.Data = fileBytes
		pbAvatar.ContentType = contentType
		pbAvatar.FileName = data.Avatar.FileName
	}

	if pbAvatar.Data == nil {
		pbAvatar = nil
	}

	profile, err := s.userClient.UpdateProfile(r.Context(), &userPb.UpdateProfileRequest{
		Id:         r.Context().Value("user_id").(string),
		FirstName:  data.FirstName,
		MiddleName: data.MiddleName,
		LastName:   data.LastName,
		About:      data.About,
		Avatar:     pbAvatar,
	})
	if err != nil {
		response.Fail(w, "Ошибка при изменении профиля")
		return
	}
	response.Success(w, "Профиль успешно изменен", profile)
}
func (s *Server) UserUpdateTheme(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Theme string `json:"theme" valid:"required,oneof=light dark system"`
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

	_, err = s.userClient.UpdateTheme(r.Context(), &userPb.UpdateThemeRequest{
		Id:    r.Context().Value("user_id").(string),
		Theme: data.Theme,
	})
	if err != nil {
		log.Printf("Error updating user theme: %v\n", err)
		response.Fail(w, "Ошибка при изменении темы")
		return
	}
	response.Success(w, "Тема успешно изменен", nil)
}
func (s *Server) UserUpdateLang(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Lang string `json:"lang" valid:"required,oneof=ru en kz"`
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
	_, err = s.userClient.UpdateLang(r.Context(), &userPb.UpdateLangRequest{
		Id:   r.Context().Value("user_id").(string),
		Lang: data.Lang,
	})
	if err != nil {
		log.Printf("Error updating user language: %v\n", err)
		response.Fail(w, "Ошибка при изменении языка")
		return
	}
	response.Success(w, "Язык успешно изменен", nil)

}
func (s *Server) UserUpdateEmail(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Code string `json:"code" valid:"required"`
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
	_, err = s.userClient.UpdateEmail(r.Context(), &userPb.UpdateEmailRequest{
		UserId: r.Context().Value("user_id").(string),
		Code:   data.Code,
	})
	if err != nil {
		log.Printf("Error updating user email: %v\n", err)
		response.Fail(w, "Ошибка при изменении электронной почты")
		return
	}
	response.Success(w, "Электронная почта успешно изменена", nil)
}
func (s *Server) UserUpdateEmailSendCode(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Email string `json:"email" valid:"required,email"`
	}
	var data Data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		response.Fail(w, "Ошибка при десериализации тела запроса")
		return
	}
	validate := validator.New()
	err = validate.Struct(data)

	_, err = s.userClient.UpdateEmailSendCode(r.Context(), &userPb.UpdateEmailSendCodeRequest{
		UserId: r.Context().Value("user_id").(string),
		Email:  data.Email,
	})
	if err != nil {
		log.Printf("Error sending email change code: %v\n", err)
		response.Fail(w, "Ошибка при отправке кода для изменения электронной почты")
		return
	}
	response.Success(w, "Код для изменения электронной почты успешно отправлен", nil)
}
