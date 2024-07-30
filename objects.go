package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type ObjectsService struct {
	store Store
}

var errTypeRequired error = errors.New("object type is required")
var errDescriptionRequired error = errors.New("object dedscription is required")

func NewObjectsService(s Store) *ObjectsService {
	return &ObjectsService{
		store: s,
	}
}

func (s *ObjectsService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/objects", WithJWTAuth(s.handleCreateObject, s.store)).Methods("POST")
	r.HandleFunc("/objects", WithJWTAuth(s.handleDeleteObjectsTest, s.store)).Methods("DELETE")
	r.HandleFunc("/objects/{id}", WithJWTAuth(s.handleGetObject, s.store)).Methods("GET")
}

func (s *ObjectsService) handleCreateObject(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}
	defer r.Body.Close()

	objPayload := &CreateObjectRequest{}
	err = json.Unmarshal(body, objPayload)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	if err = validateObjectPayload(objPayload); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	obj := NewObject(objPayload.Type, objPayload.Description)
	err = s.store.CreateObject(obj)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error during creating object"})
		return
	}

	WriteJSON(w, http.StatusCreated, CreateObjectResponse{Ref: obj.Ref})
}

func (s *ObjectsService) handleGetObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	obj, err := s.store.GetObjectByRef(id)
	if err != nil {
		if err.Error() == "no object found" {
			WriteJSON(w, http.StatusOK, ErrorResponse{Error: err.Error()})
			return
		}
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	WriteJSON(w, http.StatusOK, GetObjectResponse{
		obj.Ref,
		obj.Type,
		obj.Description,
		obj.CreatedAt,
	})
}

func (s *ObjectsService) handleDeleteObjectsTest(w http.ResponseWriter, req *http.Request) {
	err := s.store.DeleteAllObjects()
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error during delete all objects"})
		return
	}

	WriteJSON(w, http.StatusCreated, SuccessResponse{Message: "successfully delete all objects"})
}

func validateObjectPayload(req *CreateObjectRequest) error {
	if req.Type == "" {
		return errTypeRequired
	}

	if req.Description == "" {
		return errDescriptionRequired
	}

	return nil
}
