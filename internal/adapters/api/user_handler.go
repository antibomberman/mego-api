package api

import (
	"encoding/json"
	"github.com/antibomberman/mego-api/pkg/file"
	"github.com/antibomberman/mego-api/pkg/response"
	userPb "github.com/antibomberman/mego-protos/gen/go/user"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func (s *Server) UserLoginSendCode(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) UserLogin(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) UserShow(w http.ResponseWriter, r *http.Request) {
	userDetail, err := s.userClient.GetById(r.Context(), &userPb.Id{Id: chi.URLParam(r, "id")})
	if err != nil {
		response.Fail(w, 400, "Ошибка получения пользовательской информации")
		return
	}
	response.Success(w, 200, "Пользователь успешно получен", userDetail)
}

func (s *Server) UserUpdate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	type userUpdate struct {
		FirstName  string `json:"first_name"`
		MiddleName string `json:"middle_name"`
		LastName   string `json:"last_name"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		About      string `json:"about"`
		Theme      string `json:"theme"`
		Lang       string `json:"lang"`
	}
	var update userUpdate
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	userUpdateData := &userPb.UpdateUserRequest{
		Id:         id,
		FirstName:  update.FirstName,
		MiddleName: update.MiddleName,
		LastName:   update.LastName,
		Email:      update.Email,
		Phone:      update.Phone,
		About:      update.About,
		Theme:      update.Theme,
		Lang:       update.Lang,
	}

	//check exist file
	fileName, fileData, err := file.GetFile(r, "avatar")
	if err != nil {
		response.Fail(w, http.StatusBadRequest, err.Error())
		return
	}
	userUpdateData.Avatar = &userPb.NewAvatar{
		FileName: fileName,
		Data:     fileData,
	}
	_, err = s.userClient.Update(r.Context(), userUpdateData)

	if err != nil {
		log.Println("Ошибка при изменении пользовательской информации:", err)
		response.Fail(w, 400, "Ошибка при изменении пользовательской информации")
		return
	}

	response.Success(w, 200, "Пользователь успешно изменен", nil)

}

func (s *Server) UserMe(w http.ResponseWriter, r *http.Request) {

}
