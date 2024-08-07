package api

import (
	"github.com/antibomberman/mego-api/pkg/response"
	pb "github.com/antibomberman/mego-protos/gen/go/auth"
	"net/http"
)

func (s *Server) AuthLoginSendCode(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		response.Fail(w, http.StatusBadRequest, "Поле email обязательно для заполнения")
		return
	}
	_, err := s.authClient.LoginByEmailSendCode(r.Context(), &pb.LoginByEmailSendCodeRequest{
		Email: email,
	})
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, "Код отправлен на вашу электронную почту", nil)
}

func (s *Server) AuthLogin(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	if email == "" {
		response.Fail(w, http.StatusBadRequest, "Поле email обязательно для заполнения")
		return
	}
	code := r.FormValue("code")
	if code == "" {
		response.Fail(w, http.StatusBadRequest, "Поле код обязательно для заполнения")
		return
	}
	body, err := s.authClient.LoginByEmail(r.Context(), &pb.LoginByEmailRequest{
		Email: email,
		Code:  code,
	})
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, "Авторизация успешна", body)

}
