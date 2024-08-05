package api

import (
	"github.com/antibomberman/mego-api/pkg/response"
	userPb "github.com/antibomberman/mego-protos/gen/go/user"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (s *Server) UserRegister(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) UserLoginSendCode(w http.ResponseWriter, r *http.Request) {

}
func (s *Server) UserLogin(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) UserShow(w http.ResponseWriter, r *http.Request) {
	userDetail, err := s.userClient.User.GetById(r.Context(), &userPb.Id{Id: chi.URLParam(r, "id")})
	if err != nil {
		response.Fail(w, 400, "Ошибка получения пользовательской информации")
		return
	}
	response.Success(w, 200, "Пользователь успешно получен", userDetail)
}
func (s *Server) UserUpdate(w http.ResponseWriter, r *http.Request) {

}
func (s *Server) UserMe(w http.ResponseWriter, r *http.Request) {

}
