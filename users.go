package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	store Store
}

func NewUsersService(s Store) *UsersService {
	return &UsersService{
		store: s,
	}
}

func (s *UsersService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/register", s.handleUserRegister).Methods(http.MethodPost)
	r.HandleFunc("/users/login", s.handleUserLogin).Methods(http.MethodPost)
}

func (s *UsersService) handleUserRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("Execute create user")
	usrReq := &CreateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(usrReq); err != nil {
		log.Printf("Could not create user: %v", err.Error())
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "can not parse request payload"})
		return
	}

	hashpw, err := bcrypt.GenerateFromPassword([]byte(usrReq.PW), 10)
	if err != nil {
		log.Printf("Could not create user: %v", err.Error())
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "can not create user"})
		return
	}

	usr := NewUser(usrReq.FN, usrReq.LN, usrReq.USER, hashpw)

	if err := s.store.CreateUser(usr); err != nil {
		log.Printf("Could not create user: %v\n", err.Error())
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "server can not create user"})
		return
	}

	log.Printf("hash: %s", bcrypt.CompareHashAndPassword(hashpw, []byte(usrReq.PW)))

	WriteJSON(w, http.StatusCreated, SuccessResponse{Message: "succesfully created"})
}

func (s *UsersService) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Execute logging user in")
	loginReq := &LoginUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(loginReq); err != nil {
		log.Printf("Could not login user: %v", err.Error())
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "can not parse request payload"})
		return
	}

	usr, err := s.store.GetUserByUserName(loginReq.USER)

	if err != nil {
		log.Printf("Could not get user info: %v\n", err.Error())
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "can not get user info"})
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword(usr.Password, []byte(loginReq.PW))
	if err != nil {
		log.Printf("could not authenticate user: %v\n", err.Error())
		WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "could not authenticate user"})
		return
	}

	jwt, err := CreateJWTFromUser(usr)
	if err != nil {
		log.Printf("Could not authorize user: %v\n", err.Error())
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "can not authorize user"})
	}
	usr.Token = jwt

	log.Printf("user logged in: %v", usr.UserName)
	log.Printf("user logged token: %v", usr.Token)

	WriteJSON(w, http.StatusOK, LoginSuccess{Token: usr.Token})
}
