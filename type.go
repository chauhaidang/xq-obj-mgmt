package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int64     `json:"id"`
	UserName  string    `json:"userName"`
	Password  []byte    `json:"password"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
}

type Object struct {
	ID          int64     `json:"id"`
	Ref         string    `json:"ref"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type OwnerShip struct {
	ID     int64 `json:"id"`
	UserId int64 `json:"userId"`
	ObjId  int64 `json:"objId"`
}

type ErrorResponse struct {
	Error string `json:"message"`
}

type CreateUserRequest struct {
	FN   string `json:"firstName"`
	LN   string `json:"lastName"`
	USER string `json:"userName"`
	PW   string `json:"password"`
}

type LoginUserRequest struct {
	USER string `json:"userName"`
	PW   string `json:"password"`
}

type CreateObjectRequest struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type CreateObjectResponse struct {
	Ref string `json:"ref"`
}

type GetObjectResponse struct {
	Ref         string    `json:"ref"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type LoginSuccess struct {
	Token string `json:"token"`
}

func NewUser(firstName string, lastName string, userName string, password []byte) *User {
	return &User{
		UserName:  userName,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: time.Now().UTC(),
	}
}

func NewObject(objType string, objDesc string) *Object {
	return &Object{
		Ref:         uuid.New().String(),
		Type:        objType,
		Description: objDesc,
		CreatedAt:   time.Now().UTC(),
	}
}
