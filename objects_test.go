package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCreateObject(t *testing.T) {
	ms := &MockStore{}
	svc := NewObjectsService(ms)
	type Test struct {
		input *CreateObjectRequest
	}
	createObjectCases := map[string]Test{
		"type is empty":        {input: &CreateObjectRequest{Type: "", Description: "D"}},
		"desc is empty":        {input: &CreateObjectRequest{Type: "D", Description: ""}},
		"type & desc is empty": {input: &CreateObjectRequest{Type: "", Description: ""}},
	}

	for k, v := range createObjectCases {
		t.Run(fmt.Sprintf("should return an error if %s", k), func(t *testing.T) {
			b, err := json.Marshal(v.input)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/objects", bytes.NewReader(b))
			r := mux.NewRouter()
			r.HandleFunc("/objects", svc.handleCreateObject)
			r.ServeHTTP(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Error("invalid status code, it must be 400")
			}
		})
	}

}

func TestGetObject(t *testing.T) {
	ms := &MockStore{}
	svc := NewObjectsService(ms)
	type Test struct {
		input  string
		output int
	}
	getObjectCases := map[string]Test{
		"OK":              {input: "123-123", output: http.StatusOK},
		"Not Found":       {input: "", output: http.StatusNotFound},
		"500":             {input: "mock-500", output: http.StatusInternalServerError},
		"No Object Found": {input: "mock-no-object", output: http.StatusOK},
	}
	for k, v := range getObjectCases {
		t.Run(fmt.Sprintf("Should return %s when ref id is %s", k, v.input), func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/objects/%s", v.input), nil)
			r := mux.NewRouter()
			r.HandleFunc("/objects/{id}", svc.handleGetObject)
			r.ServeHTTP(rr, req)

			if rr.Code != v.output {
				t.Errorf("invalid status code, want %d, got %d", v.output, rr.Code)
			}
		})

	}
}
