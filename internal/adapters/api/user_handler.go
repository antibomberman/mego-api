package api

import (
	"fmt"
	"github.com/antibomberman/mego-api/pkg/response"
	userPb "github.com/antibomberman/mego-protos/gen/go/user"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"strings"
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
	pbAvatar := &userPb.NewAvatar{}

	file, header, err := r.FormFile("avatar")
	if err == nil {

		if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
			response.Fail(w, "Недопустимый тип файла для аватара")
			return
		}
		if header.Size > 5*1024*1024 {
			response.Fail(w, "Размер файла превышен")
			return
		}

		avatarData, err := io.ReadAll(file)
		if err != nil {
			response.Fail(w, "Ошибка при чтении файла аватара")
			return
		}
		pbAvatar.FileName = header.Filename
		pbAvatar.Data = avatarData
		pbAvatar.ContentType = header.Header.Get("Content-Type")

		fmt.Println("avatar ContentType : " + pbAvatar.ContentType)
		defer file.Close()
	}

	if pbAvatar.Data == nil {
		pbAvatar = nil
	}

	profile, err := s.userClient.UpdateProfile(r.Context(), &userPb.UpdateProfileRequest{
		Id:         r.Context().Value("user_id").(string),
		FirstName:  r.FormValue("first_name"),
		MiddleName: r.FormValue("middle_name"),
		LastName:   r.FormValue("last_name"),
		About:      r.FormValue("about"),
		Avatar:     pbAvatar,
	})
	if err != nil {
		log.Printf("Error updating user profile: %v\n", err)
		response.Fail(w, "Ошибка при изменении профиля")
		return
	}
	response.Success(w, "Профиль успешно изменен", profile)
}
func (s *Server) UserUpdateTheme(w http.ResponseWriter, r *http.Request) {
	theme := r.FormValue("theme")
	if theme == "" {
		response.Fail(w, "Не указан тема")
		return
	}
	if theme != "light" && theme != "dark" && theme != "system" {
		response.Fail(w, "Недопустимое значение темы")
		return
	}
	_, err := s.userClient.UpdateTheme(r.Context(), &userPb.UpdateThemeRequest{
		Id:    r.Context().Value("user_id").(string),
		Theme: theme,
	})
	if err != nil {
		log.Printf("Error updating user theme: %v\n", err)
		response.Fail(w, "Ошибка при изменении темы")
		return
	}
	response.Success(w, "Тема успешно изменен", nil)
}
func (s *Server) UserUpdateLang(w http.ResponseWriter, r *http.Request) {
	lang := r.FormValue("lang")
	if lang == "" {
		response.Fail(w, "Не указан язык")
		return
	}
	if lang != "ru" && lang != "en" && lang != "kz" {
		response.Fail(w, "Недопустимое значение языка")
		return
	}
	_, err := s.userClient.UpdateLang(r.Context(), &userPb.UpdateLangRequest{
		Id:   r.Context().Value("user_id").(string),
		Lang: lang,
	})
	if err != nil {
		log.Printf("Error updating user language: %v\n", err)
		response.Fail(w, "Ошибка при изменении языка")
		return
	}
	response.Success(w, "Язык успешно изменен", nil)

}
func (s *Server) UserUpdateEmail(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	if code == "" {
		response.Fail(w, "Не указан код")
		return
	}
	_, err := s.userClient.UpdateEmail(r.Context(), &userPb.UpdateEmailRequest{
		UserId: r.Context().Value("user_id").(string),
		Code:   code,
	})
	if err != nil {
		log.Printf("Error updating user email: %v\n", err)
		response.Fail(w, "Ошибка при изменении электронной почты")
		return
	}
	response.Success(w, "Электронная почта успешно изменена", nil)
}
func (s *Server) UserUpdateEmailSendCode(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		response.Fail(w, "Не указана электронная почта")
		return
	}
	_, err := s.userClient.UpdateEmailSendCode(r.Context(), &userPb.UpdateEmailSendCodeRequest{
		UserId: r.Context().Value("user_id").(string),
		Email:  email,
	})
	if err != nil {
		log.Printf("Error sending email change code: %v\n", err)
		response.Fail(w, "Ошибка при отправке кода для изменения электронной почты")
		return
	}
	response.Success(w, "Код для изменения электронной почты успешно отправлен", nil)
}
