package main

import (
	"errors"
	"time"
)

type MockStore struct {
}

func (m *MockStore) CreateUser(u *User) error {
	return nil
}

func (m *MockStore) CreateObject(u *Object) error {
	return nil
}

func (m *MockStore) DeleteAllObjects() error {
	return nil
}

func (m *MockStore) GetObjectByRef(ref string) (*Object, error) {
	if ref == "mock-500" {
		return nil, errors.New("error from mock")
	}
	if ref == "mock-no-object" {
		return nil, errors.New("no object found")
	}
	return &Object{
		Ref:         ref,
		Type:        "mocktype",
		Description: "mockdesc",
		CreatedAt:   time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
	}, nil
}
