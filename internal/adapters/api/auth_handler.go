package api

import (
	"encoding/json"
	"github.com/antibomberman/mego-api/pkg/response"
	pb "github.com/antibomberman/mego-protos/gen/go/auth"
	"io"
	"net/http"
)

func (s *Server) AuthLoginSendCode(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Email string `json:"email"`
	}
	var data Data
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Fail(w, "Ошибка при чтении тела запроса")
		return
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		response.Fail(w, "Ошибка при разборе JSON")
		return
	}
	if data.Email == "" {
		response.Fail(w, "Поле email обязательно для заполнения")
		return
	}
	_, err = s.authClient.LoginByEmailSendCode(r.Context(), &pb.LoginByEmailSendCodeRequest{
		Email: data.Email,
	})
	if err != nil {
		response.Fail(w, err.Error())
		return
	}
	response.Success(w, "Код отправлен на вашу электронную почту", nil)
}

func (s *Server) AuthLogin(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	var data Data
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Fail(w, "Ошибка при чтении тела запроса")
		return
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		response.Fail(w, "Ошибка при разборе JSON")
		return
	}
	if data.Email == "" {
		response.Fail(w, "Поле email обязательно для заполнения")
		return
	}
	if data.Code == "" {
		response.Fail(w, "Поле код обязательно для заполнения")
		return
	}
	clientResponse, err := s.authClient.LoginByEmail(r.Context(), &pb.LoginByEmailRequest{
		Email: data.Email,
		Code:  data.Code,
	})
	if err != nil {
		response.Fail(w, err.Error())
		return
	}
	response.Success(w, "Авторизация успешна", clientResponse)

}
