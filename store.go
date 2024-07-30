package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Store interface {
	CreateUser(*User) error
	GetUserByID(int64) (*User, error)
	GetUserByUserName(string) (*User, error)
	CreateObject(*Object) error
	DeleteAllObjects() error
	GetObjectByRef(string) (*Object, error)
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateObject(obj *Object) error {
	insertQuery := `
	INSERT INTO xobjects 
	(ref, type, description, created_at) VALUES 
	($1, $2, $3, $4)
	;
	`
	_, err := s.db.Query(insertQuery, obj.Ref, obj.Type, obj.Description, obj.CreatedAt)
	if err != nil {
		return err
	}
	log.Printf("Inserted to xobjects table a record with ref: %v", obj.Ref)
	return nil
}

func (s *Storage) DeleteAllObjects() error {
	if _, err := s.db.Query("truncate xobjects restart identity cascade;"); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetObjectByRef(refId string) (*Object, error) {
	query := `
	SELECT *
	FROM xobjects
	WHERE ref = $1
	;
	`
	rows, err := s.db.Query(query, refId)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return scanToObject(rows)
	}

	return nil, errors.New("No object found")
}

func (s *Storage) GetUserByID(userId int64) (*User, error) {
	query := `
	SELECT *
	FROM xusers
	WHERE id = $1
	`
	rows, err := s.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return scanToUser(rows)
	}

	return nil, errors.New("No user found")
}

func (s *Storage) GetUserByUserName(userName string) (*User, error) {
	query := `
	SELECT *
	FROM xusers
	WHERE user_name = $1
	`
	rows, err := s.db.Query(query, userName)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return scanToUser(rows)
	}

	return nil, errors.New("No user found")
}

func (s *Storage) CreateUser(usr *User) error {
	insertQuery := `
	insert into xusers 
	(user_name, password, first_name, last_name, token, created_at) values
	($1, $2, $3, $4, $5, $6)
	;
	`
	_, err := s.db.Query(insertQuery,
		usr.UserName, usr.Password, usr.FirstName,
		usr.LastName, usr.Token, usr.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func scanToObject(rows *sql.Rows) (*Object, error) {
	obj := Object{}
	err := rows.Scan(
		&obj.ID,
		&obj.Ref,
		&obj.Type,
		&obj.Description,
		&obj.CreatedAt,
	)
	return &obj, err
}

func scanToUser(rows *sql.Rows) (*User, error) {
	usr := User{}
	err := rows.Scan(
		&usr.ID,
		&usr.UserName,
		&usr.Password,
		&usr.FirstName,
		&usr.LastName,
		&usr.Token,
		&usr.CreatedAt,
	)
	return &usr, err
}

func scanToOwnership(rows *sql.Rows) (*OwnerShip, error) {
	owns := OwnerShip{}
	err := rows.Scan(
		&owns.ID,
		&owns.UserId,
		&owns.ObjId,
	)
	return &owns, err
}
